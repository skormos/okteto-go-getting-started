package main

import (
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logCtx := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Stack()

	mainLogger := logCtx.Str("module", "main").
		Logger()

	port := "8080"
	mainLogger.Info().Msgf("Starting http server on port %s", port)

	http.Handle("/", helloHandler(logCtx))
	if err := http.ListenAndServe(net.JoinHostPort("", port), nil); err != nil {
		mainLogger.Panic().Err(err).Msg("while starting http server")
	}

}
