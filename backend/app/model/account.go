package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"errors"
	"fmt"
	"github.com/leekchan/accounting"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AccountStatus string

const (
	AccountStatusCreated AccountStatus = "created"
	AccountStatusActive  AccountStatus = "active"
	AccountStatusBlocked AccountStatus = "blocked"
)

var accountStatuses = []AccountStatus{
	AccountStatusCreated,
	AccountStatusActive,
	AccountStatusBlocked,
}

func (as AccountStatus) IsValid() bool {
	for _, opt := range accountStatuses {
		if opt == as {
			return true
		}
	}

	return false
}

type AccountBalance struct {
	UpdatedAt time.Time `bson:"updated_at" json:"-"`
	Balance   int       `bson:"balance" json:"balance"`
	Currency  Currency  `bson:"currency" json:"currency"`
}

type AccountAddress struct {
	UpdatedAt  time.Time `bson:"updated_at" json:"-"`
	Address    string    `bson:"address" json:"address"`
	City       string    `bson:"city" json:"city"`
	State      string    `bson:"state" json:"state"`
	Country    string    `bson:"country" json:"country"`
	PostalCode string    `bson:"postal_code" json:"postal_code"`
}

type AccountBillingType string

const (
	AccountBillingTypePrepay  AccountBillingType = "prepay"
	AccountBillingTypePostpay AccountBillingType = "postpay"
)

type AccountBillingPeriod string

const (
	AccountBillingPeriodDay        AccountBillingPeriod = "day"
	AccountBillingPeriodWeek       AccountBillingPeriod = "week"
	AccountBillingPeriodTwoWeeks   AccountBillingPeriod = "two_weeks"
	AccountBillingPeriodMonthTwice AccountBillingPeriod = "month_twice"
	AccountBillingPeriodMonthOnce  AccountBillingPeriod = "month_once"
)

type AccountBilling struct {
	UpdatedAt time.Time            `bson:"updated_at" json:"-"`
	Currency  Currency             `bson:"currency" json:"-"`
	Type      AccountBillingType   `bson:"type" json:"-"`
	Period    AccountBillingPeriod `bson:"period" json:"-"`
	Delay     int                  `bson:"delay" json:"-"`
}

type AccountModeration string

const (
	AccountModerationOn  AccountModeration = "on"
	AccountModerationOff AccountModeration = "off"
)

var accountModerationOptions = []AccountModeration{
	AccountModerationOn,
	AccountModerationOff,
}

func (am AccountModeration) IsValid() bool {
	for _, opt := range accountModerationOptions {
		if opt == am {
			return true
		}
	}

	return false
}

type AccountMemberStatus string

const (
	AccountMemberStatusInvited  AccountMemberStatus = "invited"
	AccountMemberStatusActive   AccountMemberStatus = "active"
	AccountMemberStatusRejected AccountMemberStatus = "rejected"
)

type AccountMemberRole string

const (
	AccountMemberRoleOwner AccountMemberRole = "owner"
	AccountMemberRoleAdmin AccountMemberRole = "admin"
	AccountMemberRoleUser  AccountMemberRole = "user"
)

var accountMemberRoles = []AccountMemberRole{
	AccountMemberRoleOwner,
	AccountMemberRoleAdmin,
	AccountMemberRoleUser,
}

func (amr AccountMemberRole) IsAccountMemberRoleValid() bool {
	for _, opt := range accountMemberRoles {
		if opt == amr {
			return true
		}
	}

	return false
}

var AllowedAccountMemberRoles = []AccountMemberRole{
	AccountMemberRoleAdmin,
	AccountMemberRoleUser,
}

type AccountMemberUser struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	FirstName string `bson:"first_name" json:"fist_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
}

type AccountMember struct {
	ID        string `bson:"_id" json:"id"`
	AccountId string `bson:"-" json:"-"`

	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"-"`

	UserId string             `bson:"user_id,omitempty" json:"-"`
	User   *AccountMemberUser `bson:"user,omitempty" json:"user,omitempty"`

	Status AccountMemberStatus `bson:"status" json:"status"`
	Role   AccountMemberRole   `bson:"role" json:"role"`

	InvitationExpiresAt time.Time  `bson:"invitation_expires_at" json:"-"`
	InvitationEmail     string     `bson:"invitation_email" json:"-"`
	InvitationSentAt    *time.Time `bson:"invitation_sent_at,omitempty" json:"-"`
}

