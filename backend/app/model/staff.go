package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type StaffRole struct {
	ID    string `bson:"_id" json:"-"`
	Name  string `bson:"name" json:"name"`
	Title string `bson:"title" json:"title"`
}

type StaffRoleNew string

const (
	StaffRoleNewAdmin   StaffRoleNew = "admin"
	StaffRoleNewManager StaffRoleNew = "manager"
)

type StaffShort struct {
	ID        string     `bson:"_id" json:"id"`
	FirstName string     `bson:"first_name" json:"fist_name"`
	LastName  string     `bson:"last_name" json:"last_name"`
	Photo     string     `bson:"photo,omitempty" json:"-"`
	PhotoURL  *string    `bson:"-" json:"photo,omitempty"`
	Email     string     `bson:"email" json:"email"`
	Role      *StaffRole `bson:"role" json:"role"`
	IsActive  bool       `bson:"is_active" json:"is_active"`
}

func (s *StaffShort) UpdatePhotoURL(host string) {
	if len(s.Photo) > 0 {
		photoUrl := fmt.Sprintf("%s/%s", host, s.Photo)
		s.PhotoURL = &photoUrl
	}
}

type Staff struct {
	ID                      string     `bson:"_id" json:"id"`
	FirstName               string     `bson:"first_name" json:"fist_name"`
	LastName                string     `bson:"last_name" json:"last_name"`
	Photo                   string     `bson:"photo,omitempty" json:"-"`
	PhotoURL                *string    `bson:"-" json:"photo,omitempty"`
	Email                   string     `bson:"email" json:"email"`
	IsEmailConfirmed        bool       `bson:"is_email_confirmed" json:"is_email_confirmed"`
	EmailConfirmedAt        *time.Time `bson:"email_confirmed_at" json:"email_confirmed_at,omitempty"`
	EmailConfirmationSentAt *time.Time `bson:"email_confirmation_sent_at" json:"email_confirmation_sent_at,omitempty"`
	Login                   string     `bson:"login" json:"-"`
	Password                string     `bson:"-" json:"-"`
	PasswordHash            string     `bson:"password_hash" json:"-"`
	Role                    *StaffRole `bson:"role" json:"role"`
	IsActive                bool       `bson:"is_active" json:"is_active"`
	DeactivatedAt           *time.Time `bson:"deactivated_at" json:"deactivated_at,omitempty"`
}

func NewStaff(email, firstName, lastName string, role *StaffRole) (*Staff, error) {
	staffPass, err := GenerateStaffPassword()
	if err != nil {
		return nil, fmt.Errorf("fail to create staff - %v", err)
	}

	staffPassHash, err := bcrypt.GenerateFromPassword([]byte(staffPass), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("fail to create staff - %v", err)
	}

	return &Staff{
		ID:           uuid.NewV4().String(),
		Login:        email,
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		Password:     staffPass,
		PasswordHash: string(staffPassHash),
		IsActive:     true,
	}, nil
}

func (s *Staff) IsEmailConfirmationRequired() bool {
	return !s.IsEmailConfirmed
}

func (s *Staff) UpdatePhotoURL(host string) {
	if len(s.Photo) > 0 {
		photoUrl := fmt.Sprintf("%s/%s", host, s.Photo)
		s.PhotoURL = &photoUrl
	}
}

func (s *Staff) NewUpdate() *StaffUpdate {
	return &StaffUpdate{
		ID:    s.ID,
		staff: s,
	}
}

func (s *Staff) ResetPassword() (*StaffUpdate, error) {
	newPass, err := GenerateStaffPassword()
	if err != nil {
		return nil, fmt.Errorf("fail to reset staff password - %v", err)
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("fail to reset staff password - %v", err)
	}

	s.Password = newPass
	s.PasswordHash = string(newPassHash)

	return s.NewUpdate().SetPasswordHash(s.PasswordHash), nil
}

func (s *Staff) MarshalJson(host string) ([]byte, error) {
	s.UpdatePhotoURL(host)
	return json.Marshal(s)
}

const KCtxKeyStaffMe = "StaffMe"
const KCtxKeyStaff = "Staff"

