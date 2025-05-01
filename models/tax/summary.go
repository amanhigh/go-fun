package tax

// Summary contains all processed tax records for a given financial year.
type Summary struct {
	INRGains     []INRGains    // Processed capital gains in INR
	INRDividends []INRDividend // Processed dividends in INR
	INRInterest  []INRInterest // Placeholder for future integration
	// INRPositions  []INRPosition  // Placeholder for future integration
}
