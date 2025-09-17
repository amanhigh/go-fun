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
	AlphaBaseURL string `env:"ALPHAVANTAGE_BASE_URL" envDefault:"https://www.alphavantage.co/query"`
	AlphaAPIKey  string `env:"ALPHAVANTAGE_API_KEY" envDefault:"DUMMY_KEY_FOR_E2E"`

	// File System Configuration
	// TaxDir is the base directory for all tax-related files and subdirectories
	TaxDir string `env:"TAX_DIR" envDefault:"${HOME}/Downloads/FACompute" envExpand:"true"`
	// TickerCacheDir stores downloaded ticker data, separate from tax input files
	TickerCacheDir   string `env:"TICKER_CACHE_DIR" envDefault:"${HOME}/Downloads/FACompute/Tickers" envExpand:"true"`
	TTRateFilePath   string `env:"TTRATE_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/sbi_rates.csv" envExpand:"true"`
	TradesPath       string `env:"FA_TRADE_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/trades.csv" envExpand:"true"`
	DividendFilePath string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/dividends.csv" envExpand:"true"`
	GainsFilePath    string `env:"FA_GAINS_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/gains.csv" envExpand:"true"`
	InterestFilePath string `env:"FA_INTEREST_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/interest.csv" envExpand:"true"`
	DriveWealthPath  string `env:"VESTED_PATH" envDefault:"${HOME}/Downloads/FACompute/vested.xlsx" envExpand:"true"`
}

func NewKohanConfig() (config KohanConfig, err error) {
	if err = env.Parse(&config); err != nil {
		err = fmt.Errorf("error parsing kohan config: %w", err)
		return
	}

	return
}
