package tax

// Summary contains all processed tax records for a given financial year.
type Summary struct {
	INRGains      []INRGains    // Processed capital gains in INR
	INRDividends  []INRDividend // Processed dividends in INR
	INRInterest   []INRInterest
	INRValuations []INRValuation
	// FIXME: Repurpose Valuations or Add Format for accounts.csv
}
