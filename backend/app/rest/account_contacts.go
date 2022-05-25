package api

import (
	contacts_import "amifactory.team/sequel/coton-app-backend/app/import"
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

type userAccountContacts struct {
	contactsStore model.ContactsStore

	secureCookiesDisabled bool
	maxImportBodySize     int64

	importFilesDir string

	// TODO temp solution
	tasksLock sync.Mutex
	tasks     map[string]*contacts_import.ContactImportTask
}

func (u *userAccountContacts) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)
	group, _ := model.ContactGroupFromContext(ctx)

	form := NewPageForm(r, 3)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	contacts, total, err := u.contactsStore.ListContacts(ctx, account.ID, group, form.Offset(), form.PageSize, form.SearchPhrase)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
		} else {
			u.getLog(ctx).Errorf("Fail to fetch contacts list - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	if total == 0 && form.HasFilter() {
		httpResponseError(w, http.StatusNotFound, "result__not_found")
		return
	}

	resp := NewPageResp(total, form.PageSize, contacts)
	httpJsonResponse(w, resp)
}

func (u *userAccountContacts) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)
	group, _ := model.ContactGroupFromContext(ctx)

	form, err := NewContactForm(group, r)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	errs := form.Validate()
	if len(errs) > 0 {
		httpResponseErrors(w, http.StatusBadRequest, errs)
		return
	}

	contact := model.NewContact(
		account.ID,
		group.ID,
		*form.Phone,
		form.Vars)
	err = u.contactsStore.AddContact(ctx, contact)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to add contact - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, contact)
}

func (u *userAccountContacts) importPreviewCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)

	form := newContactImportPreviewForm(r, u.maxImportBodySize)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	csvFile, err := form.ImportFile()
	if err != nil {
		httpResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	// 1. Read header and try match row names to exists group vars
	header, err := csvReader.Read()
	if err != nil {
		// TODO check error
		httpResponseError(w, http.StatusBadRequest, "contacts__invalid")
		return
	}

	columnsCount := len(header)

	if columnsCount == 0 {
		httpResponseError(w, http.StatusBadRequest, "contacts__invalid")
		return
	}

	// 2. Read up to next 9 lines
	records := make([][]string, 1)
	records[0] = header

	for lineIdx := 1; lineIdx < 9; lineIdx++ {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			httpResponseError(w, http.StatusBadRequest, "contacts__invalid")
			break
		}

		records = append(records, record)
	}

	// 3. Try to find group variable match by column title
	type column struct {
		Variable *model.ContactVariable `json:"variable"`
	}

	columns := make([]*column, len(header))
	for idx, h := range header {
		colName := model.VariableNameFromTitle(h)
		variable := group.FindVariableByName(colName)
		if variable != nil {
			columns[idx] = &column{
				Variable: variable,
			}
		}
	}

	resp := struct {
		Vars    []*column  `json:"columns"`
		Records [][]string `json:"data"`
	}{
		Vars:    columns,
		Records: records,
	}

	httpJsonResponse(w, resp)
}

func (u *userAccountContacts) importCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)
	group, _ := model.ContactGroupFromContext(ctx)

	u.tasksLock.Lock()
	currentTask, found := u.tasks[group.ID]
	u.tasksLock.Unlock()

	if found && currentTask.Progress().Status != contacts_import.Finished {
		httpResponseError(w, http.StatusForbidden, "import__not_finished")
		return
	}

	form := newContactImportForm(r, u.maxImportBodySize, group.Vars)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	file, err := form.ImportFile()
	if err != nil {
		httpResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	importTask := contacts_import.NewContactImportTask(u.contactsStore, account.ID, group.ID)
	importTask.VarIds = form.VarIds
	importTask.DuplicateOption = form.DuplicateOption
	err = importTask.File(file)

	u.tasksLock.Lock()
	u.tasks[group.ID] = importTask
	u.tasksLock.Unlock()

	go importTask.Run()

	progress := importTask.Progress()
	httpJsonResponse(w, progress)
}

func (u *userAccountContacts) importProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)

	u.tasksLock.Lock()
	currentTask, found := u.tasks[group.ID]
	u.tasksLock.Unlock()

	if found {
		taskProgress := currentTask.Progress()
		switch taskProgress.Status {
		case contacts_import.Scheduled, contacts_import.Processing:
			httpJsonResponse(w, taskProgress)
			return

		case contacts_import.Finished:
			sinceStatusUpdated := time.Now().Sub(taskProgress.StatusUpdatedAt)
			if sinceStatusUpdated.Seconds() > 60.0 {
				httpResponseError(w, http.StatusNotFound, "import_progress__not_found")
			} else {
				httpJsonResponse(w, taskProgress)
			}
			return
		}
	}

	httpResponseError(w, http.StatusNotFound, "import_progress__not_found")
}

