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
	Tax    TaxConfig
	Barkat BarkatConfig
}

// BarkatConfig defines configuration for the Barkat Journal Explorer
// Database: SQLite file path for journal entries
type BarkatConfig struct {
	DbPath         string `env:"BARKAT_DB_PATH" envDefault:"~/Downloads/barkat.db"`
	ScreenshotPath string `env:"BARKAT_IMAGE_PATH" envDefault:"~/Downloads/Screenshots"`
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

	// Ticker data start year (avoids sparse/missing data from very old periods)
	TickerDataStartYear int `env:"TICKER_DATA_START_YEAR" envDefault:"2020"`

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
		return config, fmt.Errorf("error parsing kohan config: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("failed to get user home directory: %w", err)
	}

	applyTaxPaths(&config.Tax, homeDir)
	config.Barkat.DbPath = replaceHome(config.Barkat.DbPath, homeDir)
	config.Barkat.ScreenshotPath = replaceHome(config.Barkat.ScreenshotPath, homeDir)

	return
}

func applyTaxPaths(tax *TaxConfig, homeDir string) {
	// TODO: #C Remove this Hack.
	tax.TaxDir = replaceHome(tax.TaxDir, homeDir)
	tax.DriveWealthBase = replaceHome(tax.DriveWealthBase, homeDir)
	tax.IBKRBase = replaceHome(tax.IBKRBase, homeDir)
	tax.TickerCacheDir = replaceHome(tax.TickerCacheDir, homeDir)
	tax.TTRateFilePath = replaceHome(tax.TTRateFilePath, homeDir)
	tax.ParsedDir = replaceHome(tax.ParsedDir, homeDir)
	tax.TradesPath = replaceHome(tax.TradesPath, homeDir)
	tax.DividendFilePath = replaceHome(tax.DividendFilePath, homeDir)
	tax.InterestFilePath = replaceHome(tax.InterestFilePath, homeDir)
	tax.GainsFilePath = replaceHome(tax.GainsFilePath, homeDir)
	tax.AccountsDir = replaceHome(tax.AccountsDir, homeDir)
	tax.ReportsDir = replaceHome(tax.ReportsDir, homeDir)
	tax.ComputedDir = replaceHome(tax.ComputedDir, homeDir)
}

func replaceHome(path, homeDir string) string {
	return strings.Replace(path, "~", homeDir, 1)
}
