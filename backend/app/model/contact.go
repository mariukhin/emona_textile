package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"amifactory.team/sequel/coton-app-backend/app/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

type ContactVariableType string

var (
	VarTypeString = ContactVariableType("string")
	VarTypeNumber = ContactVariableType("number")
	VarTypePhone  = ContactVariableType("phone")
)

func Contains(list []ContactVariableType, item ContactVariableType) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

type ContactVariable struct {
	ID                 string              `bson:"id" json:"id"`
	Type               ContactVariableType `bson:"type" json:"type"`
	Title              string              `bson:"title" json:"title"`
	Name               string              `bson:"name" json:"name"`
	DefaultValueString *string             `bson:"default_value_string,omitempty" json:"default_value_string,omitempty"`
	DefaultValueNumber *int                `bson:"default_value_number,omitempty" json:"default_value_number,omitempty"`
	Min                int                 `bson:"min_length" json:"min_length"`
	Max                int                 `bson:"max_length" json:"max_length"`
	Required           bool                `bson:"required" json:"-"`
	CanDelete          bool                `bson:"can_delete" json:"can_delete"`
}

func NewContactVariable(title string, t ContactVariableType, defaultValueString *string, defaultValueNumber *int) ContactVariable {
	return ContactVariable{
		ID:                 uuid.NewV4().String(),
		Type:               t,
		Title:              title,
		Name:               VariableNameFromTitle(title),
		DefaultValueString: defaultValueString,
		DefaultValueNumber: defaultValueNumber,
		CanDelete:          true,
	}
}

func VariableNameFromTitle(title string) string {
	return strings.ReplaceAll(strings.ToLower(title), " ", "_")
}

func (v ContactVariable) Parse(value string) (interface{}, error) {
	switch v.Type {
	case VarTypeString:
		return value, nil

	case VarTypeNumber:
		i, err := strconv.Atoi(value)
		if err != nil {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}

			return f, nil
		}

		return i, nil

	case VarTypePhone:
		return ParsePhone(value)

	default:
		return nil, fmt.Errorf("unexpected var type '%s'", v.Type)
	}
}

func (v ContactVariable) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, kCtxKeyContactGroupVar, &v)
	return ctx
}

func ContactVariableFromContext(ctx context.Context) (*ContactVariable, error) {
	variable, _ := ctx.Value(kCtxKeyContactGroupVar).(*ContactVariable)
	if variable == nil {
		return nil, ErrNotFound
	}

	return variable, nil
}

type ContactGroupShort struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type ContactGroup struct {
	ID string `bson:"_id" json:"id"`

	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"-"`

	AccountId string `bson:"account" json:"-"`

	Name          string `bson:"name" json:"name"`
	Description   string `bson:"description,omitempty" json:"description,omitempty"`
	ContactsCount int    `bson:"contacts" json:"contacts"`

	Vars []ContactVariable `bson:"variables,omitempty" json:"variables,omitempty"`
}

func NewContactGroup(accountId, name, description string) *ContactGroup {
	defaultContactVars := []ContactVariable{
		{
			ID:        uuid.NewV4().String(),
			Type:      VarTypePhone,
			Title:     "Phone",
			Name:      "phone",
			Required:  true,
			CanDelete: false,
		},
		{
			ID:        uuid.NewV4().String(),
			Type:      VarTypeString,
			Title:     "Name",
			Name:      "name",
			CanDelete: true,
		},
	}

	return &ContactGroup{
		ID: uuid.NewV4().String(),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		AccountId: accountId,

		Name:        name,
		Description: description,

		Vars: defaultContactVars,
	}
}

func (g ContactGroup) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, kCtxKeyContactGroup, &g)
	return ctx
}

func ContactGroupFromContext(ctx context.Context) (*ContactGroup, error) {
	group, _ := ctx.Value(kCtxKeyContactGroup).(*ContactGroup)
	if group == nil {
		return nil, ErrNotFound
	}

	return group, nil
}

func (g *ContactGroup) FindVariableByName(name string) *ContactVariable {
	for _, v := range g.Vars {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

func (g *ContactGroup) FindVariablesByName(names []string) []ContactVariable {
	vars := make([]ContactVariable, 0)
	for _, name := range names {
		for _, v := range g.Vars {
			if v.Name == name {
				vars = append(vars, v)
				break
			}
		}
	}

	return vars
}

func (g *ContactGroup) FindVariableByID(varID string) *ContactVariable {
	for _, v := range g.Vars {
		if v.ID == varID {
			return &v
		}
	}

	return nil
}

// Attention!!! Each time you update bson serializable fields
// don't forget to update updateBson() method. Thx!
type Contact struct {
	ID string `bson:"_id" json:"id"`

	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"-"`

	AccountId string `bson:"account" json:"-"`

	GroupID   string `bson:"group" json:"-"`
	Phone     `bson:",inline" json:",inline"`
	VarValues map[string]interface{} `bson:"vars" json:"-"`
}