const kCtxKeyAccountMember = "account-member"

func (m AccountMember) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, kCtxKeyAccountMember, &m)
	return ctx
}

func MemberFromContext(ctx context.Context) (*AccountMember, error) {
	member, _ := ctx.Value(kCtxKeyAccountMember).(*AccountMember)
	if member == nil {
		return nil, ErrNotFound
	}

	return member, nil
}

func (m *AccountMember) NewUpdate() *AccountMemberUpdate {
	return &AccountMemberUpdate{
		member: m,
	}
}

func (m *AccountMember) IsInvitationExpired() bool {
	return time.Now().After(m.InvitationExpiresAt)
}

type AccountMemberUpdate struct {
	UpdatedAt time.Time `json:"-"`

	UserID              *string
	Status              *AccountMemberStatus `json:"-"`
	Role                *AccountMemberRole   `json:"-"`
	InvitationSentAt    *time.Time           `json:"-"`
	InvitationExpiresAt *time.Time           `json:"-"`

	member *AccountMember
}

func (update *AccountMemberUpdate) SetStatus(status AccountMemberStatus) *AccountMemberUpdate {
	update.Status = &status
	return update
}

func (update *AccountMemberUpdate) SetUser(user *User) *AccountMemberUpdate {
	update.UserID = &user.ID
	return update
}

func (update *AccountMemberUpdate) SetRole(role AccountMemberRole) *AccountMemberUpdate {
	update.Role = &role
	return update
}

func (update *AccountMemberUpdate) SetInvitationSent() *AccountMemberUpdate {
	sentAt := time.Now()
	expiresAt := time.Now().Add(time.Hour * 24)
	update.InvitationSentAt = &sentAt
	update.InvitationExpiresAt = &expiresAt
	return update
}

func (update *AccountMemberUpdate) setUpdatedDate() *AccountMemberUpdate {
	update.UpdatedAt = time.Now()
	return update
}

func (update *AccountMemberUpdate) Bson() bson.D {
	res := bson.D{}

	res = append(res, bson.E{Key: "members.$.updated_at", Value: update.UpdatedAt})

	if update.UserID != nil {
		res = append(res, bson.E{Key: "members.$.user_id", Value: update.UserID})
	}

	if update.Status != nil {
		res = append(res, bson.E{Key: "members.$.status", Value: update.Status})
	}

	if update.Role != nil {
		res = append(res, bson.E{Key: "members.$.role", Value: update.Role})
	}

	if update.InvitationSentAt != nil {
		res = append(res, bson.E{Key: "members.$.invitation_sent_at", Value: update.InvitationSentAt})
	}

	if update.InvitationExpiresAt != nil {
		res = append(res, bson.E{Key: "members.$.invitation_expires_at", Value: update.InvitationExpiresAt})
	}

	return res
}

type Account struct {
	ID string `bson:"_id" json:"id"`

	Name string `bson:"name" json:"name"`

	Status AccountStatus `bson:"status" json:"status"`
	Email  string        `bson:"email" json:"main_email"`

	Balance          AccountBalance `bson:"balance" json:"-"`
	FormattedBalance string         `bson:"-" json:"balance"`
}

func (a *Account) FormatBalance() {
	// TODO
	ac := accounting.Accounting{Symbol: "€", Precision: 2, Thousand: ".", Decimal: ","}
	a.FormattedBalance = ac.FormatMoneyFloat64(float64(a.Balance.Balance) / 100.0)
}

type UserAccount struct {
	ID string `bson:"_id" json:"id"`

	Name string `bson:"name" json:"name"`

	Status AccountStatus     `bson:"status" json:"status"`
	Role   AccountMemberRole `bson:"role" json:"role"`

	Current bool `bson:"-" json:"current"`
}

const kCtxKeyAccount = "account"

