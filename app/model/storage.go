package model

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

var ErrNotFound = errors.New("document not found")
var ErrPageOutOfBounds = errors.New("page out of bounds")
var ErrDuplicate = errors.New("document duplication")

type StorageOptions struct {
	MongoURI      string
	MongoDatabase string
}

func NewStorageOptions() *StorageOptions {
	return &StorageOptions{
		MongoURI:      "mongodb://localhost:27017/",
		MongoDatabase: "EmonaDB",
	}
}

type Storage interface {
	Collection(name string) *mongo.Collection
	FetchCollectionIndexes(context context.Context, collection string) ([]bson.M, error)
	StartSession() (mongo.Session, error)
	DefaultTransactionOpts() *options.TransactionOptions
}

func (opt *StorageOptions) Storage() (Storage, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
	clientOpt := options.Client().ApplyURI(opt.MongoURI)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		return nil, fmt.Errorf("fail to create storage - %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create storage - %v", err)
	}

	db := client.Database(opt.MongoDatabase)

	return &storage{
		client: client,
		db:     db,
	}, nil
}

type storage struct {
	client *mongo.Client
	db     *mongo.Database
}

func (s *storage) Collection(name string) *mongo.Collection {
	return s.db.Collection(name)
}

func (s *storage) FetchCollectionIndexes(context context.Context, collection string) ([]bson.M, error) {
	opts := options.ListIndexes()
	indexesCur, err := s.Collection(collection).Indexes().List(context, opts)
	if err != nil {
		return nil, err
	}

	defer indexesCur.Close(context)

	var indexes []bson.M
	err = indexesCur.All(context, &indexes)
	if err != nil {
		return nil, err
	}

	return indexes, nil
}

func (s *storage) StartSession() (mongo.Session, error) {
	return s.client.StartSession()
}

func (s *storage) DefaultTransactionOpts() *options.TransactionOptions {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	return options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
}

func IsErrDuplication(err error) bool {
	if err != nil {
		writeErr, ok := err.(mongo.WriteException)
		if ok && len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 {
			return true
		}

		commandErr, ok := err.(mongo.CommandError)
		if ok && commandErr.Code == 11000 {
			return true
		}
	}

	return false
}
