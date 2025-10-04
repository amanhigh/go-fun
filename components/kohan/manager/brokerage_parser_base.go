package manager

import (
	"context"
	"fmt"
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
)

// BrokerageParserBase provides common CSV generation functionality for all brokerage parsers.
// This base struct eliminates code duplication across DriveWealth, Interactive Brokers, and future parsers.
type BrokerageParserBase struct {
	config       config.TaxConfig
	gainsManager GainsComputationManager
}

// NewBrokerageParserBase creates a new BrokerageParserBase instance.
func NewBrokerageParserBase(config config.TaxConfig, gainsManager GainsComputationManager) BrokerageParserBase {
	return BrokerageParserBase{
		config:       config,
		gainsManager: gainsManager,
	}
}

// GenerateCsv creates all CSV output files from parsed brokerage information.
// It generates: trades.csv, dividends.csv, gains.csv, and interest.csv (even if empty).
func (b *BrokerageParserBase) GenerateCsv(ctx context.Context, info tax.BrokerageInfo) error {
	if err := b.createInterestFile(info.Interests); err != nil {
		return err
	}

	if err := b.createTradeFile(info.Trades); err != nil {
		return err
	}

	if err := b.createDividendFile(info.Dividends); err != nil {
		return err
	}

	return b.createGainsFile(ctx, info.Trades)
}

func (b *BrokerageParserBase) createInterestFile(interests []tax.Interest) error {
	file, err := os.Create(b.config.InterestFilePath)
	if err != nil {
		return fmt.Errorf("failed to create interest file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&interests, file); err != nil {
		return fmt.Errorf("failed to marshal interest data: %w", err)
	}
	return nil
}

func (b *BrokerageParserBase) createTradeFile(trades []tax.Trade) error {
	file, err := os.Create(b.config.TradesPath)
	if err != nil {
		return fmt.Errorf("failed to create trades file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&trades, file); err != nil {
		return fmt.Errorf("failed to marshal trades data: %w", err)
	}
	return nil
}

func (b *BrokerageParserBase) createDividendFile(dividends []tax.Dividend) error {
	file, err := os.Create(b.config.DividendFilePath)
	if err != nil {
		return fmt.Errorf("failed to create dividends file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&dividends, file); err != nil {
		return fmt.Errorf("failed to marshal dividends data: %w", err)
	}
	return nil
}

func (b *BrokerageParserBase) createGainsFile(ctx context.Context, trades []tax.Trade) error {
	gains, httpErr := b.gainsManager.ComputeGainsFromTrades(ctx, trades)
	if httpErr != nil {
		return httpErr
	}

	file, err := os.Create(b.config.GainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create gains file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&gains, file); err != nil {
		return fmt.Errorf("failed to marshal gains data: %w", err)
	}
	return nil
}

// MatchDividendWithTax applies withholding tax from taxMap to dividend and calculates net.
// This is a common pattern used by all brokers: match dividend to tax by symbol+date, then compute net.
func (b *BrokerageParserBase) MatchDividendWithTax(dividend *tax.Dividend, taxMap map[string]map[string]float64) {
	if dateTaxes, ok := taxMap[dividend.Symbol]; ok {
		if taxAmount, ok := dateTaxes[dividend.Date]; ok {
			dividend.Tax = taxAmount
			delete(dateTaxes, dividend.Date)
		}
	}
	dividend.Net = dividend.Amount - dividend.Tax
}
