package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/haninamaryia/tax-calculator/internal/core"
	"github.com/haninamaryia/tax-calculator/internal/logger"
	"github.com/haninamaryia/tax-calculator/internal/storage"
)

const SupportedYears = "2019, 2020, 2021, 2022"

// Interface for the tax calculator service
type TaxService interface {
	CalculateTax(ctx context.Context, incomeStr string, yearStr string) (core.TaxResult, error)
	ValidateTaxYear(year string) error
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

	// TODO: refactor this to be less redundant
	// Validate the year first
	if err := s.ValidateTaxYear(yearStr); err != nil {
		logger.Log.Error().Err(err).Msgf("Invalid tax year: %s", yearStr)
		return core.TaxResult{}, err
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("Error parsing year: %s", yearStr)
		return core.TaxResult{}, fmt.Errorf("error parsing year")
	}

	// Parse and validate input income
	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil || income < 0 {
		logger.Log.Error().Err(err).Msgf("Invalid income: %s", incomeStr)
		return core.TaxResult{}, fmt.Errorf("invalid income")
	}

	// Fetch tax brackets from storage
	brackets, err := s.storage.FetchTaxBrackets(ctx, year)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch tax brackets")
		return core.TaxResult{}, fmt.Errorf("failed to fetch tax brackets: %w", err)
	}

	perBand := make(map[string]float64)
	var totalTax float64

	// TODO: put this in function and create unit test
	// Calculate tax for each bracket
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

	logger.Log.Info().Msgf("Calculated tax: %.2f for income: %.2f, year: %s", totalTax, income, yearStr)

	return core.TaxResult{
		TotalTax:      totalTax,
		PerBracket:    perBand,
		EffectiveRate: effectiveRate,
	}, nil
}

func (s *taxService) ValidateTaxYear(year string) error {
	supportedYears := map[string]bool{"2019": true, "2020": true, "2021": true, "2022": true}

	if !supportedYears[year] {
		logger.Log.Warn().Msgf("Unsupported tax year: %s", year)
		return fmt.Errorf("tax year %s is not supported", year)
	}

	logger.Log.Info().Msgf("Valid tax year: %s", year)
	return nil
}
