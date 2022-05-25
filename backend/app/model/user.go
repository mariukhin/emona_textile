package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStatus string

const (
	UserStatusInvited UserStatus = "invited"
	UserStatusActive  UserStatus = "active"
)

type UserShort struct {
	ID        string `bson:"_id" json:"id"`
	FirstName string `bson:"first_name" json:"fist_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
}

type User struct {
	ID                      string             `bson:"_id" json:"id"`
	FirstName               string             `bson:"first_name" json:"fist_name"`
	LastName                string             `bson:"last_name" json:"last_name"`
	Photo                   string             `bson:"photo,omitempty" json:"-"`
	PhotoURL                *string            `bson:"-" json:"photo,omitempty"`
	Email                   string             `bson:"email" json:"email"`
	IsEmailConfirmed        bool               `bson:"is_email_confirmed" json:"is_email_confirmed"`
	EmailConfirmedAt        *time.Time         `bson:"email_confirmed_at" json:"email_confirmed_at,omitempty"`
	EmailConfirmationSentAt *time.Time         `bson:"email_confirmation_sent_at" json:"email_confirmation_sent_at,omitempty"`
	Login                   string             `bson:"login" json:"-"`
	Password                string             `bson:"-" json:"-"`
	PasswordHash            string             `bson:"password_hash" json:"-"`
	DeactivatedAt           *time.Time         `bson:"deactivated_at" json:"deactivated_at,omitempty"`
	Status                  UserStatus         `bson:"status" json:"-"`
	Role                    *AccountMemberRole `bson:"-" json:"role,omitempty"`
}

func NewUser(email, firstName, lastName, password string) (*User, error) {
	userPassHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("fail to create user - %v", err)
	}

	return &User{
		ID:           uuid.NewV4().String(),
		Login:        email,
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		Password:     password,
		PasswordHash: string(userPassHash),
		Status:       UserStatusActive,
	}, nil
}

func (u *User) IsEmailConfirmationRequired() bool {
	return !u.IsEmailConfirmed
}

func (u *User) UpdateAccountMemberRole(account *AccountDetails) {
	if account != nil {
		u.Role = &account.Role
	}
}

func (u *User) UpdatePhotoURL(host string) {
	if len(u.Photo) > 0 {
		photoUrl := fmt.Sprintf("%s/%s", host, u.Photo)
		u.PhotoURL = &photoUrl
	}
}

func (u *User) NewUpdate() *UserUpdate {
	return &UserUpdate{
		user: u,
	}
}

const KCtxKeyUserMe = "user-me"

func (u *User) NewContext(ctx context.Context, key string) context.Context {
	ctx = context.WithValue(ctx, key, u)
	return ctx
}

func UserFromContext(ctx context.Context, key string) (*User, error) {
	user, _ := ctx.Value(key).(*User)

	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}

type UserUpdate struct {
	FirstName               *string    `json:"first_name"`
	LastName                *string    `json:"last_name"`
	Photo                   *string    `json:"-"`
	Email                   *string    `json:"email"`
	IsEmailConfirmed        *bool      `json:"-"`
	EmailConfirmedAt        *time.Time `json:"-"`
	EmailConfirmationSentAt *time.Time `json:"-"`
	Login                   *string    `json:"-"`
	PasswordHash            *string    `json:"-"`

	user *User
}

func (update *UserUpdate) SetLogin(login string) *UserUpdate {
	update.Login = &login
	return update
}

func (update *UserUpdate) SetPassword(password string) (*UserUpdate, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("fail to create password hash - %v", err)
	}

	passHashStr := string(passHash)

	update.PasswordHash = &passHashStr
	return update, nil
}

func (update *UserUpdate) IsEmailUpdated() bool {
	return update.Email != nil && update.user.Email != *update.Email
}

func (update *UserUpdate) SetFirstName(firstName string) *UserUpdate {
	update.FirstName = &firstName
	return update
}

func (update *UserUpdate) SetLastName(lastName string) *UserUpdate {
	update.LastName = &lastName
	return update
}

func (update *UserUpdate) SetEmailConfirmationSent() *UserUpdate {
	now := time.Now()
	update.EmailConfirmationSentAt = &now
	return update
}

func (update *UserUpdate) SetEmailConfirmed() *UserUpdate {
	v := true
	update.IsEmailConfirmed = &v

	now := time.Now()
	update.EmailConfirmedAt = &now
	return update
}

func (update *UserUpdate) SetPhoto(fileName string) *UserUpdate {
	update.Photo = &fileName
	return update
}

func (update *UserUpdate) ClearPhoto() *UserUpdate {
	empty := ""
	update.Photo = &empty
	return update
}

func (update *UserUpdate) Bson() bson.D {
	res := bson.D{}

	if update.FirstName != nil {
		res = append(res, bson.E{Key: "first_name", Value: update.FirstName})
	}

	if update.LastName != nil {
		res = append(res, bson.E{Key: "last_name", Value: update.LastName})
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

	if update.Login != nil {
		res = append(res, bson.E{Key: "login", Value: update.Login})
	}

	if update.PasswordHash != nil {
		res = append(res, bson.E{Key: "password_hash", Value: update.PasswordHash})
	}

	if update.Photo != nil {
		if len(*update.Photo) > 0 {
			res = append(res, bson.E{Key: "photo", Value: update.Photo})
		} else {
			res = append(res, bson.E{Key: "photo", Value: nil})
		}
	}

	return res
}

type UserStore interface {
	EnsureIndexes(ctx context.Context) error

	FetchList(offset, maxCount int, searchPhase *string) ([]*UserShort, int, error)
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, update *UserUpdate) (*User, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	FindByID(ctx context.Context, userId string) (*User, error)
}

func NewUserStore(storage Storage) (UserStore, error) {
	return &userStore{
		storage: storage,
	}, nil
}

type userStore struct {
	storage Storage
}

func (store *userStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "user")

	log.Info("Fetching User storage indexes")
	indexes, err := store.storage.FetchCollectionIndexes(ctx, "user")
	if err != nil {
		log.Errorf("Fail to fetch user collection indexes - %v", err)
		return err
	}

	firstNameIdxExists := false
	lastNameIdxExists := false
	emailIdxExists := false
	loginIdxExists := false

	for _, index := range indexes {
		idxName, ok := index["name"]
		if ok {
			if idxName == "first_name_idx" {
				firstNameIdxExists = true
			}

			if idxName == "last_name_idx" {
				lastNameIdxExists = true
			}

			if idxName == "email_idx" {
				emailIdxExists = true
			}

			if idxName == "login_idx" {
				loginIdxExists = true
			}
		}
	}

	if firstNameIdxExists {
		log.Info("Index first_name_idx exists")
	} else {
		log.Info("Index first_name_idx is not exist - creating")
		emailIdx := mongo.IndexModel{
			Keys: bson.M{
				"first_name": 1,
			}, Options: options.Index().SetName("first_name_idx"),
		}
		_, err = store.storage.Collection("user").Indexes().CreateOne(ctx, emailIdx)
		if err != nil {
			log.Errorf("Fail to create index first_name_idx - %v", err)
			return err
		}
	}

	if lastNameIdxExists {
		log.Info("Index last_name_idx exists")
	} else {
		log.Info("Index last_name_idx is not exist - creating")
		emailIdx := mongo.IndexModel{
			Keys: bson.M{
				"last_name": 1,
			}, Options: options.Index().SetName("last_name_idx"),
		}
		_, err = store.storage.Collection("user").Indexes().CreateOne(ctx, emailIdx)
		if err != nil {
			log.Errorf("Fail to create index last_name_idx - %v", err)
			return err
		}
	}

	if emailIdxExists {
		log.Info("Index email_idx exists")
	} else {
		log.Info("Index email_idx is not exist - creating")
		emailIdx := mongo.IndexModel{
			Keys: bson.M{
				"email": 1,
			}, Options: options.Index().SetUnique(true).SetName("email_idx"),
		}
		_, err = store.storage.Collection("user").Indexes().CreateOne(ctx, emailIdx)
		if err != nil {
			log.Errorf("Fail to create index email_idx - %v", err)
			return err
		}
	}

	if loginIdxExists {
		log.Info("Index login_idx exists")
	} else {
		log.Info("Index login_idx is not exist - creating")
		loginIdx := mongo.IndexModel{
			Keys: bson.M{
				"login": 1,
			}, Options: options.Index().SetUnique(true).SetName("login_idx"),
		}
		_, err = store.storage.Collection("user").Indexes().CreateOne(ctx, loginIdx)
		if err != nil {
			log.Errorf("Fail to create index login_idx - %v", err)
			return err
		}
	}

	log.Info("User storage indexes are up-to-date")

	return nil
}

func (store *userStore) FetchList(offset, maxCount int, searchPhase *string) ([]*UserShort, int, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	filter := bson.D{}

	if searchPhase != nil {
		searchRegex := fmt.Sprintf("^%s.*", *searchPhase)
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{"first_name", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
				bson.D{{"last_name", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
				bson.D{{"email", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
			},
		})
	}

	count, err := store.storage.Collection("user").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch users count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(maxCount))

	users := make([]*UserShort, 0)
	cur, err := store.storage.Collection("user").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch users - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &users)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch users - %v", err)
	}

	return users, int(count), nil
}

func (store *userStore) Add(ctx context.Context, user *User) error {
	_, err := store.storage.Collection("user").InsertOne(ctx, user)
	if err != nil {
		if IsErrDuplication(err) {
			return ErrDuplicate
		}

		return err
	}

	return nil
}

func (store *userStore) Update(ctx context.Context, update *UserUpdate) (*User, error) {
	filter := bson.M{"_id": update.user.ID}
	updateBson := bson.M{"$set": update.Bson()}

	var updatedUser User
	err := store.storage.Collection("user").FindOneAndUpdate(ctx, filter, updateBson, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		if IsErrDuplication(err) {
			return nil, ErrDuplicate
		}

		return nil, err
	}

	return &updatedUser, nil
}

func (store *userStore) FindByLogin(ctx context.Context, login string) (*User, error) {
	filter := bson.M{"login": login}

	var user User
	err := store.storage.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fail to fetch user - %s", err)
	}

	return &user, nil
}

func (store *userStore) FindByID(ctx context.Context, userId string) (*User, error) {
	filter := bson.M{"_id": userId}

	var user User
	err := store.storage.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fail to fetch user - %s", err)
	}

	return &user, nil
}
