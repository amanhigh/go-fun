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

	// Yahoo Finance Configuration
	YahooBaseURL string `env:"YAHOO_BASE_URL" envDefault:"https://query1.finance.yahoo.com"`

	// File System Configuration
	// TaxDir is the base directory for all tax-related files and subdirectories
	TaxDir string `env:"TAX_DIR" envDefault:"~/Downloads/FACompute"`

	// Data Layer (Layer 2: External Reference Data - Immutable)
	// TickerCacheDir stores downloaded ticker data from Yahoo Finance
	TickerCacheDir string `env:"TICKER_CACHE_DIR" envDefault:"~/Downloads/FACompute/Data/Tickers"`
	// TTRateFilePath stores SBI exchange rates (USD-INR)
	TTRateFilePath string `env:"TTRATE_FILE_PATH" envDefault:"~/Downloads/FACompute/Data/Reference/sbi_rates.csv"`

	// Input Layer (Layer 1: User-Provided Data - Immutable)
	// Input/Brokerage/ contains vested.xlsx export from DriveWealth
	DriveWealthPath string `env:"VESTED_PATH" envDefault:"~/Downloads/FACompute/Input/Brokerage/vested.xlsx"`
	// Input/Parsed/ contains CSV files parsed from vested.xlsx by 'tax vested parse' command
	TradesPath       string `env:"FA_TRADE_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/trades.csv"`
	DividendFilePath string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/dividends.csv"`
	InterestFilePath string `env:"FA_INTEREST_FILE_PATH" envDefault:"~/Downloads/FACompute/Input/Parsed/interest.csv"`

	// Output Layer (Layer 3: System-Generated Results - Mutable)
	// Output/Computed/ contains gains.csv from capital gains calculation
	GainsFilePath string `env:"FA_GAINS_FILE_PATH" envDefault:"~/Downloads/FACompute/Output/Computed/gains.csv"`
	// Output/YearEndBalance/ contains accounts_YYYY.csv computed at year-end
	AccountsDir string `env:"ACCOUNTS_DIR" envDefault:"~/Downloads/FACompute/Output/YearEndBalance"`
	// Output/Reports/ contains tax_summary_YYYY.xlsx generated for ITR filing
	ReportsDir string `env:"REPORTS_DIR" envDefault:"~/Downloads/FACompute/Output/Reports"`
	// ComputedDir stores gains.csv and other computed results
	ComputedDir string `env:"COMPUTED_DIR" envDefault:"~/Downloads/FACompute/Output/Computed"`
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
	config.Tax.TickerCacheDir = strings.Replace(config.Tax.TickerCacheDir, "~", homeDir, 1)
	config.Tax.TradesPath = strings.Replace(config.Tax.TradesPath, "~", homeDir, 1)
	config.Tax.DividendFilePath = strings.Replace(config.Tax.DividendFilePath, "~", homeDir, 1)
	config.Tax.TTRateFilePath = strings.Replace(config.Tax.TTRateFilePath, "~", homeDir, 1)
	config.Tax.GainsFilePath = strings.Replace(config.Tax.GainsFilePath, "~", homeDir, 1)
	config.Tax.InterestFilePath = strings.Replace(config.Tax.InterestFilePath, "~", homeDir, 1)
	config.Tax.DriveWealthPath = strings.Replace(config.Tax.DriveWealthPath, "~", homeDir, 1)
	config.Tax.AccountsDir = strings.Replace(config.Tax.AccountsDir, "~", homeDir, 1)
	config.Tax.ReportsDir = strings.Replace(config.Tax.ReportsDir, "~", homeDir, 1)
	config.Tax.ComputedDir = strings.Replace(config.Tax.ComputedDir, "~", homeDir, 1)

	return
}
