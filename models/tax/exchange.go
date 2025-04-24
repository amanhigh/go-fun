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

func (r SbiRate) GetKey() string {
	return r.Date
}

func (r SbiRate) GetDate() time.Time {
	datePart := strings.Split(r.Date, " ")[0]
	t, _ := time.Parse(time.DateOnly, datePart)
	return t
}
