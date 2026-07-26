package config

import (
	"fmt"
	"time"

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
	Tax            TaxConfig
	Barkat         BarkatConfig
	Server         HttpServerConfig
	OSWaitInterval time.Duration `env:"KOHAN_OS_WAIT_INTERVAL" envDefault:"1m"`
}

// BarkatConfig defines configuration for the Barkat Journal Explorer
// Database: SQLite file path for journal entries
type BarkatConfig struct {
	DbPath         string `env:"BARKAT_DB_PATH" envDefault:"${HOME}/Downloads/barkat.db" envExpand:"true"`
	ScreenshotPath string `env:"BARKAT_IMAGE_PATH" envDefault:"${HOME}/Downloads/Screenshots" envExpand:"true"`
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
	TaxDir string `env:"TAX_DIR" envDefault:"${HOME}/Downloads/FACompute" envExpand:"true"`

	// Input: Broker statements (Layer 1)
	// Base paths for broker files - year appended at runtime: {base}_{YYYY}.{ext}
	DriveWealthBase string `env:"DRIVEWEALTH_BASE" envDefault:"${HOME}/Downloads/FACompute/Input/Brokerage/vested" envExpand:"true"`
	IBKRBase        string `env:"IBKR_BASE" envDefault:"${HOME}/Downloads/FACompute/Input/Brokerage/ibkr" envExpand:"true"`

	// Input: Parsed data (Layer 2)
	ParsedDir        string `env:"PARSED_DIR" envDefault:"${HOME}/Downloads/FACompute/Input/Parsed" envExpand:"true"`
	TradesPath       string `env:"FA_TRADE_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/Input/Parsed/trades.csv" envExpand:"true"`
	DividendFilePath string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/Input/Parsed/dividends.csv" envExpand:"true"`
	InterestFilePath string `env:"FA_INTEREST_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/Input/Parsed/interest.csv" envExpand:"true"`

	// Data: Reference data (Layer 3)
	TickerCacheDir string `env:"TICKER_CACHE_DIR" envDefault:"${HOME}/Downloads/FACompute/Data/Tickers" envExpand:"true"`
	TTRateFilePath string `env:"TTRATE_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/Data/Reference/sbi_rates.csv" envExpand:"true"`

	// Output: Computed results (Layer 4)
	GainsFilePath string `env:"FA_GAINS_FILE_PATH" envDefault:"${HOME}/Downloads/FACompute/Output/Computed/gains.csv" envExpand:"true"`
	AccountsDir   string `env:"ACCOUNTS_DIR" envDefault:"${HOME}/Downloads/FACompute/Output/YearEndBalance" envExpand:"true"`
	ReportsDir    string `env:"REPORTS_DIR" envDefault:"${HOME}/Downloads/FACompute/Output/Reports" envExpand:"true"`
	ComputedDir   string `env:"COMPUTED_DIR" envDefault:"${HOME}/Downloads/FACompute/Output/Computed" envExpand:"true"`
}

func NewKohanConfig() (config KohanConfig, err error) {
	if err = env.Parse(&config); err != nil {
		return config, fmt.Errorf("error parsing kohan config: %w", err)
	}

	config.Server.Name = "kohan"

	return
}
