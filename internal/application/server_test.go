package application_test

import (
	"bytes"
	"encoding/json"
	"go_calculator/internal/application"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculatorHandler(t *testing.T) {
	tc := []struct {
		name         string
		method       string
		body         []byte
		responseCode int
		result       float64
	}{
		{
			name:         "Valid body",
			method:       http.MethodPost,
			body:         []byte(`{"expression": "3 + 2*4"}`),
			responseCode: http.StatusOK,
			result:       11,
		},
		{
			name:         "Invalid expression",
			method:       http.MethodPost,
			body:         []byte(`{"expression": "3 + abc"}`),
			responseCode: http.StatusUnprocessableEntity,
			result:       0,
		},
		{
			name:         "Invalid json",
			method:       http.MethodPost,
			body:         []byte(`{"expression: `),
			responseCode: http.StatusBadRequest,
			result:       0,
		},
		{
			name:         "Invalid body structure",
			method:       http.MethodPost,
			body:         []byte(`{"exp": "3 + 2*4"}`),
			responseCode: http.StatusBadRequest,
			result:       0,
		},
		{
			name:         "Invalid body method",
			method:       http.MethodGet,
			body:         []byte(`{"expression": "3 + 2*4"}`),
			responseCode: http.StatusMethodNotAllowed,
			result:       0,
		},
	}

	for _, tt := range tc {

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewReader(tt.body))
			defer req.Body.Close()
			w := httptest.NewRecorder()

			application.CalculatorHandler(w, req)
			res := w.Result()

			if res.StatusCode != tt.responseCode {
				t.Errorf("expected response code %d, got %d", tt.responseCode, res.StatusCode)
			}

			if res.StatusCode == http.StatusOK {
				response := application.SuccessResponse{}
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if response.Result != tt.result {
					t.Errorf("expected result %f, got %f", tt.result, response.Result)
				}
			}
		})
	}
}
