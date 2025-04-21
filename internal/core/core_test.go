package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaxBracket_JSON(t *testing.T) {
	tests := []struct {
		name    string
		bracket TaxBracket
	}{
		{
			name: "Standard bracket",
			bracket: TaxBracket{
				Min:  0,
				Max:  50000,
				Rate: 0.15,
			},
		},
		{
			name: "Zero rate bracket",
			bracket: TaxBracket{
				Min:  100000,
				Max:  200000,
				Rate: 0.0,
			},
		},
		{
			name: "Max is zero (open-ended)",
			bracket: TaxBracket{
				Min:  200000,
				Max:  0,
				Rate: 0.33,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.bracket)
			assert.NoError(t, err)

			var decoded TaxBracket
			err = json.Unmarshal(data, &decoded)
			assert.NoError(t, err)

			assert.Equal(t, tt.bracket.Min, decoded.Min)
			assert.Equal(t, tt.bracket.Max, decoded.Max)
			assert.Equal(t, tt.bracket.Rate, decoded.Rate)
		})
	}
}

func TestTaxResult_Fields(t *testing.T) {
	tests := []struct {
		name         string
		result       TaxResult
		expectedRate float64
	}{
		{
			name: "Standard case",
			result: TaxResult{
				TotalTax: 17739.17,
				PerBracket: map[string]float64{
					"0-50000":      7500,
					"50001-100000": 10239.17,
				},
				EffectiveRate: 0.1774,
			},
			expectedRate: 0.1774,
		},
		{
			name: "Zero tax case",
			result: TaxResult{
				TotalTax:      0,
				PerBracket:    map[string]float64{},
				EffectiveRate: 0,
			},
			expectedRate: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result.TotalTax, tt.result.TotalTax)
			assert.Equal(t, tt.result.PerBracket, tt.result.PerBracket)

			if tt.expectedRate == 0 {
				assert.Equal(t, 0.0, tt.result.EffectiveRate)
			} else {
				assert.InEpsilon(t, tt.expectedRate, tt.result.EffectiveRate, 0.0001)
			}
		})
	}

}

func TestTaxResult_JSON(t *testing.T) {
	tests := []struct {
		name   string
		result TaxResult
	}{
		{
			name: "Two brackets",
			result: TaxResult{
				TotalTax: 10000,
				PerBracket: map[string]float64{
					"0-50000": 7500,
					"50001+":  2500,
				},
				EffectiveRate: 0.2,
			},
		},
		{
			name: "Empty brackets",
			result: TaxResult{
				TotalTax:      0,
				PerBracket:    map[string]float64{},
				EffectiveRate: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.result)
			assert.NoError(t, err)

			var decoded TaxResult
			err = json.Unmarshal(data, &decoded)
			assert.NoError(t, err)

			assert.Equal(t, tt.result.TotalTax, decoded.TotalTax)
			assert.Equal(t, tt.result.PerBracket, decoded.PerBracket)
			assert.Equal(t, tt.result.EffectiveRate, decoded.EffectiveRate)
		})
	}
}
