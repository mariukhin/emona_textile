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
	"net/http"
)

type userAccounts struct {
	tokenAuthority TokenAuthority
	emailTemplates template.EmailTemplates

	mailSender  mail.Sender
	linkBuilder LinkBuilder

	accountStore model.AccountStore
	userStore    model.UserStore

	backOfficeHost        string
	secureCookiesDisabled bool
}

func (a *userAccounts) invitationDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	invitationToken := r.URL.Query().Get("t")
	if len(invitationToken) == 0 {
		httpResponseError(w, http.StatusBadRequest, "token__required")
		return
	}

	memberID, err := a.tokenAuthority.ValidateAccountMemberInvitationToken(ctx, invitationToken)
	if err != nil {
		if errors.Is(err, ErrExpired) {
			httpResponseError(w, http.StatusForbidden, "invitation__expired")
			return
		}
	}

	account, member, err := a.accountStore.FindAccountMember(ctx, memberID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "invitation__not_found")
		} else {
			a.getLog(ctx).Errorf("Fail to find account member - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if member.Status != model.AccountMemberStatusInvited {
		httpResponseError(w, http.StatusForbidden, "invitation__already_used")
		return
	}

	if member.IsInvitationExpired() {
		httpResponseError(w, http.StatusForbidden, "invitation__expired")
		return
	}

	isUserExists := len(member.UserId) > 0

	resp := struct {
		ID           string `json:"id"`
		AccountName  string `json:"account_name"`
		IsUserExists bool   `json:"is_user_exists"`
	}{
		ID:           member.ID,
		AccountName:  account.Name,
		IsUserExists: isUserExists,
	}

	httpJsonResponse(w, resp)
}

func (a *userAccounts) invitationAccept(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	invitationID := chi.URLParam(r, "invitationID")
	if len(invitationID) == 0 {
		httpResponseError(w, http.StatusBadRequest, "invitation__required")
		return
	}

	// TODO check if invitationID valid

	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)

	_, member, err := a.accountStore.FindAccountMember(ctx, invitationID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "invitation__not_found")
		} else {
			a.getLog(ctx).Errorf("Fail to find account member - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if member.Status != model.AccountMemberStatusInvited {
		httpResponseError(w, http.StatusForbidden, "invitation__already_used")
		return
	}

	if member.IsInvitationExpired() {
		httpResponseError(w, http.StatusForbidden, "invitation__expired")
		return
	}

	if me.ID != member.UserId {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	memberUpdate := member.NewUpdate().SetStatus(model.AccountMemberStatusActive)
	err = a.accountStore.UpdateMember(ctx, memberUpdate)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account member - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *userAccounts) invitationReject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	invitationID := chi.URLParam(r, "invitationID")
	if len(invitationID) == 0 {
		httpResponseError(w, http.StatusBadRequest, "invitation__required")
		return
	}

	// TODO check if invitationID valid

	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)

	_, member, err := a.accountStore.FindAccountMember(ctx, invitationID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "invitation__not_found")
		} else {
			a.getLog(ctx).Errorf("Fail to find account member - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if member.Status != model.AccountMemberStatusInvited {
		httpResponseError(w, http.StatusForbidden, "invitation__already_used")
		return
	}

	if member.IsInvitationExpired() {
		httpResponseError(w, http.StatusForbidden, "invitation__expired")
		return
	}

	if me.ID != member.UserId {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	memberUpdate := member.NewUpdate().SetStatus(model.AccountMemberStatusRejected)
	err = a.accountStore.UpdateMember(ctx, memberUpdate)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account member - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *userAccounts) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	resp := a.accountDetailsResp(ctx, account)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Role != model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	update := account.NewUpdate()
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if update.Name == nil || len(*update.Name) == 0 {
		httpResponseError(w, http.StatusBadRequest, "name__required")
		return
	}

	if len(*update.Name) < 2 {
		httpResponseError(w, http.StatusBadRequest, "name__invalid")
		return
	}

	if update.Email == nil || len(*update.Email) == 0 {
		httpResponseError(w, http.StatusBadRequest, "email__required")
		return
	}

	if !validEmailRegex.MatchString(*update.Email) {
		httpResponseError(w, http.StatusBadRequest, "email__invalid")
		return
	}

	updatedAccount, err := a.accountStore.Update(ctx, update)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusBadRequest, "email__exists")
		} else {
			a.getLog(ctx).Errorf("Fail to update user - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if update.IsEmailUpdated() {
		me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
		go a.sendUserEmailConfirmation(NewAsyncTaskContext(ctx), me, updatedAccount)
	}

	resp := a.accountDetailsResp(ctx, updatedAccount)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) resendEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Role != model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if account.IsEmailConfirmed {
		httpResponseError(w, http.StatusForbidden, "email__already_confirmed")
		return
	}

	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	go a.sendUserEmailConfirmation(NewAsyncTaskContext(ctx), me, account)

	resp := a.accountDetailsResp(ctx, account)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) confirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Token) == 0 {
		httpResponseError(w, http.StatusBadRequest, "token__required")
		return
	}

	accountID, err := a.tokenAuthority.ValidateAccountEmailConfirmationToken(r.Context(), requestBody.Token)
	if err != nil {
		if errors.Is(err, ErrExpired) {
			httpResponseError(w, http.StatusBadRequest, "token__expired")
			return
		} else if errors.Is(err, ErrAlreadyUsed) {
			httpResponseError(w, http.StatusBadRequest, "token__already_used")
			return
		} else {
			httpResponseError(w, http.StatusBadRequest, "token__invalid")
			return
		}
	}

	account, err := a.accountStore.FindUserAccountByID(ctx, accountID)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to fetch account - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if account.EmailConfirmedAt != nil {
		httpPlainError(w, http.StatusOK, "Email already confirmed")
		return
	}

	update := account.NewUpdate().SetEmailConfirmed()
	_, err = a.accountStore.Update(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *userAccounts) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	accounts, err := a.accountStore.FetchUserAccounts(ctx, user.ID)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to fetch accounts list - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if len(accounts) > 0 {
		currentAccount, _ := model.AccountFromContext(ctx)
		if currentAccount != nil {
			for _, acc := range accounts {
				if acc.ID == currentAccount.ID {
					acc.Current = true
					a.currentAccountCookie(w, acc.ID)
					break
				}
			}
		} else {
			firstAcc := accounts[0]
			firstAcc.Current = true
			a.currentAccountCookie(w, firstAcc.ID)
		}
	}

	httpJsonResponse(w, accounts)
}

