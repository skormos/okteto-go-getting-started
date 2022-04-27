package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/okteto/go-getting-started/internal/logic/cluster"
)

const (
	podsListSortAge      ListPodsParamsSort = "age"
	podsListSortName     ListPodsParamsSort = "name"
	podsListSortRestarts ListPodsParamsSort = "restarts"
)

var validpodListSortFields = map[ListPodsParamsSort]struct{}{
	podsListSortAge:      {},
	podsListSortName:     {},
	podsListSortRestarts: {},
}

type (
	// PodsLister defines the required contract for this endpoint to get the Pods from the downstream data source.
	PodsLister interface {
		ListPods(ctx context.Context, namespace string) ([]cluster.Pod, error)
	}

	podsHandler struct {
		logger     zerolog.Logger
		podsLister PodsLister
	}
)

func newPodsHandler(logCtx zerolog.Context, podsLister PodsLister) podsHandler {
	logger := logCtx.Str("module", "http").
		Str("handler", "pods").
		Logger()

	return podsHandler{
		logger:     logger,
		podsLister: podsLister,
	}
}

func (p podsHandler) ListPods(w http.ResponseWriter, r *http.Request, namespace NamespacePath, params ListPodsParams) {
	//nolint:godox // ignore FIXMEs
	//FIXME use a validator library to validate the params schema

	sort := podsListSortName
	if params.Sort != nil {
		if _, isValid := validpodListSortFields[*params.Sort]; !isValid {
			http.Error(w, fmt.Errorf("sort query value %s is invalid", *params.Sort).Error(), http.StatusBadRequest)
			return
		}
		sort = *params.Sort
	}

	p.logger.Info().Msgf("sort called %s", sort)

	pods, err := p.podsLister.ListPods(r.Context(), string(namespace))
	if err != nil {
		p.logger.Err(err).Msg("while executing pods")
	}

	if err := respond(w, len(pods), http.StatusOK); err != nil {
		p.logger.Err(err).Msg("while responding to list pods")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
