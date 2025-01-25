package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

/* PSSH */
const CLUSTER_PATH = "/tmp/clusters"
const OUTPUT_PATH = "/tmp/output"
const ERROR_PATH = "/tmp/error"
const CONSOLE_FILE = CLUSTER_PATH + "/console.txt"

const DEFAULT_PARALELISM = 50

const DEBUG_FILE = "/tmp/kohandebug"

var KOHAN_DEBUG = false

type DariusConfig struct {
	MakeDir             string
	SelectedServiceFile string
}

type KohanConfig struct {
	Tax TaxConfig
}

type TaxConfig struct {
	// SBI Configuration
	SBIBaseURL string `env:"SBI_BASE_URL" envDefault:"https://raw.githubusercontent.com/sahilgupta/sbi-fx-ratekeeper/main/csv_files/SBI_REFERENCE_RATES_USD.csv"`

	// Alpha Vantage Configuration
	AlphaBaseURL string `env:"ALPHA_BASE_URL" envDefault:"https://www.alphavantage.co/query"`
	AlphaAPIKey  string `env:"ALPHA_API_KEY"` // required, no default

	// File System Configuration
	DownloadsDir        string `env:"FA_DOWNLOADS_DIR" envDefault:"~/Downloads/Tickers"`
	BrokerStatementFile string `env:"FA_BROKER_STATEMENT" envDefault:"broker_statement.csv"`
	DividendFile        string `env:"FA_DIVIDEND_FILE" envDefault:"dividends.csv"`
}

func NewKohanConfig() (config KohanConfig, err error) {
	if err = env.Parse(&config); err != nil {
		err = fmt.Errorf("error parsing kohan config: %w", err)
	}
	return
}
