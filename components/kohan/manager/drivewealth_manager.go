package manager

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	"github.com/xuri/excelize/v2"
)

const (
	// tradeRowLength is the expected number of columns in a trade row.
	tradeRowLength = 10
)

type DriveWealthManager interface {
	Parse() (info tax.DriveWealthInfo, err error)
	GenerateCsv(ctx context.Context, info tax.DriveWealthInfo) (err error)
}

// DriveWealthManagerImpl handles parsing of DriveWealth reports.
type DriveWealthManagerImpl struct {
	config       config.TaxConfig
	gainsManager GainsComputationManager
}

// NewDriveWealthManager creates a new DriveWealthManager.
func NewDriveWealthManager(config config.TaxConfig, gainsManager GainsComputationManager) DriveWealthManager {
	return &DriveWealthManagerImpl{
		config:       config,
		gainsManager: gainsManager,
	}
}

func (m *DriveWealthManagerImpl) GenerateCsv(ctx context.Context, info tax.DriveWealthInfo) (err error) {
	if err = m.createInterestFile(info.Interests); err != nil {
		return
	}
	if err = m.createTradeFile(info.Trades); err != nil {
		return
	}
	if err = m.createDividendFile(info.Dividends); err != nil {
		return
	}
	return m.createGainsFile(ctx, info.Trades)
}

func (m *DriveWealthManagerImpl) createInterestFile(interests []tax.Interest) error {
	interestFile, err := os.Create(m.config.InterestFilePath)
	if err != nil {
		return fmt.Errorf("failed to create interest file: %w", err)
	}
	defer interestFile.Close()

	if err := gocsv.MarshalFile(&interests, interestFile); err != nil {
		return fmt.Errorf("failed to marshal interest data: %w", err)
	}
	return nil
}

func (m *DriveWealthManagerImpl) createTradeFile(trades []tax.Trade) error {
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

func (m *DriveWealthManagerImpl) createDividendFile(dividends []tax.Dividend) error {
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

func (m *DriveWealthManagerImpl) createGainsFile(ctx context.Context, trades []tax.Trade) error {
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

// Parse orchestrates the parsing of the DriveWealth Excel file.
func (m *DriveWealthManagerImpl) Parse() (info tax.DriveWealthInfo, err error) {
	f, err := excelize.OpenFile(m.config.DriveWealthPath)
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

func (m *DriveWealthManagerImpl) parseSheets(f *excelize.File) (info tax.DriveWealthInfo, err error) {
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
