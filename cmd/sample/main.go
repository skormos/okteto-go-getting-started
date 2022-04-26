package main

import (
	"os"
	"os/signal"
	"syscall"

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

	shutdown := listenForShutdown()
	serverErrors := make(chan error, 1)

	httpOptions := newHTTPServerOptions("8080")
	httpServer := newHTTPServerWrapper(logCtx, httpOptions, helloHandler(logCtx), serverErrors)

	select {
	case err := <-serverErrors:
		mainLogger.Panic().Err(err).Msg("error received from the server")
	case sig := <-shutdown:
		if syscall.SIGSTOP == sig {
			mainLogger.Info().Msg("integrity issue has invoked a shutdown...")
		} else {
			mainLogger.Info().Msgf("%v : server shutdown requested...", sig)
		}

		httpServer.Stop()
	}
}

func listenForShutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	return shutdown
}
