package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"k8s.io/client-go/kubernetes"

	"github.com/okteto/go-getting-started/internal/logic/cluster"
)

func newClusterOps(logCtx zerolog.Context, k8sClient *kubernetes.Clientset) (*cluster.ClusterOps, error) {
	clusterOps, err := cluster.New(logCtx, k8sClient)
	if err != nil {
		return nil, fmt.Errorf("while creating a ClusterOps wrapper %w", err)
	}

	return clusterOps, nil
}
