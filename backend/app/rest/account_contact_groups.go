package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
)

type userAccountContactGroups struct {
	contactsStore model.ContactsStore

	secureCookiesDisabled bool
}

func (u *userAccountContactGroups) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)

	form := NewPageForm(r, 3)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	groups, total, err := u.contactsStore.ListGroups(ctx, account.ID, form.Offset(), form.PageSize, form.SearchPhrase)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
		} else {
			u.getLog(ctx).Errorf("Fail to fetch groups list - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if total == 0 && form.HasFilter() {
		httpResponseError(w, http.StatusNotFound, "result__not_found")
		return
	}

	resp := NewPageResp(total, form.PageSize, groups)
	httpJsonResponse(w, resp)
}

func (u *userAccountContactGroups) listAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)

	groups, err := u.contactsStore.ListGroupsAll(ctx, account.ID)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to fetch groups list - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, groups)
}

func (u *userAccountContactGroups) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account, _ := model.AccountFromContext(ctx)

	requestBody := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Name) == 0 {
		httpResponseError(w, http.StatusBadRequest, "name__required")
		return
	}

	if len(requestBody.Name) > 150 {
		httpResponseError(w, http.StatusForbidden, "name__invalid")
		return
	}

	if len(requestBody.Description) > 1000 {
		httpResponseError(w, http.StatusForbidden, "description__invalid")
		return
	}

	group := model.NewContactGroup(account.ID, requestBody.Name, requestBody.Description)
	err = u.contactsStore.AddGroup(ctx, group)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to add group - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, group)
}

func (u *userAccountContactGroups) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)

	httpJsonResponse(w, group)
}

func (u *userAccountContactGroups) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)

	requestBody := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Name) == 0 {
		httpResponseError(w, http.StatusBadRequest, "name__required")
		return
	}

	if len(requestBody.Name) > 150 {
		httpResponseError(w, http.StatusForbidden, "name__invalid")
		return
	}

	if len(requestBody.Description) > 1000 {
		httpResponseError(w, http.StatusForbidden, "description__invalid")
		return
	}

	group.Name = requestBody.Name
	group.Description = requestBody.Description

	err = u.contactsStore.UpdateGroup(ctx, group)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to update group - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, group)
}

func (u *userAccountContactGroups) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)
	err := u.contactsStore.DeleteGroup(ctx, group)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to delete group - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *userAccountContactGroups) populateGroup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		account, _ := model.AccountFromContext(ctx)

		groupID := chi.URLParam(r, "groupID")
		if len(groupID) == 0 {
			httpResponseError(w, http.StatusBadRequest, "group__required")
			return
		}

		group, err := u.contactsStore.FindGroup(ctx, account.ID, groupID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusNotFound, "group__not_found")
			} else {
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(group.NewContext(ctx)))
	})
}

func (u *userAccountContactGroups) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "account_contact_group",
	})
}