func (s *Staff) NewContext(ctx context.Context, key string) context.Context {
	ctx = context.WithValue(ctx, key, s)
	return ctx
}

func StaffFromContext(ctx context.Context, key string) (*Staff, error) {
	staff, _ := ctx.Value(key).(*Staff)

	if staff == nil {
		return nil, ErrNotFound
	}

	return staff, nil
}

type StaffUpdate struct {
	ID                      string     `json:"-"`
	FirstName               *string    `json:"first_name"`
	LastName                *string    `json:"last_name"`
	Photo                   *string    `json:"-"`
	Email                   *string    `json:"email"`
	IsEmailConfirmed        *bool      `json:"-"`
	EmailConfirmedAt        *time.Time `json:"-"`
	EmailConfirmationSentAt *time.Time `json:"-"`
	Login                   *string    `json:"-"`
	PasswordHash            *string    `json:"-"`
	Role                    *StaffRole `json:"-"`
	RoleName                *string    `json:"role"`
	IsActive                *bool      `json:"-"`

	staff *Staff
}

func (update *StaffUpdate) SetLogin(login string) *StaffUpdate {
	update.Login = &login
	return update
}

func (update *StaffUpdate) SetFirstName(firstName string) *StaffUpdate {
	update.FirstName = &firstName
	return update
}

func (update *StaffUpdate) SetLastName(lastName string) *StaffUpdate {
	update.LastName = &lastName
	return update
}

func (update *StaffUpdate) SetEmailConfirmationSent() *StaffUpdate {
	now := time.Now()
	update.EmailConfirmationSentAt = &now
	return update
}

func (update *StaffUpdate) SetEmailConfirmed() *StaffUpdate {
	v := true
	update.IsEmailConfirmed = &v

	now := time.Now()
	update.EmailConfirmedAt = &now
	return update
}

func (update *StaffUpdate) SetPhoto(fileName string) *StaffUpdate {
	update.Photo = &fileName
	return update
}

func (update *StaffUpdate) ClearPhoto() *StaffUpdate {
	empty := ""
	update.Photo = &empty
	return update
}

func (update *StaffUpdate) SetActive(active bool) *StaffUpdate {
	update.IsActive = &active
	return update
}

func (update *StaffUpdate) SetPasswordHash(passHash string) *StaffUpdate {
	update.PasswordHash = &passHash
	return update
}

func (update *StaffUpdate) IsEmailUpdated() bool {
	return update.Email != nil && update.staff.Email != *update.Email
}

func (update *StaffUpdate) Bson() bson.D {
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

	if update.Photo != nil {
		if len(*update.Photo) > 0 {
			res = append(res, bson.E{Key: "photo", Value: update.Photo})
		} else {
			res = append(res, bson.E{Key: "photo", Value: nil})
		}
	}

	if update.PasswordHash != nil {
		res = append(res, bson.E{Key: "password_hash", Value: update.PasswordHash})
	}

	if update.Role != nil {
		res = append(res, bson.E{Key: "role", Value: update.Role})
	}

	if update.IsActive != nil {
		res = append(res, bson.E{Key: "is_active", Value: *update.IsActive})
		if *update.IsActive {
			res = append(res, bson.E{Key: "deactivated_at", Value: nil})
		} else {
			res = append(res, bson.E{Key: "deactivated_at", Value: time.Now()})
		}
	}

	return res
}

type StaffStore interface {
	EnsureIndexes(ctx context.Context) error

	FetchStaffRoles() ([]*StaffRole, error)
	FindStaffRole(roleName string) (*StaffRole, error)

	FetchStaffByLogin(login string) (*Staff, error)
	FindStaffByEmailConfirmationToken(token string) (*Staff, error)
	FindStaffByID(ctx context.Context, ID string) (*Staff, error)
	AddStaff(staff *Staff) error
	UpdateStaff(ctx context.Context, update *StaffUpdate) (*Staff, error)

	FetchStaffList(offset, maxCount int, excludeId *string, searchPhase *string) ([]*StaffShort, int, error)
}