func (c *Contact) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	// marshal ID
	jsonValue, err := json.Marshal(c.ID)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"id\":%s,", string(jsonValue)))

	// marshal phone number
	jsonValue, err = json.Marshal(c.Phone)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"phone\":%s", string(jsonValue)))

	// marshal vars
	for key, value := range c.VarValues {
		buffer.WriteString(",")

		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		buffer.WriteString(fmt.Sprintf("\"%s\":%s", key, string(jsonValue)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func NewContact(accountId string, groupID string, phone Phone, values map[string]interface{}) *Contact {
	return &Contact{
		ID: uuid.NewV4().String(),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		AccountId: accountId,

		GroupID:   groupID,
		Phone:     phone,
		VarValues: values,
	}
}

func (c Contact) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, kCtxKeyContact, &c)
	return ctx
}

func ContactFromContext(ctx context.Context) (*Contact, error) {
	contact, _ := ctx.Value(kCtxKeyContact).(*Contact)
	if contact == nil {
		return nil, ErrNotFound
	}

	return contact, nil
}

func (c *Contact) UpdatePhone(phone *Phone) {
	if phone != nil {
		c.Number = phone.Number
		c.Mcc = phone.Mcc
		c.Mnc = phone.Mnc
	}
}

func (c *Contact) updateBson() bson.D {
	return bson.D{
		{"updated_at", c.UpdatedAt},
		{"phone", c.Phone.Number},
		{"mcc", c.Phone.Mcc},
		{"mnc", c.Phone.Mnc},
		{"vars", c.VarValues},
	}
}

type ContactsDuplicateOption string

const (
	Update ContactsDuplicateOption = "update"
	Ignore ContactsDuplicateOption = "ignore"
)

var ContactsDuplicateOptions = []ContactsDuplicateOption{
	Update, Ignore,
}

type ContactsStore interface {
	EnsureIndexes(ctx context.Context) error

	AddGroup(ctx context.Context, group *ContactGroup) error
	UpdateGroup(ctx context.Context, group *ContactGroup) error
	FindGroup(ctx context.Context, accountId string, groupId string) (*ContactGroup, error)
	FindGroups(ctx context.Context, accountId string, groupId []string) ([]*ContactGroupShort, error)
	ListGroups(ctx context.Context, accountId string, offset, maxCount int, searchPhase *string) ([]*ContactGroup, int, error)
	ListGroupsAll(ctx context.Context, accountId string) ([]*ContactGroupShort, error)
	DeleteGroup(ctx context.Context, group *ContactGroup) error

	AddGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error
	UpdateGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error
	DeleteGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error

	AddContact(ctx context.Context, contact *Contact) error
	AddContacts(ctx context.Context, accountId, groupId string, contacts []*Contact, duplicateOption ContactsDuplicateOption) (int, int, int, error)
	UpdateContact(ctx context.Context, contact *Contact) error
	FindContact(ctx context.Context, accountId string, contactId string) (*Contact, error)
	FindContactDuplicates(ctx context.Context, accountId, groupID string, contacts []*Contact) (int64, error)
	ListContacts(ctx context.Context, accountId string, group *ContactGroup, offset, maxCount int, searchPhase *string) ([]*Contact, int, error)
	ListContactsAll(ctx context.Context, accountId string, group ContactGroup) ([]*Contact, error)
	DeleteContact(ctx context.Context, contact *Contact) error
}

func NewContactsStore(storage Storage) (ContactsStore, error) {
	return &contactsStore{
		storage: storage,
	}, nil
}

type contactsStore struct {
	storage Storage
}

func (store *contactsStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "contact")

	log.Info("Fetching Contacts storage indexes")

	indexes, err := store.storage.FetchCollectionIndexes(ctx, "contacts")
	if err != nil {
		log.Errorf("Fail to fetch contacts repo indexes - %v", err)
		return err
	}

	accountFirstNameIdxExists := false
	accountLastNameIdxExists := false
	accountPhoneIdxExists := false

	for _, index := range indexes {
		idxName, ok := index["name"]
		if ok {
			if idxName == "account_phone_idx" {
				accountPhoneIdxExists = true
			} else if idxName == "account_first_name_idx" {
				accountFirstNameIdxExists = true
			} else if idxName == "account_last_name_idx" {
				accountLastNameIdxExists = true
			}
		}
	}

	if accountFirstNameIdxExists {
		log.Info("Index account_first_name_idx exists")
	} else {
		log.Info("Index account_first_name_idx is not exist - creating")
		firstNameIdx := mongo.IndexModel{
			Keys: bson.M{
				"first_name": 1,
			}, Options: options.Index().SetName("account_first_name_idx"),
		}
		_, err = store.storage.Collection("contacts").Indexes().CreateOne(ctx, firstNameIdx)
		if err != nil {
			log.Errorf("Fail to create index account_first_name_idx - %v", err)
			return err
		}
	}

	if accountLastNameIdxExists {
		log.Info("Index account_last_name_idx exists")
	} else {
		log.Info("Index account_last_name_idx is not exist - creating")
		firstNameIdx := mongo.IndexModel{
			Keys: bson.M{
				"last_name": 1,
			}, Options: options.Index().SetName("account_last_name_idx"),
		}
		_, err = store.storage.Collection("contacts").Indexes().CreateOne(ctx, firstNameIdx)
		if err != nil {
			log.Errorf("Fail to create index account_last_name_idx - %v", err)
			return err
		}
	}

	if accountPhoneIdxExists {
		log.Info("Index account_phone_idx exists")
	} else {
		log.Info("Index account_phone_idx is not exist - creating")
		phoneIdx := mongo.IndexModel{
			Keys: bson.M{
				"account": 1,
				"phone":   1,
			}, Options: options.Index().SetName("account_phone_idx"),
		}
		_, err = store.storage.Collection("contacts").Indexes().CreateOne(ctx, phoneIdx)
		if err != nil {
			log.Errorf("Fail to create index account_phone_idx - %v", err)
			return err
		}
	}

	log.Info("Contacts storage indexes are up-to-date")

	return nil
}

