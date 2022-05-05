package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type (
	apiPodsGetter interface {
		ListPods(w http.ResponseWriter, r *http.Request, namespace NamespacePath, params ListPodsParams)
	}

	api struct {
		helloSayer
		apiPodsGetter
	}
)

// New creates a new router to handle requests to the api.
func New(logCtx zerolog.Context, podsLister PodsLister) http.Handler {
	logger := logCtx.Str("module", "http").Logger()

	si := api{
		helloSayer:    newSayHelloHandler(logCtx),
		apiPodsGetter: newPodsHandler(logCtx, podsLister),
	}

	//nolint:godox // ignoring task comments for linting
	// FIXME Middleware: CORS, RequestLogger, etc.
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter:  chi.NewRouter(),
		Middlewares: nil,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Err(err).Msgf("while calling %s", r.RequestURI)

			//nolint:godox // ignoring task comments for linting
			// FIXME this needs much more granular Error handling.
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	})
}

func respond(writer http.ResponseWriter, input interface{}, status int) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("while marshalling %v for http response: %w", input, err)
	}

	writer.WriteHeader(status)
	if _, err := writer.Write(bytes); err != nil {
		return fmt.Errorf("while writing %v as bytes to Response: %w", input, err)
	}

	return nil
}
