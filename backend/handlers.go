package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
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

const (
	defaultMongoURI      = "mongodb://localhost:27017"
	defaultMongoDatabase = "EmonaDB"
)

// mongoConfig holds the MongoDB connection settings. The values come from the
// environment and default to a local instance, so no credentials ever have to
// be stored in the repository.
type mongoConfig struct {
	URI      string
	Database string
}

func mongoConfigFromEnv() mongoConfig {
	return mongoConfig{
		URI:      envOrDefault("MONGO_URI", defaultMongoURI),
		Database: envOrDefault("MONGO_DB", defaultMongoDatabase),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// app carries the dependencies shared by the HTTP handlers.
type app struct {
	mongo mongoConfig
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

// fetchCollection connects to MongoDB and decodes every document of the given
// collection into items, which must be a pointer to a slice.
func (a *app) fetchCollection(ctx context.Context, collection string, items interface{}) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(a.mongo.URI))
	if err != nil {
		return fmt.Errorf("fail to connect to mongo - %w", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("[ERROR] fail to disconnect from mongo - %v", err)
		}
	}()

	cursor, err := client.Database(a.mongo.Database).Collection(collection).Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("fail to query %s - %w", collection, err)
	}

	if err := cursor.All(ctx, items); err != nil {
		return fmt.Errorf("fail to decode %s - %w", collection, err)
	}

	return nil
}

func (a *app) carousel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	carouselItems := make([]*Carousel, 0)
	if err := a.fetchCollection(ctx, "carousel", &carouselItems); err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	httpJsonResponse(w, carouselItems)
}

func (a *app) catalog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	catalogItems := make([]*Catalog, 0)
	if err := a.fetchCollection(ctx, "catalog", &catalogItems); err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
