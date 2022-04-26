package main

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/okteto/go-getting-started/internal/api"
)

func apiHandler(logCtx zerolog.Context) http.Handler {
	return api.New(logCtx)
}