func NewStaffStore(storage Storage) (StaffStore, error) {
	return &staffStore{
		storage: storage,
	}, nil
}

type staffStore struct {
	storage Storage
}

func (store *staffStore) FetchStaffRoles() ([]*StaffRole, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	staffRoles := make([]*StaffRole, 0)
	cur, err := store.storage.Collection("staff_role").Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("fail to fetch staff roles - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &staffRoles)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch staff roles - %v", err)
	}

	return staffRoles, nil
}

func (store *staffStore) FindStaffRole(roleName string) (*StaffRole, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	filter := bson.D{{
		"name", roleName,
	}}

	res := store.storage.Collection("staff_role").FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("fail to find staff roles - %v", res.Err())
		}
	}

	var staffRole StaffRole
	err := res.Decode(&staffRole)
	if err != nil {
		return nil, fmt.Errorf("fail to find staff roles - %v", res.Err())
	}

	return &staffRole, nil
}

func (store *staffStore) FetchStaffByLogin(login string) (*Staff, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
	filter := bson.M{"login": login}

	var staff Staff
	err := store.storage.Collection("staff").FindOne(ctx, filter).Decode(&staff)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fail to fetch staff - %s", err)
	}

	return &staff, nil
}

func (store *staffStore) FindStaffByEmailConfirmationToken(token string) (*Staff, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
	filter := bson.M{"email_confirmation_token": token}

	var staff Staff
	err := store.storage.Collection("staff").FindOne(ctx, filter).Decode(&staff)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fail to fetch staff - %s", err)
	}

	return &staff, nil
}

func (store *staffStore) FindStaffByID(ctx context.Context, ID string) (*Staff, error) {
	filter := bson.M{"_id": ID}

	var staff Staff
	err := store.storage.Collection("staff").FindOne(ctx, filter).Decode(&staff)
	if err != nil {
		return nil, ErrNotFound
	}

	return &staff, nil
}

func (store *staffStore) AddStaff(staff *Staff) error {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	_, err := store.storage.Collection("staff").InsertOne(ctx, staff)
	if err != nil {
		if IsErrDuplication(err) {
			return ErrDuplicate
		}

		return err
	}

	return nil
}

func (store *staffStore) UpdateStaff(ctx context.Context, update *StaffUpdate) (*Staff, error) {
	filter := bson.M{"_id": update.ID}
	updateBson := bson.M{"$set": update.Bson()}

	var updatedStaff Staff
	err := store.storage.Collection("staff").FindOneAndUpdate(ctx, filter, updateBson, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedStaff)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		if IsErrDuplication(err) {
			return nil, ErrDuplicate
		}

		return nil, err
	}

	return &updatedStaff, nil
}

func (store *staffStore) FetchStaffList(offset, maxCount int, excludeId *string, searchPhase *string) ([]*StaffShort, int, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	filter := bson.D{}
	if excludeId != nil {
		filter = append(filter, bson.E{
			Key: "_id", Value: bson.D{{
				"$ne", *excludeId,
			}},
		})
	}

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

	count, err := store.storage.Collection("staff").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch staff count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(maxCount))

	staffList := make([]*StaffShort, 0)
	cur, err := store.storage.Collection("staff").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch staff - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &staffList)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch staff - %v", err)
	}

	return staffList, int(count), nil
}

func (store *staffStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "staff")

	log.Info("Fetching Staff storage indexes")
	indexes, err := store.storage.FetchCollectionIndexes(ctx, "staff")
	if err != nil {
		log.Errorf("Fail to fetch staff repo indexes - %v", err)
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
		_, err = store.storage.Collection("staff").Indexes().CreateOne(ctx, emailIdx)
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
		_, err = store.storage.Collection("staff").Indexes().CreateOne(ctx, emailIdx)
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
		_, err = store.storage.Collection("staff").Indexes().CreateOne(ctx, emailIdx)
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
		_, err = store.storage.Collection("staff").Indexes().CreateOne(ctx, loginIdx)
		if err != nil {
			log.Errorf("Fail to create index login_idx - %v", err)
			return err
		}
	}

	log.Info("Staff storage indexes are up-to-date")

	return nil
}
