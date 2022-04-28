package main

import (
	"fmt"

	"github.com/okteto/go-getting-started/internal/logic/cluster"
	"github.com/rs/zerolog"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
)

func newClusterOps(logCtx zerolog.Context, k8sCoreClient core.CoreV1Interface) (*cluster.ClusterOps, error) {
	clusterOps, err := cluster.New(logCtx, k8sCoreClient)
	if err != nil {
		return nil, fmt.Errorf("while creating a ClusterOps wrapper %w", err)
	}

	return clusterOps, nil
}
