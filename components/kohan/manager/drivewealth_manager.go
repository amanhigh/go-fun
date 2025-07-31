package manager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/tax"
	"github.com/xuri/excelize/v2"
)

// DriveWealthManager handles parsing of DriveWealth reports.
type DriveWealthManager struct {
	filePath string
}

// NewDriveWealthManager creates a new DriveWealthManager.
func NewDriveWealthManager(filePath string) *DriveWealthManager {
	return &DriveWealthManager{filePath: filePath}
}

// Parse orchestrates the parsing of the DriveWealth Excel file.
func (m *DriveWealthManager) Parse() (interests []tax.Interest, dividends []tax.Dividend, trades []tax.Trade, err error) {
	f, err := excelize.OpenFile(m.filePath)
	if err != nil {
		err = fmt.Errorf("failed to open excel file: %w", err)
		return
	}
	defer f.Close()

	// Check if "Income" sheet exists
	sheetExists := false
	for _, sheet := range f.GetSheetList() {
		if sheet == "Income" {
			sheetExists = true
			break
		}
	}

	if !sheetExists {
		err = fmt.Errorf("sheet 'Income' not found in the Excel file")
		return
	}

	// Check if "Trades" sheet exists
	sheetExists = false
	for _, sheet := range f.GetSheetList() {
		if sheet == "Trades" {
			sheetExists = true
			break
		}
	}

	if !sheetExists {
		err = fmt.Errorf("sheet 'Trades' not found in the Excel file")
		return
	}

	rows, err := f.GetRows("Income")
	if err != nil {
		err = fmt.Errorf("failed to get rows from 'Income' sheet: %w", err)
		return
	}

	interests, err = m.parseInterest(rows)
	if err != nil {
		return
	}

	dividends, err = m.parseDividends(rows)
	if err != nil {
		return
	}

	tradeRows, err := f.GetRows("Trades")
	if err != nil {
		err = fmt.Errorf("failed to get rows from 'Trades' sheet: %w", err)
		return
	}
	trades, err = m.parseTrades(tradeRows)
	if err != nil {
		return
	}

	return
}

func (m *DriveWealthManager) parseTrades(rows [][]string) ([]tax.Trade, error) {
	var trades []tax.Trade
	if len(rows) > 0 {
		for _, row := range rows[1:] { // Skip header row
			if len(row) >= 9 {
				quantity, err := strconv.ParseFloat(row[6], 64)
				if err != nil {
					continue
				}
				price, err := strconv.ParseFloat(row[7], 64)
				if err != nil {
					continue
				}
				value, err := strconv.ParseFloat(row[8], 64)
				if err != nil {
					continue
				}
				commission, err := strconv.ParseFloat(row[9], 64)
				if err != nil {
					continue
				}

				trade := tax.Trade{
					Symbol:     row[3],
					Date:       strings.Split(row[0], " ")[0],
					Type:       row[4],
					Quantity:   quantity,
					USDPrice:   price,
					USDValue:   value,
					Commission: commission,
				}
				trades = append(trades, trade)
			}
		}
	}
	return trades, nil
}

// parseInterest extracts interest entries from the "Income" sheet rows.
func (m *DriveWealthManager) parseInterest(rows [][]string) ([]tax.Interest, error) {
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

func (m *DriveWealthManager) parseDividends(rows [][]string) ([]tax.Dividend, error) {
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

	var dividendEntries []tax.Dividend
	// Second pass: Process dividend entries and associate them with the collected taxes.
	for _, row := range rows[1:] { // Skip header
		if len(row) >= 5 && row[2] == "Dividend" {
			amount, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				continue // Skip row if dividend amount is not a valid number.
			}

			symbol := row[3]
			date := strings.Split(row[0], " ")[0]

			entry := tax.Dividend{
				Symbol: symbol,
				Date:   date,
				Amount: amount,
			}

			// Look for a matching tax in the map using the dividend's symbol and date.
			if dateTaxes, ok := taxMap[symbol]; ok {
				if taxAmount, ok := dateTaxes[date]; ok {
					entry.Tax = taxAmount
					// Remove the tax from the map to ensure it's not used again.
					delete(dateTaxes, date)
				}
			}

			// Calculate the net amount after deducting tax.
			entry.Net = entry.Amount - entry.Tax
			dividendEntries = append(dividendEntries, entry)
		}
	}

	return dividendEntries, nil
}
