package tax

import (
	"fmt"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/models/common"
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

func (r SbiRate) GetDate() (time.Time, common.HttpError) {
	datePart := strings.Split(r.Date, " ")[0]
	t, err := time.Parse(time.DateOnly, datePart)
	if err != nil {
		return time.Time{}, NewInvalidDateError(fmt.Sprintf("failed to parse date '%s': %v", r.Date, err))
	}
	return t, nil
}
