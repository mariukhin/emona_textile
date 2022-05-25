package main

import (
	"fmt"
	"log"
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

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	//register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/category", category)
	mux.HandleFunc("/create", createSnippet)
	// Use the http.ListenAndServe() function to start a new web server. We pass in
	//two parameters: the TCP network address to listen on (in this case ":4000")
	//and the servemux we just created. If http.ListenAndServe() returns an error
	//we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
