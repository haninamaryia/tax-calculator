package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/haninamaryia/tax-calculator/internal/service"
	"github.com/stretchr/testify/assert"
)

// Mock storage
type mockStorage struct {
	brackets []core.TaxBracket
	err      error
}

func (m *mockStorage) FetchTaxBrackets(ctx context.Context, year int) ([]core.TaxBracket, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.brackets, nil
}

func TestCalculateTax(t *testing.T) {

	tests := []struct {
		name        string
		incomeStr   string
		yearStr     string
		brackets    []core.TaxBracket
		mockErr     error
		expectErr   bool
		expectTotal float64
	}{
		{
			name:      "valid income and year",
			incomeStr: "60000",
			yearStr:   "2021",
			brackets: []core.TaxBracket{
				{Min: 0, Max: 10000, Rate: 0.1},
				{Min: 10000, Max: 50000, Rate: 0.2},
				{Min: 50000, Max: 0, Rate: 0.3}, // Max 0 means no upper limit
			},
			expectErr:   false,
			expectTotal: 1000 + 8000 + 3000, // 12000
		},
		{
			name:      "invalid income format",
			incomeStr: "abc",
			yearStr:   "2021",
			expectErr: true,
		},
		{
			name:      "invalid year format",
			incomeStr: "50000",
			yearStr:   "20xx",
			expectErr: true,
		},
		{
			name:      "unsupported year",
			incomeStr: "50000",
			yearStr:   "2010",
			expectErr: true,
		},
		{
			name:      "error from storage",
			incomeStr: "50000",
			yearStr:   "2021",
			mockErr:   errors.New("failed to fetch"),
			expectErr: true,
		},
		{
			name:      "zero income",
			incomeStr: "0",
			yearStr:   "2021",
			brackets: []core.TaxBracket{
				{Min: 0, Max: 10000, Rate: 0.1},
			},
			expectErr:   false,
			expectTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStorage{
				brackets: tt.brackets,
				err:      tt.mockErr,
			}

			svc := service.NewTaxService(mock)
			result, err := svc.CalculateTax(context.Background(), tt.incomeStr, tt.yearStr)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.TotalTax != tt.expectTotal {
				t.Errorf("expected total tax %v, got %v", tt.expectTotal, result.TotalTax)
			}
		})
	}
}

func TestValidateTaxYear(t *testing.T) {

	tests := []struct {
		name          string
		year          string
		expectedError string
	}{
		{
			name:          "Valid tax year 2019",
			year:          "2019",
			expectedError: "",
		},
		{
			name:          "Valid tax year 2020",
			year:          "2020",
			expectedError: "",
		},
		{
			name:          "Valid tax year 2021",
			year:          "2021",
			expectedError: "",
		},
		{
			name:          "Valid tax year 2022",
			year:          "2022",
			expectedError: "",
		},
		{
			name:          "Invalid tax year 2023",
			year:          "2023",
			expectedError: "tax year 2023 is not supported",
		},
		{
			name:          "Invalid tax year 2025",
			year:          "2025",
			expectedError: "tax year 2025 is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStorage{}

			svc := service.NewTaxService(mock)
			err := svc.ValidateTaxYear(tt.year)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
