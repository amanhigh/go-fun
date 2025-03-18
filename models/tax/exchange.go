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

// FIXME: #B Remove Filename Changed to File Path and Inject joins filename and DownloadPath.
// File name constant for SBI Rate CSV
// DATE,PDF FILE,TT BUY,TT SELL,BILL BUY,BILL SELL,FOREX TRAVEL CARD BUY,FOREX TRAVEL CARD SELL,CN BUY,CN SELL
// 2020-01-04 09:00,https://github.com/sahilgupta/sbi_forex_rates/blob/main/pdf_files/2020/1/2020-01-04.pdf,0.00,0.00,71.29,72.34,70.70,72.55,70.40,72.70
const SBI_RATES_FILENAME = "SBI_REFERENCE_RATES_USD.csv"

func (r SbiRate) GetKey() string {
	return r.Date
}

func (r SbiRate) GetDate() time.Time {
	datePart := strings.Split(r.Date, " ")[0]
	t, _ := time.Parse(time.DateOnly, datePart)
	return t
}
