package tax

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// Account represents year-end account position
type Account struct {
	Symbol      string  `csv:"Symbol"`
	Quantity    float64 `csv:"Quantity"`
	Cost        float64 `csv:"Cost"`
	MarketValue float64 `csv:"MarketValue"`
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
	return a.Symbol != "" && a.Quantity != 0 && a.Cost != 0 && a.MarketValue != 0
}

// FromValuations converts a slice of Valuation to a slice of Account.
func FromValuations(valuations []Valuation) []Account {
	accounts := make([]Account, len(valuations))
	for i, valuation := range valuations {
		accounts[i] = Account{
			Symbol:      valuation.Ticker,
			Quantity:    valuation.YearEndPosition.Quantity,
			Cost:        valuation.YearEndPosition.USDValue(),
			MarketValue: valuation.YearEndPosition.USDValue(),
		}
	}
	return accounts
}
