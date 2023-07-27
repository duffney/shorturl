package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "available",
		"version": version,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Endpoint not found"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonResponse)
}

func (app *application) shortenHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	/*
		psudo code:
		1. get the url from the request body
		2. validate the url (later)
		3. generate a unique id
		4. save the url and the id in a map
		5. pass the id to the shortener function
		6. return the shortened url
	*/

	var input struct {
		Url string `json:"url"`
	}
	// convert the JSON request body to a struct
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	generator, err := NewIDGenerator(1)
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	id := generator.GenerateID()
	// change print to include ID:
	fmt.Println("Generated ID:", id)
	hash := DecimalToBase62(id)
	fmt.Println("Hashed ID:", hash)

	// save the url and the id in a map
	s := Shorten{
		// Id:        id,
		Long_url:  input.Url,
		Short_url: shortenerAddress + hash, // combine hash with shorturl address
	}

	// store Shorten in a map
	app.db[id] = s

	// marshal the database to JSON
	js, err := json.MarshalIndent(app.db, "", "\t")
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}
	// print the js to standard out
	fmt.Println("Database:", string(js))

	fmt.Fprintf(w, "%+v\n", app.db[id].Short_url)
}

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the hash from the url
	hash := r.URL.Path[len("/v1/"):]
	// reverse hash to id
	id := Base62ToDecimal(hash)
	// get the url from the map
	url := app.db[id].Long_url
	// 301 redirect to the url
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
