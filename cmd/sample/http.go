package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type httpServerWrapper struct {
	logger          zerolog.Logger
	server          *http.Server
	shutdownTimeout time.Duration
}

type httpServerOptions struct {
	host            string
	port            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

func newHTTPServerOptions(port string) *httpServerOptions {
	return &httpServerOptions{
		port:            port,
		readTimeout:     120 * time.Second,
		writeTimeout:    120 * time.Second,
		shutdownTimeout: 30 * time.Second,
	}
}

// Host allows the caller to set the host string for the server.
func (o *httpServerOptions) Host(host string) *httpServerOptions {
	o.host = host
	return o
}

// ReadTimeout overrides ReadTimeout for the underlying http server.
func (o *httpServerOptions) ReadTimeout(timeout time.Duration) *httpServerOptions {
	o.readTimeout = timeout
	return o
}

// ShutdownTimeout set the timeout on the cancel context for shutting down the server.
func (o *httpServerOptions) ShutdownTimeout(timeout time.Duration) *httpServerOptions {
	o.shutdownTimeout = timeout
	return o
}

// WriteTimeout overrides the WriteTime for the underlying http server.
func (o *httpServerOptions) WriteTimeout(timeout time.Duration) *httpServerOptions {
	o.writeTimeout = timeout
	return o
}

// Stop attempts to shut down the underlying http server.
func (h *httpServerWrapper) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), h.shutdownTimeout)
	defer cancel()

	err := h.server.Shutdown(ctx)
	if err != nil {
		h.logger.Info().Err(err).Msgf("Graceful shutdown did not complete in duration: %s. Attempting close...", h.shutdownTimeout)
		if err := h.server.Close(); err != nil {
			h.logger.Info().Err(err).Msg("While closing server...")
		} else {
			h.logger.Info().Msg("Successfully closed server.")
		}
	} else {
		h.logger.Info().Msg("Successfully completed graceful shutdown.")
	}
}

func newHTTPServerWrapper(logCtx zerolog.Context, options *httpServerOptions, handler http.Handler, errors chan<- error) *httpServerWrapper {
	addr := net.JoinHostPort(options.host, options.port)

	logger := logCtx.Str("module", "http").Logger()

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  options.readTimeout,
		WriteTimeout: options.writeTimeout,
	}

	wrapper := &httpServerWrapper{
		logger:          logger,
		server:          server,
		shutdownTimeout: options.shutdownTimeout,
	}

	go func() {
		logger.Info().Msgf("starting http server at: %s", addr)
		errors <- server.ListenAndServe()
	}()

	return wrapper
}
