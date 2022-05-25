package model

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"sync"
	"unicode"
)

type RefBookShort struct {
	Country  string `bson:"country"`
	Operator string `bson:"operator"`
	Mcc      int    `bson:"mcc"`
	Mnc      int    `bson:"mnc"`
}

type RefBookRecord struct {
	ID               string `bson:"_id"`
	CountryPhoneCode string `bson:"country_code"`
	CountryCode2     string `bson:"cc2"`
	Country          string `bson:"country"`
	Operator         string `bson:"operator"`
	PhonePrefix      string `bson:"prefix"`
	Mcc              int    `bson:"mcc"`
	Mnc              int    `bson:"mnc"`
	MinLen           int    `bson:"min_len"`
	MaxLen           int    `bson:"max_len"`
}

type RefBookStore interface {
	EnsureIndexes(ctx context.Context) error

	Save(ctx context.Context, record *RefBookRecord) error
	List(ctx context.Context) ([]*RefBookRecord, error)
	FindByMccMnc(ctx context.Context, mccMnc MccMnc) (*RefBookShort, error)
}

func NewRefBookStore(storage Storage) (RefBookStore, error) {
	return &refBookStore{
		storage: storage,
	}, nil
}

type refBookStore struct {
	storage Storage
}

func (store *refBookStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "refbook")

	log.Info("Fetching refbook storage indexes")
	indexes, err := store.storage.FetchCollectionIndexes(ctx, "refbook")
	if err != nil {
		log.Errorf("Fail to fetch refbook storage indexes - %v", err)
		return err
	}

	phonePrefixIdxExists := false

	for _, index := range indexes {
		idxName, ok := index["name"]
		if ok {
			if idxName == "prefix_idx" {
				phonePrefixIdxExists = true
			}
		}
	}

	if phonePrefixIdxExists {
		log.Info("Index prefix_idx exists")
	} else {
		log.Info("Index prefix_idx is not exist - creating")
		prefixIdx := mongo.IndexModel{
			Keys: bson.M{
				"prefix": 1,
			}, Options: options.Index().SetName("prefix_idx"),
		}
		_, err = store.storage.Collection("refbook").Indexes().CreateOne(ctx, prefixIdx)
		if err != nil {
			log.Errorf("Fail to create index prefix_idx - %v", err)
			return err
		}
	}

	log.Info("Refbook storage indexes are up-to-date")

	return nil
}

func (store *refBookStore) Save(ctx context.Context, record *RefBookRecord) error {
	filter := bson.D{{"_id", record.ID}}
	update := bson.D{{"$set", record}}
	opts := options.Update().SetUpsert(true)
	_, err := store.storage.Collection("refbook").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		if IsErrDuplication(err) {
			return ErrDuplicate
		}

		return err
	}

	return nil
}

func (store *refBookStore) List(ctx context.Context) ([]*RefBookRecord, error) {
	records := make([]*RefBookRecord, 0)
	cur, err := store.storage.Collection("refbook").Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("fail to fetch ref book records - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &records)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch ref book records - %v", err)
	}

	return records, nil
}

func (store *refBookStore) FindByMccMnc(ctx context.Context, mccMnc MccMnc) (*RefBookShort, error) {
	filter := bson.D{
		{"mcc", mccMnc.Mcc},
		{"mnc", mccMnc.Mnc},
	}
	res := store.storage.Collection("refbook").FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fail to find ref book record - %v", res.Err())
	}

	var rec RefBookShort
	err := res.Decode(&rec)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

type RefBookService interface {
	FetchByMccMnc(ctx context.Context, mccMnc MccMnc) (RefBookShort, error)
}

type refBookService struct {
	store RefBookStore
}

type MccMnc struct {
	Mcc int `bson:"mcc" json:"mcc"`
	Mnc int `bson:"mnc" json:"mnc"`
}

type Phone struct {
	Number string `bson:"phone" json:"phone"`
	MccMnc
}

func ParsePhone(rawPhone string) (*Phone, error) {
	return SharedPhoneService.ParsePhone(rawPhone)
}

func NewPhone(number string, mcc, mnc int) *Phone {
	return &Phone{
		Number: number,
		MccMnc: MccMnc{
			Mcc: mcc,
			Mnc: mnc,
		},
	}
}

func (p Phone) MarshalText() ([]byte, error) {
	return []byte(p.Number), nil
}

func (p *Phone) UnmarshalJSON(b []byte) error {
	var phoneRaw string
	if err := json.Unmarshal(b, &phoneRaw); err != nil {
		return err
	}

	p.Number = phoneRaw
	p.Mcc = 255
	p.Mnc = 6

	return nil
}

