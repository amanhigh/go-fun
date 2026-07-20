package tax

import "time"

// BrokerageInfo contains all processed data from brokerage reports (DriveWealth, Interactive Brokers, etc.).
type BrokerageInfo struct {
	// CoverageThrough is the latest date covered by this broker result.
	// When merging multiple brokerage results, the later non-zero CoverageThrough is preserved.
	CoverageThrough time.Time
	Interests       []Interest
	Dividends       []Dividend
	Trades          []Trade
}
