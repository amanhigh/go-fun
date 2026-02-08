package tax

// BrokerageInfo contains all processed data from brokerage reports (DriveWealth, Interactive Brokers, etc.).
type BrokerageInfo struct {
	Interests []Interest
	Dividends []Dividend
	Trades    []Trade
}
