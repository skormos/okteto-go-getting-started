package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newHandler(api, metrics, ping http.Handler) http.Handler {
	router := chi.NewRouter()

	router.Mount("/api", api)

	router.Handle("/metrics", metrics)
	router.Handle("/ping", ping)

	return router
}
