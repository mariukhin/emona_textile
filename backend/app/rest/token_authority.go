package api

import (
	"amifactory.team/sequel/coton-app-backend/app/model"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// Token validation error
var (
	ErrInvalid     = errors.New("token authority: token is malformed, unverifiable or signature is not valid")
	ErrIssuer      = errors.New("token authority: issuer is not valid")
	ErrAudience    = errors.New("token authority: audience is not valid")
	ErrSubject     = errors.New("token authority: subject is not valid or missing")
	ErrTarget      = errors.New("token authority: target is not valid")
	ErrExpired     = errors.New("token authority: token expired")
	ErrAlreadyUsed = errors.New("token authority: token is already used")
	ErrUnexpected  = errors.New("token authority: unexpected err")
)

// Common claims
const (
	kClaimsIssuer         = "iss"
	kClaimsSubject        = "sub"
	kClaimsAudience       = "aud"
	kClaimsExpirationTime = "exp"
	kClaimsJWTId          = "jti"
)

// App claims
const (
	kClaimsTarget                    = "trg"
	kClaimsEmailOwner                = "eow"
	kClaimsEmailConfirmationRequired = "ecr"
	kClaimsEmail                     = "ead"
	kClaimsAccountName               = "acn"
)

const (
	taIssuer = "Coton"

	taAudienceAdminPanel = "Coton Admin Panel"
	taAudienceBackOffice = "Coton Back-office"

	taTargetApi               = "api"
	taTargetEmailConfirmation = "email"
	taTargetInvitation        = "invitation"
	taTargetPasswordReset     = "pass-reset"
	taTargetRefresh           = "refresh"

	taEmailOwnerUser    = "user"
	taEmailOwnerAccount = "account"
)

type jwtToken struct {
	AccessToken        string    `json:"access_token"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpire time.Time `json:"-"`
}

type TokenAuthority interface {
	IssueAccountEmailConfirmationToken(ctx context.Context, a *model.AccountDetails) (*string, error)
	IssueAccountMemberInvitationToken(ctx context.Context, a *model.AccountDetails, m *model.AccountMember) (*string, error)

	IssueStaffAPIToken(ctx context.Context, s *model.Staff) (*jwtToken, error)
	IssueStaffEmailConfirmationToken(ctx context.Context, c *model.Staff) (*string, error)

	IssueUserAPIToken(ctx context.Context, c *model.User) (*jwtToken, error)
	IssueUserEmailConfirmationToken(ctx context.Context, c *model.User) (*string, error)
	IssueUserPasswordResetToken(ctx context.Context, c *model.User) (*string, error)

	ValidateAccountEmailConfirmationToken(ctx context.Context, token string) (string, error)
	ValidateAccountMemberInvitationToken(ctx context.Context, token string) (string, error)

	ValidateStaffApiToken(ctx context.Context, token string) (string, error)
	ValidateStaffRefreshToken(ctx context.Context, token string) (string, error)
	ValidateStaffEmailConfirmationToken(ctx context.Context, token string) (string, error)

	ValidateUserApiToken(ctx context.Context, token string) (string, error)
	ValidateUserRefreshToken(ctx context.Context, token string) (string, error)
	ValidateUserEmailConfirmationToken(ctx context.Context, token string) (string, error)
	ValidateUserPasswordResetToken(ctx context.Context, token string) (string, error)
}

type TokenAuthorityConf struct {
	Secret                         string
	AccessTokenDuration            time.Duration
	RefreshTokenDuration           time.Duration
	EmailConfirmationTokenDuration time.Duration
	InvitationTokenDuration        time.Duration
	PasswordResetTokenDuration     time.Duration
}

func NewTokenAuthorityConf(secret string) TokenAuthorityConf {
	return TokenAuthorityConf{
		Secret:                         secret,
		AccessTokenDuration:            time.Minute * 10,
		RefreshTokenDuration:           time.Hour * 24,
		EmailConfirmationTokenDuration: time.Hour * 12,
		InvitationTokenDuration:        time.Hour * 24,
		PasswordResetTokenDuration:     time.Minute * 10,
	}
}

func NewTokenAuthority(conf TokenAuthorityConf, store model.Storage) (TokenAuthority, error) {
	return &tokenAuthority{
		accessTokenDuration:            conf.AccessTokenDuration,
		refreshTokenDuration:           conf.RefreshTokenDuration,
		emailConfirmationTokenDuration: conf.EmailConfirmationTokenDuration,
		invitationTokenDuration:        conf.InvitationTokenDuration,
		passwordResetTokenDuration:     conf.PasswordResetTokenDuration,
		jwtAuth:                        jwtauth.New("HS256", []byte(conf.Secret), nil),
		storage: &tokenStorage{
			storage: store,
		},
	}, nil
}

type tokenAuthority struct {
	accessTokenDuration            time.Duration
	refreshTokenDuration           time.Duration
	emailConfirmationTokenDuration time.Duration
	invitationTokenDuration        time.Duration
	passwordResetTokenDuration     time.Duration

	jwtAuth *jwtauth.JWTAuth
	storage *tokenStorage
}

func (ta *tokenAuthority) IssueAccountEmailConfirmationToken(ctx context.Context, a *model.AccountDetails) (*string, error) {
	tokenMeta := NewTokenMeta(taTargetEmailConfirmation, ta.emailConfirmationTokenDuration)

	var err error

	ecClaims := jwt.MapClaims{}

	// Common claims
	ecClaims[kClaimsIssuer] = taIssuer
	ecClaims[kClaimsAudience] = taAudienceBackOffice
	ecClaims[kClaimsSubject] = a.ID
	ecClaims[kClaimsExpirationTime] = tokenMeta.ExpiresAt.Unix()
	ecClaims[kClaimsJWTId] = tokenMeta.ID

	// App specific claims
	ecClaims[kClaimsTarget] = taTargetEmailConfirmation
	ecClaims[kClaimsEmail] = a.Email
	ecClaims[kClaimsEmailOwner] = taEmailOwnerAccount
	ecClaims[kClaimsAccountName] = a.Name

	_, ecToken, err := ta.jwtAuth.Encode(ecClaims)
	if err != nil {
		return nil, err
	}

	err = ta.storage.AddToken(ctx, tokenMeta)
	if err != nil {
		return nil, err
	}

	return &ecToken, nil
}

func (ta *tokenAuthority) IssueAccountMemberInvitationToken(ctx context.Context, a *model.AccountDetails, m *model.AccountMember) (*string, error) {
	tokenMeta := NewTokenMeta(taTargetInvitation, ta.emailConfirmationTokenDuration)

	var err error

	ecClaims := jwt.MapClaims{}

	// Common claims
	ecClaims[kClaimsIssuer] = taIssuer
	ecClaims[kClaimsAudience] = taAudienceBackOffice
	ecClaims[kClaimsSubject] = m.ID
	ecClaims[kClaimsExpirationTime] = tokenMeta.ExpiresAt.Unix()
	ecClaims[kClaimsJWTId] = tokenMeta.ID

	// App specific claims
	ecClaims[kClaimsTarget] = taTargetInvitation
	ecClaims[kClaimsAccountName] = a.Name

	_, ecToken, err := ta.jwtAuth.Encode(ecClaims)
	if err != nil {
		return nil, err
	}

	return &ecToken, nil
}

func (ta *tokenAuthority) IssueUserAPIToken(ctx context.Context, c *model.User) (*jwtToken, error) {
	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}

	// Common claims
	atClaims[kClaimsIssuer] = taIssuer
	atClaims[kClaimsAudience] = taAudienceBackOffice
	atClaims[kClaimsSubject] = c.ID
	atClaims[kClaimsExpirationTime] = time.Now().Add(ta.accessTokenDuration).Unix()

	// App specific claims
	atClaims[kClaimsTarget] = taTargetApi
	if c.IsEmailConfirmationRequired() {
		atClaims[kClaimsEmailConfirmationRequired] = true
		atClaims[kClaimsEmail] = c.Email
	}

	_, at, err := ta.jwtAuth.Encode(atClaims)
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}

	rtExpire := time.Now().Add(ta.refreshTokenDuration)

	// Common claims
	rtClaims[kClaimsIssuer] = taIssuer
	rtClaims[kClaimsAudience] = taAudienceBackOffice
	rtClaims[kClaimsSubject] = c.ID
	rtClaims[kClaimsExpirationTime] = rtExpire.Unix()

	// App specific claims
	rtClaims[kClaimsTarget] = taTargetRefresh

	_, rt, err := ta.jwtAuth.Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return &jwtToken{
		AccessToken:        at,
		RefreshToken:       rt,
		RefreshTokenExpire: rtExpire,
	}, nil
}

func (ta *tokenAuthority) IssueUserEmailConfirmationToken(ctx context.Context, c *model.User) (*string, error) {
	tokenMeta := NewTokenMeta(taTargetEmailConfirmation, ta.emailConfirmationTokenDuration)

	var err error

	ecClaims := jwt.MapClaims{}

	// Common claims
	ecClaims[kClaimsIssuer] = taIssuer
	ecClaims[kClaimsAudience] = taAudienceBackOffice
	ecClaims[kClaimsSubject] = c.ID
	ecClaims[kClaimsExpirationTime] = tokenMeta.ExpiresAt.Unix()
	ecClaims[kClaimsJWTId] = tokenMeta.ID

	// App specific claims
	ecClaims[kClaimsTarget] = taTargetEmailConfirmation
	ecClaims[kClaimsEmail] = c.Email
	ecClaims[kClaimsEmailOwner] = taEmailOwnerUser

	_, ecToken, err := ta.jwtAuth.Encode(ecClaims)
	if err != nil {
		return nil, err
	}

	err = ta.storage.AddToken(ctx, tokenMeta)
	if err != nil {
		return nil, err
	}

	return &ecToken, nil
}

func (ta *tokenAuthority) IssueUserPasswordResetToken(ctx context.Context, c *model.User) (*string, error) {
	tokenMeta := NewTokenMeta(taTargetPasswordReset, ta.passwordResetTokenDuration)

	var err error

	ecClaims := jwt.MapClaims{}

	// Common claims
	ecClaims[kClaimsIssuer] = taIssuer
	ecClaims[kClaimsAudience] = taAudienceBackOffice
	ecClaims[kClaimsSubject] = c.ID
	ecClaims[kClaimsExpirationTime] = tokenMeta.ExpiresAt.Unix()
	ecClaims[kClaimsJWTId] = tokenMeta.ID

	// App specific claims
	ecClaims[kClaimsTarget] = taTargetPasswordReset

	_, ecToken, err := ta.jwtAuth.Encode(ecClaims)
	if err != nil {
		return nil, err
	}

	err = ta.storage.AddToken(ctx, tokenMeta)
	if err != nil {
		return nil, err
	}

	return &ecToken, nil
}

func (ta *tokenAuthority) IssueStaffAPIToken(ctx context.Context, s *model.Staff) (*jwtToken, error) {
	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}

	// Common claims
	atClaims[kClaimsIssuer] = taIssuer
	atClaims[kClaimsAudience] = taAudienceAdminPanel
	atClaims[kClaimsSubject] = s.ID
	atClaims[kClaimsExpirationTime] = time.Now().Add(ta.accessTokenDuration).Unix()

	// App specific claims
	atClaims[kClaimsTarget] = taTargetApi
	if s.IsEmailConfirmationRequired() {
		atClaims[kClaimsEmailConfirmationRequired] = true
		atClaims[kClaimsEmail] = s.Email
	}

	_, at, err := ta.jwtAuth.Encode(atClaims)
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}

	// Common claims
	rtClaims[kClaimsIssuer] = taIssuer
	rtClaims[kClaimsAudience] = taAudienceAdminPanel
	rtClaims[kClaimsSubject] = s.ID
	rtClaims[kClaimsExpirationTime] = time.Now().Add(ta.refreshTokenDuration).Unix()

	// App specific claims
	rtClaims[kClaimsTarget] = taTargetRefresh

	_, rt, err := ta.jwtAuth.Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return &jwtToken{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (ta *tokenAuthority) IssueStaffEmailConfirmationToken(ctx context.Context, s *model.Staff) (*string, error) {
	tokenMeta := NewTokenMeta(taTargetEmailConfirmation, ta.emailConfirmationTokenDuration)

	var err error

	ecClaims := jwt.MapClaims{}

	// Common claims
	ecClaims[kClaimsIssuer] = taIssuer
	ecClaims[kClaimsAudience] = taAudienceAdminPanel
	ecClaims[kClaimsSubject] = s.ID
	ecClaims[kClaimsExpirationTime] = tokenMeta.ExpiresAt.Unix()
	ecClaims[kClaimsJWTId] = tokenMeta.ID

	// App specific claims
	ecClaims[kClaimsTarget] = taTargetEmailConfirmation
	ecClaims[kClaimsEmail] = s.Email

	_, ecToken, err := ta.jwtAuth.Encode(ecClaims)
	if err != nil {
		return nil, err
	}

	err = ta.storage.AddToken(ctx, tokenMeta)
	if err != nil {
		return nil, err
	}

	return &ecToken, nil
}

func (ta *tokenAuthority) ValidateAccountEmailConfirmationToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetEmailConfirmation {
		return "", ErrTarget
	}

	emailOwner, ok := tokenClaims[kClaimsEmailOwner].(string)
	if !ok || emailOwner != taEmailOwnerAccount {
		return "", ErrUnexpected
	}

	tokenID, ok := tokenClaims[kClaimsJWTId].(string)
	if !ok {
		return "", ErrUnexpected
	}

	err = ta.storage.RemoveTokenById(ctx, tokenID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", ErrAlreadyUsed
		}
	}

	accountID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return accountID, nil
}

func (ta *tokenAuthority) ValidateAccountMemberInvitationToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetInvitation {
		return "", ErrTarget
	}

	memberID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return memberID, nil
}

func (ta *tokenAuthority) ValidateStaffApiToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceAdminPanel, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetApi {
		return "", ErrTarget
	}

	staffID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return staffID, nil
}

func (ta *tokenAuthority) ValidateStaffRefreshToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceAdminPanel, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetRefresh {
		return "", ErrTarget
	}

	staffID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return staffID, nil
}

func (ta *tokenAuthority) ValidateStaffEmailConfirmationToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceAdminPanel, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetEmailConfirmation {
		return "", ErrTarget
	}

	tokenID, ok := tokenClaims[kClaimsJWTId].(string)
	if !ok {
		return "", ErrUnexpected
	}

	err = ta.storage.RemoveTokenById(ctx, tokenID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", ErrAlreadyUsed
		}
	}

	staffID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return staffID, nil
}

func (ta *tokenAuthority) ValidateUserApiToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetApi {
		return "", ErrTarget
	}

	customerID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return customerID, nil
}

func (ta *tokenAuthority) ValidateUserRefreshToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetRefresh {
		return "", ErrTarget
	}

	customerID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return customerID, nil
}

func (ta *tokenAuthority) ValidateUserEmailConfirmationToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetEmailConfirmation {
		return "", ErrTarget
	}

	emailOwner, ok := tokenClaims[kClaimsEmailOwner].(string)
	if !ok || emailOwner != taEmailOwnerUser {
		return "", ErrUnexpected
	}

	tokenID, ok := tokenClaims[kClaimsJWTId].(string)
	if !ok {
		return "", ErrUnexpected
	}

	err = ta.storage.RemoveTokenById(ctx, tokenID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", ErrAlreadyUsed
		}
	}

	customerID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return customerID, nil
}

func (ta *tokenAuthority) ValidateUserPasswordResetToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := ta.jwtAuth.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return "", ErrExpired
			}
		}

		return "", ErrUnexpected
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnexpected
	}

	if !tokenClaims.VerifyIssuer(taIssuer, true) {
		return "", ErrIssuer
	}

	if !tokenClaims.VerifyAudience(taAudienceBackOffice, true) {
		return "", ErrAudience
	}

	target, ok := tokenClaims[kClaimsTarget].(string)
	if !ok || target != taTargetPasswordReset {
		return "", ErrTarget
	}

	tokenID, ok := tokenClaims[kClaimsJWTId].(string)
	if !ok {
		return "", ErrUnexpected
	}

	err = ta.storage.RemoveTokenById(ctx, tokenID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", ErrAlreadyUsed
		}
	}

	customerID, ok := tokenClaims[kClaimsSubject].(string)
	if !ok {
		return "", ErrSubject
	}

	return customerID, nil
}

type tokenStorage struct {
	storage model.Storage
}

func (store *tokenStorage) AddToken(ctx context.Context, token TokenMeta) error {
	_, err := store.storage.Collection("token").InsertOne(ctx, token)
	if err != nil {
		// TODO check duplication
		//if IsErrDuplication(err) {
		//	return ErrDuplicate
		//}

		return err
	}

	return nil
}

func (store *tokenStorage) RemoveTokenById(ctx context.Context, tokenID string) error {
	filter := bson.M{"_id": tokenID}

	res, err := store.storage.Collection("token").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount != 1 {
		return model.ErrNotFound
	}

	return nil
}

type TokenMeta struct {
	ID        string    `bson:"_id"`
	Type      string    `bson:"type"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func NewTokenMeta(ttype string, duration time.Duration) TokenMeta {
	return TokenMeta{
		ID:        uuid.NewV4().String(),
		Type:      ttype,
		ExpiresAt: time.Now().Add(duration),
	}
}
