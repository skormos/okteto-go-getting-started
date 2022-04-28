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
}

func (t testPodsLister) ListPods(_ context.Context, namespace string) ([]cluster.Pod, error) {
	assert.Equal(t.t, t.expectedNamespace, namespace)
	return []cluster.Pod{
		{
			Name:     "first",
			Age:      (3 * time.Minute).String(),
			Restarts: 0,
		},
	}, nil
}

func TestPodsHandler_ListPods(t *testing.T) {
	expectedNamespace := "test-namespace"

	tests := map[string]struct {
		expectedStatusCode int
		expectedBody       map[string]interface{}
		params             ListPodsParams
	}{
		"Single pod with no params returns successfully": {
			expectedStatusCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"pods": []interface{}{
					map[string]interface{}{
						"age":      "3m0s",
						"name":     "first",
						"restarts": float64(0),
					},
				},
				"limit":  float64(100),
				"offset": float64(0),
				"total":  float64(1),
			},
			params: ListPodsParams{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			target := newPodsHandler(testLoggerContext(tt), testPodsLister{tt, expectedNamespace})

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
