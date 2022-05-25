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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"unicode/utf8"
)

type userAuth struct {
	tokenAuthority TokenAuthority
	emailTemplates template.EmailTemplates

	userStore    model.UserStore
	accountStore model.AccountStore

	mailSender  mail.Sender
	linkBuilder LinkBuilder

	backOfficeHost        string
	secureCookiesDisabled bool
}

var validEmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u *userAuth) signInEmail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	form := newSignInForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	// TODO add password validation

	user, err := u.userStore.FindByLogin(r.Context(), form.Email)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to fetch user - %v", err)
		httpResponseError(w, http.StatusForbidden, "credential__invalid")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password))
	if err != nil {
		u.getLog(ctx).Errorf("Fail to generate token - %v", err)
		httpResponseError(w, http.StatusForbidden, "credential__invalid")
		return
	}

	//if !user.IsActive {
	//	restResponseError(w, http.StatusForbidden, "account__suspended")
	//	return
	//}

	jwtToken, err := u.tokenAuthority.IssueUserAPIToken(r.Context(), user)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to issue api token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	u.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (u *userAuth) signUpEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	form := newSignUpForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	// TODO transaction

	newUser, err := model.NewUser(form.Email, form.FirstName, form.LastName, form.Password)
	err = u.userStore.Add(ctx, newUser)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusBadRequest, "email__exists")
		} else {
			u.getLog(ctx).Errorf("Fail to create user - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	if len(form.AccountInvitationToken) > 0 {
		invitationToken := form.AccountInvitationToken
		memberID, err := u.tokenAuthority.ValidateAccountMemberInvitationToken(ctx, invitationToken)
		if err != nil {
			u.getLog(ctx).Errorf("Account member invitation failed - %v", err)
		} else {
			_, member, err := u.accountStore.FindAccountMember(ctx, memberID)
			if err != nil {
				u.getLog(ctx).Errorf("Account member invitation failed - %v", err)
			} else {
				if member.Status == model.AccountMemberStatusInvited {
					memberUpdate := member.NewUpdate().SetStatus(model.AccountMemberStatusActive).SetUser(newUser)
					err = u.accountStore.UpdateMember(ctx, memberUpdate)
					if err != nil {
						u.getLog(ctx).Errorf("Account member invitation failed - %v", err)
					}
				} else {
					u.getLog(ctx).Errorf("Account member invitation failed - invalid status %v", member.Status)
				}
			}
		}
	}

	go u.sendUserSignUpEmailConfirmation(NewAsyncTaskContext(ctx), newUser)

	jwtToken, err := u.tokenAuthority.IssueUserAPIToken(r.Context(), newUser)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to generate token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	u.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (u *userAuth) signOut(w http.ResponseWriter, r *http.Request) {
	u.clearRefreshTokenCookie(w)
	http.StatusText(http.StatusOK)
}

func (u *userAuth) requestPasswordReset(w http.ResponseWriter, r *http.Request) {


	ctx := r.Context()

	form := newRequestPasswordResetForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	user, err := u.userStore.FindByLogin(r.Context(), form.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			u.getLog(ctx).Warn("User with provided login/email is not found")
			w.WriteHeader(http.StatusOK)
		} else {
			u.getLog(ctx).Errorf("Fail to find user by login - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	go u.sendUserPasswordResetLink(NewAsyncTaskContext(ctx), user)

	w.WriteHeader(http.StatusOK)
}

func (u *userAuth) passwordResetWithToken(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	form := newPasswordResetWithTokenForm(r)
	if form.HasValidationErrors() {
		httpResponseErrors(w, http.StatusBadRequest, form.ValidationErrors)
		return
	}

	userID, err := u.tokenAuthority.ValidateUserPasswordResetToken(r.Context(), form.Token)
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

	user, err := u.userStore.FindByID(r.Context(), userID)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to fetch user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	update, err := user.NewUpdate().SetPassword(form.Password)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to update user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	user, err = u.userStore.Update(r.Context(), update)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to update user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	jwtToken, err := u.tokenAuthority.IssueUserAPIToken(r.Context(), user)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to generate token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	u.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (u *userAuth) refreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	ctx := r.Context()

	refreshToken := cookie.Value
	userId, err := u.tokenAuthority.ValidateUserRefreshToken(ctx, refreshToken)
	if err != nil {
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	user, err := u.userStore.FindByID(r.Context(), userId)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to fetch user - %v", err)
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	jwtToken, err := u.tokenAuthority.IssueUserAPIToken(r.Context(), user)
	if err != nil {
		u.getLog(ctx).Errorf("Fail to issue token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	u.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (u *userAuth) refreshTokenCookie(w http.ResponseWriter, jwtToken *jwtToken) {
	cookie := &http.Cookie{
		Name:  "refresh-token",
		Value: jwtToken.RefreshToken,

		Domain:  u.backOfficeHost,
		Path:    "/customer/api/v1/refresh",
		Expires: jwtToken.RefreshTokenExpire,

		Secure:   !u.secureCookiesDisabled,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (u *userAuth) clearRefreshTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "refresh-token",
		Value: "",

		Domain: u.backOfficeHost,
		Path:   "/customer/api/v1/refresh",
		MaxAge: -1,

		Secure:   !u.secureCookiesDisabled,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (u *userAuth) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := tokenFromHeader(r)
		if len(tokenStr) == 0 {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		ctx := r.Context()

		userId, err := u.tokenAuthority.ValidateUserApiToken(ctx, tokenStr)
		if err != nil {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		user, err := u.userStore.FindByID(ctx, userId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
				return
			}

			u.getLog(ctx).Errorf("Fail to find user by ID - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}

		// User is found, pass it through
		next.ServeHTTP(w, r.WithContext(user.NewContext(r.Context(), model.KCtxKeyUserMe)))
	})
}

func (u *userAuth) ensureConfirmed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		me, _ := model.UserFromContext(r.Context(), model.KCtxKeyUserMe)
		if me.IsEmailConfirmationRequired() {
			httpResponseError(w, http.StatusForbidden, "email__not_confirmed")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (u *userAuth) sendUserSignUpEmailConfirmation(c context.Context, newUser *model.User) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send user sign up email confirmation")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Welcome to Coton SMS!")
	mailOpts.SetRecipients(newUser.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := u.tokenAuthority.IssueUserEmailConfirmationToken(ctx, newUser)
	if err != nil {
		log.Errorf("Fail to create confirmation token - %v", err)
		return
	}

	confirmationLink := u.linkBuilder.UserEmailConfirmation(*ecToken)

	tmpl, err := u.emailTemplates.UserSignUpEmail(confirmationLink, fmt.Sprintf("%s %s", newUser.FirstName, newUser.LastName))
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = u.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := newUser.NewUpdate()
	update.SetEmailConfirmationSent()
	_, err = u.userStore.Update(ctx, update)
	if err != nil {
		log.Errorf("Fail to update user - %v", err)
		return
	}

	log.Info("Successfully sent user sign up email confirmation")
}

func (u *userAuth) sendUserPasswordResetLink(c context.Context, user *model.User) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send user password reset link")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Coton SMS password reset")
	mailOpts.SetRecipients(user.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := u.tokenAuthority.IssueUserPasswordResetToken(ctx, user)
	if err != nil {
		log.Errorf("Fail to create confirmation token - %v", err)
		return
	}

	passwordResetLink := u.linkBuilder.UserPasswordReset(*ecToken)

	tmpl, err := u.emailTemplates.UserPasswordReset(passwordResetLink, fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = u.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	log.Info("Successfully sent user password reset link")
}

func (u *userAuth) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "user_auth",
	})
}


/// Forms

type signUpForm struct {
	*BaseForm

	FirstName              string
	LastName               string
	Email                  string
	Password               string
	AccountInvitationToken string
}

func newSignUpForm(r *http.Request) signUpForm {
	body := struct {
		FirstName              string `json:"first_name"`
		LastName               string `json:"last_name"`
		Email                  string `json:"email"`
		Password               string `json:"password"`
		AccountInvitationToken string `json:"account_invitation_token"`
	}{}

	f := signUpForm{
		BaseForm: &BaseForm{},
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}
	if len(body.FirstName) != 0 {
		if utf8.RuneCountInString(body.FirstName) < 1 {
			f.AddValidationError("first_name__invalid", nil)
		} else if utf8.RuneCountInString(body.FirstName) > 50 {
			f.AddValidationError("first_name__invalid", nil)
		} else {
			f.FirstName = body.FirstName
		}
	} else {
		f.AddValidationError("first_name__required", nil)
	}
	if len(body.LastName) != 0 {
		if utf8.RuneCountInString(body.LastName) < 1 {
			f.AddValidationError("last_name__invalid", nil)
		} else if utf8.RuneCountInString(body.LastName) > 50 {
			f.AddValidationError("last_name__invalid", nil)
		} else {
			f.LastName = body.LastName
		}
	} else {
		f.AddValidationError("last_name__required", nil)
	}

	if len(body.Email) != 0 {
		if validEmailRegex.MatchString(body.Email){
			f.Email = body.Email
		} else {
			f.AddValidationError("email__invalid", nil)
		}
	} else {
		f.AddValidationError("email__required", nil)
	}

	if len(body.Password) != 0 {
		if model.IsStrongPassword(body.Password){
			f.Password = body.Password
		}else {
			f.AddValidationError("password__invalid", nil)
		}
	} else {
		f.AddValidationError("password__required", nil)
	}
	f.AccountInvitationToken = body.AccountInvitationToken
	return f
}

type signInForm struct {
	*BaseForm

	Email                  string
	Password               string

}

func newSignInForm(r *http.Request) signInForm {
	body := struct {
		Email                  string `json:"email"`
		Password               string `json:"password"`
	}{}

	f := signInForm{
		BaseForm: &BaseForm{},
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}

	if len(body.Email) != 0 {
		if validEmailRegex.MatchString(body.Email){
			f.Email = body.Email
		} else {
			f.AddValidationError("email__invalid", nil)
		}
	} else {
		f.AddValidationError("email__required", nil)
	}

	if len(body.Password) != 0 {
		f.Password = body.Password
	} else {
		f.AddValidationError("password__required", nil)
	}

	return f
}

type requestPasswordResetForm struct {
	*BaseForm

	Email                  string
}

func newRequestPasswordResetForm(r *http.Request) requestPasswordResetForm {
	body := struct {
		Email                  string `json:"email"`
	}{}

	f := requestPasswordResetForm{
		BaseForm: &BaseForm{},
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}

	if len(body.Email) != 0 {
		if validEmailRegex.MatchString(body.Email){
			f.Email = body.Email
		} else {
			f.AddValidationError("email__invalid", nil)
		}
	} else {
		f.AddValidationError("email__required", nil)
	}

	return f
}

type passwordResetWithTokenForm struct {
	*BaseForm

	Token    string
	Password string
}

func newPasswordResetWithTokenForm(r *http.Request) passwordResetWithTokenForm {
	body := struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}{}

	f := passwordResetWithTokenForm{
		BaseForm: &BaseForm{},
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		f.AddValidationError("body_json__invalid", err)
		return f
	}

	if len(body.Password) != 0 {
		if model.IsStrongPassword(body.Password){
			f.Password = body.Password
		}else {
			f.AddValidationError("password__invalid", nil)
		}
	} else {
		f.AddValidationError("password__required", nil)
	}

	if len(body.Token) != 0 {
		f.Token = body.Token
	} else {
		f.AddValidationError("token__required", nil)
	}

	return f
}