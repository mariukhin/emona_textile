package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
)

type userAccountContactGroupVars struct {
	contactsStore model.ContactsStore

	secureCookiesDisabled bool
}

func (u *userAccountContactGroupVars) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)
	vars := group.Vars
	if vars == nil {
		vars = make([]model.ContactVariable, 0)
	}

	httpJsonResponse(w, vars)
}

func (u *userAccountContactGroupVars) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)

	form, err := NewContactVariableForm(r, true)
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

	newVariable := model.NewContactVariable(form.Title, form.Type, form.DefaultValueString, form.DefaultValueNumber)
	err = u.contactsStore.AddGroupVar(ctx, group, &newVariable)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to add group var - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, newVariable)
}

func (u *userAccountContactGroupVars) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	variable, _ := model.ContactVariableFromContext(ctx)

	httpJsonResponse(w, variable)
}

func (u *userAccountContactGroupVars) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)
	variable, _ := model.ContactVariableFromContext(ctx)

	form, err := NewContactVariableForm(r, false)
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

	switch variable.Type {
	case model.VarTypeNumber:
		variable.DefaultValueNumber = form.DefaultValueNumber
	case model.VarTypeString:
		variable.DefaultValueString = form.DefaultValueString
	}

	err = u.contactsStore.UpdateGroupVar(ctx, group, variable)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to update group var - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, variable)
}

func (u *userAccountContactGroupVars) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	group, _ := model.ContactGroupFromContext(ctx)
	variable, _ := model.ContactVariableFromContext(ctx)

	err := u.contactsStore.DeleteGroupVar(ctx, group, variable)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to delete group var - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *userAccountContactGroupVars) populateVariable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		group, _ := model.ContactGroupFromContext(ctx)

		varID := chi.URLParam(r, "varID")
		if len(varID) == 0 {
			httpResponseError(w, http.StatusBadRequest, "var__required")
			return
		}

		if !uuidRegex.MatchString(varID) {
			httpResponseError(w, http.StatusBadRequest, "var__invalid")
			return
		}

		var variable *model.ContactVariable
		for _, v := range group.Vars {
			if v.ID == varID {
				variable = &v
				break
			}
		}

		if variable == nil {
			httpResponseError(w, http.StatusNotFound, "var__not_found")
			return
		}

		next.ServeHTTP(w, r.WithContext(variable.NewContext(ctx)))
	})
}

func (u *userAccountContactGroupVars) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "account_contact_group_var",
	})
}
