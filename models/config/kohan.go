package config

import (
	"fmt"
	"os"
	"strings"

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
	TaxDir string `env:"TAX_DIR" envDefault:"~/Downloads/FACompute"`
	// TickerCacheDir stores downloaded ticker data, separate from tax input files
	TickerCacheDir   string `env:"TICKER_CACHE_DIR" envDefault:"~/Downloads/FACompute/Tickers"`
	TTRateFilePath   string `env:"TTRATE_FILE_PATH" envDefault:"~/Downloads/FACompute/sbi_rates.csv"`
	TradesPath       string `env:"FA_TRADE_FILE_PATH" envDefault:"~/Downloads/FACompute/trades.csv"`
	DividendFilePath string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"~/Downloads/FACompute/dividends.csv"`
	GainsFilePath    string `env:"FA_GAINS_FILE_PATH" envDefault:"~/Downloads/FACompute/gains.csv"`
	InterestFilePath string `env:"FA_INTEREST_FILE_PATH" envDefault:"~/Downloads/FACompute/interest.csv"`
	DriveWealthPath  string `env:"VESTED_PATH" envDefault:"~/Downloads/FACompute/vested.xlsx"`
}

func NewKohanConfig() (config KohanConfig, err error) {
	if err = env.Parse(&config); err != nil {
		err = fmt.Errorf("error parsing kohan config: %w", err)
		return
	}

	// Expand home directory in file paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// HACK: #C Remove this Hack.
	config.Tax.TaxDir = strings.Replace(config.Tax.TaxDir, "~", homeDir, 1)
	config.Tax.TickerCacheDir = strings.Replace(config.Tax.TickerCacheDir, "~", homeDir, 1)
	config.Tax.TradesPath = strings.Replace(config.Tax.TradesPath, "~", homeDir, 1)
	config.Tax.DividendFilePath = strings.Replace(config.Tax.DividendFilePath, "~", homeDir, 1)
	config.Tax.TTRateFilePath = strings.Replace(config.Tax.TTRateFilePath, "~", homeDir, 1)
	config.Tax.GainsFilePath = strings.Replace(config.Tax.GainsFilePath, "~", homeDir, 1)
	config.Tax.InterestFilePath = strings.Replace(config.Tax.InterestFilePath, "~", homeDir, 1)
	config.Tax.DriveWealthPath = strings.Replace(config.Tax.DriveWealthPath, "~", homeDir, 1)

	return
}
