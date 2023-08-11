package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/v1/shorten", app.shortenMuxHandler) // GET|POST Return all URLs or create a new one
	mux.HandleFunc("/v1/", app.redirectUrlHandler)       // GET redirects to the original URL
	// mux.HandleFunc("/", app.NotFoundHandler)             // GET returns 404 <--- This is not working
	return mux
}
