package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type ExcelManager interface {
	GenerateTaxSummaryExcel(ctx context.Context, year int, summary tax.Summary) error
}

type ExcelManagerImpl struct {
	outputDir string // Directory path where Excel files will be saved
}

func NewExcelManager(outputDir string) ExcelManager {
	return &ExcelManagerImpl{
		outputDir: outputDir,
	}
}

// generateYearlyFilePath creates the year-specific filepath for tax summary
func (e *ExcelManagerImpl) generateYearlyFilePath(year int) string {
	filename := fmt.Sprintf("tax_summary_%d.xlsx", year)
	return filepath.Join(e.outputDir, filename)
}

func (e *ExcelManagerImpl) GenerateTaxSummaryExcel(ctx context.Context, year int, summary tax.Summary) (err error) {
	outputFilePath := e.generateYearlyFilePath(year)

	if err = e.ensureDirectoryExists(ctx, outputFilePath); err != nil {
		return
	}

	f := excelize.NewFile()
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Ctx(ctx).Error().Err(closeErr).Msg("Error closing Excel file")
			if err == nil {
				err = closeErr
			}
		}
	}()

	if err = e.writeSheets(ctx, f, summary); err != nil {
		return
	}

	// Delete the default "Sheet1"
	if err = f.DeleteSheet("Sheet1"); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to delete default sheet")
	}

	if err = f.SaveAs(outputFilePath); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("path", outputFilePath).Msg("Failed to save Excel file")
		return fmt.Errorf("failed to save excel file to %s: %w", outputFilePath, err)
	}

	log.Ctx(ctx).Info().Str("path", outputFilePath).Msg("Excel summary generated successfully.")
	return
}

func (e *ExcelManagerImpl) writeSheets(ctx context.Context, f *excelize.File, summary tax.Summary) (err error) {
	if err = e.writeGainsSheet(ctx, f, summary.INRGains); err != nil {
		return
	}

	if err = e.writeDividendsSheet(ctx, f, summary.INRDividends); err != nil {
		return
	}

	if err = e.writeValuationsSheet(ctx, f, summary.INRValuations); err != nil {
		return
	}

	if err = e.writeInterestSheet(ctx, f, summary.INRInterest); err != nil {
		return
	}
	return
}

// ensureDirectoryExists creates the directory for the output file if it doesn't exist.
func (e *ExcelManagerImpl) ensureDirectoryExists(ctx context.Context, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("directory", dir).Msg("Failed to create output directory")
		return fmt.Errorf("failed to create output directory %s: %w", dir, err)
	}
	return nil
}

// createSheetWithHeaders creates a new sheet and writes the headers.
func (e *ExcelManagerImpl) createSheetWithHeaders(ctx context.Context, f *excelize.File, sheetName string, headers []string) error {
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("sheet", sheetName).Msg("Failed to create sheet")
		return fmt.Errorf("failed to create sheet %s: %w", sheetName, err)
	}
	f.SetActiveSheet(index)

	return e.writeHeaders(f, sheetName, headers)
}

// writeGainsSheet handles the creation and population of the "Gains" sheet.
func (e *ExcelManagerImpl) writeGainsSheet(ctx context.Context, f *excelize.File, gains []tax.INRGains) error {
	sheetName := "Gains"
	headers := []string{
		"Symbol", "BuyDate", "SellDate", "Quantity", "PNL (USD)",
		"Commission (USD)", "Type", "TTDate", "TTRate", "PNL (INR)",
	}
	if err := e.createSheetWithHeaders(ctx, f, sheetName, headers); err != nil {
		return err
	}

	for idx, gainRecord := range gains {
		rowNum := idx + 2 // Data starts from row 2

		// Write data columns (A-I) - all except PNL (INR)
		rowData := []interface{}{
			gainRecord.Symbol,
			gainRecord.BuyDate,  // Assuming string
			gainRecord.SellDate, // Assuming string
			gainRecord.Quantity,
			gainRecord.PNL, // USD PNL (column E)
			gainRecord.Commission,
			gainRecord.Type,
			e.formatDateForExcel(gainRecord.TTDate), // Format date
			gainRecord.TTRate,                       // Column I
			"",                                      // Placeholder for PNL (INR) - will be formula
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
			return err
		}

		// Write formula for PNL (INR) column
		// J: PNL (INR) = E * I
		const pnlINRColumn = 10 // Column J
		formula := fmt.Sprintf("=E%d*I%d", rowNum, rowNum)
		if err := e.writeFormulaCell(f, sheetName, rowNum, pnlINRColumn, formula); err != nil {
			return err
		}
	}

	return nil
}

