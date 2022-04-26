package api

import (
	"net/http"

	"github.com/rs/zerolog"
)

// NewHelloHandler creates a new http Handler that responds with Hello World, and an error if something is wrong.
func NewHelloHandler(logCtx zerolog.Context) http.Handler {
	logger := logCtx.Str("module", "http").
		Str("handler", "hello").
		Logger()

	return http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		if _, err := writer.Write([]byte("Hello Worldeee!")); err != nil {
			logger.Err(err).Msg("while trying to server Hello World request")

			http.Error(writer, "Could not complete hello world request.", http.StatusInternalServerError)
			return
		}
	})
}
