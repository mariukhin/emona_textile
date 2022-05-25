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
	"regexp"
	"unicode/utf8"
)

type staffAccounts struct {
	accountStore model.AccountStore
	staffStore   model.StaffStore

	secureCookiesDisabled bool
}

func (a *staffAccounts) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	form := newAccountsPageForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	accountsList, total, err := a.accountStore.FetchAccounts(ctx, form.Offset(), form.PageSize, form.SearchPhrase, form.StatusParam)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
		} else {
			a.getLog(ctx).Errorf("Fail to fetch accounts list - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	resp := NewPageResp(total, form.PageSize, accountsList)

	httpJsonResponse(w, resp)
}
func (a *staffAccounts) listAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	form := newAllAccountsPageForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	accountsList, err := a.accountStore.FetchAllAccounts(ctx, form.SrcParam, form.StatusParam)
	if err != nil {
		if errors.Is(err, model.ErrPageOutOfBounds) {
			httpResponseError(w, http.StatusBadRequest, "page__out_of_bounds")
		} else {
			a.getLog(ctx).Errorf("Fail to fetch accounts list - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	httpJsonResponse(w, accountsList)
}

func (a *staffAccounts) fetchAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)
	httpJsonResponse(w, account.Details())
}

func (a *staffAccounts) updateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	requestBody := struct {
		Moderation model.AccountModeration `json:"moderation"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Moderation) == 0 {
		httpResponseError(w, http.StatusBadRequest, "moderation__required")
		return
	}

	if !requestBody.Moderation.IsValid() {
		httpResponseError(w, http.StatusForbidden, "moderation__invalid")
		return
	}

	update := account.NewUpdate().SetModeration(requestBody.Moderation)
	updatedAccount, err := a.accountStore.Update(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, updatedAccount.Details())
}

func (a *staffAccounts) activateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Status == model.AccountStatusActive {
		httpResponseError(w, http.StatusForbidden, "status__invalid")
		return
	}

	update := account.NewUpdate().SetStatus(model.AccountStatusActive)
	updatedAccount, err := a.accountStore.Update(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, updatedAccount.Details())
}

func (a *staffAccounts) blockAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Status == model.AccountStatusBlocked {
		httpResponseError(w, http.StatusForbidden, "status__invalid")
		return
	}

	update := account.NewUpdate().SetStatus(model.AccountStatusBlocked)
	updatedAccount, err := a.accountStore.Update(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	httpJsonResponse(w, updatedAccount.Details())
}

func (a *staffAccounts) accountMembers(w http.ResponseWriter, r *http.Request) {
	account, _ := model.AccountFromContext(r.Context())
	httpJsonResponse(w, account.Members)
}

func (a *staffAccounts) accountAddress(w http.ResponseWriter, r *http.Request) {
	account, _ := model.AccountFromContext(r.Context())
	httpJsonResponse(w, account.ServiceAddress)
}

func (a *staffAccounts) accountMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	member, _ := model.MemberFromContext(ctx)
	httpJsonResponse(w, a.accountMemberDetailsResp(ctx, member))
}

func (a *staffAccounts) updateAccountMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	member, _ := model.MemberFromContext(ctx)

	if member.Role == model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	requestBody := struct {
		Role model.AccountMemberRole `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %sv", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Role) == 0 {
		httpResponseError(w, http.StatusBadRequest, "role__required")
		return
	}

	if !requestBody.Role.IsAccountMemberRoleValid() {
		httpResponseError(w, http.StatusForbidden, "role__invalid")
		return
	}

	update := member.NewUpdate().SetRole(requestBody.Role)
	err = a.accountStore.UpdateMember(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account member - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	// TODO hack
	member.Role = requestBody.Role

	httpJsonResponse(w, a.accountMemberDetailsResp(ctx, member))
}

func (a *staffAccounts) deleteAccountMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	member, _ := model.MemberFromContext(ctx)

	if member.Role == model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	err := a.accountStore.DeleteAccountMember(ctx, member)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "member__not_found")
		} else {
			a.getLog(ctx).Errorf("Fail to delete account member -%v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *staffAccounts) accountMemberDetailsResp(ctx context.Context, member *model.AccountMember) interface{} {
	me, _ := model.StaffFromContext(ctx, model.KCtxKeyStaffMe)

	canEdit := false
	canRemove := false

	if me.Role.Name == string(model.StaffRoleNewAdmin) && member.Role != model.AccountMemberRoleOwner {
		canEdit = true
		canRemove = true
	}

	resp := struct {
		*model.AccountMember
		CanEdit   bool `json:"can_edit"`
		CanRemove bool `json:"can_remove"`
	}{
		AccountMember: member,
		CanEdit:       canEdit,
		CanRemove:     canRemove,
	}

	return resp
}

func (a *staffAccounts) populateAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accountId := chi.URLParam(r, "accountID")
		if len(accountId) == 0 {
			httpResponseError(w, http.StatusBadRequest, "account__invalid")
			return
		}

		account, err := a.accountStore.FindAccountByID(ctx, accountId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusNotFound, "account__not_found")
			} else {
				a.getLog(ctx).Errorf("Fail to fetch account by ID - %v", err)
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}

			return
		}

		next.ServeHTTP(w, r.WithContext(account.NewContext(r.Context())))
	})
}

func (a *staffAccounts) populateAccountMember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		memberID := chi.URLParam(r, "memberID")
		if len(memberID) == 0 {
			httpResponseError(w, http.StatusBadRequest, "member__required")
			return
		}

		account, err := model.AccountFromContext(r.Context())
		if err != nil {
			httpResponseError(w, http.StatusBadRequest, "account__required")
			return
		}

		for _, m := range account.Members {
			if m.ID == memberID {
				next.ServeHTTP(w, r.WithContext(m.NewContext(r.Context())))
				return
			}
		}

		httpResponseError(w, http.StatusNotFound, "member__not_found")
	})
}

func (a *staffAccounts) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "accounts",
	})
}

/// Forms

type accountsPageForm struct {
	*PageForm
	StatusParam *model.AccountStatus
}

func newAccountsPageForm(r *http.Request) *accountsPageForm {
	pageForm := NewPageForm(r, 2)
	form := accountsPageForm{
		PageForm: pageForm,
	}

	statusParam := r.URL.Query().Get("status")

	if len(statusParam) > 0 {
		statusFilter := ParseAccountStatusFilter(statusParam)
		if statusFilter != nil {
			form.StatusParam = statusFilter.accountStatus
		}
	}

	return &form
}

type accountsAllPageForm struct {
	*BaseForm

	StatusParam *model.AccountStatus
	SrcParam    *string
}

func newAllAccountsPageForm(r *http.Request) *accountsAllPageForm {

	form := accountsAllPageForm{
		BaseForm: &BaseForm{},
	}

	statusParam := r.URL.Query().Get("status")
	srcParam := r.URL.Query().Get("src")

	if len(statusParam) > 0 {
		statusFilter := ParseAccountStatusFilter(statusParam)
		if statusFilter != nil {
			form.StatusParam = statusFilter.accountStatus
		}
	}
	if len(srcParam) > 0 {
		if utf8.RuneCountInString(srcParam) >= 3 {
			srcEscaped := regexp.QuoteMeta(srcParam)
			form.SrcParam = &srcEscaped
		} else {
			form.AddValidationError("src__invalid", nil)
		}
	} else {
		form.SrcParam = nil
	}

	return &form
}
