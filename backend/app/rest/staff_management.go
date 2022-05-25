package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/mail"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"amifactory.team/sequel/coton-app-backend/app/template"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

var ErrEmailConfirmationTooOften = errors.New("email confirmation too often")

type staffManagement struct {
	tokenAuthority TokenAuthority
	emailTemplates template.EmailTemplates

	store model.StaffStore

	mailSender  mail.Sender
	linkBuilder LinkBuilder

	mediaDir  string
	mediaHost string

	staffProfilePhotoDir string

	maxEmailConfirmationCount  int
	maxEmailConfirmationPeriod time.Duration
}

func (s *staffManagement) getStaffList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if staff.Role == nil || staff.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	srcParam := r.URL.Query().Get("src")
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("page_size")

	page := 0
	if len(pageParam) > 0 {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			s.getLog(ctx).Errorf("Fail to get parse 'page' param - %v", err)
			httpResponseError(w, http.StatusBadRequest, "page__invalid")
			return
		}

		if page < 0 {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
			return
		}
	}

	pageSize := 10
	if len(pageSizeParam) > 0 {
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil {
			s.getLog(ctx).Errorf("Fail to get parse 'page_size' param - %v", err)
			httpResponseError(w, http.StatusBadRequest, "page_size__invalid")
			return
		}

		if pageSize <= 0 || 100 < pageSize {
			httpResponseError(w, http.StatusBadRequest, "page_size__out_of_bounds")
			return
		}
	}

	var searchPhase *string
	if len(srcParam) > 0 {
		if len(srcParam) >= 3 {
			srcEscaped := regexp.QuoteMeta(srcParam)
			searchPhase = &srcEscaped
		} else {
			httpResponseError(w, http.StatusBadRequest, "src__invalid")
			return
		}
	} else {
		searchPhase = nil
	}

	staffList, total, err := s.store.FetchStaffList(page*pageSize, pageSize, &staff.ID, searchPhase)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
			return
		}

		s.getLog(ctx).Errorf("Fail to fetch staff list - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	for _, item := range staffList {
		item.UpdatePhotoURL(s.mediaHost)
	}

	resp := struct {
		Pages  int                 `json:"pages"`
		Total  int                 `json:"total"`
		Result []*model.StaffShort `json:"result"`
	}{
		Pages:  int(math.Ceil(float64(total) / float64(pageSize))),
		Total:  total,
		Result: staffList,
	}

	httpJsonResponse(w, resp)
}

func (s *staffManagement) createStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if staff.Role == nil || staff.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	requestBody := struct {
		Email     *string `json:"email"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Role      *string `json:"role"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if requestBody.Email == nil {
		httpResponseError(w, http.StatusBadRequest, "email__required")
		return
	}

	if !validEmailRegex.MatchString(*requestBody.Email) {
		httpResponseError(w, http.StatusBadRequest, "email__invalid")
		return
	}

	if requestBody.FirstName == nil {
		httpResponseError(w, http.StatusBadRequest, "first_name__required")
		return
	}

	if len(*requestBody.FirstName) == 0 || len(*requestBody.FirstName) > 50 {
		httpResponseError(w, http.StatusBadRequest, "first_name__invalid")
		return
	}

	if requestBody.LastName == nil {
		httpResponseError(w, http.StatusBadRequest, "last_name__required")
		return
	}

	if len(*requestBody.LastName) == 0 || len(*requestBody.LastName) > 50 {
		httpResponseError(w, http.StatusBadRequest, "last_name__invalid")
		return
	}

	if requestBody.Role == nil {
		httpResponseError(w, http.StatusBadRequest, "role__required")
		return
	}

	if len(*requestBody.Role) == 0 || len(*requestBody.Role) > 15 {
		httpResponseError(w, http.StatusBadRequest, "role__invalid")
		return
	}

	role, err := s.store.FindStaffRole(*requestBody.Role)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusBadRequest, "role__invalid")
		} else {
			s.getLog(ctx).Errorf("Fail to find role - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	newStaff, err := model.NewStaff(
		*requestBody.Email,
		*requestBody.FirstName,
		*requestBody.LastName,
		role,
	)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to create staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	err = s.store.AddStaff(newStaff)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusBadRequest, "email__exists")
		} else {
			s.getLog(ctx).Errorf("Fail to create staff - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	go s.sendStaffSignUpEmailConfirmation(NewAsyncTaskContext(ctx), newStaff)

	httpJsonResponse(w, newStaff)
}

func (s *staffManagement) fetchStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	staff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, staff)
}

