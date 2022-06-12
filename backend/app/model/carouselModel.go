package model

import (
	"backend/app/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Carousel struct {
	ID string `bson:"_id" json:"_id"`

	Title      string `bson:"title" json:"title"`
	ButtonText string `bson:"btnText" json:"btnText"`
	ImageUrl   string `bson:"imageUrl" json:"imageUrl"`
	IsCurrent  bool   `bson:"isCurrent" json:"isCurrent"`
}

type CarouselStore interface {
	EnsureIndexes(ctx context.Context) error

	FetchList() ([]*Carousel, int, error)
}

func NewCarouselStore(store Storage) (CarouselStore, error) {
	return &carouselStore{
		store: store,
	}, nil
}

type carouselStore struct {
	store Storage
}

func (cs *carouselStore) EnsureIndexes(ctx context.Context) error {
	log := logger.GetLogger(ctx).WithField("store", "carousel")

	log.Info("Fetching storage indexes")

	//filter := bson.D{}
	//
	//cur, err := cs.store.Collection("carousel").Find(context.TODO(), filter)
	//if err != nil {
	//	return fmt.Errorf("fail to fetch carousel - %v", err)
	//}
	//
	//var results []bson.M
	//if err = cur.All(context.TODO(), &results); err != nil {
	//	panic(err)
	//}
	//for _, result := range results {
	//	output, err := json.MarshalIndent(result, "", "    ")
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("%s\n", output)
	//}

	log.Info("Storage indexes are up-to-date")

	return nil
}

func (cs *carouselStore) FetchList() ([]*Carousel, int, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	filter := bson.D{}

	count, err := cs.store.Collection("carousel").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch carousel count - %v", err)
	}

	carouselItems := make([]*Carousel, 0)
	cur, err := cs.store.Collection("carousel").Find(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch carousel - %v", err)
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &carouselItems)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to fetch users - %v", err)
	}

	return carouselItems, int(count), nil
}