type AccountDetails struct {
	ID string `bson:"_id" json:"id"`

	CreatedAt time.Time  `bson:"created_at" json:"-"`
	UpdatedAt time.Time  `bson:"updated_at" json:"-"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"-"`

	Name string `bson:"name" json:"name"`

	Email                   string     `bson:"email" json:"main_email"`
	IsEmailConfirmed        bool       `bson:"is_email_confirmed" json:"is_email_confirmed"`
	EmailConfirmedAt        *time.Time `bson:"email_confirmed_at" json:"email_confirmed_at,omitempty"`
	EmailConfirmationSentAt *time.Time `bson:"email_confirmation_sent_at" json:"email_confirmation_sent_at,omitempty"`

	Status AccountStatus     `bson:"status" json:"status"`
	Role   AccountMemberRole `bson:"-" json:"role,omitempty"`

	Members []*AccountMember `bson:"members" json:"-"`

	Balance        AccountBalance `bson:"balance" json:"-"`
	ServiceAddress AccountAddress `bson:"service_address" json:"-"`
	Billing        AccountBilling `bson:"billing" json:"-"`

	Moderation AccountModeration `bson:"moderation" json:"-"`
}

func (a *AccountDetails) NewMember(user *User, invitationEmail string, role AccountMemberRole) *AccountMember {
	var userID string
	if user != nil {
		userID = user.ID
	}

	member := AccountMember{
		ID:                  uuid.NewV4().String(),
		AccountId:           a.ID,
		UserId:              userID,
		Status:              AccountMemberStatusInvited,
		Role:                role,
		InvitationExpiresAt: time.Now().Add(time.Hour * 24),
		InvitationEmail:     invitationEmail,
	}

	return &member
}

func NewAccount(owner *User) *AccountDetails {
	account := AccountDetails{
		ID:   uuid.NewV4().String(),
		Name: fmt.Sprintf("%s %s", owner.FirstName, owner.LastName),

		Email:            owner.Email,
		IsEmailConfirmed: true,

		Status: AccountStatusCreated,

		Balance: AccountBalance{
			Currency: Euro,
		},
		Billing: AccountBilling{
			Currency: Euro,
			Type:     AccountBillingTypePrepay,
			Period:   AccountBillingPeriodMonthOnce,
			Delay:    0,
		},
		Moderation: AccountModerationOn,
	}

	member := AccountMember{
		ID:        uuid.NewV4().String(),
		AccountId: account.ID,
		UserId:    owner.ID,
		Status:    AccountMemberStatusActive,
		Role:      AccountMemberRoleOwner,
	}
	account.Members = append(account.Members, &member)

	return &account
}

func (a *AccountDetails) NewUpdate() *UserAccountUpdate {
	return &UserAccountUpdate{
		account: a,
	}
}

func (a *AccountDetails) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, kCtxKeyAccount, a)
	return ctx
}

func (a *AccountDetails) Details() interface{} {
	details := struct {
		ID string `json:"id"`

		CreatedAt time.Time  `json:"-"`
		UpdatedAt time.Time  `json:"-"`
		DeletedAt *time.Time `json:"-"`

		Name string `json:"name"`

		Email                   string     `json:"main_email"`
		IsEmailConfirmed        bool       `json:"is_email_confirmed"`
		EmailConfirmedAt        *time.Time `json:"email_confirmed_at,omitempty"`
		EmailConfirmationSentAt *time.Time `json:"email_confirmation_sent_at,omitempty"`

		Status AccountStatus `json:"status"`

		Moderation AccountModeration `json:"moderation"`
	}{
		ID:                      a.ID,
		CreatedAt:               a.CreatedAt,
		UpdatedAt:               a.UpdatedAt,
		DeletedAt:               a.DeletedAt,
		Name:                    a.Name,
		Email:                   a.Email,
		IsEmailConfirmed:        a.IsEmailConfirmed,
		EmailConfirmedAt:        a.EmailConfirmedAt,
		EmailConfirmationSentAt: a.EmailConfirmationSentAt,
		Status:                  a.Status,
		Moderation:              a.Moderation,
	}

	return &details
}

func AccountFromContext(ctx context.Context) (*AccountDetails, error) {
	account, _ := ctx.Value(kCtxKeyAccount).(*AccountDetails)

	if account == nil {
		return nil, ErrNotFound
	}

	return account, nil
}

type UserAccountUpdate struct {
	UpdatedAt time.Time `json:"-"`

	Name *string `json:"name"`

	Email                   *string    `json:"email"`
	IsEmailConfirmed        *bool      `json:"-"`
	EmailConfirmedAt        *time.Time `json:"-"`
	EmailConfirmationSentAt *time.Time `json:"-"`

	ServiceAddress *AccountAddress    `json:"-"`
	Status         *AccountStatus     `json:"-"`
	Moderation     *AccountModeration `json:"-"`

	account *AccountDetails
}

