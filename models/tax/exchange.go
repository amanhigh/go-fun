package tax

import (
	"strings"
	"time"
)

type SbiRate struct {
	Date   string  `csv:"DATE"`
	TTBuy  float64 `csv:"TT BUY"`
	TTSell float64 `csv:"TT SELL"`
}

// IsValid checks if the rate has all required fields populated
func (r SbiRate) IsValid() bool {
	return r.Date != "" && r.TTBuy != 0 && r.TTSell != 0
}

// BUG: #B Remove Filename Changed to File Path
// File name constant for SBI Rate CSV
const SBI_RATES_FILENAME = "SBI_REFERENCE_RATES_USD.csv"

// BUG: Change To Key?
func (r SbiRate) GetSymbol() string {
	return r.Date
}

// ParseDate implementation for SbiRate
func (r SbiRate) ParseDate() (time.Time, error) {
	// Parse only date part as file has dates in format "2024-01-23 Wednesday"
	datePart := strings.Split(r.Date, " ")[0]
	return time.Parse(time.DateOnly, datePart)
}
