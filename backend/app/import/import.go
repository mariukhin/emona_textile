package contacts_import

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"encoding/csv"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path"
	"sync"
	"time"
)

type ImportErr struct {
	Line       int
	ErrMessage string
}

func (e ImportErr) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.ErrMessage)
}

func (e ImportErr) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

type TaskStatus string

var (
	Scheduled  = TaskStatus("scheduled")
	Processing = TaskStatus("processing")
	Finished   = TaskStatus("finished")
)

type ContactImportTask struct {
	ID              string
	AccountID       string
	GroupID         string
	DuplicateOption model.ContactsDuplicateOption
	VarIds          []*string

	importFilePath string

	progressLock sync.Mutex
	progress     *ContactImportProgress

	store model.ContactsStore
}

type ContactImportProgress struct {
	Status          TaskStatus `json:"status"`
	StatusUpdatedAt time.Time  `json:"-"`
	LinesProcessed  int        `json:"lines_processed"`
	ContactsAdded   int        `json:"contacts_added"`
	DuplicatesFound int        `json:"duplicates_found"`
	ErrorsCount     int        `json:"phone_errors_count"`
	Errors          []error    `json:"phone_errors,omitempty"`
}

func (p *ContactImportProgress) Update(lines, contacts int) {
	p.LinesProcessed = lines
	p.ContactsAdded = contacts
}

func (p *ContactImportProgress) AddError(e error) {
	p.Errors = append(p.Errors, e)
	p.ErrorsCount++
}

func NewContactImportTask(store model.ContactsStore, accountId, groupId string) *ContactImportTask {
	return &ContactImportTask{
		ID:        uuid.NewV4().String(),
		AccountID: accountId,
		GroupID:   groupId,
		progress: &ContactImportProgress{
			Status:          Scheduled,
			StatusUpdatedAt: time.Now(),
		},
		store: store,
	}
}

func (t *ContactImportTask) File(file multipart.File) error {
	tmpDir := "tmp"
	err := os.MkdirAll(tmpDir, 0770)
	if err != nil {
		return fmt.Errorf("fail to create tmp dir for import file - %v", err)
	}

	tmpFilePath := path.Join(tmpDir, fmt.Sprintf("contacts_%s.csv", t.ID))
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("fail to create tmp import file - %v", err)
	}

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		os.Remove(tmpFilePath)
		return fmt.Errorf("fail to copy input import file to tmp import file - %v", err)
	}

	t.importFilePath = tmpFilePath
	return nil
}

func (t *ContactImportTask) UpdateProgress(lines, contacts int) {
	t.progressLock.Lock()
	t.progress.Update(lines, contacts)
	t.progressLock.Unlock()
}

func (t *ContactImportTask) UpdateStatus(status TaskStatus) {
	t.progressLock.Lock()
	t.progress.Status = status
	t.progress.StatusUpdatedAt = time.Now()
	t.progressLock.Unlock()
}

func (t *ContactImportTask) UpdateDuplicates(count int) {
	t.progressLock.Lock()
	t.progress.DuplicatesFound = count
	t.progressLock.Unlock()
}

func (t *ContactImportTask) AddError(e error) {
	t.progressLock.Lock()
	t.progress.AddError(e)
	t.progressLock.Unlock()
}

func (t *ContactImportTask) AddFatalError(e error) {
	t.progressLock.Lock()
	t.progress.AddError(e)
	t.progress.Status = Finished
	t.progressLock.Unlock()
}

func (t *ContactImportTask) Run() {
	ctx := context.Background()
	log := logger.GetLogger(ctx)
	log = log.WithFields(logrus.Fields{
		"task_id": t.ID,
		"name":    "import_contacts",
		"account": t.AccountID,
		"group":   t.GroupID,
	})

	t.UpdateStatus(Processing)

	defer os.Remove(t.importFilePath)

	group, err := t.store.FindGroup(ctx, t.AccountID, t.GroupID)
	if err != nil {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.Errorf("Fail to find group - %v", err)
		return
	}

	columnsCount := len(t.VarIds)
	varsCount := 0
	vars := make([]*model.ContactVariable, columnsCount)
	for varIdx, varId := range t.VarIds {
		if varId == nil {
			continue
		}

		v := group.FindVariableByID(*varId)
		if v != nil {
			vars[varIdx] = v
			varsCount++
		}
	}

	csvFile, err := os.Open(t.importFilePath)
	if err != nil {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.WithError(err).Errorf("Fail open contacts import csv '%s'", t.importFilePath, err)
		return
	}

	csvReader := csv.NewReader(csvFile)

	header, err := csvReader.Read()
	if err != nil {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.WithError(err).Errorf("Fail to read file header - %v", t.importFilePath)
		return
	}

	actualColumnsCount := len(header)

	if actualColumnsCount == 0 {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.Error("Fail to read contacts csv file - no columns found")
		return
	}

	if columnsCount > actualColumnsCount {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.Errorf("Fail to read contacts csv file - expected columns count %d - actual %d", columnsCount, actualColumnsCount)
		return
	}

	// looks like csv is ok, start reading records

	lineNo := 1

	contacts := make([]*model.Contact, 0)

	for {
		lineNo++

		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			t.AddError(ImportErr{Line: lineNo, ErrMessage: "File read error"})
			log.Errorf("File read error - %v", err)
			break
		}

		values := make(map[string]interface{})
		for idx, cell := range record {
			if idx >= columnsCount {
				// skip cells out of selected scope
				break
			}

			variable := vars[idx]
			if variable == nil {
				// skip cells without selected variable
				continue
			}

			value, err := variable.Parse(cell)
			if err != nil && variable.Required {
				t.AddError(ImportErr{Line: lineNo, ErrMessage: fmt.Sprintf("Fail to parse %s", variable.Name)})
				continue
			}

			values[variable.Name] = value
		}

		phoneVal, found := values["phone"]
		phone, typeOk := phoneVal.(*model.Phone)
		if found && typeOk {
			// TODO
			delete(values, "phone")
			newContact := model.NewContact(t.AccountID, t.GroupID, *phone, values)
			contacts = append(contacts, newContact)
		}

		//time.Sleep(2 * time.Millisecond)

		// TODO
		t.UpdateProgress(lineNo-1, len(contacts))
	}

	if len(contacts) == 0 {
		t.AddFatalError(ImportErr{ErrMessage: "No contacts parsed"})
		return
	}

	duplicatesNumber, err := t.store.FindContactDuplicates(ctx, t.AccountID, t.GroupID, contacts)
	t.UpdateDuplicates(int(duplicatesNumber))

	added, updated, ignored, err := t.store.AddContacts(ctx, t.AccountID, t.GroupID, contacts, t.DuplicateOption)
	if err != nil {
		t.AddFatalError(ImportErr{ErrMessage: "Unexpected error"})
		log.Errorf("Fail to save imported contacts - %v", err)
		return
	}

	log.Infof("Contacts added: %d, updated: %d, ignored: %d", added, updated, ignored)
	t.UpdateStatus(Finished)
}

func (t *ContactImportTask) Progress() *ContactImportProgress {
	t.progressLock.Lock()
	defer t.progressLock.Unlock()
	return t.progress
}
