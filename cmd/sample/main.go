package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	k8sCoreClient, err := k8sInClusterCoreClient()
	if err != nil {
		mainLogger.Panic().Err(err).Msg("while initializing kubernetes client")
	}
	clusterOps, err := newClusterOps(logCtx, k8sCoreClient)
	if err != nil {
		mainLogger.Panic().Err(err).Msg("while initializing ClusterOps logic")
	}

	ctx := context.Background()
	podsMetricsCancel, err := initPodsMetricsGauge(ctx, logCtx, k8sCoreClient, "skormos", 5*time.Second)
	defer func() {
		podsMetricsCancel()
	}()
	if err != nil {
		mainLogger.Panic().Err(err).Msg("while initializing Pods Gauge metric")
	}

	apiHandler := newAPIHandler(logCtx, clusterOps)
	metricsHandler := newMetricsHandler()
	pingHandler, err := newPingHandler()
	if err != nil {
		mainLogger.Panic().Err(err).Msg("while creating ping Handler")
	}
	rootHandler := newHandler(apiHandler, metricsHandler, pingHandler)

	httpOptions := newHTTPServerOptions("8080")
	httpServer := newHTTPServerWrapper(logCtx, httpOptions, rootHandler, serverErrors)

	select {
	case err := <-serverErrors:
		podsMetricsCancel()
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
