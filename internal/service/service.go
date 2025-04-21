package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/haninamaryia/tax-calculator/internal/storage"
)

// Interface for the tax calculator service
type TaxService interface {
	CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
}

// Struct implementing the interface
type taxService struct {
	storage storage.TaxStorage
}

// Constructor
func NewTaxService(s storage.TaxStorage) TaxService {
	return &taxService{
		storage: s,
	}
}

// Business logic to calculate tax
func (s *taxService) CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error) {
	// Parse and validate input
	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil || income < 0 {
		return core.TaxResult{}, fmt.Errorf("invalid income")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2019 || year > 2022 {
		return core.TaxResult{}, fmt.Errorf("unsupported year")
	}

	// Fetch tax brackets from storage
	brackets, err := s.storage.FetchTaxBrackets(ctx, year)
	if err != nil {
		return core.TaxResult{}, fmt.Errorf("failed to fetch tax brackets: %w", err)
	}

	perBand := make(map[string]float64)
	var totalTax float64

	for _, b := range brackets {
		if income <= b.Min {
			break
		}

		upper := b.Max
		if upper == 0 || income < upper {
			upper = income
		}

		taxable := upper - b.Min
		if taxable < 0 {
			taxable = 0
		}

		tax := taxable * b.Rate
		totalTax += tax
		perBand[fmt.Sprintf("%.2f-%.2f", b.Min, upper)] = tax
	}

	effectiveRate := 0.0
	if income > 0 {
		effectiveRate = totalTax / income
	}

	return core.TaxResult{
		TotalTax:      totalTax,
		PerBracket:    perBand,
		EffectiveRate: effectiveRate,
	}, nil
}
