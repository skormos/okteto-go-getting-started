package main

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/okteto/go-getting-started/internal/api"
	"github.com/okteto/go-getting-started/internal/logic/cluster"
)

func apiHandler(logCtx zerolog.Context, clusterOps *cluster.ClusterOps) http.Handler {
	return api.New(logCtx, clusterOps)
}
