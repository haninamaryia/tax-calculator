package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/haninamaryia/tax-calculator/internal/core"
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
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// Fetch the tax brackets from the API for the specified year
func (t *taxAPIClient) FetchTaxBrackets(ctx context.Context, year int) ([]core.TaxBracket, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", t.baseURL, year)

	// Prepare the request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the API call
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tax brackets: %w", err)
	}
	defer resp.Body.Close()

	// Handle unexpected HTTP status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Decode the JSON response
	var response struct {
		TaxBrackets []core.TaxBracket `json:"tax_brackets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.TaxBrackets, nil
}
