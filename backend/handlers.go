package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Emona backend"))
}

func category(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	//not found response.
	name, err := strconv.Atoi(r.URL.Query().Get("name"))
	if err != nil || name < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the name value with our response
	//and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a catehory name from URL %d...", name)
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