func (update *UserAccountUpdate) IsEmailUpdated() bool {
	return update.Email != nil && update.account.Email != *update.Email
}

func (update *UserAccountUpdate) SetEmailConfirmationSent() *UserAccountUpdate {
	now := time.Now()
	update.EmailConfirmationSentAt = &now
	return update
}

func (update *UserAccountUpdate) SetEmailConfirmed() *UserAccountUpdate {
	v := true
	update.IsEmailConfirmed = &v

	now := time.Now()
	update.EmailConfirmedAt = &now
	return update
}

func (update *UserAccountUpdate) SetStatus(status AccountStatus) *UserAccountUpdate {
	update.Status = &status
	return update
}

func (update *UserAccountUpdate) SetModeration(moderation AccountModeration) *UserAccountUpdate {
	update.Moderation = &moderation
	return update
}

func (update *UserAccountUpdate) SetServiceAddress(address AccountAddress) *UserAccountUpdate {
	update.ServiceAddress = &address
	return update
}

func (update *UserAccountUpdate) setUpdatedDate() *UserAccountUpdate {
	update.UpdatedAt = time.Now()
	if update.ServiceAddress != nil {
		update.ServiceAddress.UpdatedAt = update.UpdatedAt
	}
	return update
}

func (update *UserAccountUpdate) Bson() bson.D {
	res := bson.D{}

	res = append(res, bson.E{Key: "updated_at", Value: update.UpdatedAt})

	if update.Name != nil {
		res = append(res, bson.E{Key: "name", Value: update.Name})
	}

	if update.Status != nil {
		res = append(res, bson.E{Key: "status", Value: update.Status})
	}

	if update.Moderation != nil {
		res = append(res, bson.E{Key: "moderation", Value: update.Moderation})
	}

	if update.IsEmailUpdated() {
		// TODO
		res = append(res, bson.E{Key: "email", Value: update.Email})
		res = append(res, bson.E{Key: "is_email_confirmed", Value: false})
		res = append(res, bson.E{Key: "email_confirmed_at", Value: nil})
		res = append(res, bson.E{Key: "email_confirmation_sent_at", Value: nil})
	}

	if update.IsEmailConfirmed != nil {
		res = append(res, bson.E{Key: "is_email_confirmed", Value: update.IsEmailConfirmed})
	}

	if update.EmailConfirmedAt != nil {
		res = append(res, bson.E{Key: "email_confirmed_at", Value: *update.EmailConfirmedAt})
	}

	if update.EmailConfirmationSentAt != nil {
		res = append(res, bson.E{Key: "email_confirmation_sent_at", Value: *update.EmailConfirmationSentAt})
	}

	if update.ServiceAddress != nil {
		res = append(res, bson.E{Key: "service_address", Value: *update.ServiceAddress})
	}

	return res
}

type AccountStore interface {
	EnsureIndexes(ctx context.Context) error

	Add(ctx context.Context, account *AccountDetails) error
	Update(ctx context.Context, update *UserAccountUpdate) (*AccountDetails, error)

	FindAccountByID(ctx context.Context, accountId string) (*AccountDetails, error)
	FindUserAccountByID(ctx context.Context, accountId string) (*AccountDetails, error)
	FetchAccounts(ctx context.Context, offset, maxCount int, searchPhase *string, status *AccountStatus) ([]*Account, int, error)
	FetchAllAccounts(ctx context.Context, searchPhase *string, status *AccountStatus) ([]*Account, error)
	FetchUserAccounts(ctx context.Context, userId string) ([]*UserAccount, error)
	FindUserAccount(ctx context.Context, userId, accountId string) (*AccountDetails, error)

	AddMember(ctx context.Context, member *AccountMember) (*AccountMember, error)
	UpdateMember(ctx context.Context, update *AccountMemberUpdate) error
	FindAccountMember(ctx context.Context, memberID string) (*UserAccount, *AccountMember, error)
	FindUserAccountMember(ctx context.Context, accountId, userId string) (*AccountMember, error)
	DeleteAccountMember(ctx context.Context, member *AccountMember) error
}

func NewAccountStore(storage Storage) (AccountStore, error) {
	return &accountStore{
		storage: storage,
	}, nil
}

type accountStore struct {
	storage Storage
}

func (store *accountStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "account")

	log.Info("Fetching Account storage indexes")
	log.Info("User Account indexes are up-to-date")

	return nil
}

