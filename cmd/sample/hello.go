package main

import (
	"net/http"

	"github.com/okteto/go-getting-started/internal/api"
)

func helloHandler() http.Handler {
	return api.NewHelloHandler()
}
