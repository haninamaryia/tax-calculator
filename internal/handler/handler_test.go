package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/stretchr/testify/assert"
)

// mockTaxCalculator is a test double
type mockTaxCalculator struct {
	CalculateTaxFunc func(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
}

func (m *mockTaxCalculator) CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error) {
	return m.CalculateTaxFunc(ctx, incomeStr, yearStr)
}

func TestTaxHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           map[string]interface{}
		mockFunc       func(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
		expectedCode   int
		expectedBody   string
		expectedHeader map[string]string
		validateJSON   bool
		expectedJSON   core.TaxResult
	}{
		{
			name:   "Success case",
			method: "POST",
			body:   map[string]interface{}{"income": 10000.0, "year": 2022},
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{
					TotalTax:      1234.56,
					EffectiveRate: 0.123,
					PerBracket: map[string]float64{
						"0-50000": 1000,
						"50000+":  234.56,
					},
				}, nil
			},
			expectedCode: http.StatusOK,
			validateJSON: true,
			expectedJSON: core.TaxResult{
				TotalTax:      1234.56,
				EffectiveRate: 0.123,
				PerBracket: map[string]float64{
					"0-50000": 1000,
					"50000+":  234.56,
				},
			},
		},
		{
			name:         "Missing body parameters",
			method:       "POST",
			body:         map[string]interface{}{}, // Missing income and year
			expectedCode: http.StatusBadRequest,
			expectedBody: "Missing required fields",
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{}, nil
			},
		},

		{
			name:         "Invalid income",
			method:       "POST",
			body:         map[string]interface{}{"income": "abc", "year": 2022}, // Invalid income format
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid income",
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{}, nil
			},
		},
		{
			name:   "Internal error from service",
			method: "POST",
			body:   map[string]interface{}{"income": 10000.0, "year": 2022},
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{}, errors.New("internal error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error calculating tax",
		},
		{
			name:           "OPTIONS method",
			method:         "OPTIONS",
			expectedCode:   http.StatusNoContent,
			expectedHeader: map[string]string{"Allow": "POST, OPTIONS"},
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{}, nil
			},
		},
		{
			name:           "Method not allowed",
			method:         "GET",
			expectedCode:   http.StatusMethodNotAllowed,
			expectedHeader: map[string]string{"Allow": "POST, OPTIONS"},
			mockFunc: func(ctx context.Context, incomeStr, yearStr string) (core.TaxResult, error) {
				return core.TaxResult{}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTaxCalculator{CalculateTaxFunc: tt.mockFunc}
			handler := &TaxCalculatorHandler{tc: mock}

			var reqBody io.Reader
			if tt.body != nil {
				b, _ := json.Marshal(tt.body)
				reqBody = bytes.NewReader(b)
			}

			req := httptest.NewRequest(tt.method, "/tax", reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			for k, v := range tt.expectedHeader {
				assert.Equal(t, v, w.Header().Get(k))
			}

			if tt.validateJSON {
				var result core.TaxResult
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedJSON.TotalTax, result.TotalTax)
				assert.Equal(t, tt.expectedJSON.EffectiveRate, result.EffectiveRate)
				assert.Equal(t, tt.expectedJSON.PerBracket, result.PerBracket)
			} else if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
