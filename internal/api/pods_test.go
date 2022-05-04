package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/okteto/go-getting-started/internal/logic/cluster"
)

type testPodsLister struct {
	t                 *testing.T
	expectedNamespace string
	podsData          []cluster.Pod
}

func (l testPodsLister) ListPods(_ context.Context, namespace string) ([]cluster.Pod, error) {
	assert.Equal(l.t, l.expectedNamespace, namespace)
	return l.podsData, nil
}

func testPodsData(t *testing.T) (map[string]cluster.Pod, []cluster.Pod) {
	t.Helper()

	values := []cluster.Pod{
		{"Alpha", time.Minute, 0},
		{"Bravo", 3 * time.Minute, 4},
		{"Alphaaaa", 61 * time.Second, 1},
		{"Delta", time.Minute, 2},
		{"Charlie", time.Minute, 0},
		{"Foxtrot", 30 * time.Minute, 2},
		{"Echo", 20 * time.Second, 0},
	}

	mapped := make(map[string]cluster.Pod, len(values))
	for _, pod := range values {
		mapped[pod.Name] = pod
	}

	return mapped, values
}

func expectedBody(t *testing.T, limit, offset, total int, expectedPods []interface{}) map[string]interface{} {
	t.Helper()

	response := map[string]interface{}{
		"limit":  float64(limit),
		"offset": float64(offset),
		"total":  float64(total),
	}

	pods := make([]interface{}, 0, len(expectedPods))
	for _, pod := range expectedPods {
		pods = append(pods, pod)
	}

	response["pods"] = pods

	return response
}

func expectedPods(t *testing.T, in ...cluster.Pod) []interface{} {
	t.Helper()

	out := make([]interface{}, 0, len(in))
	for _, item := range in {
		out = append(out, map[string]interface{}{
			"age":      item.Age.String(),
			"name":     item.Name,
			"restarts": float64(item.Restarts),
		})
	}
	return out
}

func expectedListParams(t *testing.T, limit, offset int, sort ListPodsParamsSort) ListPodsParams {
	lmt := RecordLimit(limit)
	off := RecordOffset(offset)

	return ListPodsParams{
		Limit:  &lmt,
		Offset: &off,
		Sort:   &sort,
	}
}

func TestPodsHandler_ListPods(t *testing.T) {
	expectedNamespace := "test-namespace"

	podsByName, podsData := testPodsData(t)

	tests := map[string]struct {
		expectedStatusCode int
		expectedBody       map[string]interface{}
		podsData           []cluster.Pod
		params             ListPodsParams
	}{
		"Zero pods with no params returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody:       expectedBody(t, 100, 0, 0, []interface{}{}),
			podsData:           []cluster.Pod{},
			params:             ListPodsParams{},
		},
		"Single pod with no params returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 100, 0, 1,
				expectedPods(t, podsByName["Alpha"]),
			),
			podsData: []cluster.Pod{podsByName["Alpha"]},
			params:   ListPodsParams{},
		},
		"All pods under limit and no params returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 100, 0, len(podsData),
				expectedPods(t,
					podsByName["Alpha"],
					podsByName["Alphaaaa"],
					podsByName["Bravo"],
					podsByName["Charlie"],
					podsByName["Delta"],
					podsByName["Echo"],
					podsByName["Foxtrot"],
				),
			),
			podsData: podsData,
			params:   ListPodsParams{},
		},
		"All pods with no offset and a low limit returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 3, 0, len(podsData),
				expectedPods(t,
					podsByName["Alpha"],
					podsByName["Alphaaaa"],
					podsByName["Bravo"],
				),
			),
			podsData: podsData,
			params:   expectedListParams(t, 3, 0, podsListSortName),
		},
		"All pods with an in-range offset and a low limit returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 4, 2, len(podsData),
				expectedPods(t,
					podsByName["Bravo"],
					podsByName["Charlie"],
					podsByName["Delta"],
					podsByName["Echo"],
				),
			),
			podsData: podsData,
			params:   expectedListParams(t, 4, 2, podsListSortName),
		},
		"All pods under limit sort by Age returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 100, 0, len(podsData),
				expectedPods(t,
					podsByName["Foxtrot"],
					podsByName["Bravo"],
					podsByName["Alphaaaa"],
					podsByName["Alpha"],
					podsByName["Charlie"],
					podsByName["Delta"],
					podsByName["Echo"],
				),
			),
			podsData: podsData,
			params:   expectedListParams(t, 100, 0, podsListSortAge),
		},
		"All pods under limit sort by Restarts returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 100, 0, len(podsData),
				expectedPods(t,
					podsByName["Bravo"],
					podsByName["Delta"],
					podsByName["Foxtrot"],
					podsByName["Alphaaaa"],
					podsByName["Alpha"],
					podsByName["Charlie"],
					podsByName["Echo"],
				),
			),
			podsData: podsData,
			params:   expectedListParams(t, 100, 0, podsListSortRestarts),
		},
		"Offset last page with data returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: expectedBody(t, 4, 4, len(podsData),
				expectedPods(t,
					podsByName["Delta"],
					podsByName["Echo"],
					podsByName["Foxtrot"],
				),
			),
			podsData: podsData,
			params:   expectedListParams(t, 4, 4, podsListSortName),
		},
		"Offset past last records returns no records": {
			expectedStatusCode: http.StatusOK,
			expectedBody:       expectedBody(t, 100, 100, len(podsData), []interface{}{}),
			podsData:           podsData,
			params:             expectedListParams(t, 100, 100, podsListSortName),
		},
	}

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			target := newPodsHandler(
				testLoggerContext(tt),
				testPodsLister{tt, expectedNamespace, test.podsData},
			)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			defer func() {
				_ = response.Result().Body.Close()
			}()

			target.ListPods(response, request, NamespacePath(expectedNamespace), test.params)
			assert.Equal(tt, test.expectedStatusCode, response.Result().StatusCode)

			var actualBody map[string]interface{}
			require.NoError(tt, json.NewDecoder(response.Result().Body).Decode(&actualBody))

			assert.Equal(tt, test.expectedBody, actualBody)
		})
	}
}
