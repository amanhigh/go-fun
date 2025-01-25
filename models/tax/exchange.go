package tax

type SbiRate struct {
	Date   string  `csv:"DATE"`
	TTBuy  float64 `csv:"TT BUY"`
	TTSell float64 `csv:"TT SELL"`
}

// File name constant for SBI Rate CSV
const SBI_RATES_FILENAME = "SBI_REFERENCE_RATES_USD.csv"

// BUG: #C Add Constant for other Files
