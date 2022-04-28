package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
)

type (
	// ClusterOps manages operations for a kubernetes cluster
	ClusterOps struct { //nolint:revive // yeah, this stutters. Will fix when a better name comes to mind.
		logger zerolog.Logger
		client core.CoreV1Interface
	}

	// Pod is a simple representation of a Pod
	Pod struct {
		Name     string
		Age      string
		Restarts int32
	}
)

// New creates a new instance of ClusterOps given the client set.
func New(logCtx zerolog.Context, client core.CoreV1Interface) (*ClusterOps, error) {
	logger := logCtx.Str("module", "logic").
		Str("handler", "kubernetes-cluster").
		Logger()

	return &ClusterOps{
		logger: logger,
		client: client,
	}, nil
}

// ListPods gets all the pods in the provided namespace for the given cluster.
func (o *ClusterOps) ListPods(ctx context.Context, namespace string) ([]Pod, error) {
	list, err := o.client.Pods(namespace).List(ctx, meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get pod list for namespace %s: %w", namespace, err)
	}

	pods := make([]Pod, 0, len(list.Items))
	for _, podItem := range list.Items {
		createTime := podItem.GetCreationTimestamp().Time
		ageDuration := time.Since(createTime)

		restartCount := int32(0)

		for _, container := range podItem.Status.ContainerStatuses {
			o.logger.Info().Msgf("Pod %s Container %s", podItem.GetName(), container.Name)
			if strings.HasPrefix(podItem.GetName(), container.Name) {
				restartCount = container.RestartCount
			}
		}

		pods = append(pods, Pod{
			Name:     podItem.GetName(),
			Age:      ageDuration.String(),
			Restarts: restartCount,
		})
	}

	return pods, nil
}