func (store *accountStore) Add(ctx context.Context, account *AccountDetails) error {
	updatedAt := time.Now()
	account.CreatedAt = updatedAt
	account.UpdatedAt = updatedAt
	account.Balance.UpdatedAt = updatedAt
	account.ServiceAddress.UpdatedAt = updatedAt
	account.Billing.UpdatedAt = updatedAt

	for _, member := range account.Members {
		member.CreatedAt = updatedAt
		member.UpdatedAt = updatedAt
	}

	_, err := store.storage.Collection("account").InsertOne(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (store *accountStore) Update(ctx context.Context, update *UserAccountUpdate) (*AccountDetails, error) {
	filter := bson.M{"_id": update.account.ID}
	updateBson := bson.M{"$set": update.setUpdatedDate().Bson()}

	var updatedAccount AccountDetails
	err := store.storage.Collection("account").FindOneAndUpdate(ctx, filter, updateBson, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedAccount)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		if IsErrDuplication(err) {
			return nil, ErrDuplicate
		}

		return nil, err
	}

	updatedAccount.Role = update.account.Role
	updatedAccount.Members = update.account.Members

	return &updatedAccount, nil
}

func (store *accountStore) FindAccountByID(ctx context.Context, accountId string) (*AccountDetails, error) {
	// TODO simplify
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{
				Key:   "_id",
				Value: accountId,
			},
		},
	}}

	unwindStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			"path",
			"$members",
		}},
	}}

	lookupStage := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{
				"from",
				"user",
			},
			{
				"localField",
				"members.user_id",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"user",
			},
		},
	}}

	addFieldsStage := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"members.user",
			bson.D{{
				"$arrayElemAt", bson.A{"$user", 0},
			}},
		}},
	}}

	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{
				"_id",
				bson.D{
					{
						"_id", "$_id",
					},
					{
						"created_at", "$created_at",
					},
					{
						"updated_at", "$updated_at",
					},
					{
						"name", "$name",
					},
					{
						"email", "$email",
					},
					{
						"is_email_confirmed", "$is_email_confirmed",
					},
					{
						"email_confirmed_at", "$email_confirmed_at",
					},
					{
						"email_confirmation_sent_at", "$email_confirmation_sent_at",
					},
					{
						"status", "$status",
					},
					{
						"balance", "$balance",
					},
					{
						"service_address", "$service_address",
					},
					{
						"billing", "$billing",
					},
					{
						"moderation", "$moderation",
					},
				},
			},
			{
				"members",
				bson.D{{
					"$push", "$members",
				}},
			},
		},
	}}

	replaceWithStage := bson.D{{
		"$replaceWith",
		bson.D{
			{
				"_id", "$_id._id",
			},
			{
				"created_at", "$_id.created_at",
			},
			{
				"updated_at", "$_id.updated_at",
			},
			{
				"name", "$_id.name",
			},
			{
				"email", "$_id.email",
			},
			{
				"is_email_confirmed", "$_id.is_email_confirmed",
			},
			{
				"email_confirmed_at", "$_id.email_confirmed_at",
			},
			{
				"email_confirmation_sent_at", "$_id.email_confirmation_sent_at",
			},
			{
				"status", "$_id.status",
			},
			{
				"balance", "$_id.balance",
			},
			{
				"service_address", "$_id.service_address",
			},
			{
				"billing", "$_id.billing",
			},
			{
				"moderation", "$_id.moderation",
			},
			{
				"members", "$members",
			},
		},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		unwindStage,
		lookupStage,
		addFieldsStage,
		groupStage,
		replaceWithStage,
	}

	opts := options.Aggregate()
	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("fail to fetch account - %v", err)
		}
	}

	defer cur.Close(ctx)

	var account AccountDetails

	if cur.Next(ctx) {
		err = cur.Decode(&account)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch account - %v", err)
		}

		for _, m := range account.Members {
			m.AccountId = account.ID

			if m.User == nil {
				m.User = &AccountMemberUser{
					FirstName: "Invited",
					LastName:  "user",
					Email:     m.InvitationEmail,
				}
			}
		}

		return &account, nil
	}

	return nil, ErrNotFound
}

