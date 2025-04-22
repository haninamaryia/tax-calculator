package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/haninamaryia/tax-calculator/internal/logger"
)

type TaxCalculator interface {
	CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
}

type TaxCalculatorHandler struct {
	tc TaxCalculator
}

// HealthCheckHandler will handle requests to the health check endpoint.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Respond with a 200 OK status code
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewServer(port int, tc TaxCalculator) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", HealthCheckHandler)
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
		logger.Log.Warn().Msgf("Invalid URL path: %s", r.URL.Path) // Log invalid path
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var request struct {
			Income interface{} `json:"income"`
			Year   int         `json:"year"`
		}

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Log.Error().Err(err).Msg("Invalid JSON body") // Log error when JSON is invalid
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Check for missing required fields
		if request.Income == nil || request.Year == 0 {
			logger.Log.Warn().Msg("Missing required fields in request") // Log missing fields
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Check for invalid income type
		switch v := request.Income.(type) {
		case string:
			logger.Log.Warn().Msg("Invalid income type: expected number but received string") // Log invalid income type
			http.Error(w, "Invalid income", http.StatusBadRequest)
			return
		case float64:
			// Valid income format
			if v < 0 {
				logger.Log.Warn().Msgf("Invalid income value: %f, income must be non-negative", v) // Log invalid income value
				http.Error(w, "Income must be non-negative", http.StatusBadRequest)
				return
			}
		default:
			logger.Log.Warn().Msg("Invalid income type") // Log invalid income type
			http.Error(w, "Invalid income", http.StatusBadRequest)
			return
		}

		// Convert income and year to strings as before
		incomeStr := fmt.Sprintf("%.2f", request.Income)
		yearStr := fmt.Sprintf("%d", request.Year)

		// Call the service
		result, err := t.tc.CalculateTax(r.Context(), incomeStr, yearStr)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error calculating tax") // Log error calculating tax
			http.Error(w, "Error calculating tax: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to encode response") // Log error encoding response
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		logger.Log.Warn().Msgf("Method %s not allowed for %s", r.Method, r.URL.Path) // Log method not allowed
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
