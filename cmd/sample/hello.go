package main

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/okteto/go-getting-started/internal/api"
)

func helloHandler(logCtx zerolog.Context) http.Handler {
	return api.NewHelloHandler(logCtx)
}
