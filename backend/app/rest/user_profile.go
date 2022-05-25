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
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type userProfile struct {
	tokenAuthority TokenAuthority
	emailTemplates template.EmailTemplates

	mailSender  mail.Sender
	linkBuilder LinkBuilder

	userStore    model.UserStore
	accountStore model.AccountStore

	mediaDir  string
	mediaHost string
}

func (up *userProfile) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	currentAccount, _ := model.AccountFromContext(ctx)

	user.UpdateAccountMemberRole(currentAccount)
	user.UpdatePhotoURL(up.mediaHost)
	httpJsonResponse(w, user)
}

func (up *userProfile) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	update := user.NewUpdate()
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if update.FirstName == nil || len(*update.FirstName) == 0 {
		httpResponseError(w, http.StatusBadRequest, "first_name__required")
		return
	}

	if len(*update.FirstName) < 2 {
		httpResponseError(w, http.StatusBadRequest, "first_name__invalid")
		return
	}

	if update.LastName == nil || len(*update.LastName) == 0 {
		httpResponseError(w, http.StatusBadRequest, "last_name__required")
		return
	}

	if len(*update.LastName) < 2 {
		httpResponseError(w, http.StatusBadRequest, "last_name__invalid")
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

	updatedUser, err := up.userStore.Update(ctx, update)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			httpResponseError(w, http.StatusBadRequest, "email__exists")
		} else {
			up.getLog(ctx).Errorf("Fail to update user - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}
		return
	}

	if update.IsEmailUpdated() {
		go up.sendUserEmailConfirmation(NewAsyncTaskContext(ctx), updatedUser)
	}

	currentAccount, _ := model.AccountFromContext(ctx)

	updatedUser.UpdateAccountMemberRole(currentAccount)
	updatedUser.UpdatePhotoURL(up.mediaHost)
	httpJsonResponse(w, updatedUser)
}

func (up *userProfile) updatePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get photo from body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "photo__invalid")
		return
	}

	if !IsImage(bodyBytes) {
		httpResponseError(w, http.StatusBadRequest, "photo__invalid")
		return
	}

	err = os.MkdirAll(up.mediaDir, 0770)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to make photo dir - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	photoFullPath := path.Join(up.mediaDir, "user", "photo")
	err = os.MkdirAll(photoFullPath, 0770)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to make photo dir - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	userPhotoName := path.Join(photoFullPath, user.ID)
	err = ioutil.WriteFile(userPhotoName, bodyBytes, 0770)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to write photo - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	update := user.NewUpdate()
	update.SetPhoto(path.Join("user", "photo", user.ID))
	updatedUser, err := up.userStore.Update(ctx, update)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to update user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	currentAccount, _ := model.AccountFromContext(ctx)

	updatedUser.UpdateAccountMemberRole(currentAccount)
	updatedUser.UpdatePhotoURL(up.mediaHost)
	httpJsonResponse(w, updatedUser)
}

func (up *userProfile) deletePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	if len(user.Photo) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	photoFullPath := path.Join(up.mediaDir, user.Photo)
	_, err = os.Stat(photoFullPath)
	if err != nil {
		if !os.IsNotExist(err) {
			up.getLog(ctx).Errorf("Fail to remove photo file - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	} else {
		err = os.Remove(photoFullPath)
		if err != nil {
			up.getLog(ctx).Errorf("Fail to remove photo file - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	update := user.NewUpdate().ClearPhoto()
	_, err = up.userStore.Update(ctx, update)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to update user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (up *userProfile) confirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Token) == 0 {
		httpResponseError(w, http.StatusBadRequest, "token__required")
		return
	}

	userID, err := up.tokenAuthority.ValidateUserEmailConfirmationToken(r.Context(), requestBody.Token)
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

	user, err := up.userStore.FindByID(ctx, userID)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to fetch user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if user.EmailConfirmedAt != nil {
		httpPlainError(w, http.StatusOK, "Email already confirmed")
		return
	}

	// TODO transaction

	update := user.NewUpdate().SetLogin(user.Email).SetEmailConfirmed()
	updatedUser, err := up.userStore.Update(ctx, update)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to update user - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	// check if user has account. If not - create one
	account, err := up.accountStore.FetchUserAccounts(ctx, updatedUser.ID)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to fetch user accounts - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if len(account) == 0 {
		newAccount := model.NewAccount(updatedUser)
		err = up.accountStore.Add(ctx, newAccount)
		if err != nil {
			up.getLog(ctx).Errorf("Fail to create user account  - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (up *userProfile) sendEmailConfirmationUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	me, err := model.UserFromContext(ctx, model.KCtxKeyUserMe)
	if err != nil {
		up.getLog(ctx).Errorf("Fail to get user from ctx - %v", err)
		httpResponseError(w, http.StatusForbidden, "request__unauthorized")
		return
	}

	if me.IsEmailConfirmed {
		httpResponseError(w, http.StatusForbidden, "email__already_confirmed")
		return
	}

	go up.sendUserEmailConfirmation(NewAsyncTaskContext(ctx), me)

	currentAccount, _ := model.AccountFromContext(ctx)

	me.UpdateAccountMemberRole(currentAccount)
	me.UpdatePhotoURL(up.mediaHost)
	httpJsonResponse(w, me)
}

func (up *userProfile) sendUserEmailConfirmation(c context.Context, user *model.User) {
	ctx, cancel := context.WithTimeout(c, sendEmailTaskTimeout)
	log := logger.GetLogger(ctx)

	log.Info("Try send user email confirmation")

	defer cancel()

	mailOpts := mail.NewSendOpt()
	mailOpts.SetSubject("Your Coton SMS account has been updated")
	mailOpts.SetRecipients(user.Email)

	// TODO make
	//mailOpts.SetPlainText(fmt.Sprintf("Your password is %s", newStaff.Password))

	ecToken, err := up.tokenAuthority.IssueUserEmailConfirmationToken(ctx, user)
	if err != nil {
		log.Errorf("Fail to create confirmation token - %v", err)
		return
	}

	confirmationLink := up.linkBuilder.UserEmailConfirmation(*ecToken)

	tmpl, err := up.emailTemplates.UserEmailConfirmation(confirmationLink, fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	if err != nil {
		log.Errorf("Fail to create html email - %v", err)
		return
	} else {
		mailOpts.SetHtml(tmpl)
	}

	err = up.mailSender.Send(ctx, mailOpts)
	if err != nil {
		log.Errorf("Fail to send mail - %v", err)
		return
	}

	update := user.NewUpdate().SetEmailConfirmationSent()
	_, err = up.userStore.Update(ctx, update)
	if err != nil {
		log.Errorf("Fail to update user - %v", err)
		return
	}

	log.Info("Successfully sent user email confirmation")
}

func (up *userProfile) getLog(ctx context.Context) *logrus.Entry {
	return logger.GetLogger(ctx).WithFields(logrus.Fields{
		"scope": "user_profile",
	})
}
