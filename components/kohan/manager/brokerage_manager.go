package manager

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type BrokerageManager interface {
	ParseAndGenerate(ctx context.Context, year int) error
}

type BrokerageManagerImpl struct {
	DriveWealth  Broker                  `container:"name"`
	IB           Broker                  `container:"name"`
	GainsManager GainsComputationManager `container:"type"`
	Config       config.TaxConfig
}

// Constructor for tests
func NewBrokerageManager(
	dwManager Broker,
	ibManager Broker,
	gainsManager GainsComputationManager,
	config config.TaxConfig,
) *BrokerageManagerImpl {
	return &BrokerageManagerImpl{
		DriveWealth:  dwManager,
		IB:           ibManager,
		GainsManager: gainsManager,
		Config:       config,
	}
}

var _ BrokerageManager = (*BrokerageManagerImpl)(nil)

func (m *BrokerageManagerImpl) ParseAndGenerate(ctx context.Context, year int) error {
	brokers := []Broker{m.DriveWealth, m.IB}
	requiredDate := time.Date(year+1, tax.COVERAGE_CUTOFF_MONTH, tax.COVERAGE_CUTOFF_DAY, 0, 0, 0, 0, time.UTC)

	var merged tax.BrokerageInfo
	for _, broker := range brokers {
		info, err := broker.Parse(year)
		if err != nil {
			return fmt.Errorf("broker '%s': parse failed: %w", broker.GetName(), err)
		}
		log.Info().Str("broker", broker.GetName()).Msg("Parsed broker successfully")

		if info.CoverageThrough.Before(requiredDate) {
			return fmt.Errorf("broker '%s': coverage through %s, required %s",
				broker.GetName(), info.CoverageThrough.Format(time.DateOnly), requiredDate.Format(time.DateOnly))
		}

		merged = mergeBrokerageInfo(merged, info)
	}

	hasData := len(merged.Trades) > 0 || len(merged.Dividends) > 0 || len(merged.Interests) > 0
	if !hasData {
		return fmt.Errorf("no data found from any broker")
	}

	return m.writeCSVs(ctx, merged)
}

func (m *BrokerageManagerImpl) writeCSVs(ctx context.Context, info tax.BrokerageInfo) error {
	if err := m.createInterestFile(info.Interests); err != nil {
		return err
	}

	// Sort trades before writing to ensure correct ordering (BUY before SELL on same date)
	sortedTrades := sortTradesByDate(info.Trades)

	// Persist raw (broker-original) trades
	if err := m.createTradeFile(sortedTrades); err != nil {
		return err
	}

	if err := m.createDividendFile(info.Dividends); err != nil {
		return err
	}

	return m.createGainsFile(ctx, sortedTrades)
}

func (m *BrokerageManagerImpl) createInterestFile(interests []tax.Interest) error {
	file, err := os.Create(m.Config.InterestFilePath)
	if err != nil {
		return fmt.Errorf("failed to create interest file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&interests, file); err != nil {
		return fmt.Errorf("failed to marshal interest data: %w", err)
	}
	return nil
}

func (m *BrokerageManagerImpl) createTradeFile(trades []tax.Trade) error {
	file, err := os.Create(m.Config.TradesPath)
	if err != nil {
		return fmt.Errorf("failed to create trades file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&trades, file); err != nil {
		return fmt.Errorf("failed to marshal trades data: %w", err)
	}
	return nil
}

func (m *BrokerageManagerImpl) createDividendFile(dividends []tax.Dividend) error {
	file, err := os.Create(m.Config.DividendFilePath)
	if err != nil {
		return fmt.Errorf("failed to create dividends file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&dividends, file); err != nil {
		return fmt.Errorf("failed to marshal dividends data: %w", err)
	}
	return nil
}

func (m *BrokerageManagerImpl) createGainsFile(ctx context.Context, trades []tax.Trade) error {
	gains, httpErr := m.GainsManager.ComputeGainsFromTrades(ctx, trades)
	if httpErr != nil {
		return httpErr
	}

	file, err := os.Create(m.Config.GainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create gains file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&gains, file); err != nil {
		return fmt.Errorf("failed to marshal gains data: %w", err)
	}
	return nil
}

// FIXME: Broker source is lost after merge — append does not tag which broker contributed each row.
// Gains, Dividends, Interest sheets and Summary all show aggregated values with no per-broker
// breakdown. Adding a broker field to models would let downstream output show DW vs IBKR splits.
func mergeBrokerageInfo(a, b tax.BrokerageInfo) tax.BrokerageInfo {
	coverage := a.CoverageThrough
	if b.CoverageThrough.After(coverage) {
		coverage = b.CoverageThrough
	}
	return tax.BrokerageInfo{
		CoverageThrough: coverage,
		Trades:          append(a.Trades, b.Trades...),
		Dividends:       append(a.Dividends, b.Dividends...),
		Interests:       append(a.Interests, b.Interests...),
	}
}

// sortTradesByDate sorts trades chronologically by date and ensures BUY trades come before SELL trades on the same date
func sortTradesByDate(trades []tax.Trade) []tax.Trade {
	sortedTrades := make([]tax.Trade, len(trades))
	copy(sortedTrades, trades)

	sort.SliceStable(sortedTrades, func(i, j int) bool {
		dateI, errI := time.Parse(time.DateOnly, sortedTrades[i].Date)
		dateJ, errJ := time.Parse(time.DateOnly, sortedTrades[j].Date)
		if errI != nil || errJ != nil {
			return false // Keep original order if date parsing fails
		}

		// If dates are different, sort by date
		if !dateI.Equal(dateJ) {
			return dateI.Before(dateJ)
		}

		// For same-day trades, ensure BUY comes before SELL
		typeI := strings.ToUpper(sortedTrades[i].Type)
		typeJ := strings.ToUpper(sortedTrades[j].Type)

		// BUY should come before SELL on the same day
		if typeI == tax.TRADE_TYPE_BUY && typeJ == tax.TRADE_TYPE_SELL {
			return true
		}
		if typeI == tax.TRADE_TYPE_SELL && typeJ == tax.TRADE_TYPE_BUY {
			return false
		}

		return false // Keep original order for same type
	})

	return sortedTrades
}