func (a *userAccounts) makeCurrent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	accountId := chi.URLParam(r, "accountID")
	if len(accountId) == 0 {
		httpResponseError(w, http.StatusBadRequest, "account__invalid")
		return
	}

	// TODO check if account allowed for the user

	a.currentAccountCookie(w, accountId)
}

func (a *userAccounts) membersList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)
	httpJsonResponse(w, account.Members)
}

func (a *userAccounts) inviteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	requestBody := struct {
		Email string                  `json:"email"`
		Role  model.AccountMemberRole `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Email) == 0 {
		httpResponseError(w, http.StatusBadRequest, "email__required")
		return
	}

	if !validEmailRegex.MatchString(requestBody.Email) {
		httpResponseError(w, http.StatusBadRequest, "email__invalid")
		return
	}

	if len(requestBody.Role) == 0 {
		httpResponseError(w, http.StatusBadRequest, "role__required")
		return
	}

	roleFound := false
	for _, r := range model.AllowedAccountMemberRoles {
		if r == requestBody.Role {
			roleFound = true
			break
		}
	}

	if !roleFound {
		httpResponseError(w, http.StatusBadRequest, "role__invalid")
		return
	}

	// TODO transaction
	existsUser, err := a.userStore.FindByLogin(ctx, requestBody.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// do nothing
		} else {
			a.getLog(ctx).Errorf("Fail to find user - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	newMember := account.NewMember(existsUser, requestBody.Email, requestBody.Role)
	updatedMember, err := a.accountStore.AddMember(ctx, newMember)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusForbidden, "member__already_exist")
		} else {
			a.getLog(ctx).Errorf("Fail to save member - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	go a.sendUserInvitation(NewAsyncTaskContext(ctx), account, newMember)

	resp := a.accountMemberDetailsResp(ctx, updatedMember)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) getMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestedMember, _ := model.MemberFromContext(ctx)

	resp := a.accountMemberDetailsResp(ctx, requestedMember)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) updateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	requestedMember, _ := model.MemberFromContext(ctx)

	if me.ID == requestedMember.UserId {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if requestedMember.Role == model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	requestBody := struct {
		Role model.AccountMemberRole `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Role) == 0 {
		httpResponseError(w, http.StatusBadRequest, "role__required")
		return
	}

	roleFound := false
	for _, r := range model.AllowedAccountMemberRoles {
		if r == requestBody.Role {
			roleFound = true
			break
		}
	}

	if !roleFound {
		httpResponseError(w, http.StatusBadRequest, "role__invalid")
		return
	}

	memberUpdate := requestedMember.NewUpdate().SetRole(requestBody.Role)
	err = a.accountStore.UpdateMember(ctx, memberUpdate)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update member - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	// TODO hack
	requestedMember.Role = requestBody.Role

	resp := a.accountMemberDetailsResp(ctx, requestedMember)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) deleteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	requestedMember, _ := model.MemberFromContext(ctx)

	if me.ID == requestedMember.UserId {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	if requestedMember.Role == model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	err := a.accountStore.DeleteAccountMember(ctx, requestedMember)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httpResponseError(w, http.StatusNotFound, "member__not_found")
		} else {
			a.getLog(ctx).Errorf("Fail to delete account member - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *userAccounts) resendMemberInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)
	accountMember, _ := model.MemberFromContext(ctx)

	if accountMember.Status != model.AccountMemberStatusInvited {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	go a.sendUserInvitation(NewAsyncTaskContext(ctx), account, accountMember)

	w.WriteHeader(http.StatusNoContent)
}

func (a *userAccounts) serviceAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	resp := a.accountAddressResp(ctx, account)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) serviceAddressUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account, _ := model.AccountFromContext(ctx)

	if account.Role != model.AccountMemberRoleOwner {
		httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		return
	}

	address := model.AccountAddress{}
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(address.Address) == 0 {
		httpResponseError(w, http.StatusBadRequest, "address__required")
		return
	}

	if len(address.Address) < 3 {
		httpResponseError(w, http.StatusBadRequest, "address__invalid")
		return
	}

	if len(address.City) == 0 {
		httpResponseError(w, http.StatusBadRequest, "city__required")
		return
	}

	if len(address.City) < 2 {
		httpResponseError(w, http.StatusBadRequest, "city__invalid")
		return
	}

	if len(address.State) == 0 {
		httpResponseError(w, http.StatusBadRequest, "state__required")
		return
	}

	if len(address.State) < 2 {
		httpResponseError(w, http.StatusBadRequest, "state__invalid")
		return
	}

	if len(address.Country) == 0 {
		httpResponseError(w, http.StatusBadRequest, "country__required")
		return
	}

	if len(address.Country) < 2 {
		httpResponseError(w, http.StatusBadRequest, "country__invalid")
		return
	}

	if len(address.PostalCode) == 0 {
		httpResponseError(w, http.StatusBadRequest, "postal_code__required")
		return
	}

	// TODO check postal code

	update := account.NewUpdate().SetServiceAddress(address)
	updatedAccount, err := a.accountStore.Update(ctx, update)
	if err != nil {
		a.getLog(ctx).Errorf("Fail to update account service address - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	resp := a.accountAddressResp(ctx, updatedAccount)

	httpJsonResponse(w, resp)
}

func (a *userAccounts) currentAccountCookie(w http.ResponseWriter, accountId string) {
	cookie := &http.Cookie{
		Name:  "account-id",
		Value: accountId,

		Domain: a.backOfficeHost,
		Path:   "/customer/api/v1",

		Secure: !a.secureCookiesDisabled,
	}

	http.SetCookie(w, cookie)
}

func (a *userAccounts) accountDetailsResp(ctx context.Context, account *model.AccountDetails) interface{} {
	canEdit := false

	if account.Role == model.AccountMemberRoleOwner {
		canEdit = true
	}

	resp := struct {
		*model.AccountDetails
		CanEdit bool `json:"can_edit"`
	}{
		AccountDetails: account,
		CanEdit:        canEdit,
	}

	return resp
}

func (a *userAccounts) accountMemberDetailsResp(ctx context.Context, member *model.AccountMember) interface{} {
	me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	account, _ := model.AccountFromContext(ctx)

	canEdit := false
	canRemove := false
	canResendInvite := false

	if account.Role == model.AccountMemberRoleOwner {
		if me.ID != member.UserId {
			canEdit = true
			canRemove = true
		}
	}

	if account.Role == model.AccountMemberRoleAdmin {
		if me.ID != member.UserId && member.Role != model.AccountMemberRoleOwner {
			canEdit = true
			canRemove = true
		}
	}

	if account.Role == model.AccountMemberRoleOwner || account.Role == model.AccountMemberRoleAdmin {
		if member.Status == model.AccountMemberStatusInvited {
			canResendInvite = true
		}
	}

	resp := struct {
		*model.AccountMember
		CanEdit         bool `json:"can_edit"`
		CanRemove       bool `json:"can_remove"`
		CanResendInvite bool `json:"can_resend_invite"`
	}{
		AccountMember:   member,
		CanEdit:         canEdit,
		CanRemove:       canRemove,
		CanResendInvite: canResendInvite,
	}

	return resp
}

func (a *userAccounts) accountAddressResp(ctx context.Context, account *model.AccountDetails) interface{} {
	canEdit := false

	if account.Role == model.AccountMemberRoleOwner {
		canEdit = true
	}

	resp := struct {
		model.AccountAddress
		CanEdit bool `json:"can_edit"`
	}{
		AccountAddress: account.ServiceAddress,
		CanEdit:        canEdit,
	}

	return resp
}

func (a *userAccounts) sendUserEmailConfirmation(c context.Context, initiator *model.User, account *model.AccountDetails) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send user email confirmation")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Your Coton SMS account has been updated")
	mailOpts.SetRecipients(account.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := a.tokenAuthority.IssueAccountEmailConfirmationToken(ctx, account)
	if err != nil {
		log.Errorf("Fail to create confirmation token - %v", err)
		return
	}

	confirmationLink := a.linkBuilder.AccountEmailConfirmation(*ecToken)

	tmpl, err := a.emailTemplates.AccountEmailConfirmation(confirmationLink, fmt.Sprintf("%s %s", initiator.FirstName, initiator.LastName))
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = a.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := account.NewUpdate().SetEmailConfirmationSent()
	_, err = a.accountStore.Update(ctx, update)
	if err != nil {
		log.Errorf("Fail to update account - %v", err)
		return
	}

	log.Info("Successfully sent user email confirmation")
}

func (a *userAccounts) sendUserInvitation(c context.Context, account *model.AccountDetails, member *model.AccountMember) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send user invitation email")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("CotonSMS account invitation")
	mailOpts.SetRecipients(member.InvitationEmail)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %a", newStaff.Password))

	ecToken, err := a.tokenAuthority.IssueAccountMemberInvitationToken(ctx, account, member)
	if err != nil {
		log.Errorf("Fail to create invitation token - %v", err)
		return
	}

	invitationLink := a.linkBuilder.AccountMemberInvitation(*ecToken)
	userName := "user"
	if member.User != nil {
		userName = fmt.Sprintf("%s %s", member.User.FirstName, member.User.LastName)
	}

	tmpl, err := a.emailTemplates.AccountMemberInvitation(invitationLink, account.Name, userName)
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = a.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := member.NewUpdate().SetInvitationSent()
	err = a.accountStore.UpdateMember(ctx, update)
	if err != nil {
		log.Errorf("Fail to update account member - %v", err)
		return
	}

	log.Info("Successfully sent user invitation email")
}

