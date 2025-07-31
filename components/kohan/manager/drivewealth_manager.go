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
func (m *DriveWealthManager) Parse() ([]tax.Interest, error) {
	f, err := excelize.OpenFile(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
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
		return nil, fmt.Errorf("sheet 'Income' not found in the Excel file")
	}

	rows, err := f.GetRows("Income")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from 'Income' sheet: %w", err)
	}

	interestEntries, err := m.parseInterest(rows)
	if err != nil {
		return nil, err
	}

	// TODO: Parse Dividends from the rows
	// TODO: Parse Taxes from the rows

	return interestEntries, nil
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
