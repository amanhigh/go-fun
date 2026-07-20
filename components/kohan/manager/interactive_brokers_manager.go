package manager

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/tax"
)

// InteractiveBrokersManagerImpl handles parsing of IB reports and implements Broker interface.
type InteractiveBrokersManagerImpl struct {
	basePath string
}

// ibRecordTypeData is the IB CSV record type for data rows (vs header/summary rows).
const ibRecordTypeData = "Data"

// ibRecordTypeInterest is the IB CSV record type for interest data rows.
const ibRecordTypeInterest = "Interest"

func NewInteractiveBrokersManagerImpl(basePath string) Broker {
	return &InteractiveBrokersManagerImpl{
		basePath: basePath,
	}
}

// GetName returns the broker name.
func (m *InteractiveBrokersManagerImpl) GetName() string {
	return "Interactive Brokers"
}

// discoverFiles finds every file matching <basePath>_YYYY.csv pattern in deterministic lexical order.
func (m *InteractiveBrokersManagerImpl) discoverFiles() ([]string, error) {
	pattern := m.basePath + "_[0-9][0-9][0-9][0-9].csv"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("failed to open CSV file")
	}
	return files, nil
}

func (m *InteractiveBrokersManagerImpl) Parse(_ int) (tax.BrokerageInfo, error) {
	files, err := m.discoverFiles()
	if err != nil {
		return tax.BrokerageInfo{}, err
	}

	var merged tax.BrokerageInfo
	for _, filePath := range files {
		records, err := m.readCSVRecords(filePath)
		if err != nil {
			return tax.BrokerageInfo{}, err
		}

		info, err := m.parseRecords(records)
		if err != nil {
			return tax.BrokerageInfo{}, err
		}

		merged = mergeBrokerageInfo(merged, info)
	}

	return merged, nil
}

// parseRecords parses all record types (interest, trades, dividends) from the given CSV records
// and returns a complete BrokerageInfo. Each call builds its own withholding-tax map from the
// records, so per-file tax matching is preserved.
func (m *InteractiveBrokersManagerImpl) parseRecords(records [][]string) (tax.BrokerageInfo, error) {
	interests, err := m.parseInterest(records)
	if err != nil {
		return tax.BrokerageInfo{}, err
	}

	trades, err := m.parseTrades(records)
	if err != nil {
		return tax.BrokerageInfo{}, err
	}

	dividends, err := m.parseDividends(records)
	if err != nil {
		return tax.BrokerageInfo{}, err
	}

	return tax.BrokerageInfo{
		Interests: interests,
		Trades:    trades,
		Dividends: dividends,
	}, nil
}

// readCSVRecords reads and parses all CSV records from the given file path.
// Extracted from Parse to keep statement count within funlen limit.
func (m *InteractiveBrokersManagerImpl) readCSVRecords(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	return records, nil
}

func (m *InteractiveBrokersManagerImpl) parseInterest(records [][]string) ([]tax.Interest, error) {
	var interests []tax.Interest

	for _, record := range records {
		if !m.isValidInterestRecord(record) {
			continue
		}

		date := record[3]
		amount, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}

		interests = append(interests, tax.Interest{
			Symbol: "CASH",
			Date:   date,
			Amount: amount,
			Tax:    0,
			Net:    amount,
		})
	}

	return interests, nil
}

func (m *InteractiveBrokersManagerImpl) isValidInterestRecord(record []string) bool {
	return len(record) >= 6 && record[0] == ibRecordTypeInterest && record[1] == ibRecordTypeData && record[2] == "USD"
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
	return len(record) >= 14 && record[0] == "Trades" && record[1] == ibRecordTypeData
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

	// Skip C. Price at record[9], read Proceeds at record[10]
	proceeds, err := strconv.ParseFloat(record[10], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse proceeds: %w", err)
	}

	// Read Comm/Fee at record[11]
	commission, err := strconv.ParseFloat(record[11], 64)
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
	return len(record) >= 6 && record[0] == "Dividends" && record[1] == ibRecordTypeData
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

	MatchDividendWithTax(dividend, taxMap)

	return *dividend, nil
}

func (m *InteractiveBrokersManagerImpl) buildTaxMap(records [][]string) map[string]map[string]float64 {
	taxMap := make(map[string]map[string]float64)

	for _, record := range records {
		if len(record) < 6 || record[0] != "Withholding Tax" || record[1] != ibRecordTypeData {
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
	if before, _, ok := strings.Cut(description, "("); ok {
		return before
	}
	return ""
}
