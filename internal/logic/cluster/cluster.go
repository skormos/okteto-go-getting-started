package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// ClusterOps manages operations for a kubernetes cluster
	ClusterOps struct { //nolint:revive // yeah, this stutters. Will fix when a better name comes to mind.
		logger zerolog.Logger
		client *kubernetes.Clientset
	}

	// Pod is a simple representation of a Pod
	Pod struct {
		Name     string
		Age      string
		Restarts int
	}
)

// New creates a new instance of ClusterOps given the client set.
func New(logCtx zerolog.Context, client *kubernetes.Clientset) (*ClusterOps, error) {
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
	list, err := o.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get pod list for namespace %s: %w", namespace, err)
	}

	pods := make([]Pod, 0, len(list.Items))
	for _, podItem := range list.Items {
		createTime := podItem.GetCreationTimestamp().Time
		ageDuration := time.Since(createTime)

		for _, container := range podItem.Spec.Containers {
			o.logger.Info().Msgf("Pod %s Container %s", podItem.GetName(), container.Name)
		}

		pods = append(pods, Pod{
			Name:     podItem.GetName(),
			Age:      ageDuration.String(),
			Restarts: 3000,
		})
	}

	return pods, nil
}
