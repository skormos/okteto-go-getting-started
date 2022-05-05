package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Implements the Prometheus quickstart, for testing and validation with the metrics handler in the Okteto cluster
func newPingHandler() (http.Handler, error) {
	pingCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "Number of requests handled by the ping handler.",
		})
	if err := prometheus.Register(pingCounter); err != nil {
		return nil, fmt.Errorf("while registering pingCounter %w", err)
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		pingCounter.Inc()
		_, _ = fmt.Fprint(writer, "PONG")
	}), nil
}

// convenience function to not have to add another import to main.go
func newMetricsHandler() http.Handler {
	return promhttp.Handler()
}

func initPodsMetricsGauge(
	ctx context.Context,
	logCtx zerolog.Context,
	k8sClient core.CoreV1Interface,
	namespace string,
	sleep time.Duration) (func(), error) {

	logger := logCtx.Str("module", "metrics").Str("metric", "kubernetes_pods_gauge").Logger()
	podGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kubernetes_pods_gauge",
		Help: fmt.Sprintf("The number of pods in the %s namespace.", namespace),
	})

	if err := prometheus.Register(podGauge); err != nil {
		return func() {}, fmt.Errorf("while registering pod gauge for namespace %s: %w", namespace, err)
	}

	cancelCtx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
	checkLoop:
		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("Metrics gauge was requested to exit.")
				break checkLoop
			default:
				if pods, err := k8sClient.Pods(namespace).List(ctx, meta.ListOptions{}); err != nil {
					logger.Err(err).Msg("while listing pods for namespace " + namespace)
				} else {
					podGauge.Set(float64(len(pods.Items)))
				}
			}

			time.Sleep(sleep)
		}

	}(cancelCtx)

	return cancel, nil
}
