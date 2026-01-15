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

// TaxConfig defines all paths and URLs for tax computation
// Directory structure:
// FACompute/
// ├── Input/
// │   ├── Brokerage/
// │   │   ├── vested_YYYY.xlsx         (DriveWealth exports by year)
// │   │   ├── ibkr_YYYY.csv            (Interactive Brokers by year)
// │   └── Parsed/
// │       ├── trades.csv               (Merged from all brokers)
// │       ├── dividends.csv
// │       └── interest.csv
// ├── Data/
// │   ├── Tickers/                     (Yahoo Finance cache)
// │   └── Reference/
// │       └── sbi_rates.csv            (Exchange rates)
// └── Output/
//
//	├── Computed/
//	│   └── gains.csv                (Capital gains)
//	├── YearEndBalance/
//	│   └── accounts_YYYY.csv        (Year-end positions)
//	├── Reports/
//	│   └── tax_summary_YYYY.xlsx
type TaxConfig struct {
	// External APIs
	SBIBaseURL   string `env:"SBI_BASE_URL" envDefault:"https://raw.githubusercontent.com/sahilgupta/sbi-fx-ratekeeper/main/csv_files/SBI_REFERENCE_RATES_USD.csv"`
	YahooBaseURL string `env:"YAHOO_BASE_URL" envDefault:"https://query1.finance.yahoo.com"`

	// Root directory
	TaxDir string `env:"TAX_DIR" envDefault:"~/Downloads/FACompute"`

	// Input: Broker statements (Layer 1)
	// Base paths for broker files - year appended at runtime: {base}_{YYYY}.{ext}
	DriveWealthBase string `env:"DRIVEWEALTH_BASE" envDefault:"~/Downloads/FACompute/Input/Brokerage/vested"`
	IBKRBase        string `env:"IBKR_BASE" envDefault:"~/Downloads/FACompute/Input/Brokerage/ibkr"`

	// Input: Parsed data (Layer 2)
	ParsedDir        string `env:"PARSED_DIR" envDefault:"~/Downloads/FACompute/Input/Parsed"`
	TradesPath       string `env:"FA_TRADE_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/trades.csv"`
	DividendFilePath string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/dividends.csv"`
	InterestFilePath string `env:"FA_INTEREST_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/interest.csv"`

	// Data: Reference data (Layer 3)
	TickerCacheDir string `env:"TICKER_CACHE_DIR" envDefault:"~/Downloads/FACompute/Data/Tickers"`
	TTRateFilePath string `env:"TTRATE_FILE_PATH" envDefault:"~/Downloads/FACompute/Data/Reference/sbi_rates.csv"`

	// Output: Computed results (Layer 4)
	GainsFilePath string `env:"FA_GAINS_FILE_PATH" envDefault:"~/Downloads/FACompute/Output/Computed/gains.csv"`
	AccountsDir   string `env:"ACCOUNTS_DIR" envDefault:"~/Downloads/FACompute/Output/YearEndBalance"`
	ReportsDir    string `env:"REPORTS_DIR" envDefault:"~/Downloads/FACompute/Output/Reports"`
	ComputedDir   string `env:"COMPUTED_DIR" envDefault:"~/Downloads/FACompute/Output/Computed"`
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
	// Expand home directory (~) in all file paths
	config.Tax.TaxDir = strings.Replace(config.Tax.TaxDir, "~", homeDir, 1)
	config.Tax.DriveWealthBase = strings.Replace(config.Tax.DriveWealthBase, "~", homeDir, 1)
	config.Tax.IBKRBase = strings.Replace(config.Tax.IBKRBase, "~", homeDir, 1)
	config.Tax.TickerCacheDir = strings.Replace(config.Tax.TickerCacheDir, "~", homeDir, 1)
	config.Tax.TTRateFilePath = strings.Replace(config.Tax.TTRateFilePath, "~", homeDir, 1)
	config.Tax.ParsedDir = strings.Replace(config.Tax.ParsedDir, "~", homeDir, 1)
	config.Tax.TradesPath = strings.Replace(config.Tax.TradesPath, "~", homeDir, 1)
	config.Tax.DividendFilePath = strings.Replace(config.Tax.DividendFilePath, "~", homeDir, 1)
	config.Tax.InterestFilePath = strings.Replace(config.Tax.InterestFilePath, "~", homeDir, 1)
	config.Tax.GainsFilePath = strings.Replace(config.Tax.GainsFilePath, "~", homeDir, 1)
	config.Tax.AccountsDir = strings.Replace(config.Tax.AccountsDir, "~", homeDir, 1)
	config.Tax.ReportsDir = strings.Replace(config.Tax.ReportsDir, "~", homeDir, 1)
	config.Tax.ComputedDir = strings.Replace(config.Tax.ComputedDir, "~", homeDir, 1)

	return
}