// writeDividendsSheet handles the creation and population of the "Dividends" sheet.
// It assumes tax.INRDividend has fields: Symbol, Date, Amount, TTDate, TTRate.
func (e *ExcelManagerImpl) writeDividendsSheet(ctx context.Context, f *excelize.File, dividends []tax.INRDividend) error {
	sheetName := "Dividends"
	headers := []string{
		"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)", "TTDate", "TTRate",
		"Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := e.createSheetWithHeaders(ctx, f, sheetName, headers); err != nil {
		return err
	}

	for idx, dividendRecord := range dividends {
		rowNum := idx + 2 // Data starts from row 2
		rowData := []interface{}{
			dividendRecord.Symbol,
			dividendRecord.Date,
			dividendRecord.Amount,
			dividendRecord.Tax,
			dividendRecord.Net,
			e.formatDateForExcel(dividendRecord.TTDate),
			dividendRecord.TTRate,
			"", // Placeholder for Amount (INR) - will be formula
			"", // Placeholder for Tax (INR) - will be formula
			"", // Placeholder for Net (INR) - will be formula
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
			return err
		}

		// Write formulas for INR columns (Amount, Tax, Net all converted to INR)
		if err := e.writeTaxWithheldINRFormulas(f, sheetName, rowNum); err != nil {
			return err
		}
	}

	return nil
}

// writeValuationsSheet handles the creation and population of the "Valuations" sheet.
func (e *ExcelManagerImpl) writeValuationsSheet(ctx context.Context, f *excelize.File, valuations []tax.INRValuation) error {
	sheetName := "Valuations"
	if err := e.createSheetWithHeaders(ctx, f, sheetName, e.getValuationHeaders()); err != nil {
		return err
	}

	for idx, valuationRecord := range valuations {
		rowNum := idx + 2 // Data starts from row 2
		rowData := e.buildValuationRow(valuationRecord)
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
			return err
		}

		// Write formulas for all 3 positions (6 formulas total)
		// Each position: ValUSD = Qty * Price, ValINR = ValUSD * TTRate
		formulas := map[int]string{
			// First Position: Columns E (ValUSD), H (ValINR)
			5: fmt.Sprintf("=C%d*D%d", rowNum, rowNum), // E: ValUSD = Qty * Price
			8: fmt.Sprintf("=E%d*G%d", rowNum, rowNum), // H: ValINR = ValUSD * TTRate (uses E!)

			// Peak Position: Columns L (ValUSD), O (ValINR)
			12: fmt.Sprintf("=J%d*K%d", rowNum, rowNum), // L: ValUSD = Qty * Price
			15: fmt.Sprintf("=L%d*N%d", rowNum, rowNum), // O: ValINR = ValUSD * TTRate (uses L!)

			// YearEnd Position: Columns S (ValUSD), V (ValINR)
			19: fmt.Sprintf("=Q%d*R%d", rowNum, rowNum), // S: ValUSD = Qty * Price
			22: fmt.Sprintf("=S%d*U%d", rowNum, rowNum), // V: ValINR = ValUSD * TTRate (uses S!)
		}
		if err := e.writeFormulaRange(f, sheetName, rowNum, formulas); err != nil {
			return err
		}
	}

	return nil
}

func (e *ExcelManagerImpl) getValuationHeaders() []string {
	return []string{
		"Symbol",
		"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
		"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
		"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
		"AmountPaid (INR)",
	}
}

func (e *ExcelManagerImpl) buildValuationRow(valuation tax.INRValuation) []interface{} {
	rowData := []interface{}{valuation.Ticker}
	rowData = append(rowData, e.getPositionRowData(&valuation.FirstPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.PeakPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.YearEndPosition)...)
	rowData = append(rowData, valuation.AmountPaid)
	return rowData
}

func (e *ExcelManagerImpl) getPositionRowData(pos *tax.INRPosition) []interface{} {
	return []interface{}{
		e.formatDateForExcel(pos.Date),
		pos.Quantity,
		pos.RoundedUSDPrice(),
		"", // Placeholder for ValUSD - will be formula: Qty * Price
		e.formatDateForExcel(pos.TTDate),
		pos.TTRate,
		"", // Placeholder for ValINR - will be formula: ValUSD * TTRate
	}
}