func (store *contactsStore) AddGroup(ctx context.Context, group *ContactGroup) error {
	_, err := store.storage.Collection("groups").InsertOne(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func (store *contactsStore) UpdateGroup(ctx context.Context, group *ContactGroup) error {
	filter := bson.D{{"_id", group.ID}}
	update := bson.D{{
		Key:   "$set",
		Value: group,
	}}

	group.UpdatedAt = time.Now()

	_, err := store.storage.Collection("groups").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (store *contactsStore) FindGroup(ctx context.Context, accountId string, groupId string) (*ContactGroup, error) {
	filter := bson.D{
		{"_id", groupId},
		{"account", accountId},
	}

	groups := make([]*ContactGroup, 0)
	cur, err := store.storage.Collection("groups").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("fail to find group - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &groups)
	if err != nil {
		return nil, fmt.Errorf("fail to find group - %v", err)
	}

	if len(groups) == 0 {
		return nil, ErrNotFound
	}

	return groups[0], nil
}

func (store *contactsStore) FindGroups(ctx context.Context, accountId string, groupId []string) ([]*ContactGroupShort, error) {
	filter := bson.D{
		{
			Key:   "account",
			Value: accountId,
		},
		{
			Key: "_id",
			Value: bson.D{{
				Key:   "$in",
				Value: groupId,
			}},
		},
	}

	opts := options.Find().SetSort(bson.D{{"name", 1}})

	groups := make([]*ContactGroupShort, 0)
	cur, err := store.storage.Collection("groups").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("fail to find groups - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &groups)
	if err != nil {
		return nil, fmt.Errorf("fail to find groups - %v", err)
	}

	return groups, nil
}

func (store *contactsStore) ListGroups(ctx context.Context, accountId string, offset, maxCount int, searchPhase *string) ([]*ContactGroup, int, error) {
	filter := bson.D{
		{
			Key:   "account",
			Value: accountId,
		},
	}

	if searchPhase != nil {
		searchRegex := fmt.Sprintf("^%s.*", *searchPhase)
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{"name", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
			},
		})
	}

	count, err := store.storage.Collection("groups").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch groups count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(maxCount)).SetSort(bson.D{{"name", 1}})
	cur, err := store.storage.Collection("groups").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch groups - %v", err)
	}

	defer cur.Close(ctx)

	groups := make([]*ContactGroup, 0)
	err = cur.All(ctx, &groups)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch groups - %v", err)
	}

	return groups, int(count), nil
}