func (store *accountStore) FindUserAccountByID(ctx context.Context, accountId string) (*AccountDetails, error) {
	// TODO simplify
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{
				Key:   "_id",
				Value: accountId,
			},
			{
				Key: "$or",
				Value: bson.A{
					bson.D{{
						Key: "status", Value: AccountStatusCreated,
					}},
					bson.D{{
						Key: "status", Value: AccountStatusActive,
					}},
				},
			},
		},
	}}

	addFieldsStage1 := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"m", bson.D{{
				"$arrayElemAt", bson.A{
					"$members",
					bson.D{{
						"$indexOfArray",
						bson.A{"$members", bson.D{
							{
								"status", "active",
							},
						},
						},
					}},
				},
			}},
		}},
	}}

	addFieldsStage2 := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"role", "$m.role",
		}},
	}}

	projectStage := bson.D{{
		Key: "$project",
		Value: bson.D{{
			"m", 0,
		}},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		addFieldsStage1,
		addFieldsStage2,
		projectStage,
	}

	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	defer cur.Close(ctx)

	var account AccountDetails
	if cur.Next(ctx) {
		err = cur.Decode(&account)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch account - %v", err)
		}

		for _, m := range account.Members {
			m.AccountId = account.ID

			if m.User == nil {
				m.User = &AccountMemberUser{
					FirstName: "Invited",
					LastName:  "user",
					Email:     m.InvitationEmail,
				}
			}
		}

		return &account, nil
	}

	return nil, ErrNotFound
}

func (store *accountStore) FetchAccounts(ctx context.Context, offset, maxCount int, searchPhase *string, status *AccountStatus) ([]*Account, int, error) {
	filter := bson.D{}
	if searchPhase != nil {
		searchRegex := fmt.Sprintf("^%s.*", *searchPhase)
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{"name", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
				bson.D{{"email", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
			},
		})
	}

	if status != nil {
		filter = append(filter, bson.E{
			Key: "status", Value: string(*status),
		})
	}

	count, err := store.storage.Collection("account").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch accounts count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(maxCount))

	accounts := make([]*Account, 0)
	cur, err := store.storage.Collection("account").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	for _, a := range accounts {
		a.FormatBalance()
	}

	return accounts, int(count), nil
}

func (store *accountStore) FetchAllAccounts(ctx context.Context, searchPhase *string, status *AccountStatus) ([]*Account, error) {
	filter := bson.D{}
	if searchPhase != nil {
		searchRegex := fmt.Sprintf("^%s.*", *searchPhase)
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{"name", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
				bson.D{{"email", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
			},
		})
	}

	if status != nil {
		filter = append(filter, bson.E{
			Key: "status", Value: string(*status),
		})
	}

	opts := options.Find()
	accounts := make([]*Account, 0)
	cur, err := store.storage.Collection("account").Find(ctx, filter, opts)
	if err != nil {
		return nil,  fmt.Errorf("fail to fetch accounts - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	for _, a := range accounts {
		a.FormatBalance()
	}

	return accounts, nil
}

func (store *accountStore) FetchUserAccounts(ctx context.Context, userId string) ([]*UserAccount, error) {
	// TODO simplify
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{
				Key: "$or",
				Value: bson.A{
					bson.D{{
						Key: "status", Value: AccountStatusCreated,
					}},
					bson.D{{
						Key: "status", Value: AccountStatusActive,
					}},
				},
			},
			{
				Key: "members",
				Value: bson.D{{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "user_id",
							Value: userId,
						},
						{
							Key:   "status",
							Value: AccountMemberStatusActive,
						},
					},
				}},
			},
		},
	}}

	addFieldsStage1 := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"m", bson.D{{
				"$arrayElemAt", bson.A{
					"$members",
					bson.D{{
						"$indexOfArray",
						bson.A{"$members", bson.D{
							{
								"user_id", userId,
							},
							{
								"status", "active",
							},
						},
						},
					}},
				},
			}},
		}},
	}}

	addFieldsStage2 := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"role", "$m.role",
		}},
	}}

	projectStage := bson.D{{
		Key: "$project",
		Value: bson.D{{
			"m", 0,
		}},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		addFieldsStage1,
		addFieldsStage2,
		projectStage,
	}

	accounts := make([]*UserAccount, 0)
	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	return accounts, nil
}