func (s *staffManagement) updateStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	staffUpdate := staff.NewUpdate()
	err = json.NewDecoder(r.Body).Decode(&staffUpdate)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if staffUpdate.RoleName != nil {
		staffUpdate.Role, err = s.store.FindStaffRole(*staffUpdate.RoleName)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusBadRequest, "role__invalid")
			} else {
				s.getLog(ctx).Errorf("Fail to find role - %v", err)
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}

			return
		}
	}

	updatedStaff, err := s.store.UpdateStaff(ctx, staffUpdate)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusBadRequest, "email__exists")
		} else {
			s.getLog(ctx).Errorf("Not found - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if staffUpdate.IsEmailUpdated() {
		go s.sendStaffEmailConfirmation(NewAsyncTaskContext(ctx), updatedStaff)
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffManagement) updateStaffPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get photo from body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "photo__invalid")
		return
	}

	if !IsImage(bodyBytes) {
		httpResponseError(w, http.StatusBadRequest, "photo__invalid")
		return
	}

	photoFullPath := path.Join(s.mediaDir, s.staffProfilePhotoDir)
	err = os.MkdirAll(photoFullPath, 0770)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to make photo dir - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	staffPhotoName := path.Join(photoFullPath, staff.ID)
	err = ioutil.WriteFile(staffPhotoName, bodyBytes, 0770)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to write photo - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	update := staff.NewUpdate().SetPhoto(path.Join(s.staffProfilePhotoDir, staff.ID))
	updatedStaff, err := s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffManagement) deleteStaffPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	if len(staff.Photo) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	photoFullPath := path.Join(s.mediaDir, staff.Photo)
	_, err = os.Stat(photoFullPath)
	if err != nil {
		if !os.IsNotExist(err) {
			s.getLog(ctx).Errorf("Fail to remove photo file - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	} else {
		err = os.Remove(photoFullPath)
		if err != nil {
			s.getLog(ctx).Errorf("Fail to remove photo file - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	update := staff.NewUpdate().ClearPhoto()
	_, err = s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *staffManagement) activateStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if staff.ID == me.ID {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__active")
		return
	}

	update := staff.NewUpdate().SetActive(true)
	updatedStaff, err := s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffManagement) deactivateStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if staff.ID == me.ID {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	update := staff.NewUpdate().SetActive(false)
	updatedStaff, err := s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffManagement) resetPasswordStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if staff.ID == me.ID {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	if !staff.IsEmailConfirmed {
		httpResponseError(w, http.StatusForbidden, "email__not_confirmed")
		return
	}

	update, err := staff.ResetPassword()
	if err != nil {
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	_, err = s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to reset staff password - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	go s.sendStaffPasswordResetLink(NewAsyncTaskContext(ctx), staff)

	w.WriteHeader(http.StatusOK)
}

func (s *staffManagement) sendEmailConfirmationStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	// TODO update later
	if me.Role == nil || me.Role.Name != "admin" {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaff)
	if err != nil {
		httpResponseError(w, http.StatusForbidden, "staff__invalid")
		return
	}

	if staff.ID == me.ID {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__deactivated")
		return
	}

	if staff.IsEmailConfirmed {
		httpResponseError(w, http.StatusForbidden, "email__already_confirmed")
		return
	}

	go s.sendStaffEmailConfirmation(NewAsyncTaskContext(ctx), staff)

	w.WriteHeader(http.StatusOK)
}

func (s *staffManagement) staffById(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		me, _ := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)

		staffId := chi.URLParam(r, "staffID")
		if len(staffId) == 0 {
			httpResponseError(w, http.StatusBadRequest, "staff__invalid")
			return
		}

		var staff *model.Staff
		var err error
		if me != nil && staffId == me.ID {
			staff = me
		} else {
			staff, err = s.store.FindStaffByID(ctx, staffId)
			if err != nil {
				if errors.Is(err, model.ErrNotFound) {
					httpResponseError(w, http.StatusNotFound, "staff__not_found")
				} else {
					s.getLog(ctx).Errorf("Fail to fetch staff by ID - %v", err)
					httpResponseError(w, http.StatusInternalServerError, "internal")
				}

				return
			}
		}

		next.ServeHTTP(w, r.WithContext(staff.NewContext(ctx, model.KCtxKeyStaff)))
	})
}

func (s *staffManagement) sendStaffSignUpEmailConfirmation(c context.Context, newStaff *model.Staff) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send staff sign-up email confirmation")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Welcome to Coton Admin panel!")
	mailOpts.SetRecipients(newStaff.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := s.tokenAuthority.IssueStaffEmailConfirmationToken(ctx, newStaff)
	if err != nil {
		log.Errorf("Fail to issue confirmation token - %v", err)
		return
	}

	confirmationLink := s.linkBuilder.StaffEmailConfirmation(*ecToken)
	loginLink := s.linkBuilder.StaffLogin()

	tmpl, err := s.emailTemplates.StaffSignUpEmail(confirmationLink, loginLink, fmt.Sprintf("%s %s", newStaff.FirstName, newStaff.LastName), newStaff.Email, newStaff.Password)
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = s.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := newStaff.NewUpdate().SetEmailConfirmationSent()
	_, err = s.store.UpdateStaff(ctx, update)
	if err != nil {
		log.Errorf("Fail to update staff - %v", err)
		return
	}

	log.Info("Successfully sent staff sign-up email confirmation")
}

func (s *staffManagement) sendStaffPasswordResetLink(c context.Context, staff *model.Staff) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send staff password reset link")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Your Coton Admin panel password has been changed")
	mailOpts.SetRecipients(staff.Email)
	// TODO make
	mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", staff.Password))

	loginLink := s.linkBuilder.StaffLogin()

	tmpl, err := s.emailTemplates.StaffPasswordUpdated(loginLink, fmt.Sprintf("%s %s", staff.FirstName, staff.LastName), staff.Email, staff.Password)
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = s.mailSender.Send(context.Background(), mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	log.Info("Successfully sent staff password reset link")
}

func (s *staffManagement) sendStaffEmailConfirmation(c context.Context, staff *model.Staff) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send staff email confirmation")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Your Coton Admin panel account has been updated")
	mailOpts.SetRecipients(staff.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := s.tokenAuthority.IssueStaffEmailConfirmationToken(ctx, staff)
	if err != nil {
		log.Errorf("Fail to create confirmation token - %v", err)
		return
	}

	confirmationLink := s.linkBuilder.StaffEmailConfirmation(*ecToken)

	tmpl, err := s.emailTemplates.StaffEmailConfirmation(confirmationLink, fmt.Sprintf("%s %s", staff.FirstName, staff.LastName))
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = s.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := staff.NewUpdate().SetEmailConfirmationSent()
	_, err = s.store.UpdateStaff(ctx, update)
	if err != nil {
		log.Errorf("Fail to update staff - %v", err)
		return
	}

	log.Info("Successfully sent staff email confirmation")
}

func (s *staffManagement) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "staff_management",
	})
}