func (store *contactsStore) ListGroupsAll(ctx context.Context, accountId string) ([]*ContactGroupShort, error) {
	filter := bson.D{
		{
			Key:   "account",
			Value: accountId,
		},
	}

	groups := make([]*ContactGroupShort, 0)
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cur, err := store.storage.Collection("groups").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch groups - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &groups)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch groups - %v", err)
	}

	return groups, nil
}

func (store *contactsStore) DeleteGroup(ctx context.Context, group *ContactGroup) error {
	sess, err := store.storage.StartSession()
	if err != nil {
		return fmt.Errorf("fail to delete group - %v", err)
	}
	defer sess.EndSession(ctx)

	transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Group group
		groupFilter := bson.D{{"_id", group.ID}}
		_, err := store.storage.Collection("groups").DeleteOne(ctx, groupFilter)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}

			return nil, err
		}

		// 2. Remove all group contacts
		contactsFilter := bson.D{{"group", group.ID}}
		_, err = store.storage.Collection("contacts").DeleteMany(ctx, contactsFilter)
		// group might not have contacts, so we ignore ErrNoDocuments
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
	return err
}

func (store *contactsStore) AddGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error {
	filter := bson.D{{"_id", group.ID}}
	update := bson.D{{
		Key:   "$push",
		Value: bson.D{{"variables", variable}},
	}}

	group.UpdatedAt = time.Now()

	_, err := store.storage.Collection("groups").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (store *contactsStore) UpdateGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error {
	filter := bson.D{
		{"_id", group.ID},
		{"variables.id", variable.ID},
	}
	update := bson.D{{
		Key:   "$set",
		Value: bson.D{{"variables.$", variable}},
	}}

	group.UpdatedAt = time.Now()

	_, err := store.storage.Collection("groups").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (store *contactsStore) DeleteGroupVar(ctx context.Context, group *ContactGroup, variable *ContactVariable) error {
	sess, err := store.storage.StartSession()
	if err != nil {
		return fmt.Errorf("fail to add contact - %v", err)
	}
	defer sess.EndSession(ctx)

	transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Remove Group variable
		variableFilter := bson.D{{"_id", group.ID}}
		variableUpdate := bson.D{
			{"$pull", bson.D{
				{"variables",
					bson.D{{"id", variable.ID}},
				},
			}},
		}

		group.UpdatedAt = time.Now()

		_, err := store.storage.Collection("groups").UpdateOne(ctx, variableFilter, variableUpdate)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}

			return nil, err
		}

		// 2. Remove variable field from Group contacts
		contactsFilter := bson.D{{"group", group.ID}}
		contactsUpdate := bson.D{{"$unset", bson.D{{fmt.Sprintf("vars.%s", variable.Name), ""}}}}
		_, err = store.storage.Collection("contacts").UpdateMany(ctx, contactsFilter, contactsUpdate)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
	return err
}

func (store *contactsStore) AddContact(ctx context.Context, contact *Contact) error {
	sess, err := store.storage.StartSession()
	if err != nil {
		return fmt.Errorf("fail to add contact - %v", err)
	}
	defer sess.EndSession(ctx)

	transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Add contact to Group
		_, err := store.storage.Collection("contacts").InsertOne(ctx, contact)
		if err != nil {
			return nil, err
		}

		// 2. Increment Group contacts counter
		filter := bson.D{{"_id", contact.GroupID}}
		update := bson.D{{"$inc", bson.D{{"contacts", 1}}}}
		_, err = store.storage.Collection("groups").UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
	return err
}

