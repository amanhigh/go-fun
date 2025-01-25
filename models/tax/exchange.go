package tax

// BUG: #A Use CSV Tag
type SbiRate struct {
	Date   string  `json:"DATE"`
	TTBuy  float64 `json:"TT BUY"`
	TTSell float64 `json:"TT SELL"`
}

// File name constant for SBI Rate CSV
const SBI_RATES_FILENAME = "SBI_REFERENCE_RATES_USD.csv"

// BUG: #C Add Constant for other Files
