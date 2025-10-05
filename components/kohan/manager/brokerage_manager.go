package manager

import (
	"context"
	"fmt"
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type BrokerageManager interface {
	ParseAndGenerate(ctx context.Context) error
}

type BrokerageManagerImpl struct {
	brokers      []Broker
	gainsManager GainsComputationManager
	config       config.TaxConfig
}

func NewBrokerageManager(
	dwManager Broker,
	ibManager Broker,
	gainsManager GainsComputationManager,
	config config.TaxConfig,
) BrokerageManager {
	return &BrokerageManagerImpl{
		brokers:      []Broker{dwManager, ibManager},
		gainsManager: gainsManager,
		config:       config,
	}
}

func (m *BrokerageManagerImpl) ParseAndGenerate(ctx context.Context) error {
	var merged tax.BrokerageInfo

	for _, broker := range m.brokers {
		info, err := broker.Parse()
		if err != nil {
			log.Warn().Str("broker", broker.GetName()).Err(err).Msg("Broker parse failed, skipping")
			continue
		}
		merged = mergeBrokerageInfo(merged, info)
		log.Info().Str("broker", broker.GetName()).Msg("Parsed broker successfully")
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

	if err := m.createTradeFile(info.Trades); err != nil {
		return err
	}

	if err := m.createDividendFile(info.Dividends); err != nil {
		return err
	}

	return m.createGainsFile(ctx, info.Trades)
}

func (m *BrokerageManagerImpl) createInterestFile(interests []tax.Interest) error {
	file, err := os.Create(m.config.InterestFilePath)
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
	file, err := os.Create(m.config.TradesPath)
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
	file, err := os.Create(m.config.DividendFilePath)
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
	gains, httpErr := m.gainsManager.ComputeGainsFromTrades(ctx, trades)
	if httpErr != nil {
		return httpErr
	}

	file, err := os.Create(m.config.GainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create gains file: %w", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&gains, file); err != nil {
		return fmt.Errorf("failed to marshal gains data: %w", err)
	}
	return nil
}

func mergeBrokerageInfo(a, b tax.BrokerageInfo) tax.BrokerageInfo {
	return tax.BrokerageInfo{
		Trades:    append(a.Trades, b.Trades...),
		Dividends: append(a.Dividends, b.Dividends...),
		Interests: append(a.Interests, b.Interests...),
	}
}