func (store *contactsStore) AddContacts(ctx context.Context, accountId, groupId string, contacts []*Contact, duplicateOption ContactsDuplicateOption) (int, int, int, error) {
	type result struct {
		added   int
		updated int
		ignored int
	}

	switch duplicateOption {
	case Update:
		var operations []mongo.WriteModel
		for _, c := range contacts {
			contactUpdate := mongo.NewUpdateOneModel()
			contactUpdate.SetFilter(bson.M{"phone": c.Phone.Number, "account": accountId, "group": groupId})
			contactUpdate.SetUpdate(bson.M{
				"$set":         c.updateBson(),
				"$setOnInsert": bson.D{{"_id", c.ID}},
			})
			contactUpdate.SetUpsert(true)
			operations = append(operations, contactUpdate)
		}

		sess, err := store.storage.StartSession()
		if err != nil {
			return 0, 0, 0, fmt.Errorf("fail to add contacts - %v", err)
		}
		defer sess.EndSession(ctx)

		transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
			// 1. Add contacts
			r, err := store.storage.Collection("contacts").BulkWrite(ctx, operations)
			if err != nil {
				return result{}, err
			}

			res := result{
				added:   int(r.UpsertedCount),
				updated: int(r.ModifiedCount),
			}

			// 2. Increment Group contacts counter
			filter := bson.D{{"_id", groupId}}
			update := bson.D{{"$inc", bson.D{{"contacts", res.added}}}}
			_, err = store.storage.Collection("groups").UpdateOne(ctx, filter, update)
			if err != nil {
				return result{}, err
			}

			return res, nil
		}

		r, err := sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
		res := r.(result)
		return res.added, res.updated, res.ignored, err

	case Ignore:
		// find duplicate IDs
		duplicateIds, err := store.findContactDuplicatePhoneNumbers(ctx, accountId, groupId, contacts)
		if err != nil {
			return 0, 0, 0, err
		}

		// filter out contacts with duplicate IDs
		var iContacts []interface{}
		for _, c := range contacts {
			if _, found := duplicateIds[c.Phone.Number]; !found {
				iContacts = append(iContacts, c)
			}
		}

		if len(iContacts) == 0 {
			return 0, 0, len(duplicateIds), nil
		}

		sess, err := store.storage.StartSession()
		if err != nil {
			return 0, 0, 0, fmt.Errorf("fail to add contacts - %v", err)
		}
		defer sess.EndSession(ctx)

		transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
			// 1. Add contacts - insert new contacts only
			_, err = store.storage.Collection("contacts").InsertMany(ctx, iContacts)
			if err != nil {
				return result{}, err
			}

			res := result{
				added:   len(iContacts),
				ignored: len(duplicateIds),
			}

			// 2. Increment Group contacts counter
			filter := bson.D{{"_id", groupId}}
			update := bson.D{{"$inc", bson.D{{"contacts", res.added}}}}
			_, err = store.storage.Collection("groups").UpdateOne(ctx, filter, update)
			if err != nil {
				return result{}, err
			}

			return res, nil
		}

		r, err := sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
		res := r.(result)
		return res.added, res.updated, res.ignored, err

	default:
		return 0, 0, 0, fmt.Errorf("unknown duplicate option '%s'", duplicateOption)
	}
}

func (store *contactsStore) UpdateContact(ctx context.Context, contact *Contact) error {
	filter := bson.D{{"_id", contact.ID}}
	update := bson.D{{
		Key:   "$set",
		Value: contact,
	}}

	contact.UpdatedAt = time.Now()

	_, err := store.storage.Collection("contacts").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (store *contactsStore) FindContact(ctx context.Context, accountId string, contactId string) (*Contact, error) {
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{Key: "_id", Value: contactId},
				{Key: "account", Value: accountId},
			},
		},
	}

	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "groups"},
				{Key: "localField", Value: "group_ids"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "groups"},
			},
		},
	}

	pipeline := mongo.Pipeline{
		matchStage,
		lookupStage,
	}

	cur, err := store.storage.Collection("contacts").Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return nil, fmt.Errorf("fail to find contact - cursor is empty")
	}

	var contact Contact
	err = cur.Decode(&contact)
	if err != nil {
		return nil, fmt.Errorf("fail to find contact - %v", err)
	}

	return &contact, nil
}

func (store *contactsStore) FindContactDuplicates(ctx context.Context, accountId, groupId string, contacts []*Contact) (int64, error) {
	splitSize := 100000
	numberOfSplits := len(contacts) / splitSize
	remainder := len(contacts) - numberOfSplits*splitSize
	if remainder > 0 {
		numberOfSplits++
	}

	var duplicatesCount int64

	for splitIdx := 0; splitIdx < numberOfSplits; splitIdx++ {
		from := splitSize * splitIdx
		to := utils.Min(from+splitSize, len(contacts))
		chunkSize := to - from
		chunk := make([]string, chunkSize)

		chunkIdx := 0
		for idx := from; idx < to; idx++ {
			chunk[chunkIdx] = contacts[idx].Phone.Number
			chunkIdx++
		}

		filter := bson.D{
			{
				"phone", bson.D{{
					"$in", chunk,
				}},
			},
			{
				"account", accountId,
			},
			{
				"group", groupId,
			},
		}

		count, err := store.storage.Collection("contacts").CountDocuments(ctx, filter)
		if err != nil {
			return 0, err
		}

		duplicatesCount += count
	}

	return duplicatesCount, nil
}

