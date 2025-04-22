package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/haninamaryia/tax-calculator/internal/logger"
)

type TaxStorage interface {
	FetchTaxBrackets(ctx context.Context, year int) ([]core.TaxBracket, error)
}

type taxAPIClient struct {
	baseURL string
	client  *http.Client
}

// Constructor to initialize the taxAPIClient
func NewTaxAPIClient(baseURL string) TaxStorage {
	return &taxAPIClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second}, // Increased timeout for robustness
	}
}

// Fetch the tax brackets from the API for the specified year
func (t *taxAPIClient) FetchTaxBrackets(ctx context.Context, year int) ([]core.TaxBracket, error) {

	//TODO: put this in config
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", t.baseURL, year)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to make HTTP request")
		return nil, fmt.Errorf("failed to fetch tax brackets: %w", err)
	}
	defer resp.Body.Close()

	logger.Log.Info().Msgf("Received response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log.Warn().Msgf("Unexpected status code %d, response body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected response status: %s. Response body: %s", resp.Status, string(body))
	}

	var response struct {
		TaxBrackets []core.TaxBracket `json:"tax_brackets"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to decode response body")
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(response.TaxBrackets) == 0 {
		logger.Log.Warn().Msgf("Missing or invalid tax brackets in response: %s", string(body))
		return nil, fmt.Errorf("missing or invalid tax brackets in the response: %s", string(body))
	}

	logger.Log.Info().Msgf("Fetched tax brackets for year %d successfully", year)
	return response.TaxBrackets, nil
}
