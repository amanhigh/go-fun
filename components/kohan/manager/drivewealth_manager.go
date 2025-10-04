package manager

import (
	"fmt"
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
	BrokerageParserHelper
	excelPath string
}

// NewDriveWealthManagerImpl creates a new DriveWealth broker parser.
func NewDriveWealthManagerImpl(excelPath string) Broker {
	return &DriveWealthManagerImpl{
		excelPath: excelPath,
	}
}

// Parse orchestrates the parsing of the DriveWealth Excel file.
func (m *DriveWealthManagerImpl) Parse() (info tax.BrokerageInfo, err error) {
	f, err := excelize.OpenFile(m.excelPath)
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
	info.Trades, err = m.parseTrades(tradeRows)
	if err != nil {
		return
	}
	return
}

func (m *DriveWealthManagerImpl) checkSheetExists(f *excelize.File, sheetName string) error {
	for _, sheet := range f.GetSheetList() {
		if sheet == sheetName {
			return nil
		}
	}
	return fmt.Errorf("sheet '%s' not found in the Excel file", sheetName)
}

func (m *DriveWealthManagerImpl) parseTrades(rows [][]string) ([]tax.Trade, error) {
	var trades []tax.Trade
	if len(rows) > 0 {
		for _, row := range rows[1:] { // Skip header row
			if len(row) >= tradeRowLength {
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

			m.MatchDividendWithTax(dividend, taxMap)
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
