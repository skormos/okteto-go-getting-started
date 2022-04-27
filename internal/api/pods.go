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

var validPodListSortFields = map[ListPodsParamsSort]struct{}{
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

	listPodsParamsWrapper ListPodsParams
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

	paramsWrapper := listPodsParamsWrapper(params)
	sort, err := paramsWrapper.sort()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.logger.Info().Msgf("sort called %s", sort)

	pods, err := p.podsLister.ListPods(r.Context(), string(namespace))
	if err != nil {
		p.logger.Err(err).Msg("while executing pods")
	}

	if err := respond(w, listPodsResponse(paramsWrapper, pods), http.StatusOK); err != nil {
		p.logger.Err(err).Msg("while responding to list pods")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func listPodsResponse(params listPodsParamsWrapper, input []cluster.Pod) ListPodsResponse {

	pods := make([]Pod, 0, len(input))
	for _, p := range input {
		pods = append(pods, Pod{
			Age:      p.Age,
			Name:     p.Name,
			Restarts: int(p.Restarts),
		})
	}

	return ListPodsResponse{
		Limit:  params.limit(),
		Offset: params.offset(),
		Pods:   pods,
		Total:  TotalRecords(len(pods)),
	}
}

func (p listPodsParamsWrapper) offset() RecordOffset {
	if p.Offset != nil {
		return *p.Offset
	}

	return 0
}

func (p listPodsParamsWrapper) limit() RecordLimit {
	if p.Limit != nil {
		return *p.Limit
	}

	return 100
}

func (p listPodsParamsWrapper) sort() (ListPodsParamsSort, error) {
	if p.Sort != nil {
		if _, isValid := validPodListSortFields[*p.Sort]; !isValid {
			return "", fmt.Errorf("sort query value %s is invalid", *p.Sort)
		}
		return *p.Sort, nil
	}

	return podsListSortName, nil
}
