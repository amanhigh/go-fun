package tax

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// FIXME: #A Create Test Data for Integration Test for all CSV Models
type Gains struct {
	Symbol     string  `csv:"Symbol"`
	BuyDate    string  `csv:"BuyDate"`
	SellDate   string  `csv:"SellDate"`
	Quantity   float64 `csv:"Quantity"`
	PNL        float64 `csv:"PNL"`
	Commission float64 `csv:"Commission"`
	Type       string  `csv:"Type"`
}

func (g Gains) GetKey() string {
	return g.Symbol
}

func (g Gains) IsValid() bool {
	return g.Symbol != "" && g.BuyDate != "" && g.SellDate != ""
}

func (g Gains) GetDate() (time.Time, common.HttpError) {
	t, err := time.Parse(time.DateOnly, g.SellDate)
	if err != nil {
		return time.Time{}, NewInvalidDateError(fmt.Sprintf("failed to parse sell date '%s': %v", g.SellDate, err))
	}
	return t, nil
}

func (g Gains) ParseBuyDate() (time.Time, error) {
	return time.Parse(time.DateOnly, g.BuyDate)
}

func (g Gains) ParseSellDate() (time.Time, error) {
	return time.Parse(time.DateOnly, g.SellDate)
}

// FIXME: #A Create TaxSummary model and Wire up TaxManager.
// INRGains adds exchange rate details to basic gains
type INRGains struct {
	Gains            // Embed original gains
	TTDate time.Time // Sell date for exchange rate
	TTRate float64   // Exchange rate on sell date
}

// Implement Exchangeable interface
func (g *INRGains) GetDate() (time.Time, common.HttpError) {
	// Call the embedded Gains's GetDate method to avoid infinite recursion
	return g.Gains.GetDate()
}

func (g *INRGains) GetUSDAmount() float64 {
	return g.PNL
}

func (g *INRGains) SetTTRate(rate float64) {
	g.TTRate = rate
}

func (g *INRGains) SetTTDate(date time.Time) {
	g.TTDate = date
}

// INRValue computes the PNL value in INR
func (g *INRGains) INRValue() float64 {
	return g.PNL * g.TTRate
}
