package core

const DateFormat = "2006-01-02"

type (
	TaxBracket struct {
		Min  float64 `json:"min"`
		Max  float64 `json:"max,omitempty"`
		Rate float64 `json:"rate"`
	}

	TaxResult struct {
		TotalTax      float64            `json:"total_tax"`
		PerBracket    map[string]float64 `json:"per_bracket"`
		EffectiveRate float64            `json:"effective_rate"`
	}
)