// writeInterestSheet handles the creation and population of the "Interest" sheet.
func (e *ExcelManagerImpl) writeInterestSheet(ctx context.Context, f *excelize.File, interest []tax.INRInterest) error {
	sheetName := "Interest"
	headers := []string{
		"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
		"TTDate", "TTRate", "Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := e.createSheetWithHeaders(ctx, f, sheetName, headers); err != nil {
		return err
	}

	for idx, interestRecord := range interest {
		rowNum := idx + 2 // Data starts from row 2
		rowData := []interface{}{
			interestRecord.Symbol,
			interestRecord.Date,
			interestRecord.Amount,
			interestRecord.Tax,
			interestRecord.Net,
			e.formatDateForExcel(interestRecord.TTDate),
			interestRecord.TTRate,
			"", // Placeholder for Amount (INR) - will be formula
			"", // Placeholder for Tax (INR) - will be formula
			"", // Placeholder for Net (INR) - will be formula
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
			return err
		}

		// Write formulas for INR columns (Amount, Tax, Net all converted to INR)
		if err := e.writeTaxWithheldINRFormulas(f, sheetName, rowNum); err != nil {
			return err
		}
	}

	return nil
}

// writeHeaders writes the header row and applies basic styling.
func (e *ExcelManagerImpl) writeHeaders(f *excelize.File, sheetName string, headers []string) error {
	for i, header := range headers {
		cellName, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return fmt.Errorf("failed to get cell name for header col %d: %w", i+1, err)
		}
		if err := f.SetCellValue(sheetName, cellName, header); err != nil {
			return fmt.Errorf("failed to set header '%s': %w", header, err)
		}
	}
	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}
	if err := f.SetRowStyle(sheetName, 1, 1, style); err != nil {
		log.Warn().Err(err).Msg("Failed to set header row style")
	}
	return nil
}

// writeRow writes a slice of interface{} data to a specific row.
func (e *ExcelManagerImpl) writeRow(f *excelize.File, sheetName string, rowNum int, data []interface{}) error {
	for i, val := range data {
		cellName, err := excelize.CoordinatesToCellName(i+1, rowNum)
		if err != nil {
			return fmt.Errorf("failed to get cell name for data col %d, row %d: %w", i+1, rowNum, err)
		}
		if err := f.SetCellValue(sheetName, cellName, val); err != nil {
			return fmt.Errorf("failed to set cell value at %s: %w", cellName, err)
		}
	}
	return nil
}

// formatDateForExcel formats a time.Time for Excel.
// If time.Time is zero, return an empty string.
func (e *ExcelManagerImpl) formatDateForExcel(t time.Time) interface{} {
	if t.IsZero() {
		return "" // Return empty string for zero/uninitialized dates
	}
	return t.Format(time.DateOnly) // "2006-01-02"
}

// writeFormulaCell writes a formula to a specific cell
func (e *ExcelManagerImpl) writeFormulaCell(f *excelize.File, sheetName string, rowNum, colNum int, formula string) error {
	cellName, err := excelize.CoordinatesToCellName(colNum, rowNum)
	if err != nil {
		return fmt.Errorf("failed to get cell name for col %d, row %d: %w", colNum, rowNum, err)
	}

	if err := f.SetCellFormula(sheetName, cellName, formula); err != nil {
		return fmt.Errorf("failed to set formula at %s: %w", cellName, err)
	}
	return nil
}

// writeFormulaRange writes multiple formulas for a single row
// formulas is a map of columnIndex -> formulaString
func (e *ExcelManagerImpl) writeFormulaRange(f *excelize.File, sheetName string, rowNum int, formulas map[int]string) error {
	for colNum, formula := range formulas {
		if err := e.writeFormulaCell(f, sheetName, rowNum, colNum, formula); err != nil {
			return err
		}
	}
	return nil
}

// writeTaxWithheldINRFormulas writes three formulas for tax-withheld items (Dividends/Interest)
// Formulas:
//   - Column H: Amount(INR) = Amount(USD) * TTRate
//   - Column I: Tax(INR) = Tax(USD) * TTRate
//   - Column J: Net(INR) = Net(USD) * TTRate
func (e *ExcelManagerImpl) writeTaxWithheldINRFormulas(f *excelize.File, sheetName string, rowNum int) error {
	formulas := map[int]string{
		8:  fmt.Sprintf("=C%d*G%d", rowNum, rowNum), // Column H: Amount(INR) = Amount(USD) * TTRate
		9:  fmt.Sprintf("=D%d*G%d", rowNum, rowNum), // Column I: Tax(INR) = Tax(USD) * TTRate
		10: fmt.Sprintf("=E%d*G%d", rowNum, rowNum), // Column J: Net(INR) = Net(USD) * TTRate
	}
	return e.writeFormulaRange(f, sheetName, rowNum, formulas)
}
