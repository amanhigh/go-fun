package tax

// SecurityInfo represents basic security information from Yahoo Finance.
// This is the minimal model for Phase 1 covering symbol, name, exchange, and type.
type SecurityInfo struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Type     string `json:"type"`
}
