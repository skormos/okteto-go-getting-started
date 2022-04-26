package api

import "github.com/go-chi/chi/v5"

// New creates a new router to handle requests to the api.
func New() chi.Router {
	return chi.NewRouter()
}
