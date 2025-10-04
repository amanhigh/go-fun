package manager

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
)

type InteractiveBrokersManager interface {
	Parse() (info tax.BrokerageInfo, err error)
	GenerateCsv(ctx context.Context, info tax.BrokerageInfo) (err error)
}

type InteractiveBrokersManagerImpl struct {
	config       config.TaxConfig
	gainsManager GainsComputationManager
}

func NewInteractiveBrokersManager(config config.TaxConfig, gainsManager GainsComputationManager) InteractiveBrokersManager {
	return &InteractiveBrokersManagerImpl{
		config:       config,
		gainsManager: gainsManager,
	}
}

func (m *InteractiveBrokersManagerImpl) GenerateCsv(ctx context.Context, info tax.BrokerageInfo) (err error) {
	// TODO: Add Interest file generation when IB Activity Statement parser is implemented
	// Interest income is available in IB Activity Statement, not in Realized.csv
	// if err = m.createInterestFile(info.Interests); err != nil {
	// 	return
	// }

	if err = m.createTradeFile(info.Trades); err != nil {
		return
	}
	if err = m.createDividendFile(info.Dividends); err != nil {
		return
	}
	return m.createGainsFile(ctx, info.Trades)
}

func (m *InteractiveBrokersManagerImpl) createTradeFile(trades []tax.Trade) error {
	tradeFile, err := os.Create(m.config.TradesPath)
	if err != nil {
		return fmt.Errorf("failed to create trades file: %w", err)
	}
	defer tradeFile.Close()

	if err := gocsv.MarshalFile(&trades, tradeFile); err != nil {
		return fmt.Errorf("failed to marshal trades data: %w", err)
	}
	return nil
}

func (m *InteractiveBrokersManagerImpl) createDividendFile(dividends []tax.Dividend) error {
	dividendFile, err := os.Create(m.config.DividendFilePath)
	if err != nil {
		return fmt.Errorf("failed to create dividends file: %w", err)
	}
	defer dividendFile.Close()

	if err := gocsv.MarshalFile(&dividends, dividendFile); err != nil {
		return fmt.Errorf("failed to marshal dividends data: %w", err)
	}
	return nil
}

func (m *InteractiveBrokersManagerImpl) createGainsFile(ctx context.Context, trades []tax.Trade) error {
	gains, httpErr := m.gainsManager.ComputeGainsFromTrades(ctx, trades)
	if httpErr != nil {
		return httpErr
	}

	gainsFile, err := os.Create(m.config.GainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to create gains file: %w", err)
	}
	defer gainsFile.Close()

	if err := gocsv.MarshalFile(&gains, gainsFile); err != nil {
		return fmt.Errorf("failed to marshal gains data: %w", err)
	}
	return nil
}

func (m *InteractiveBrokersManagerImpl) Parse() (info tax.BrokerageInfo, err error) {
	file, err := os.Open(m.config.IBPath)
	if err != nil {
		err = fmt.Errorf("failed to open CSV file: %w", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		err = fmt.Errorf("failed to read CSV: %w", err)
		return
	}

	// TODO: Parse Interest from IB Activity Statement (separate from Realized.csv)
	// info.Interests, err = m.parseInterest(records)
	// if err != nil {
	// 	return
	// }

	info.Trades, err = m.parseTrades(records)
	if err != nil {
		return
	}

	info.Dividends, err = m.parseDividends(records)
	if err != nil {
		return
	}

	return
}

func (m *InteractiveBrokersManagerImpl) parseTrades(records [][]string) ([]tax.Trade, error) {
	var trades []tax.Trade

	for _, record := range records {
		if !m.isValidTradeRecord(record) {
			continue
		}

		trade, err := m.parseTradeRecord(record)
		if err != nil {
			continue
		}

		trades = append(trades, trade)
	}

	return trades, nil
}

func (m *InteractiveBrokersManagerImpl) isValidTradeRecord(record []string) bool {
	return len(record) >= 14 && record[0] == "Trades" && record[1] == "Data"
}

func (m *InteractiveBrokersManagerImpl) parseTradeRecord(record []string) (tax.Trade, error) {
	symbol := record[5]
	dateTime := record[6]
	date := strings.Split(dateTime, ",")[0]

	quantity, err := strconv.ParseFloat(record[7], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse quantity: %w", err)
	}

	price, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse price: %w", err)
	}

	proceeds, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse proceeds: %w", err)
	}

	commission, err := strconv.ParseFloat(record[10], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse commission: %w", err)
	}

	tradeType := m.determineTradeType(quantity)

	return tax.Trade{
		Symbol:     symbol,
		Date:       date,
		Type:       tradeType,
		Quantity:   math.Abs(quantity),
		USDPrice:   price,
		USDValue:   math.Abs(proceeds),
		Commission: math.Abs(commission),
	}, nil
}

func (m *InteractiveBrokersManagerImpl) parseDividends(records [][]string) ([]tax.Dividend, error) {
	taxMap := m.buildTaxMap(records)

	var dividends []tax.Dividend

	for _, record := range records {
		if !m.isValidDividendRecord(record) {
			continue
		}

		dividend, err := m.parseDividendRecord(record, taxMap)
		if err != nil {
			continue
		}

		dividends = append(dividends, dividend)
	}

	return dividends, nil
}

func (m *InteractiveBrokersManagerImpl) isValidDividendRecord(record []string) bool {
	return len(record) >= 6 && record[0] == "Dividends" && record[1] == "Data"
}

func (m *InteractiveBrokersManagerImpl) parseDividendRecord(record []string, taxMap map[string]map[string]float64) (tax.Dividend, error) {
	date := record[3]
	description := record[4]
	symbol := extractSymbol(description)

	if symbol == "" {
		return tax.Dividend{}, fmt.Errorf("empty symbol")
	}

	amount, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return tax.Dividend{}, fmt.Errorf("failed to parse dividend amount: %w", err)
	}

	dividend := tax.Dividend{
		Symbol: symbol,
		Date:   date,
		Amount: amount,
	}

	m.applyWithholdingTax(&dividend, taxMap)
	dividend.Net = dividend.Amount - dividend.Tax

	return dividend, nil
}

func (m *InteractiveBrokersManagerImpl) applyWithholdingTax(dividend *tax.Dividend, taxMap map[string]map[string]float64) {
	if dateTaxes, ok := taxMap[dividend.Symbol]; ok {
		if taxAmount, ok := dateTaxes[dividend.Date]; ok {
			dividend.Tax = taxAmount
			delete(dateTaxes, dividend.Date)
		}
	}
}

func (m *InteractiveBrokersManagerImpl) buildTaxMap(records [][]string) map[string]map[string]float64 {
	taxMap := make(map[string]map[string]float64)

	for _, record := range records {
		if len(record) < 6 || record[0] != "Withholding Tax" || record[1] != "Data" {
			continue
		}

		date := record[3]
		description := record[4]
		symbol := extractSymbol(description)

		if symbol == "" {
			continue
		}

		taxAmount, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}

		if _, ok := taxMap[symbol]; !ok {
			taxMap[symbol] = make(map[string]float64)
		}
		taxMap[symbol][date] += math.Abs(taxAmount)
	}

	return taxMap
}

func (m *InteractiveBrokersManagerImpl) determineTradeType(quantity float64) string {
	if quantity > 0 {
		return tax.TRADE_TYPE_BUY
	}
	return tax.TRADE_TYPE_SELL
}

func extractSymbol(description string) string {
	if idx := strings.Index(description, "("); idx != -1 {
		return description[:idx]
	}
	return ""
}
