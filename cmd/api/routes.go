package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	return mux
}
