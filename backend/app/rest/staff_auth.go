package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/model"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type staffAuth struct {
	tokenAuthority TokenAuthority

	store model.StaffStore

	adminPanelHost        string
	secureCookiesDisabled bool
}

func (s *staffAuth) signInEmail(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	log := logger.GetLogger(r.Context())
	log.Info("Sign-in request")

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// do not print json body here, because it contains password
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Email) == 0 {
		httpResponseError(w, http.StatusBadRequest, "email__required")
		return
	}

	if len(requestBody.Password) == 0 {
		httpResponseError(w, http.StatusBadRequest, "password__required")
		return
	}

	staff, err := s.store.FetchStaffByLogin(requestBody.Email)
	if err != nil {
		log.Errorf("fail to fetch staff - %v", err)
		httpResponseError(w, http.StatusForbidden, "credential__invalid")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(requestBody.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			httpResponseError(w, http.StatusForbidden, "credential__invalid")
		} else {
			log.Errorf("fail to compare hash and password - %v", err)
			httpResponseError(w, http.StatusInternalServerError, "internal")
		}

		return
	}

	if !staff.IsActive {
		httpResponseError(w, http.StatusForbidden, "account__suspended")
		return
	}

	if !staff.IsEmailConfirmed {
		httpResponseError(w, http.StatusForbidden, "email__not_confirmed")
		return
	}

	jwtToken, err := s.tokenAuthority.IssueStaffAPIToken(r.Context(), staff)
	if err != nil {
		log.Errorf("fail to generate token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	s.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (s *staffAuth) signOut(w http.ResponseWriter, r *http.Request) {
	s.clearRefreshTokenCookie(w)
	http.StatusText(http.StatusOK)
}

func (s *staffAuth) confirmEmail(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Token string `json:"token"`
	}{}

	log := logger.GetLogger(r.Context())

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Errorf("fail to decode body - %v", err)
		httpResponseError(w, http.StatusBadRequest, "body_json__invalid")
		return
	}

	if len(requestBody.Token) == 0 {
		httpResponseError(w, http.StatusBadRequest, "token__required")
		return
	}

	staffID, err := s.tokenAuthority.ValidateStaffEmailConfirmationToken(r.Context(), requestBody.Token)
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

	staff, err := s.store.FindStaffByID(r.Context(), staffID)
	if err != nil {
		log.Errorf("fail to fetch staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	if staff.EmailConfirmedAt != nil {
		httpPlainError(w, http.StatusOK, "Email already confirmed")
		return
	}

	update := staff.NewUpdate().SetLogin(staff.Email).SetEmailConfirmed()
	_, err = s.store.UpdateStaff(r.Context(), update)
	if err != nil {
		log.Errorf("fail to update staff - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *staffAuth) refreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	ctx := r.Context()
	log := logger.GetLogger(r.Context())

	refreshToken := cookie.Value
	staffId, err := s.tokenAuthority.ValidateStaffRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Errorf("staff refresh token invalid - %v", err)
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	staff, err := s.store.FindStaffByID(r.Context(), staffId)
	if err != nil {
		log.Errorf("fail to fetch staff - %v", err)
		httpResponseError(w, http.StatusBadRequest, "refresh__failed")
		return
	}

	jwtToken, err := s.tokenAuthority.IssueStaffAPIToken(r.Context(), staff)
	if err != nil {
		log.Errorf("fail to generate token - %v", err)
		httpResponseError(w, http.StatusInternalServerError, "internal")
		return
	}

	s.refreshTokenCookie(w, jwtToken)
	httpJsonResponse(w, jwtToken)
}

func (s *staffAuth) refreshTokenCookie(w http.ResponseWriter, jwtToken *jwtToken) {
	cookie := &http.Cookie{
		Name:  "refresh-token",
		Value: jwtToken.RefreshToken,

		Domain:  s.adminPanelHost,
		Path:    "/admin/api/v1/refresh",
		Expires: jwtToken.RefreshTokenExpire,

		Secure:   !s.secureCookiesDisabled,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (s *staffAuth) clearRefreshTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "refresh-token",
		Value: "",

		Domain: s.adminPanelHost,
		MaxAge: -1,

		Secure:   !s.secureCookiesDisabled,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (s *staffAuth) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := tokenFromHeader(r)
		if len(tokenStr) == 0 {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		ctx := r.Context()

		staffId, err := s.tokenAuthority.ValidateStaffApiToken(ctx, tokenStr)
		if err != nil {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		staff, err := s.store.FindStaffByID(ctx, staffId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
				return
			}

			httpResponseError(w, http.StatusInternalServerError, "internal")
			return
		}

		// Staff is found, pass it through
		next.ServeHTTP(w, r.WithContext(staff.NewContext(r.Context(), model.KCtxKeyStaffMe)))
	})
}
