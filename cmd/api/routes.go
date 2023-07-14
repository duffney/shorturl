package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/v1/data/shorten", app.shortenHandler) // POST passes data to be shortened
	mux.HandleFunc("/v1", app.redirectHandler)             // GET redirects to the original URL
	return mux
}
