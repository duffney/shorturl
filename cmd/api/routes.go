package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/v1/shorten", app.shortenHandler) // POST passes data to be shortened
	mux.HandleFunc("/v1/", app.redirectHandler)       // GET redirects to the original URL
	mux.HandleFunc("/v1/shortlinks", app.listHandler) // GET returns a list of shortened URLs
	mux.HandleFunc("/", app.NotFoundHandler)          // GET returns 404
	return mux
}
