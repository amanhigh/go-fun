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

	// Move Summary sheet to the first position
	e.moveSummarySheetToFirst(f)

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

	// Write Summary sheet LAST (after all detail sheets exist)
	return e.writeSummarySheet(ctx, f, summary)
}

// moveSummarySheetToFirst moves the Summary sheet to the first position in the workbook
func (e *ExcelManagerImpl) moveSummarySheetToFirst(f *excelize.File) {
	summaryIndex, getErr := f.GetSheetIndex("Summary")
	if getErr == nil && summaryIndex > 0 {
		_ = f.MoveSheet("Summary", "Gains")
	}
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

	// Write totals section with STCG/LTCG breakdown
	if len(gains) > 0 {
		lastDataRow := len(gains) + 1
		if err := e.writeGainsTotals(f, sheetName, lastDataRow); err != nil {
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

	// Write totals section
	if len(dividends) > 0 {
		lastDataRow := len(dividends) + 1
		if err := e.writeDividendsTotals(f, sheetName, lastDataRow); err != nil {
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

	// Write totals section - ONLY AmountPaid (INR) in column W
	if len(valuations) > 0 {
		lastDataRow := len(valuations) + 1
		if err := e.writeValuationsTotals(f, sheetName, lastDataRow); err != nil {
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

	// Write totals section
	if len(interest) > 0 {
		lastDataRow := len(interest) + 1
		if err := e.writeInterestTotals(f, sheetName, lastDataRow); err != nil {
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

// writeTotalsLabel writes "TOTALS", "STCG", "LTCG" etc. to column A with bold styling
func (e *ExcelManagerImpl) writeTotalsLabel(f *excelize.File, sheetName string, rowNum int, label string) error {
	cellName, err := excelize.CoordinatesToCellName(1, rowNum) // Column A
	if err != nil {
		return fmt.Errorf("failed to get cell name for label row %d: %w", rowNum, err)
	}
	if err = f.SetCellValue(sheetName, cellName, label); err != nil {
		return fmt.Errorf("failed to set label '%s': %w", label, err)
	}

	// Apply bold style to label
	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return fmt.Errorf("failed to create bold style: %w", err)
	}
	if err = f.SetCellStyle(sheetName, cellName, cellName, style); err != nil {
		return fmt.Errorf("failed to set style for label: %w", err)
	}
	return nil
}

// writeSimpleTotals writes a single TOTALS row with SUM formulas for specified columns
// columns: map of columnIndex -> columnLetter for SUM formulas (e.g., {3: "C", 4: "D"})
// This generic helper reduces duplication across Dividends, Interest, and Valuations sheets
func (e *ExcelManagerImpl) writeSimpleTotals(
	f *excelize.File,
	sheetName string,
	lastDataRow int,
	columns map[int]string,
) error {
	totalsRow := lastDataRow + 2

	// Write TOTALS label in column A
	if err := e.writeTotalsLabel(f, sheetName, totalsRow, "TOTALS"); err != nil {
		return err
	}

	// Build SUM formulas for each specified column
	formulas := make(map[int]string)
	for colIdx, colLetter := range columns {
		formulas[colIdx] = fmt.Sprintf("=SUM(%s2:%s%d)", colLetter, colLetter, lastDataRow)
	}

	return e.writeFormulaRange(f, sheetName, totalsRow, formulas)
}

// writeGainsTotals writes TOTALS row and STCG/LTCG breakdown for Gains sheet
// Columns: E (PNL USD), F (Commission USD), J (PNL INR)
// STCG/LTCG use SUMIF based on Type column (G)
func (e *ExcelManagerImpl) writeGainsTotals(f *excelize.File, sheetName string, lastDataRow int) error {
	// Calculate row numbers
	totalsRow := lastDataRow + 2
	stcgRow := totalsRow + 2
	ltcgRow := stcgRow + 1

	// Write TOTALS row
	if err := e.writeTotalsLabel(f, sheetName, totalsRow, "TOTALS"); err != nil {
		return err
	}

	// Write TOTALS formulas for columns E, F, J
	totalsFormulas := map[int]string{
		5:  fmt.Sprintf("=SUM(E2:E%d)", lastDataRow), // Column E: PNL USD
		6:  fmt.Sprintf("=SUM(F2:F%d)", lastDataRow), // Column F: Commission USD
		10: fmt.Sprintf("=SUM(J2:J%d)", lastDataRow), // Column J: PNL INR
	}
	if err := e.writeFormulaRange(f, sheetName, totalsRow, totalsFormulas); err != nil {
		return err
	}

	// Write STCG row
	if err := e.writeTotalsLabel(f, sheetName, stcgRow, "STCG"); err != nil {
		return err
	}

	// Write STCG formulas (SUMIF for Type="STCG")
	stcgFormulas := map[int]string{
		5:  fmt.Sprintf("=SUMIF(G2:G%d,\"STCG\",E2:E%d)", lastDataRow, lastDataRow), // Column E: PNL USD
		10: fmt.Sprintf("=SUMIF(G2:G%d,\"STCG\",J2:J%d)", lastDataRow, lastDataRow), // Column J: PNL INR
	}
	if err := e.writeFormulaRange(f, sheetName, stcgRow, stcgFormulas); err != nil {
		return err
	}

	// Write LTCG row
	if err := e.writeTotalsLabel(f, sheetName, ltcgRow, "LTCG"); err != nil {
		return err
	}

	// Write LTCG formulas (SUMIF for Type="LTCG")
	ltcgFormulas := map[int]string{
		5:  fmt.Sprintf("=SUMIF(G2:G%d,\"LTCG\",E2:E%d)", lastDataRow, lastDataRow), // Column E: PNL USD
		10: fmt.Sprintf("=SUMIF(G2:G%d,\"LTCG\",J2:J%d)", lastDataRow, lastDataRow), // Column J: PNL INR
	}
	return e.writeFormulaRange(f, sheetName, ltcgRow, ltcgFormulas)
}

// writeDividendsTotals writes TOTALS row for Dividends sheet
// Columns: C (Amount USD), D (Tax USD), E (Net USD), H (Amount INR), I (Tax INR), J (Net INR)
func (e *ExcelManagerImpl) writeDividendsTotals(f *excelize.File, sheetName string, lastDataRow int) error {
	return e.writeSimpleTotals(f, sheetName, lastDataRow, map[int]string{
		3: "C", 4: "D", 5: "E", // USD columns: Amount, Tax, Net
		8: "H", 9: "I", 10: "J", // INR columns: Amount, Tax, Net
	})
}

// writeInterestTotals writes TOTALS row for Interest sheet
// Columns: C (Amount USD), D (Tax USD), E (Net USD), H (Amount INR), I (Tax INR), J (Net INR)
func (e *ExcelManagerImpl) writeInterestTotals(f *excelize.File, sheetName string, lastDataRow int) error {
	return e.writeSimpleTotals(f, sheetName, lastDataRow, map[int]string{
		3: "C", 4: "D", 5: "E", // USD columns: Amount, Tax, Net
		8: "H", 9: "I", 10: "J", // INR columns: Amount, Tax, Net
	})
}

// writeValuationsTotals writes TOTALS row for Valuations sheet - ONLY AmountPaid (INR)
// Column: W (AmountPaid INR) - No totals for position valuations as they're not meaningful
func (e *ExcelManagerImpl) writeValuationsTotals(f *excelize.File, sheetName string, lastDataRow int) error {
	return e.writeSimpleTotals(f, sheetName, lastDataRow, map[int]string{
		23: "W", // Column W: AmountPaid (INR)
	})
}

// writeSummarySheet creates the Summary sheet with cross-referenced formulas to detail sheets
// This sheet must be created AFTER all detail sheets to know their TOTALS row positions
func (e *ExcelManagerImpl) writeSummarySheet(ctx context.Context, f *excelize.File, summary tax.Summary) error {
	const dividendsSectionStartRow = 3

	sheetName := "Summary"
	if err := e.createSheetWithHeaders(ctx, f, sheetName, []string{}); err != nil {
		return err
	}

	// Write header (Row 1: "SUMMARY")
	if err := e.writeSummaryHeader(f, sheetName); err != nil {
		return err
	}

	// Calculate Dividends TOTALS row position
	// lastDataRow = len(dividends) + 1
	// totalsRow = lastDataRow + 2
	dividendsDataCount := len(summary.INRDividends)
	var dividendsTotalsRow int
	if dividendsDataCount > 0 {
		dividendsTotalsRow = (dividendsDataCount + 1) + 2
	}

	// Write Dividends section (starts at row 3)
	return e.writeDividendsSection(f, sheetName, dividendsSectionStartRow, dividendsTotalsRow)
}

// writeSummaryHeader writes the title row (Row 1: "SUMMARY" - bold)
func (e *ExcelManagerImpl) writeSummaryHeader(f *excelize.File, sheetName string) error {
	return e.writeTotalsLabel(f, sheetName, 1, "SUMMARY")
}

// writeDividendsSection writes the Dividends section with cross-referenced formulas
// startRow: where to start writing section header (typically row 3)
// dividendsTotalsRow: row number in Dividends sheet where TOTALS row exists (0 if no data)
func (e *ExcelManagerImpl) writeDividendsSection(f *excelize.File, sheetName string, startRow, dividendsTotalsRow int) error {
	// Row startRow: "Dividends" (bold section header)
	if err := e.writeTotalsLabel(f, sheetName, startRow, "Dividends"); err != nil {
		return err
	}

	// Row startRow+1: Column headers (USD first, then INR) - bold
	headerRow := startRow + 1
	headers := []interface{}{
		"Amount (USD)", "Tax (USD)", "Net (USD)",
		"Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := e.writeRow(f, sheetName, headerRow, headers); err != nil {
		return err
	}

	// Apply bold style to headers
	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}
	if err := f.SetRowStyle(sheetName, headerRow, headerRow, style); err != nil {
		return fmt.Errorf("failed to set header row style: %w", err)
	}

	// Row startRow+2: Formulas referencing Dividends TOTALS row
	valuesRow := startRow + 2
	formulas := map[int]string{
		1: fmt.Sprintf("=Dividends!C%d", dividendsTotalsRow), // A: Amount USD
		2: fmt.Sprintf("=Dividends!D%d", dividendsTotalsRow), // B: Tax USD
		3: fmt.Sprintf("=Dividends!E%d", dividendsTotalsRow), // C: Net USD
		4: fmt.Sprintf("=Dividends!H%d", dividendsTotalsRow), // D: Amount INR
		5: fmt.Sprintf("=Dividends!I%d", dividendsTotalsRow), // E: Tax INR
		6: fmt.Sprintf("=Dividends!J%d", dividendsTotalsRow), // F: Net INR
	}

	return e.writeFormulaRange(f, sheetName, valuesRow, formulas)
}
