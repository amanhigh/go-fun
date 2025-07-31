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
	// DownloadsDir stores downloaded ticker data, separate from tax input files.
	DownloadsDir        string `env:"FA_DOWNLOADS_DIR" envDefault:"~/Downloads/FACompute"`
	BrokerStatementPath string `env:"FA_BROKER_STATEMENT_PATH" envDefault:"~/Downloads/FACompute/trades.csv"`
	DividendFilePath    string `env:"FA_DIVIDEND_FILE_PATH" envDefault:"~/Downloads/FACompute/dividends.csv"`
	SBIFilePath         string `env:"SBI_FILE_PATH" envDefault:"~/Downloads/FACompute/sbi_rates.csv"`
	AccountFilePath     string `env:"ACCOUNT_FILE_PATH" envDefault:"~/Downloads/FACompute/accounts.csv"`
	GainsFilePath       string `env:"FA_GAINS_FILE_PATH" envDefault:"~/Downloads/FACompute/gains.csv"`
	InterestFilePath    string `env:"FA_INTEREST_FILE_PATH" envDefault:"~/Downloads/FACompute/interest.csv"`
	YearlySummaryPath   string `env:"YEARLY_SUMMARY_PATH" envDefault:"~/Downloads/FACompute/tax_summary.xlsx"`
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

	config.Tax.DownloadsDir = strings.Replace(config.Tax.DownloadsDir, "~", homeDir, 1)
	config.Tax.BrokerStatementPath = strings.Replace(config.Tax.BrokerStatementPath, "~", homeDir, 1)
	config.Tax.DividendFilePath = strings.Replace(config.Tax.DividendFilePath, "~", homeDir, 1)
	config.Tax.SBIFilePath = strings.Replace(config.Tax.SBIFilePath, "~", homeDir, 1)
	config.Tax.AccountFilePath = strings.Replace(config.Tax.AccountFilePath, "~", homeDir, 1)
	config.Tax.GainsFilePath = strings.Replace(config.Tax.GainsFilePath, "~", homeDir, 1)
	config.Tax.InterestFilePath = strings.Replace(config.Tax.InterestFilePath, "~", homeDir, 1)
	config.Tax.YearlySummaryPath = strings.Replace(config.Tax.YearlySummaryPath, "~", homeDir, 1)

	return
}