func (store *accountStore) FindUserAccount(ctx context.Context, userId, accountId string) (*AccountDetails, error) {
	// TODO simplify
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{
				Key:   "_id",
				Value: accountId,
			},
			{
				Key: "members",
				Value: bson.D{{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "user_id",
							Value: userId,
						},
						{
							Key:   "status",
							Value: AccountMemberStatusActive,
						},
					},
				}},
			},
		},
	}}

	unwindStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			"path",
			"$members",
		}},
	}}

	lookupStage := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{
				"from",
				"user",
			},
			{
				"localField",
				"members.user_id",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"user",
			},
		},
	}}

	addFieldsStage := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"members.user",
			bson.D{{
				"$arrayElemAt", bson.A{"$user", 0},
			}},
		}},
	}}

	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{
				"_id",
				bson.D{
					{
						"_id", "$_id",
					},
					{
						"created_at", "$created_at",
					},
					{
						"updated_at", "$updated_at",
					},
					{
						"name", "$name",
					},
					{
						"email", "$email",
					},
					{
						"is_email_confirmed", "$is_email_confirmed",
					},
					{
						"email_confirmed_at", "$email_confirmed_at",
					},
					{
						"email_confirmation_sent_at", "$email_confirmation_sent_at",
					},
					{
						"status", "$status",
					},
					{
						"balance", "$balance",
					},
					{
						"service_address", "$service_address",
					},
					{
						"billing", "$billing",
					},
					{
						"moderation", "$moderation",
					},
				},
			},
			{
				"members",
				bson.D{{
					"$push", "$members",
				}},
			},
		},
	}}

	replaceWithStage := bson.D{{
		"$replaceWith",
		bson.D{
			{
				"_id", "$_id._id",
			},
			{
				"created_at", "$_id.created_at",
			},
			{
				"updated_at", "$_id.updated_at",
			},
			{
				"name", "$_id.name",
			},
			{
				"email", "$_id.email",
			},
			{
				"is_email_confirmed", "$_id.is_email_confirmed",
			},
			{
				"email_confirmed_at", "$_id.email_confirmed_at",
			},
			{
				"email_confirmation_sent_at", "$_id.email_confirmation_sent_at",
			},
			{
				"status", "$_id.status",
			},
			{
				"balance", "$_id.balance",
			},
			{
				"service_address", "$_id.service_address",
			},
			{
				"billing", "$_id.billing",
			},
			{
				"moderation", "$_id.moderation",
			},
			{
				"members", "$members",
			},
		},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		unwindStage,
		lookupStage,
		addFieldsStage,
		groupStage,
		replaceWithStage,
	}

	opts := options.Aggregate()
	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("fail to fetch account - %v", err)
		}
	}

	defer cur.Close(ctx)

	var account AccountDetails

	if cur.Next(ctx) {
		err = cur.Decode(&account)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch account - %v", err)
		}

		for _, m := range account.Members {
			m.AccountId = account.ID

			if m.UserId == userId {
				account.Role = m.Role
			}

			if m.User == nil {
				m.User = &AccountMemberUser{
					FirstName: "Invited",
					LastName:  "user",
					Email:     m.InvitationEmail,
				}
			}
		}

		return &account, nil
	}

	return nil, ErrNotFound
}

func (store *accountStore) AddMember(ctx context.Context, member *AccountMember) (*AccountMember, error) {
	filter := bson.D{{"_id", member.AccountId}}
	if member.User != nil {
		filter = append(filter, bson.E{Key: "members", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{"user_id", member.User.ID}}}}})
	} else {
		filter = append(filter, bson.E{Key: "members", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{"invitation_email", member.InvitationEmail}}}}})
	}

	res := store.storage.Collection("account").FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			// ok, do update
		} else {
			return nil, ErrDuplicate
		}
	} else {
		return nil, ErrDuplicate
	}

	updatedAt := time.Now()
	member.CreatedAt = updatedAt
	member.UpdatedAt = updatedAt

	updateFilter := bson.D{{"_id", member.AccountId}}
	update := bson.D{
		{
			"$push", bson.D{{
				"members", member,
			}},
		},
	}

	opts := options.Update()
	_, err := store.storage.Collection("account").UpdateOne(ctx, updateFilter, update, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if member.User == nil {
		member.User = &AccountMemberUser{
			FirstName: "Invited",
			LastName:  "user",
			Email:     member.InvitationEmail,
		}
	}

	return member, nil
}

