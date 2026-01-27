package tax

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// Account represents year-end account position with original acquisition metadata.
// Preserves FirstPosition for Schedule FA reporting across years.
type Account struct {
	Symbol      string  `csv:"Symbol"`
	Quantity    float64 `csv:"Quantity"`    // Year-end quantity
	MarketValue float64 `csv:"MarketValue"` // Year-end market value
	OriginDate  string  `csv:"OriginDate"`  // ISO format: YYYY-MM-DD (original acquisition date)
	OriginQty   float64 `csv:"OriginQty"`   // Original acquisition quantity
	OriginPrice float64 `csv:"OriginPrice"` // Original cost basis per unit (USD)
}

// GetKey implements CSVRecord interface
func (a Account) GetKey() string {
	return a.Symbol
}

// GetDate implements CSVRecord interface.
// For an Account record from accounts.csv, the date is implicit (year-end).
// This method returns a zero time; the effective date is set by the consuming logic.
func (a Account) GetDate() (time.Time, common.HttpError) {
	// Return zero time and no error, as account snapshot date is context-dependent.
	return time.Time{}, nil
}

// IsValid implements CSVRecord interface
func (a Account) IsValid() bool {
	return a.Symbol != "" && a.Quantity != 0 && a.MarketValue != 0
}

// FromValuations converts a slice of Valuation to a slice of Account.
// Preserves FirstPosition metadata in origin fields for Schedule FA reporting.
func FromValuations(valuations []Valuation) []Account {
	accounts := make([]Account, len(valuations))
	for i, valuation := range valuations {
		// Only include OriginDate if FirstPosition has a valid date
		var originDate string
		if !valuation.FirstPosition.Date.IsZero() {
			originDate = valuation.FirstPosition.Date.Format(time.DateOnly)
		}

		accounts[i] = Account{
			Symbol:      valuation.Ticker,
			Quantity:    valuation.YearEndPosition.Quantity,
			MarketValue: valuation.YearEndPosition.USDValue(),

			// Preserve FirstPosition metadata for carryover
			OriginDate:  originDate,
			OriginQty:   valuation.FirstPosition.Quantity,
			OriginPrice: valuation.FirstPosition.RoundedUSDPrice(),
		}
	}
	return accounts
}
