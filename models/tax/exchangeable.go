package tax

import "time"

// Exchangeable represents any entity that needs USD to INR exchange rate conversion
type Exchangeable interface {
	// GetDate returns the original transaction date for exchange lookup
	GetDate() time.Time

	// GetUSDAmount returns the USD amount to be converted
	GetUSDAmount() float64

	// SetTTRate stores the exchange rate used for conversion
	SetTTRate(float64)

	// SetTTDate stores the date for which exchange rate was used
	SetTTDate(time.Time)
}
