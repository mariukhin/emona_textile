package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

// Define a home handler function which writes a byte slice containing
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Emona backend"))
}

func httpJsonResponse(w http.ResponseWriter, resp interface{}) {
	httpJsonResponseCode(w, resp, http.StatusOK)
}

func httpJsonResponseCode(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"error": [
				"internal"
			]
		}`)
		return
	}
}

type MongoDBCommand struct {
	MongoURI      string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb+srv://zorkiy:admin@emonacluster.udns5gz.mongodb.net/?retryWrites=true&w=majority"`
	MongoDatabase string `long:"mongo-database" env:"MONGO_DB" default:"EmonaDB" description:"MongoDB database name"`
}

type Carousel struct {
	ID string `bson:"_id" json:"id"`

	Title      string `bson:"title" json:"title"`
	ButtonText string `bson:"btnText" json:"btnText"`
	ImageUrl   string `bson:"imageUrl" json:"imageUrl"`
	IsCurrent  bool   `bson:"isCurrent" json:"isCurrent"`
}

type Catalog struct {
	ID string `bson:"_id" json:"id"`

	Title      string `bson:"title" json:"title"`
	ImageUrl   string `bson:"imageUrl" json:"imageUrl"`
}

func carousel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://zorkiy:admin@emonacluster.udns5gz.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("EmonaDB").Collection("carousel")

	carouselItems := make([]*Carousel, 0)
	cursor, err := coll.Find(context.TODO(), bson.D{})

	err = cursor.All(ctx, &carouselItems)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	httpJsonResponse(w, carouselItems)
}

func catalog(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://zorkiy:admin@emonacluster.udns5gz.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("EmonaDB").Collection("catalog")

	catalogItems := make([]*Catalog, 0)
	cursor, err := coll.Find(context.TODO(), bson.D{})

	err = cursor.All(ctx, &catalogItems)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	httpJsonResponse(w, catalogItems)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not. Note that
	//http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodPost {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		//response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		w.Header().Set("Allow", http.MethodPost)
		// Set a new cache-control header. If an existing "Cache-Control" header exists // it will be overwritten.
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		// In contrast, the Add() method appends a new "Cache-Control" header and can // be called multiple times.
		w.Header().Add("Cache-Control", "public")
		w.Header().Add("Cache-Control", "max-age=31536000")
		// Delete all values for the "Cache-Control" header.
		w.Header().Del("Cache-Control")
		// Retrieve the first value for the "Cache-Control" header.
		w.Header().Get("Cache-Control")
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"name":"Max"}`))
	w.Write([]byte("Create a new snippet..."))
}
