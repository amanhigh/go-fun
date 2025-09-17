package tax

// DriveWealthInfo contains all processed data from a DriveWealth report.
type DriveWealthInfo struct {
	Interests []Interest
	Dividends []Dividend
	Trades    []Trade
}
