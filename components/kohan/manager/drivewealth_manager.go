package manager

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/tax"
	"github.com/xuri/excelize/v2"
)

const (
	// tradeRowLength is the expected number of columns in a trade row.
	tradeRowLength = 10
)

// DriveWealthManagerImpl handles parsing of DriveWealth reports and implements Broker interface.
type DriveWealthManagerImpl struct {
	basePath string
}

// NewDriveWealthManagerImpl creates a new DriveWealth broker parser.
// basePath should be the base path without year or extension (e.g., ~/path/to/vested)
func NewDriveWealthManagerImpl(basePath string) Broker {
	return &DriveWealthManagerImpl{
		basePath: basePath,
	}
}

// GetName returns the broker name.
func (m *DriveWealthManagerImpl) GetName() string {
	return "DriveWealth"
}

// resolveFilePath constructs the year-specific file path
func (m *DriveWealthManagerImpl) resolveFilePath(year int) string {
	return fmt.Sprintf("%s_%d.xlsx", m.basePath, year)
}

// Parse orchestrates the parsing of the DriveWealth Excel file for the given year.
func (m *DriveWealthManagerImpl) Parse(year int) (info tax.BrokerageInfo, err error) {
	filePath := m.resolveFilePath(year)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		err = fmt.Errorf("failed to open excel file: %w", err)
		return
	}
	defer f.Close()

	// Verify that the required sheets exist.
	if err = m.checkSheetExists(f, "Income"); err != nil {
		return
	}
	if err = m.checkSheetExists(f, "Trades"); err != nil {
		return
	}

	// Parse the sheets.
	if info, err = m.parseSheets(f); err != nil {
		return
	}

	return
}

func (m *DriveWealthManagerImpl) parseSheets(f *excelize.File) (info tax.BrokerageInfo, err error) {
	rows, err := f.GetRows("Income")
	if err != nil {
		err = fmt.Errorf("failed to get rows from 'Income' sheet: %w", err)
		return
	}

	info.Interests, err = m.parseInterest(rows)
	if err != nil {
		return
	}

	info.Dividends, err = m.parseDividends(rows)
	if err != nil {
		return
	}

	tradeRows, err := f.GetRows("Trades")
	if err != nil {
		err = fmt.Errorf("failed to get rows from 'Trades' sheet: %w", err)
		return
	}

	// Load commission map from All Transactions sheet if available
	commissionMap := make(map[string]float64)
	allTransRows, errAT := f.GetRows("All Transactions")
	if errAT == nil {
		commissionMap = m.parseCommissions(allTransRows)
	}

	info.Trades, err = m.parseTrades(tradeRows, commissionMap)
	return
}

func (m *DriveWealthManagerImpl) checkSheetExists(f *excelize.File, sheetName string) error {
	if slices.Contains(f.GetSheetList(), sheetName) {
		return nil
	}
	return fmt.Errorf("sheet '%s' not found in the Excel file", sheetName)
}