var ErrPhoneInvalidSymbols = errors.New("phone number contains invalid symbols")
var ErrPhoneMccMncNotFound = errors.New("mcc and mnc codes for phone number not found")
var ErrPhoneWrongLen = errors.New("phone number has wrong len")

type PhoneService interface {
	PrefetchRefBook(ctx context.Context) error
	ParsePhone(phoneRaw string) (*Phone, error)
	RefBookRecordForMccMnc(ctx context.Context, mccMnc MccMnc) (RefBookShort, error)
}

var SharedPhoneService PhoneService

func InitPhoneService(ctx context.Context, store RefBookStore) error {
	var err error
	SharedPhoneService, err = NewPhoneService(store)
	if err != nil {
		return err
	}
	return SharedPhoneService.PrefetchRefBook(ctx)
}

func NewPhoneService(store RefBookStore) (PhoneService, error) {
	return &phoneService{
		store:          store,
		refBookRecords: make(map[string]*RefBookRecord),
		rawPhoneRegex:  regexp.MustCompile(`^\+?[\d\-\s()]+$`),
		mccMncRecords:  make(map[MccMnc]*RefBookShort),
	}, nil
}

type phoneService struct {
	store RefBookStore

	maxPrefixLen       int
	refBookRecords     map[string]*RefBookRecord
	refBookRecordsLock sync.RWMutex

	mccMncRecords     map[MccMnc]*RefBookShort
	mccMncRecordsLock sync.Mutex

	rawPhoneRegex *regexp.Regexp
}

func (srv *phoneService) PrefetchRefBook(ctx context.Context) error {
	srv.refBookRecordsLock.Lock()
	defer srv.refBookRecordsLock.Unlock()

	records, err := srv.store.List(ctx)
	if err != nil {
		return fmt.Errorf("fail to fetch refbook records - %v", err)
	}

	srv.maxPrefixLen = 0
	srv.refBookRecords = make(map[string]*RefBookRecord)

	for _, rec := range records {
		srv.refBookRecords[rec.PhonePrefix] = rec
		if srv.maxPrefixLen < len(rec.PhonePrefix) {
			srv.maxPrefixLen = len(rec.PhonePrefix)
		}
	}

	return nil
}

func (srv *phoneService) ParsePhone(phoneRaw string) (*Phone, error) {
	// 1. Check if phone raw is at least fit in allowed boundaries:
	// - optional + at the begin
	// - optionally contains spaces, '(', ')' or '-'
	if !srv.rawPhoneRegex.MatchString(phoneRaw) {
		return nil, ErrPhoneInvalidSymbols
	}

	// 2. Remove all symbols except digits
	phoneClearedRunes := make([]rune, 0, len(phoneRaw))
	for _, r := range []rune(phoneRaw) {
		if unicode.IsDigit(r) {
			phoneClearedRunes = append(phoneClearedRunes, r)
		}
	}

	// 3. Try find refbook record
	srv.refBookRecordsLock.RLock()
	rec := srv.findMccMncUnsafe(phoneClearedRunes)
	srv.refBookRecordsLock.RUnlock()

	if rec == nil {
		return nil, ErrPhoneMccMncNotFound
	}

	// 4. Check phone number length
	if 0 < rec.MinLen && len(phoneClearedRunes) < rec.MinLen {
		return nil, ErrPhoneWrongLen
	}

	if 0 < rec.MaxLen && rec.MaxLen < len(phoneClearedRunes) {
		return nil, ErrPhoneWrongLen
	}

	// 5. Everything is ok - return Phone
	return NewPhone(string(phoneClearedRunes), rec.Mcc, rec.Mnc), nil
}

func (srv *phoneService) RefBookRecordForMccMnc(ctx context.Context, mccMnc MccMnc) (RefBookShort, error) {
	srv.mccMncRecordsLock.Lock()
	defer srv.mccMncRecordsLock.Unlock()

	var err error

	rec, ok := srv.mccMncRecords[mccMnc]
	if !ok {
		rec, err = srv.store.FindByMccMnc(ctx, mccMnc)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				srv.mccMncRecords[mccMnc] = nil
				return RefBookShort{}, ErrNotFound
			} else {
				return RefBookShort{}, fmt.Errorf("fail to find mcc/mnc record - %v", err)
			}
		}

		srv.mccMncRecords[mccMnc] = rec
	}

	if rec != nil {
		return *rec, nil
	} else {
		return RefBookShort{}, nil
	}
}

func (srv *phoneService) findMccMncUnsafe(phone []rune) *RefBookRecord {
	phoneLen := srv.maxPrefixLen
	if len(phone) < phoneLen {
		phoneLen = len(phone)
	}

	for idx := phoneLen; idx > 0; idx-- {
		prefix := string(phone[:idx])
		rec := srv.refBookRecords[prefix]
		if rec != nil {
			return rec
		}
	}

	return nil
}
