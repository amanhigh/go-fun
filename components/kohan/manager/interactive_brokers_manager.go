package manager

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/tax"
)

// InteractiveBrokersManagerImpl handles parsing of IB reports and implements Broker interface.
type InteractiveBrokersManagerImpl struct {
	BrokerageParserHelper
	csvPath string
}

func NewInteractiveBrokersManagerImpl(csvPath string) Broker {
	return &InteractiveBrokersManagerImpl{
		csvPath: csvPath,
	}
}

func (m *InteractiveBrokersManagerImpl) Parse() (info tax.BrokerageInfo, err error) {
	file, err := os.Open(m.csvPath)
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

	dividend := &tax.Dividend{
		Symbol: symbol,
		Date:   date,
		Amount: amount,
	}

	m.MatchDividendWithTax(dividend, taxMap)

	return *dividend, nil
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
