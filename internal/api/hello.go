package api

import (
	"net/http"
)

// NewHelloHandler creates a new http Handler that responds with Hello World, and an error if something is wrong.
func NewHelloHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		if _, err := writer.Write([]byte("Hello World!")); err != nil {
			http.Error(writer, "Could not complete hello world request.", http.StatusInternalServerError)
		}
	})
}