func (a *userAccounts) populateAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		me, _ := model.UserFromContext(ctx, model.KCtxKeyUserMe)

		cookie, err := r.Cookie("account-id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
			} else {
				a.getLog(ctx).Errorf("Fail to get account id from cookie - %v", err)
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		accountId := cookie.Value
		if len(accountId) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		account, err := a.accountStore.FindUserAccount(ctx, me.ID, accountId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				next.ServeHTTP(w, r)
			} else {
				a.getLog(ctx).Errorf("Fail find account - %v", err)
				httpResponseError(w, http.StatusInternalServerError, "internal")
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(account.NewContext(r.Context())))
	})
}

func (a *userAccounts) ensureAccountPopulated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		account, _ := model.AccountFromContext(ctx)

		if account == nil {
			httpResponseError(w, http.StatusBadRequest, "account__required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *userAccounts) populateAccountMember(next http.Handler) http.Handler {
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

func (a *userAccounts) restrictAccess(allowedRoles ...model.AccountMemberRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			account, err := model.AccountFromContext(r.Context())
			if err != nil {
				httpResponseError(w, http.StatusBadRequest, "account__required")
				return
			}

			for _, role := range allowedRoles {
				if role == account.Role {
					next.ServeHTTP(w, r)
					return
				}
			}

			httpResponseError(w, http.StatusForbidden, "action__not_allowed")
		}
		return http.HandlerFunc(fn)
	}
}

func (a *userAccounts) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "account_contact_group",
	})
}
