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
	"strconv"
	"time"
)

type Price struct {
	ID string `bson:"_id" json:"id"`

	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"-"`

	AccountId *string `bson:"account,omitempty" json:"-"`

	MccMnc   `bson:",inline"`
	Country  string `bson:"country" json:"country"`
	Operator string `bson:"operator" json:"operator"`

	Enabled bool `bson:"enabled" json:"-"`

	Price    Cost     `bson:"price" json:"-"`
	Currency Currency `bson:"currency" json:"-"`
}

func NewPrice(accountID *string, mccMnc MccMnc, price Cost, cur Currency) Price {
	return Price{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
		AccountId: accountID,
		MccMnc:    mccMnc,
		Enabled:   true,
		Price:     price,
		Currency:  cur,
	}
}

type PricesFilter struct {
	SearchPhase *string
	Enabled     *bool
	AccountID   *string
	Currency    *Currency
}

func (pf *Price) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "price", pf)
	return ctx
}

func PriceItemFromContext(ctx context.Context) (*Price, error) {
	price, _ := ctx.Value("price").(*Price)
	if price == nil {
		return nil, ErrNotFound
	}

	return price, nil
}

func (pf PricesFilter) SetEnabled(enabled bool) PricesFilter {
	return PricesFilter{
		SearchPhase: pf.SearchPhase,
		Enabled:     &enabled,
		AccountID:   pf.AccountID,
		Currency:    pf.Currency,
	}
}

func (pf PricesFilter) SetCurrency(cur Currency) PricesFilter {
	return PricesFilter{
		SearchPhase: pf.SearchPhase,
		Enabled:     pf.Enabled,
		AccountID:   pf.AccountID,
		Currency:    &cur,
	}
}

func (pf PricesFilter) BsonFilter() bson.D {
	filter := bson.D{}
	if pf.Currency != nil {
		filter = append(filter, bson.E{
			Key: "currency", Value: *pf.Currency,
		})
	}

	if pf.SearchPhase != nil && len(*pf.SearchPhase) > 0 {
		// 1. If searchPhase is number (code) - filter by mcc/mnc
		code, err := strconv.Atoi(*pf.SearchPhase)
		if err != nil {
			// 2. If searchPhase is not number - filter by country/operator
			searchRegex := fmt.Sprintf("^%s.*", *pf.SearchPhase)
			filter = append(filter, bson.E{
				Key: "$or", Value: bson.A{
					bson.D{{"country", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
					bson.D{{"operator", bson.D{{"$regex", searchRegex}, {"$options", "i"}}}},
				},
			})
		} else {
			filter = append(filter, bson.E{
				Key: "$or", Value: bson.A{
					bson.M{"mcc": code},
					bson.M{"mnc": code},
				},
			})
		}
	}

	if pf.Enabled != nil {
		filter = append(filter, bson.E{
			Key: "enabled", Value: *pf.Enabled,
		})
	}

	return filter
}

type PricesStore interface {
	EnsureIndexes(ctx context.Context) error

	FetchPriceForAccount(ctx context.Context, accountId string, mccMnc MccMnc) (*Price, error)
	FetchPricesForAccount(ctx context.Context, offset, maxCount int, pf PricesFilter) ([]Price, int, error)

	AddPrice(ctx context.Context, price Price) error
	FetchPrice(ctx context.Context, priceID string) (*Price, error)
	UpdatePrice(ctx context.Context, price Price) error
	DeletePrice(ctx context.Context, price Price) error
}

func NewPricesStore(store Storage) (PricesStore, error) {
	return &pricesStore{
		store: store,
	}, nil
}

type pricesStore struct {
	store Storage
}

func (ps *pricesStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "prices")

	log.Info("Fetching storage indexes")

	indexes, err := ps.store.FetchCollectionIndexes(ctx, "prices")
	if err != nil {
		log.Errorf("Fail to fetch repo indexes - %v", err)
		return err
	}

	accountMccMncIdxExists := false

	for _, index := range indexes {
		idxName, ok := index["name"]
		if ok {
			if idxName == "account_mcc_mnc_idx" {
				accountMccMncIdxExists = true
			}
		}
	}

	if accountMccMncIdxExists {
		log.Info("Index account_mcc_mnc_idx exists")
	} else {
		log.Info("Index account_mcc_mnc_idx is not exist - creating")
		accountMccMncIdx := mongo.IndexModel{
			Keys: bson.M{
				"account":  1,
				"mcc":      1,
				"mnc":      1,
				"currency": 1,
			}, Options: options.Index().SetUnique(true).SetName("account_mcc_mnc_idx"),
		}
		_, err = ps.store.Collection("prices").Indexes().CreateOne(ctx, accountMccMncIdx)
		if err != nil {
			log.Errorf("Fail to create index account_first_name_idx - %v", err)
			return err
		}
	}

	log.Info("Storage indexes are up-to-date")

	return nil
}

func (ps *pricesStore) FetchPriceForAccount(ctx context.Context, accountId string, mccMnc MccMnc) (*Price, error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"account": accountId}, bson.M{"account": nil},
		},
		"mcc": mccMnc.Mcc,
		"mnc": mccMnc.Mnc,
	}

	opts := options.Find().SetSort(bson.M{"account": -1}).SetLimit(1)
	cur, err := ps.store.Collection("prices").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("fail to find price for account - %v", err)
	}
	defer cur.Close(ctx)

	if cur.Next(ctx) {
		var price Price
		err = cur.Decode(&price)
		if err != nil {
			return nil, fmt.Errorf("fail to find price for account - %v", err)
		}

		return &price, nil
	}

	return nil, ErrNotFound
}

func (ps *pricesStore) FetchPricesForAccount(ctx context.Context, offset, maxCount int, pf PricesFilter) ([]Price, int, error) {
	generalPriceFilter := pf.BsonFilter()
	generalPriceFilter = append(generalPriceFilter,
		bson.E{
			Key: "account", Value: nil,
		},
	)

	// 1. First we take all items for general price
	count, err := ps.store.Collection("prices").CountDocuments(ctx, generalPriceFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch prices count - %v", err)
	}

	if int(count) < offset {
		return nil, int(count), ErrPageOutOfBounds
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(maxCount)).
		SetSort(bson.M{"country": 1, "operator": 1})
	cur1, err := ps.store.Collection("prices").Find(ctx, generalPriceFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch prices - %v", err)
	}

	defer cur1.Close(ctx)

	generalPrices := make([]Price, 0)
	err = cur1.All(ctx, &generalPrices)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch prices - %v", err)
	}

	if pf.AccountID == nil {
		return generalPrices, int(count), nil
	}

	// 2. In we need to fetch price for particular account - we additionally fetch prices
	// for account and replace general price items with found items

	accountPriceFilter := bson.D{
		bson.E{
			Key: "account", Value: *pf.AccountID,
		},
	}

	if pf.Currency != nil {
		accountPriceFilter = append(accountPriceFilter,
			bson.E{
				Key: "currency", Value: *pf.Currency,
			},
		)
	}

	if pf.Enabled != nil {
		accountPriceFilter = append(accountPriceFilter,
			bson.E{
				Key: "enabled", Value: *pf.Enabled,
			},
		)
	}

	cur2, err := ps.store.Collection("prices").Find(ctx, accountPriceFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch prices - %v", err)
	}
	defer cur2.Close(ctx)

	accountPrices := make([]Price, 0)
	err = cur2.All(ctx, &accountPrices)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch prices - %v", err)
	}

	accountPricesSet := make(map[MccMnc]Price, 0)
	for _, accountPrice := range accountPrices {
		accountPricesSet[accountPrice.MccMnc] = accountPrice
	}

	// 3. Replace general price items with found account price items
	for idx, generalPrice := range generalPrices {
		accountPrice, found := accountPricesSet[generalPrice.MccMnc]
		if found {
			generalPrices[idx] = accountPrice
		}
	}

	return generalPrices, int(count), nil
}

func (ps *pricesStore) FetchPrice(ctx context.Context, priceID string) (*Price, error) {
	filter := bson.M{
		"_id": priceID,
	}

	var price Price
	err := ps.store.Collection("prices").FindOne(ctx, filter).Decode(&price)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fail to find price - %v", err)
	}

	return &price, nil
}

func (ps *pricesStore) AddPrice(ctx context.Context, price Price) error {
	price.UpdatedAt = time.Now()
	_, err := ps.store.Collection("prices").InsertOne(ctx, price)
	if err != nil {
		if IsErrDuplication(err) {
			return ErrDuplicate
		} else {
			return fmt.Errorf("fail to add price - %v", err)
		}
	}

	return nil
}

func (ps *pricesStore) UpdatePrice(ctx context.Context, price Price) error {
	filter := bson.M{
		"_id":      price.ID,
		"account":  price.AccountId,
		"mcc":      price.MccMnc.Mcc,
		"mnc":      price.MccMnc.Mnc,
		"currency": price.Currency,
	}
	price.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"updated_at": price.UpdatedAt,
			"enabled":    price.Enabled,
			"price":      price.Price,
		},
	}
	opts := options.Update()
	res, err := ps.store.Collection("prices").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("fail to update price - %v", err)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (ps *pricesStore) DeletePrice(ctx context.Context, price Price) error {
	filter := bson.M{}

	if price.AccountId != nil {
		filter["_id"] = price.ID
	} else {
		// remove all account price items with the same mcc,mnc and currency
		filter["mcc"] = price.Mcc
		filter["mnc"] = price.Mnc
		filter["currency"] = price.Currency
	}

	_, err := ps.store.Collection("prices").DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("fail to delete price - %v", err)
	}

	return nil
}
