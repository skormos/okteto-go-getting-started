package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

type (
	helloSayer interface {
		SayHello(w http.ResponseWriter, r *http.Request, params SayHelloParams)
	}

	sayHelloHandler struct {
		logger zerolog.Logger
	}
)

func newSayHelloHandler(logCtx zerolog.Context) sayHelloHandler {
	logger := logCtx.Str("module", "http").
		Str("handler", "hello").
		Logger()

	return sayHelloHandler{
		logger: logger,
	}
}

func (s sayHelloHandler) SayHello(w http.ResponseWriter, _ *http.Request, params SayHelloParams) {
	name := "World"
	if params.Name != nil && strings.TrimSpace(*params.Name) != "" {
		name = strings.TrimSpace(*params.Name)
	}

	response := Greeting{
		Greeting: fmt.Sprintf("Hello %s!", name),
	}

	if err := respond(w, &response, http.StatusOK); err != nil {
		s.logger.Err(err).Msgf("while responding with %v", response)
		http.Error(w, "Could not complete hello world request.", http.StatusInternalServerError)
	}
}
