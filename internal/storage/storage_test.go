package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestFetchTaxBrackets(t *testing.T) {
	expectedBrackets := []core.TaxBracket{
		{Min: 0, Max: 50000, Rate: 0.1},
		{Min: 50000, Max: 100000, Rate: 0.2},
	}

	tests := []struct {
		name           string
		serverBehavior func(w http.ResponseWriter, r *http.Request)
		expectedError  string
		expectedResult []core.TaxBracket
		contextTimeout time.Duration
		apiURL         string
	}{
		{
			name: "Success",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/tax-calculator/tax-year/2023", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"tax_brackets": expectedBrackets,
				})
			},
			expectedResult: expectedBrackets,
		},
		{
			name: "Bad status code",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "not found", http.StatusNotFound)
			},
			expectedError: "unexpected response status",
		},
		{
			name: "Bad JSON",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("{invalid json"))
			},
			expectedError: "failed to decode response",
		},
		{
			name:          "Request creation error",
			apiURL:        "http://[::1]:NamedPort", // invalid
			expectedError: "failed to create request",
		},
		{
			name: "Timeout",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(6 * time.Second)
			},
			contextTimeout: 1 * time.Second,
			expectedError:  "failed to fetch tax brackets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.serverBehavior != nil {
				server = httptest.NewServer(http.HandlerFunc(tt.serverBehavior))
				defer server.Close()
			}

			apiURL := tt.apiURL
			if apiURL == "" && server != nil {
				apiURL = server.URL
			}

			client := NewTaxAPIClient(apiURL)

			ctx := context.Background()
			if tt.contextTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.contextTimeout)
				defer cancel()
			}

			brackets, err := client.FetchTaxBrackets(ctx, 2023)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, brackets)
			}
		})
	}
}
