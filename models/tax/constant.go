package tax

// Tax File Names
const (
	TRADES_FILENAME    = "trades.csv"
	DIVIDENDS_FILENAME = "dividends.csv"
	// File name constant for SBI Rate CSV
	// DATE,PDF FILE,TT BUY,TT SELL,BILL BUY,BILL SELL,FOREX TRAVEL CARD BUY,FOREX TRAVEL CARD SELL,CN BUY,CN SELL
	// 2020-01-04 09:00,https://github.com/sahilgupta/sbi_forex_rates/blob/main/pdf_files/2020/1/2020-01-04.pdf,0.00,0.00,71.29,72.34,70.70,72.55,70.40,72.70
	SBI_RATES_FILENAME = "sbi_rates.csv"
	ACCOUNTS_FILENAME  = "accounts.csv"
	GAINS_FILENAME     = "gains.csv"
	INTEREST_FILENAME  = "interest.csv"

	// Trade Types
	TRADE_TYPE_BUY  = "BUY"
	TRADE_TYPE_SELL = "SELL"

	// Gain Types
	GAIN_TYPE_STCG = "STCG"
	GAIN_TYPE_LTCG = "LTCG"

	// Rounding factor for 2 decimal places
	ROUNDING_FACTOR_2_DECIMALS = 100.0
)