func (m *DriveWealthManagerImpl) parseTrades(rows [][]string, commissionMap map[string]float64) ([]tax.Trade, error) {
	var trades []tax.Trade
	if len(rows) <= 1 {
		return trades, nil
	}

	for _, row := range rows[1:] { // Skip header row
		if len(row) < tradeRowLength {
			continue
		}

		trade, err := m.parseTradeRow(row, commissionMap)
		if err != nil {
			continue
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// parseTradeRow parses a single trade row and applies commission fallback logic.
func (m *DriveWealthManagerImpl) parseTradeRow(row []string, commissionMap map[string]float64) (tax.Trade, error) {
	quantity, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse quantity: %w", err)
	}
	price, err := strconv.ParseFloat(row[7], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse price: %w", err)
	}
	value, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse value: %w", err)
	}
	commission, err := strconv.ParseFloat(row[9], 64)
	if err != nil {
		return tax.Trade{}, fmt.Errorf("failed to parse commission: %w", err)
	}

	// Apply commission fallback: if Trades sheet commission is zero, lookup from All Transactions
	if commission == 0 && len(commissionMap) > 0 {
		date := strings.Split(row[0], " ")[0]
		lookupKey := fmt.Sprintf("%s|%s|%s", date, row[3], row[4])
		if fallbackCommission, exists := commissionMap[lookupKey]; exists {
			commission = fallbackCommission
		}
	}

	date := strings.Split(row[0], " ")[0]
	return tax.Trade{
		Symbol:     row[3],
		Date:       date,
		Type:       row[4],
		Quantity:   quantity,
		USDPrice:   price,
		USDValue:   value,
		Commission: commission,
	}, nil
}

// parseInterest extracts interest entries from the "Income" sheet rows.
func (m *DriveWealthManagerImpl) parseInterest(rows [][]string) ([]tax.Interest, error) {
	var interestEntries []tax.Interest

	if len(rows) > 0 {
		for _, row := range rows[1:] { // Skip header row
			if len(row) >= 5 && row[2] == "Interest" {
				amount, err := strconv.ParseFloat(row[4], 64)
				if err != nil {
					// Potentially log the error or handle it more gracefully
					continue // Skip rows where amount is not a valid float
				}

				entry := tax.Interest{
					Symbol: "CASH",
					Date:   strings.Split(row[0], " ")[0],
					Amount: amount,
					Tax:    0,
					Net:    amount,
				}
				interestEntries = append(interestEntries, entry)
			}
		}
	}

	return interestEntries, nil
}

func (m *DriveWealthManagerImpl) parseDividends(rows [][]string) ([]tax.Dividend, error) {
	taxMap, err := m.buildTaxMap(rows)
	if err != nil {
		return nil, err
	}

	var dividendEntries []tax.Dividend
	for _, row := range rows[1:] {
		if len(row) >= 5 && row[2] == "Dividend" {
			amount, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				continue
			}

			dividend := &tax.Dividend{
				Symbol: row[3],
				Date:   strings.Split(row[0], " ")[0],
				Amount: amount,
			}

			MatchDividendWithTax(dividend, taxMap)
			dividendEntries = append(dividendEntries, *dividend)
		}
	}

	return dividendEntries, nil
}

func (m *DriveWealthManagerImpl) buildTaxMap(rows [][]string) (map[string]map[string]float64, error) {
	// taxMap stores tax amounts keyed by symbol, then by date.
	// This allows for efficient lookup of taxes for a given dividend.
	taxMap := make(map[string]map[string]float64) // symbol -> date -> taxAmount

	// First pass: Iterate through all rows to aggregate tax entries.
	// This is done first because taxes may not appear immediately after dividends.
	for _, row := range rows[1:] { // Skip header
		if len(row) >= 5 && row[2] == "Tax" {
			symbol := row[3]
			date := strings.Split(row[0], " ")[0]
			taxAmount, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				continue // Skip row if tax amount is not a valid number.
			}

			// Create the nested map if it doesn't exist for the symbol.
			if _, ok := taxMap[symbol]; !ok {
				taxMap[symbol] = make(map[string]float64)
			}
			// Add the tax amount. Report lists taxes as negative, so we negate to store as a positive value.
			taxMap[symbol][date] += -taxAmount
		}
	}
	return taxMap, nil
}

// parseCommissions extracts commission data from "All Transactions" sheet.
// Returns a map with key format "Date|Symbol|Type" -> commission amount.
// Expected comment format: "COMM Buy SYMBOL base=amount" or "COMM Sell SYMBOL base=amount"
func (m *DriveWealthManagerImpl) parseCommissions(rows [][]string) map[string]float64 {
	commissionMap := make(map[string]float64)

	if len(rows) <= 1 {
		return commissionMap // Return empty map if no data rows
	}

	for _, row := range rows[1:] { // Skip header row
		// Expected columns: Date(0), Time(1), Type(2), Amount(3), Account Balance(4), Comment(5)
		if len(row) >= 6 && row[2] == "COMM" {
			comment := row[5]
			// Parse comment: "COMM Buy SYMBOL base=amount"
			parts := strings.Fields(comment)
			if len(parts) >= 4 && parts[0] == "COMM" {
				tradeType := parts[1]  // "Buy" or "Sell"
				symbol := parts[2]     // Stock symbol
				basePrefix := parts[3] // "base=amount"

				// Extract amount from "base=amount"
				if after, ok := strings.CutPrefix(basePrefix, "base="); ok {
					commissionStr := after
					commission, err := strconv.ParseFloat(commissionStr, 64)
					if err != nil {
						continue // Skip malformed entries
					}

					date := strings.Split(row[0], " ")[0]
					lookupKey := fmt.Sprintf("%s|%s|%s", date, symbol, tradeType)
					commissionMap[lookupKey] = commission
				}
			}
		}
	}

	return commissionMap
}
