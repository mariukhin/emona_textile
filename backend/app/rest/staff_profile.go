package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type staffProfile struct {
	store model.StaffStore

	mediaDir  string
	mediaHost string
}

func (s *staffProfile) getStaffProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	staff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, staff)
}

func (s *staffProfile) updateStaffProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	requestBody := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.FirstName) == 0 {
		httpResponseError(w, http.StatusBadRequest, "first_name__required")
		return
	}

	if len(requestBody.LastName) == 0 {
		httpResponseError(w, http.StatusBadRequest, "last_name__required")
		return
	}

	update := staff.NewUpdate().
		SetFirstName(requestBody.FirstName).
		SetLastName(requestBody.LastName)

	updatedStaff, err := s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffProfile) updateStaffProfilePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
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

	staffProfilePhotoDir := path.Join("staff", "photo")

	photoFullPath := path.Join(s.mediaDir, staffProfilePhotoDir)
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

	update := staff.NewUpdate().SetPhoto(path.Join(staffProfilePhotoDir, staff.ID))
	updatedStaff, err := s.store.UpdateStaff(ctx, update)
	if err != nil {
		s.getLog(ctx).Errorf("Not found - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	updatedStaff.UpdatePhotoURL(s.mediaHost)
	httpJsonResponse(w, updatedStaff)
}

func (s *staffProfile) deleteStaffProfilePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	staff, err := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)
	if err != nil {
		s.getLog(ctx).Errorf("Fail to get staff from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
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

func (s *staffProfile) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "staff_profile",
	})
}
