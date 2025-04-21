package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/haninamaryia/tax-calculator/internal/core"
)

type TaxCalculator interface {
	CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
}

type TaxCalculatorHandler struct {
	tc TaxCalculator
}

func NewServer(port int, tc TaxCalculator) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/tax", &TaxCalculatorHandler{tc})

	return &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 Mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}

func (t *TaxCalculatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tax" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Parse query params
		incomeStr := r.URL.Query().Get("income")
		yearStr := r.URL.Query().Get("year")

		if incomeStr == "" || yearStr == "" {
			http.Error(w, "Missing required query parameters: 'income' and 'year'", http.StatusBadRequest)
			return
		}

		// Parse income and year from query params
		income, err := strconv.ParseFloat(incomeStr, 64)
		if err != nil || income < 0 {
			http.Error(w, "Invalid income", http.StatusBadRequest)
			return
		}

		// Call the service to calculate tax
		ctx := r.Context()
		result, err := t.tc.CalculateTax(ctx, incomeStr, yearStr)
		if err != nil {
			http.Error(w, "Error calculating tax: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set response content type and encode result
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
