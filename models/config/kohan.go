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
	AlphaBaseURL string `env:"ALPHA_BASE_URL" envDefault:"https://www.alphavantage.co/query"`
	AlphaAPIKey  string `env:"ALPHA_API_KEY" envDefault:"DUMMY_KEY_FOR_E2E"`

	// File System Configuration
	// TickerInfoDir stores downloaded ticker data, separate from tax input files.
	TickerInfoDir     string `env:"TICKER_INFO_DIR" envDefault:"~/Downloads/FACompute/Tickers"`
	TTRateFilePath    string `env:"TTRATE_FILE_PATH" envDefault:"~/Downloads/FACompute/sbi_rates.csv"`
	YearlySummaryPath string `env:"YEARLY_SUMMARY_PATH" envDefault:"~/Downloads/FACompute/tax_summary.xlsx"`

	AccountDir       string `env:"ACCOUNT_DIR" envDefault:"~/Downloads/FACompute"`
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
	config.Tax.TickerInfoDir = strings.Replace(config.Tax.TickerInfoDir, "~", homeDir, 1)
	config.Tax.TradesPath = strings.Replace(config.Tax.TradesPath, "~", homeDir, 1)
	config.Tax.DividendFilePath = strings.Replace(config.Tax.DividendFilePath, "~", homeDir, 1)
	config.Tax.TTRateFilePath = strings.Replace(config.Tax.TTRateFilePath, "~", homeDir, 1)
	config.Tax.AccountDir = strings.Replace(config.Tax.AccountDir, "~", homeDir, 1)
	config.Tax.GainsFilePath = strings.Replace(config.Tax.GainsFilePath, "~", homeDir, 1)
	config.Tax.InterestFilePath = strings.Replace(config.Tax.InterestFilePath, "~", homeDir, 1)
	config.Tax.YearlySummaryPath = strings.Replace(config.Tax.YearlySummaryPath, "~", homeDir, 1)
	config.Tax.DriveWealthPath = strings.Replace(config.Tax.DriveWealthPath, "~", homeDir, 1)

	return
}