func (store *contactsStore) ListContacts(ctx context.Context, accountId string, group *ContactGroup, offset, maxCount int, searchPhase *string) ([]*Contact, int, error) {
	filter := bson.D{
		{Key: "account", Value: accountId},
	}

	if group == nil {
		return nil, 0, fmt.Errorf("fail to fetch contacts - group is not provided")
	}

	filter = append(filter, bson.E{Key: "group", Value: group.ID})

	if searchPhase != nil {
		searchRegex := fmt.Sprintf(`^\+?%s.*`, *searchPhase)
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{"phone", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
			},
		})
	}

	count, err := store.storage.Collection("contacts").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch contacts count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(maxCount)).SetSort(bson.D{{"created_at", -1}})
	cur, err := store.storage.Collection("contacts").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch contacts - %v", err)
	}

	defer cur.Close(ctx)

	contacts := make([]*Contact, 0)
	err = cur.All(ctx, &contacts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch contacts - %v", err)
	}

	return contacts, int(count), nil
}

func (store *contactsStore) ListContactsAll(ctx context.Context, accountId string, group ContactGroup) ([]*Contact, error) {
	filter := bson.M{
		"account": accountId,
		"group":   group.ID,
	}

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cur, err := store.storage.Collection("contacts").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch contacts - %v", err)
	}

	defer cur.Close(ctx)

	contacts := make([]*Contact, 0)
	err = cur.All(ctx, &contacts)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch contacts - %v", err)
	}

	return contacts, nil
}

func (store *contactsStore) DeleteContact(ctx context.Context, contact *Contact) error {
	sess, err := store.storage.StartSession()
	if err != nil {
		return fmt.Errorf("fail to add contact - %v", err)
	}
	defer sess.EndSession(ctx)

	transaction := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Remove contact
		contactFilter := bson.D{{"_id", contact.ID}}
		_, err := store.storage.Collection("contacts").DeleteOne(ctx, contactFilter)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrNotFound
			}

			return nil, err
		}

		// 2. Decrement Group contacts counter
		groupFilter := bson.D{{"_id", contact.GroupID}}
		update := bson.D{{"$inc", bson.D{{"contacts", -1}}}}
		_, err = store.storage.Collection("groups").UpdateOne(ctx, groupFilter, update)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, transaction, store.storage.DefaultTransactionOpts())
	return err
}

func (store *contactsStore) findContactDuplicatePhoneNumbers(ctx context.Context, accountId, groupId string, contacts []*Contact) (map[string]string, error) {
	splitSize := 100000
	numberOfSplits := len(contacts) / splitSize
	remainder := len(contacts) - numberOfSplits*splitSize
	if remainder > 0 {
		numberOfSplits++
	}

	duplicateIds := make(map[string]string, 0)

	for splitIdx := 0; splitIdx < numberOfSplits; splitIdx++ {
		from := splitSize * splitIdx
		to := utils.Min(from+splitSize, len(contacts))
		chunkSize := to - from
		chunk := make([]string, chunkSize)

		chunkIdx := 0
		for idx := from; idx < to; idx++ {
			chunk[chunkIdx] = contacts[idx].Phone.Number
			chunkIdx++
		}

		filter := bson.D{
			{
				"phone", bson.D{{
					"$in", chunk,
				}},
			},
			{
				"account", accountId,
			},
			{
				"group", groupId,
			},
		}

		cur, err := store.storage.Collection("contacts").Find(ctx, filter)
		if err != nil {
			return map[string]string{}, err
		}

		res := struct {
			ID    string `bson:"_id"`
			Phone string `bson:"phone"`
		}{}

		for {
			if cur.Next(ctx) {
				err = cur.Decode(&res)
				if err != nil {
					cur.Close(ctx)
					return map[string]string{}, err
				}

				duplicateIds[res.Phone] = res.ID
			} else {
				break
			}
		}

		cur.Close(ctx)
	}

	return duplicateIds, nil
}