func (u *userAccountContacts) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)
	contact, _ := model.ContactFromContext(ctx)

	form, err := NewContactForm(group, r)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	errs := form.Validate()
	if len(errs) > 0 {
		httpResponseErrors(w, http.StatusBadRequest, errs)
		return
	}

	contact.UpdatePhone(form.Phone)
	contact.VarValues = form.Vars

	err = u.contactsStore.UpdateContact(ctx, contact)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to update contact - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, contact)
}

func (u *userAccountContacts) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact, _ := model.ContactFromContext(ctx)
	err := u.contactsStore.DeleteContact(ctx, contact)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to delete contact - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *userAccountContacts) populateContact(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		account, _ := model.AccountFromContext(ctx)

		contactID := chi.URLParam(r, "contactID")
		if len(contactID) == 0 {
			httpResponseError(w, http.StatusBadRequest, "contact__required")
			return
		}

		contact, err := u.contactsStore.FindContact(ctx, account.ID, contactID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusNotFound, "contact__not_found")
			} else {
				u.getLog(ctx).Errorf("Fail to find contact by ID - %v", err)
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(contact.NewContext(ctx)))
	})
}

func (u *userAccountContacts) findGroups(ctx context.Context, accountId string, groupsIds []string) ([]*model.ContactGroupShort, error) {
	if len(groupsIds) > 0 {
		filteredGroupIDs := make([]string, 0)

		for _, groupId := range groupsIds {
			if uuidRegex.MatchString(groupId) {
				filteredGroupIDs = append(filteredGroupIDs, groupId)
			}
		}

		groups, err := u.contactsStore.FindGroups(ctx, accountId, filteredGroupIDs)
		if err != nil {
			return nil, err
		}

		return groups, nil
	} else {
		return make([]*model.ContactGroupShort, 0), nil
	}
}

func (u *userAccountContacts) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "account_contact",
	})
}

//
// Contact import preview form
//
type contactImportPreviewForm struct {
	*BaseForm
	multipartForm *multipart.Form
}

func newContactImportPreviewForm(r *http.Request, maxImportBodySize int64) *contactImportPreviewForm {
	f := &contactImportPreviewForm{
		BaseForm: &BaseForm{},
	}

	err := r.ParseMultipartForm(maxImportBodySize)
	if err != nil {
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			f.AddValidationError("contacts__too_large", err)
		} else {
			f.AddValidationError("body_multipart__invalid", err)
		}
		return f
	}

	f.multipartForm = r.MultipartForm

	importFiles := f.multipartForm.File["contacts"]
	if len(importFiles) == 0 {
		f.AddValidationError("contacts__required", err)
	}

	return f
}

func (f *contactImportPreviewForm) ImportFile() (multipart.File, error) {
	importFiles := f.multipartForm.File["contacts"]
	if len(importFiles) == 0 {
		return nil, errors.New("contacts__required")
	}

	firstFile := importFiles[0]
	return firstFile.Open()
}

//
// Contact import form
//
type contactImportForm struct {
	contactImportPreviewForm
	VarIds          []*string
	DuplicateOption model.ContactsDuplicateOption
}

func newContactImportForm(r *http.Request, maxImportBodySize int64, allowedVars []model.ContactVariable) *contactImportForm {
	f := &contactImportForm{}
	f.contactImportPreviewForm = *newContactImportPreviewForm(r, maxImportBodySize)

	varsValue, ok := r.MultipartForm.Value["variables"]
	if ok && len(varsValue) > 0 {
		firstVarValue := varsValue[0]
		if len(firstVarValue) > 0 {
			varIds := strings.Split(firstVarValue, ",")
			f.VarIds = make([]*string, len(varIds))
			for idx, varId := range varIds {
				cleanedVarId := strings.ToLower(strings.TrimSpace(varId))
				if IsValidUUID(cleanedVarId) {
					f.VarIds[idx] = &cleanedVarId
				} else if cleanedVarId == "null" {
					// do nothing
				} else {
					f.AddValidationError(fmt.Sprintf("variable_%d__invalid", idx), nil)
				}
			}
		} else {
			f.AddValidationError("variables__required", nil)
		}
	} else {
		f.AddValidationError("variables__required", nil)
	}

	dupValue, ok := r.MultipartForm.Value["duplicates"]
	if ok && len(dupValue) > 0 {
		f.DuplicateOption = model.ContactsDuplicateOption(dupValue[0])
	}

	// validate variable IDs
	for idx, varId := range f.VarIds {
		if varId == nil {
			continue
		}
		found := false
		for _, allowedVar := range allowedVars {
			if *varId == allowedVar.ID {
				found = true
				break
			}
		}

		if !found {
			f.AddValidationError(fmt.Sprintf("variable_%d__invalid", idx), nil)
		}
	}

	// validate duplicate options
	found := false
	for _, option := range model.ContactsDuplicateOptions {
		if f.DuplicateOption == option {
			found = true
			break
		}
	}

	if !found {
		f.AddValidationError("duplicates__invalid", nil)
	}

	return f
}
