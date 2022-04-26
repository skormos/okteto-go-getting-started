package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSayHelloHandler_SayHello(t *testing.T) {
	tests := map[string]struct {
		inputParams      SayHelloParams
		expectedResponse Greeting
	}{
		"A nil name parameter should default successfully": {
			inputParams: SayHelloParams{},
			expectedResponse: Greeting{
				Greeting: "Hello World!",
			},
		},
		"An empty string parameter should default successfully": {
			inputParams: testCreateSayHelloParams(t, ""),
			expectedResponse: Greeting{
				Greeting: "Hello World!",
			},
		},
		"Spaces in the string should be trimmed and default successfully": {
			inputParams: testCreateSayHelloParams(t, ""),
			expectedResponse: Greeting{
				Greeting: "Hello World!",
			},
		},
		"A populated parameter should return successfully": {
			inputParams: testCreateSayHelloParams(t, "Okteto"),
			expectedResponse: Greeting{
				Greeting: "Hello Okteto!",
			},
		},
		"A populated parameter with spaces should trim and return successfully": {
			inputParams: testCreateSayHelloParams(t, "  Okteto Yo  "),
			expectedResponse: Greeting{
				Greeting: "Hello Okteto Yo!",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			target := newSayHelloHandler(testLoggerContext(tt))

			response := httptest.NewRecorder()
			target.SayHello(response, nil, test.inputParams)

			actualBytes := response.Body.Bytes()
			expectedBytes, err := json.Marshal(test.expectedResponse)
			require.NoError(tt, err)

			assert.Equal(tt, string(expectedBytes), string(actualBytes))
		})
	}
}

func testLoggerContext(tt *testing.T) zerolog.Context {
	tt.Helper()
	return zerolog.New(bytes.NewBuffer(make([]byte, 0))).With()
}

func testCreateSayHelloParams(t *testing.T, name string) SayHelloParams {
	t.Helper()

	return SayHelloParams{
		Name: &name,
	}
}
