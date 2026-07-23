package manager

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// InteractiveBrokersManagerImpl handles parsing of IB reports and implements Broker interface.
type InteractiveBrokersManagerImpl struct {
	basePath string
}

// SecurityIDProvider defines a narrow lookup for ticker-to-Security-ID from IBKR files.
type SecurityIDProvider interface {
	GetSecurityID(ctx context.Context, ticker string) (string, common.HttpError)
}

// ibRecordTypeData is the IB CSV record type for data rows (vs header/summary rows).
const ibRecordTypeData = "Data"

// ibRecordTypeInterest is the IB CSV record type for interest data rows.
const ibRecordTypeInterest = "Interest"

// ibRecordTypeStatement is the IB CSV record type for statement-level metadata rows.
const ibRecordTypeStatement = "Statement"

// ibStatementFieldPeriod is the IB CSV statement field name for the statement period.
const ibStatementFieldPeriod = "Period"

func NewInteractiveBrokersManagerImpl(basePath string) *InteractiveBrokersManagerImpl {
	return &InteractiveBrokersManagerImpl{
		basePath: basePath,
	}
}

var _ Broker = (*InteractiveBrokersManagerImpl)(nil)
var _ SecurityIDProvider = (*InteractiveBrokersManagerImpl)(nil)

const (
	ibSectionFI        = "Financial Instrument Information"
	ibFIDataTypeStocks = "Stocks"
	ibFIColSymbol      = 3
	ibFIColSecurityID  = 6
	ibMinFIELength     = 7
)

// GetSecurityID looks up the Security ID (CUSIP/ISIN) for the given ticker from IBKR
// Financial Instrument Information files. Returns a not-found error when the ticker
// has no entry, and a conflict error when multiple non-empty Security IDs disagree.
func (m *InteractiveBrokersManagerImpl) GetSecurityID(ctx context.Context, ticker string) (string, common.HttpError) {
	files, err := m.discoverFiles()
	if err != nil {
		return "", common.NewServerError(fmt.Errorf("failed to discover IBKR files: %w", err))
	}

	if err := m.checkContext(ctx); err != nil {
		return "", err
	}

	var foundID string
	for _, filePath := range files {
		fileID, conflict, err := m.readFileAndScanFI(ctx, filePath, ticker, foundID)
		if err != nil {
			return "", err
		}
		if conflict {
			return "", common.ErrEntityExists
		}
		if fileID != "" {
			foundID = fileID
		}
	}

	if foundID == "" {
		return "", common.ErrNotFound
	}
	return foundID, nil
}

// checkContext returns a server error wrapping ctx.Err() if the context is done.
func (m *InteractiveBrokersManagerImpl) checkContext(ctx context.Context) common.HttpError {
	select {
	case <-ctx.Done():
		return common.NewServerError(ctx.Err())
	default:
		return nil
	}
}

// readFileAndScanFI reads a CSV file and scans its FI records for the ticker.
// Returns the Security ID found, a conflict flag, or an error.
func (m *InteractiveBrokersManagerImpl) readFileAndScanFI(ctx context.Context, filePath, ticker, existingID string) (string, bool, common.HttpError) {
	if err := m.checkContext(ctx); err != nil {
		return "", false, err
	}

	records, err := m.readCSVRecords(filePath)
	if err != nil {
		return "", false, common.NewServerError(fmt.Errorf("failed to read IBKR file %s: %w", filePath, err))
	}

	fileID, conflict := m.scanFIRecords(records, ticker, existingID)
	return fileID, conflict, nil
}

// scanFIRecords scans Financial Instrument Information records for the given ticker.
// Returns the Security ID if found, and a conflict flag if a different non-empty ID
// is found compared to the existingID.
func (m *InteractiveBrokersManagerImpl) scanFIRecords(records [][]string, ticker, existingID string) (string, bool) {
	for _, record := range records {
		if len(record) < ibMinFIELength {
			continue
		}
		if record[0] != ibSectionFI || record[1] != ibRecordTypeData || record[2] != ibFIDataTypeStocks {
			continue
		}
		if record[ibFIColSymbol] != ticker {
			continue
		}

		secID := strings.TrimSpace(record[ibFIColSecurityID])
		if secID == "" {
			continue
		}

		if existingID == "" {
			existingID = secID
		} else if existingID != secID {
			return "", true
		}
	}
	return existingID, false
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
	periodEnd, err := m.parsePeriod(records)
	if err != nil {
		return tax.BrokerageInfo{}, err
	}

	interests := m.parseInterest(records)
	trades := m.parseTrades(records)
	dividends := m.parseDividends(records)

	return tax.BrokerageInfo{
		CoverageThrough: periodEnd,
		Interests:       interests,
		Trades:          trades,
		Dividends:       dividends,
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

// parsePeriod extracts the Statement,Data,Period row and returns the period end date.
// Format: January 1, 2024 - December 31, 2024
func (m *InteractiveBrokersManagerImpl) parsePeriod(records [][]string) (time.Time, error) {
	for _, record := range records {
		if len(record) >= 4 && record[0] == ibRecordTypeStatement && record[1] == ibRecordTypeData && record[2] == ibStatementFieldPeriod {
			periodStr := record[3]
			_, endDateStr, ok := strings.Cut(periodStr, " - ")
			if !ok {
				return time.Time{}, fmt.Errorf("invalid period format: %s", periodStr)
			}
			endDateStr = strings.TrimSpace(endDateStr)
			endDate, err := time.Parse("January 2, 2006", endDateStr)
			if err != nil {
				return time.Time{}, fmt.Errorf("failed to parse period end date %q: %w", endDateStr, err)
			}
			return endDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("period metadata not found")
}

func (m *InteractiveBrokersManagerImpl) parseInterest(records [][]string) []tax.Interest {
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

	return interests
}

func (m *InteractiveBrokersManagerImpl) isValidInterestRecord(record []string) bool {
	return len(record) >= 6 && record[0] == ibRecordTypeInterest && record[1] == ibRecordTypeData && record[2] == "USD"
}

func (m *InteractiveBrokersManagerImpl) parseTrades(records [][]string) []tax.Trade {
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

	return trades
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

func (m *InteractiveBrokersManagerImpl) parseDividends(records [][]string) []tax.Dividend {
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

	return dividends
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

	// FIXME: Decide whether an IBKR dividend without matching withholding tax should fail parsing instead of defaulting tax to zero.
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
