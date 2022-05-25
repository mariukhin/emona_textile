package api

import (
	"amifactory.team/sequel/coton-app-backend/app/model"
	"net/http"
)

type userAccountRoles struct {
}

func (u *userAccountRoles) availableRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Role != model.AccountMemberRoleOwner &&
		account.Role != model.AccountMemberRoleAdmin {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	httpJsonResponse(w, model.AllowedAccountMemberRoles)
}
