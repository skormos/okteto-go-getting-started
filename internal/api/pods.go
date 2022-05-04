package api

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/rs/zerolog"

	"github.com/okteto/go-getting-started/internal/logic/cluster"
)

const (
	podsListSortAge      ListPodsParamsSort = "age"
	podsListSortName     ListPodsParamsSort = "name"
	podsListSortRestarts ListPodsParamsSort = "restarts"

	podsListOffsetDefault RecordOffset = 0
	podsListLimitDefault  RecordLimit  = 100
)

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
	//FIXME use a validator library to validate the params schema for max, min and sort.

	pods, err := p.podsLister.ListPods(r.Context(), string(namespace))
	if err != nil {
		p.logger.Err(err).Msg("while listing pods")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusInternalServerError)
		return
	}

	paramsWrapper := listPodsParamsWrapper(params)
	response, err := listPodsResponse(paramsWrapper, pods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := respond(w, response, http.StatusOK); err != nil {
		p.logger.Err(err).Msg("while responding to list pods")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func listPodsResponse(params listPodsParamsWrapper, input []cluster.Pod) (ListPodsResponse, error) {

	offset := int(params.offset())

	// if offset is past the last record
	if offset >= len(input) {
		return ListPodsResponse{
			Limit:  params.limit(),
			Offset: params.offset(),
			Pods:   []Pod{},
			Total:  TotalRecords(len(input)),
		}, nil
	}

	// sort the pods according to the requested field
	sorted, err := sortPods(input, params.sort())
	if err != nil {
		return ListPodsResponse{}, err
	}

	endIdx := offset + int(params.limit())
	// if the limit is passed the last record, get the last record
	if endIdx > len(input) {
		endIdx = len(input)
	}

	responsePods := make([]Pod, 0, endIdx-offset)
	for i := offset; i < endIdx; i++ {
		p := sorted[i]
		responsePods = append(responsePods, Pod{
			Age:      p.Age.String(),
			Name:     p.Name,
			Restarts: int(p.Restarts),
		})
	}

	return ListPodsResponse{
		Limit:  params.limit(),
		Offset: params.offset(),
		Pods:   responsePods,
		Total:  TotalRecords(len(input)),
	}, nil
}

func (p listPodsParamsWrapper) offset() RecordOffset {
	if p.Offset != nil {
		return *p.Offset
	}

	return podsListOffsetDefault
}

func (p listPodsParamsWrapper) limit() RecordLimit {
	if p.Limit != nil {
		return *p.Limit
	}

	return podsListLimitDefault
}

func (p listPodsParamsWrapper) sort() ListPodsParamsSort {
	if p.Sort != nil {
		return *p.Sort
	}

	return podsListSortName
}

func sortPods(input []cluster.Pod, sortParam ListPodsParamsSort) ([]cluster.Pod, error) {
	switch sortParam {
	case podsListSortAge:
		sort.Slice(input, func(i, j int) bool {
			if input[i].Age == input[j].Age {
				return input[i].Name < input[j].Name
			}
			return input[i].Age > input[j].Age
		})
		return input, nil
	case podsListSortName:
		sort.Slice(input, func(i, j int) bool {
			return input[i].Name < input[j].Name
		})
		return input, nil
	case podsListSortRestarts:
		sort.Slice(input, func(i, j int) bool {
			if input[i].Restarts == input[j].Restarts {
				return input[i].Name < input[j].Name
			}
			return input[i].Restarts > input[j].Restarts
		})
		return input, nil
	default:
		return nil, fmt.Errorf("sort query value %s is invalid", sortParam)
	}
}