func (store *accountStore) UpdateMember(ctx context.Context, update *AccountMemberUpdate) error {
	filter := bson.D{
		{
			"_id", update.member.AccountId,
		},
		{
			"members._id", update.member.ID,
		},
	}
	updateBson := bson.M{"$set": update.setUpdatedDate().Bson()}
	_, err := store.storage.Collection("account").UpdateOne(ctx, filter, updateBson, options.Update())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (store *accountStore) FindAccountMember(ctx context.Context, memberID string) (*UserAccount, *AccountMember, error) {
	// TODO simplify
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{
				Key: "members",
				Value: bson.D{{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "_id",
							Value: memberID,
						},
					},
				}},
			},
		},
	}}

	unwindStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			"path",
			"$members",
		}},
	}}

	lookupStage := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{
				"from",
				"user",
			},
			{
				"localField",
				"members.user_id",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"user",
			},
		},
	}}

	addFieldsStage := bson.D{{
		Key: "$addFields",
		Value: bson.D{{
			"members.user",
			bson.D{{
				"$arrayElemAt", bson.A{"$user", 0},
			}},
		}},
	}}

	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{
				"_id",
				bson.D{
					{
						"_id", "$_id",
					},
					{
						"created_at", "$created_at",
					},
					{
						"updated_at", "$updated_at",
					},
					{
						"name", "$name",
					},
					{
						"status", "$status",
					},
				},
			},
			{
				"members",
				bson.D{{
					"$push", "$members",
				}},
			},
		},
	}}

	replaceWithStage := bson.D{{
		"$replaceWith",
		bson.D{
			{
				"_id", "$_id._id",
			},
			{
				"created_at", "$_id.created_at",
			},
			{
				"updated_at", "$_id.updated_at",
			},
			{
				"name", "$_id.name",
			},
			{
				"status", "$_id.status",
			},
			{
				"members", "$members",
			},
		},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		unwindStage,
		lookupStage,
		addFieldsStage,
		groupStage,
		replaceWithStage,
	}

	opts := options.Aggregate()
	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, ErrNotFound
		} else {
			return nil, nil, fmt.Errorf("fail to fetch account - %v", err)
		}
	}

	defer cur.Close(ctx)

	var account AccountDetails

	if cur.Next(ctx) {
		err = cur.Decode(&account)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to fetch account - %v", err)
		}

		for _, m := range account.Members {
			m.AccountId = account.ID
			if m.User == nil {
				m.User = &AccountMemberUser{
					FirstName: "Invited",
					LastName:  "user",
					Email:     m.InvitationEmail,
				}
			}

			if m.ID == memberID {
				a := &UserAccount{
					ID:     account.ID,
					Name:   account.Name,
					Status: account.Status,
					Role:   account.Role,
				}
				return a, m, nil
			}
		}

		return nil, nil, ErrNotFound
	}

	return nil, nil, ErrNotFound
}

func (store *accountStore) FindUserAccountMember(ctx context.Context, accountId, userId string) (*AccountMember, error) {
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{
			{"_id", accountId},
			{
				Key: "members",
				Value: bson.D{{
					Key: "$elemMatch",
					Value: bson.D{
						{
							Key:   "user_id",
							Value: userId,
						},
					},
				}},
			},
		},
	}}

	projectStage := bson.D{{
		Key: "$project",
		Value: bson.D{
			{"member", bson.D{{
				"$arrayElemAt", bson.A{
					"$members",
					bson.D{{
						"$indexOfArray",
						bson.A{"$members", bson.D{
							{
								"user_id", userId,
							},
						},
						},
					}},
				},
			}}},
			{"_id", 0},
		},
	}}

	pipeline := mongo.Pipeline{
		matchStage,
		projectStage,
	}

	cur, err := store.storage.Collection("account").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch accounts - %v", err)
	}

	defer cur.Close(ctx)

	res := struct {
		Member AccountMember `bson:"member"`
	}{}
	if cur.Next(ctx) {
		err = cur.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch account member - %v", err)
		}
		return &res.Member, nil
	}

	return nil, ErrNotFound
}

func (store *accountStore) DeleteAccountMember(ctx context.Context, member *AccountMember) error {
	filter := bson.D{{"_id", member.AccountId}}
	update := bson.D{
		{
			"$pull", bson.D{
				{
					"members", bson.D{
						{
							"_id", member.ID,
						},
					},
				},
			},
		},
	}

	_, err := store.storage.Collection("account").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}
