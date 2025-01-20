package config

type FAConfig struct {
	SBIBaseURL string `env:"SBI_BASE_URL" envDefault:"https://raw.githubusercontent.com/sahilgupta/sbi-fx-ratekeeper/main/csv_files/SBI_REFERENCE_RATES_USD.csv"`
	// FIXME: #A Add Remaining Config API Key, Downloads Directory. Vantage Base URL
}
