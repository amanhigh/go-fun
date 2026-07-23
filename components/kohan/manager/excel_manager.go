//nolint:dupl // Excel operations have similar patterns for different entity types
package manager

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

// Column width constants used across all sheets.
// These replace magic number literals to satisfy the mnd linter.
const (
	colWidthNarrow    = 8
	colWidthSemi      = 10
	colWidthMedium    = 12
	colWidthWide      = 14
	colWidthExtraWide = 16
	colWidthMax       = 18

	sheetNameInterest = "Interest"
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

// sortedCopy returns a defensive clone of s sorted stably by cmpFn.
// cmpFn must return -1 if a < b, 0 if a == b, +1 if a > b (cmp.Compare convention).
func sortedCopy[S ~[]E, E any](s S, cmpFn func(a, b E) int) S {
	result := slices.Clone(s)
	slices.SortStableFunc(result, cmpFn)
	return result
}

// --- Local typed comparators (ExcelManager-private) ---

func compareGainsBySellDateSymbol(a, b tax.INRGains) int {
	// SellDate descending (newest first), then Symbol ascending
	if c := cmp.Compare(a.SellDate, b.SellDate); c != 0 {
		return -c
	}
	return cmp.Compare(a.Symbol, b.Symbol)
}

func compareDividendsByDateSymbol(a, b tax.INRDividend) int {
	// Date descending (newest first), then Symbol ascending
	if c := cmp.Compare(a.Date, b.Date); c != 0 {
		return -c
	}
	return cmp.Compare(a.Symbol, b.Symbol)
}

func compareInterestByDateSymbol(a, b tax.INRInterest) int {
	// Date descending (newest first), then Symbol ascending
	if c := cmp.Compare(a.Date, b.Date); c != 0 {
		return -c
	}
	return cmp.Compare(a.Symbol, b.Symbol)
}

func compareValuationsByTicker(a, b tax.INRValuation) int {
	return cmp.Compare(a.Ticker, b.Ticker)
}

func compareTTRatesByActualDate(a, b tax.MonthEndRate) int {
	// ActualDate descending (newest first)
	if a.ActualDate.Before(b.ActualDate) {
		return 1
	}
	if a.ActualDate.After(b.ActualDate) {
		return -1
	}
	return 0
}

// generateYearlyFilePath creates the year-specific filepath for tax summary
func (e *ExcelManagerImpl) generateYearlyFilePath(year int) string {
	filename := fmt.Sprintf("%d_Tax_Summary.xlsx", year)
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

// FIXME: Add a Security Info sheet to the final workbook with resolved security metadata.
func (e *ExcelManagerImpl) writeSheets(ctx context.Context, f *excelize.File, summary tax.Summary) (err error) {
	// Defensive sorted copies preserve caller's original slices
	gains := sortedCopy(summary.INRGains, compareGainsBySellDateSymbol)
	dividends := sortedCopy(summary.INRDividends, compareDividendsByDateSymbol)
	valuations := sortedCopy(summary.INRValuations, compareValuationsByTicker)
	interest := sortedCopy(summary.INRInterest, compareInterestByDateSymbol)
	rates := sortedCopy(summary.TTMonthEndRates, compareTTRatesByActualDate)

	if err = e.writeGainsSheet(ctx, f, gains); err != nil {
		return
	}

	if err = e.writeDividendsSheet(ctx, f, dividends); err != nil {
		return
	}

	if err = e.writeValuationsSheet(ctx, f, valuations); err != nil {
		return
	}

	if err = e.writeInterestSheet(ctx, f, interest); err != nil {
		return
	}

	// Write TT Rates sheet with FY month-end data
	if err = e.writeTTRatesSheet(ctx, f, summary.Year, rates); err != nil {
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
		rowData := []any{
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
		if err := util.WriteRow(f, sheetName, rowNum, rowData); err != nil {
			return fmt.Errorf("failed to write row %d in sheet %s: %w", rowNum, sheetName, err)
		}

		// Write formula for PNL (INR) column
		// J: PNL (INR) = E * I
		const pnlINRColumn = 10 // Column J
		formula := fmt.Sprintf("=E%d*I%d", rowNum, rowNum)
		if err := util.WriteFormulaCell(f, sheetName, rowNum, pnlINRColumn, formula); err != nil {
			return fmt.Errorf("failed to write formula cell at col %d row %d in sheet %s: %w", pnlINRColumn, rowNum, sheetName, err)
		}
	}

	// Write totals section with STCG/LTCG breakdown
	if len(gains) > 0 {
		lastDataRow := len(gains) + 1
		if err := e.writeGainsTotals(f, sheetName, lastDataRow); err != nil {
			return err
		}
	}

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthNarrow, "B": colWidthMedium, "C": colWidthMedium, "D": colWidthSemi, "E": colWidthMedium,
		"F": colWidthExtraWide, "G": colWidthNarrow, "H": colWidthMedium, "I": colWidthSemi, "J": colWidthMedium,
	})
	if err := util.ApplyAutoFilter(f, sheetName, "J", len(gains)+1); err != nil {
		return fmt.Errorf("failed to apply auto filter on sheet %s: %w", sheetName, err)
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
		rowData := []any{
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
		if err := util.WriteRow(f, sheetName, rowNum, rowData); err != nil {
			return fmt.Errorf("failed to write row %d in sheet %s: %w", rowNum, sheetName, err)
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

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthNarrow, "B": colWidthMedium, "C": colWidthWide, "D": colWidthMedium, "E": colWidthMedium,
		"F": colWidthMedium, "G": colWidthSemi, "H": colWidthWide, "I": colWidthMedium, "J": colWidthMedium,
	})
	if err := util.ApplyAutoFilter(f, sheetName, "J", len(dividends)+1); err != nil {
		return fmt.Errorf("failed to apply auto filter on sheet %s: %w", sheetName, err)
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
		if err := util.WriteRow(f, sheetName, rowNum, rowData); err != nil {
			return fmt.Errorf("failed to write row %d in sheet %s: %w", rowNum, sheetName, err)
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
		if err := util.WriteFormulaRange(f, sheetName, rowNum, formulas); err != nil {
			return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", rowNum, sheetName, err)
		}
	}

	// Write totals section - ONLY AmountPaid (INR) in column W
	if len(valuations) > 0 {
		lastDataRow := len(valuations) + 1
		if err := e.writeValuationsTotals(f, sheetName, lastDataRow); err != nil {
			return err
		}
	}

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthNarrow,
		"B": colWidthWide, "C": colWidthNarrow, "D": colWidthSemi, "E": colWidthSemi, "F": colWidthMedium, "G": colWidthSemi, "H": colWidthSemi,
		"I": colWidthWide, "J": colWidthNarrow, "K": colWidthSemi, "L": colWidthSemi, "M": colWidthMedium, "N": colWidthSemi, "O": colWidthSemi,
		"P": colWidthExtraWide, "Q": colWidthNarrow, "R": colWidthSemi, "S": colWidthSemi, "T": colWidthMedium, "U": colWidthSemi, "V": colWidthSemi,
		"W": colWidthExtraWide,
	})
	if err := util.ApplyAutoFilter(f, sheetName, "W", len(valuations)+1); err != nil {
		return fmt.Errorf("failed to apply auto filter on sheet %s: %w", sheetName, err)
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

func (e *ExcelManagerImpl) buildValuationRow(valuation tax.INRValuation) []any {
	rowData := []any{valuation.Ticker}
	rowData = append(rowData, e.getPositionRowData(&valuation.FirstPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.PeakPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.YearEndPosition)...)
	rowData = append(rowData, valuation.AmountPaid)
	return rowData
}

func (e *ExcelManagerImpl) getPositionRowData(pos *tax.INRPosition) []any {
	return []any{
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
	sheetName := sheetNameInterest
	headers := []string{
		"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
		"TTDate", "TTRate", "Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := e.createSheetWithHeaders(ctx, f, sheetName, headers); err != nil {
		return err
	}

	for idx, interestRecord := range interest {
		rowNum := idx + 2 // Data starts from row 2
		rowData := []any{
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
		if err := util.WriteRow(f, sheetName, rowNum, rowData); err != nil {
			return fmt.Errorf("failed to write row %d in sheet %s: %w", rowNum, sheetName, err)
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

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthNarrow, "B": colWidthMedium, "C": colWidthWide, "D": colWidthMedium, "E": colWidthMedium,
		"F": colWidthMedium, "G": colWidthSemi, "H": colWidthWide, "I": colWidthMedium, "J": colWidthMedium,
	})
	if err := util.ApplyAutoFilter(f, sheetName, "J", len(interest)+1); err != nil {
		return fmt.Errorf("failed to apply auto filter on sheet %s: %w", sheetName, err)
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

// formatDateForExcel formats a time.Time for Excel.
// If time.Time is zero, return an empty string.
func (e *ExcelManagerImpl) formatDateForExcel(t time.Time) any {
	if t.IsZero() {
		return "" // Return empty string for zero/uninitialized dates
	}
	return t.Format(time.DateOnly) // "2006-01-02"
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
	if err := util.WriteFormulaRange(f, sheetName, rowNum, formulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", rowNum, sheetName, err)
	}
	return nil
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

	if err := util.WriteFormulaRange(f, sheetName, totalsRow, formulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", totalsRow, sheetName, err)
	}
	return nil
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
	if err := util.WriteFormulaRange(f, sheetName, totalsRow, totalsFormulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", totalsRow, sheetName, err)
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
	if err := util.WriteFormulaRange(f, sheetName, stcgRow, stcgFormulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", stcgRow, sheetName, err)
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
	if err := util.WriteFormulaRange(f, sheetName, ltcgRow, ltcgFormulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", ltcgRow, sheetName, err)
	}
	return nil
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

// writeValuationsTotals writes TOTALS row for Valuations sheet
// Columns: V (YearEnd ValINR), W (AmountPaid INR)
func (e *ExcelManagerImpl) writeValuationsTotals(f *excelize.File, sheetName string, lastDataRow int) error {
	return e.writeSimpleTotals(f, sheetName, lastDataRow, map[int]string{
		22: "V", // Column V: YearEnd ValINR = Qty * Price * TTRate
		23: "W", // Column W: AmountPaid (INR)
	})
}

// setColumnWidths sets custom column widths for the given sheet.
// widths is a map of column letter to desired width.
func (e *ExcelManagerImpl) setColumnWidths(f *excelize.File, sheetName string, widths map[string]float64) {
	for col, width := range widths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			log.Warn().Err(err).Str("sheet", sheetName).Str("col", col).Msg("Failed to set column width")
		}
	}
}

// writeTTRatesSheet creates the "TT Rates" sheet with FY month-end rate reference data.
// Rows are sorted by ActualDate descending (latest rates first). Month/Year labels are derived
// from each rate's ActualDate.AddDate(0, 1, 0), so the applicable month follows the SBI rate date.
// Each row shows the actual rate date, the TT buy rate, a clickable PDF link when available,
// and the day of the week.
func (e *ExcelManagerImpl) writeTTRatesSheet(ctx context.Context, f *excelize.File, year int, rates []tax.MonthEndRate) error {
	sheetName := "TT Rates"
	headers := []string{"Month", "Year", "TTDate", "TTRate", "PDF Link", "DayOfWeek"}
	if err := e.createSheetWithHeaders(ctx, f, sheetName, headers); err != nil {
		return err
	}

	for idx, rate := range rates {
		rowNum := idx + 2 // Data starts from row 2

		// Derive Month/Year label from ActualDate's next month.
		// This ensures the latest rate correctly reflects its applicable month
		// (e.g., a Feb 20 rate is labeled MAR).
		// Fall back to index-based Apr→Mar order for zero ActualDate values.
		var labelDate time.Time
		if rate.ActualDate.IsZero() {
			labelDate = time.Date(year, time.April, 1, 0, 0, 0, 0, time.UTC).AddDate(0, idx, 0)
		} else {
			labelDate = rate.ActualDate.AddDate(0, 1, 0)
		}
		monthLabel := strings.ToUpper(labelDate.Format("Jan"))
		fyYear := labelDate.Year()

		// Compute day of week from the actual SBI rate date
		dayOfWeek := rate.ActualDate.Format("Monday")

		rowData := []any{
			monthLabel,
			fyYear,
			e.formatDateForExcel(rate.ActualDate),
			rate.Rate,
			"", // Placeholder for PDF Link - written with hyperlink below
			dayOfWeek,
		}
		if err := util.WriteRow(f, sheetName, rowNum, rowData); err != nil {
			return fmt.Errorf("failed to write row %d in sheet %s: %w", rowNum, sheetName, err)
		}

		// Write PDF Link in column E with clickable hyperlink if URL is present
		if err := e.writePDFLink(f, sheetName, rowNum, rate.PDFFile); err != nil {
			return err
		}
	}

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthNarrow, "B": colWidthNarrow, "C": colWidthWide, "D": colWidthSemi, "E": colWidthMedium, "F": colWidthWide,
	})
	if err := util.ApplyAutoFilter(f, sheetName, "F", len(rates)+1); err != nil {
		return fmt.Errorf("failed to apply auto filter on sheet %s: %w", sheetName, err)
	}
	return nil
}

// writePDFLink writes the PDF Link cell value and sets an external hyperlink if the
// pdfFile string is an HTTP/HTTPS URL. Non-URL values are written as-is.
func (e *ExcelManagerImpl) writePDFLink(f *excelize.File, sheetName string, rowNum int, pdfFile string) error {
	const pdfLinkColumn = 5 // Column E

	cellName, err := excelize.CoordinatesToCellName(pdfLinkColumn, rowNum)
	if err != nil {
		return fmt.Errorf("failed to get cell name for PDF link at row %d: %w", rowNum, err)
	}

	if strings.HasPrefix(pdfFile, "http://") || strings.HasPrefix(pdfFile, "https://") {
		display := "PDF"
		if err := f.SetCellValue(sheetName, cellName, display); err != nil {
			return fmt.Errorf("failed to set PDF link display at %s: %w", cellName, err)
		}
		if err := f.SetCellHyperLink(sheetName, cellName, pdfFile, "External", excelize.HyperlinkOpts{
			Display: &display,
		}); err != nil {
			return fmt.Errorf("failed to set hyperlink at %s: %w", cellName, err)
		}
	} else {
		if err := f.SetCellValue(sheetName, cellName, pdfFile); err != nil {
			return fmt.Errorf("failed to set PDF link value at %s: %w", cellName, err)
		}
	}

	return nil
}

// writeSummarySheet creates the Summary sheet with cross-referenced formulas to detail sheets
// This sheet must be created AFTER all detail sheets to know their TOTALS row positions
func (e *ExcelManagerImpl) writeSummarySheet(ctx context.Context, f *excelize.File, summary tax.Summary) error {
	sheetName := "Summary"
	if err := e.createSheetWithHeaders(ctx, f, sheetName, []string{}); err != nil {
		return err
	}

	// Write header (Row 1: "SUMMARY")
	if err := e.writeSummaryHeader(f, sheetName); err != nil {
		return err
	}

	return e.writeSummarySections(f, sheetName, summary)
}

// writeSummarySections writes Gains, Dividends, and Interest sections into the Summary sheet.
// This is extracted from writeSummarySheet to keep statement count within the funlen limit.
func (e *ExcelManagerImpl) writeSummarySections(f *excelize.File, sheetName string, summary tax.Summary) error {
	// Calculate TOTALS row positions for all sheets (re-calculated from summary)
	gainsTotalsRow, gainsSTCGRow, gainsLTCGRow := e.calculateGainsRows(summary.INRGains)
	dividendsTotalsRow := e.calculateTotalsRow(summary.INRDividends)
	interestTotalsRow := e.calculateTotalsRow(summary.INRInterest)
	currentRow := 3 // Start after header and empty row

	// Gains section (Short Term + Long Term) - only if data exists
	if len(summary.INRGains) > 0 {
		if err := e.writeGainsSection(f, sheetName, currentRow, gainsSTCGRow, gainsLTCGRow, gainsTotalsRow); err != nil {
			return err
		}
		currentRow += 10 // Short Term (4 rows) + Long Term (4 rows) + 2 empty
	}

	// Dividends section - only if data exists
	if len(summary.INRDividends) > 0 {
		if err := e.writeDividendsSection(f, sheetName, currentRow, dividendsTotalsRow); err != nil {
			return err
		}
		currentRow += 5 // Section header + headers + values + 2 empty
	}

	// Interest section - only if data exists
	if len(summary.INRInterest) > 0 {
		if err := e.writeInterestSection(f, sheetName, currentRow, interestTotalsRow); err != nil {
			return err
		}
	}

	e.setColumnWidths(f, sheetName, map[string]float64{
		"A": colWidthMax, "B": colWidthMax, "C": colWidthMax, "D": colWidthMax, "E": colWidthMax, "F": colWidthMax,
	})
	return nil
}

// calculateGainsRows computes TOTALS, STCG, and LTCG row positions for Gains sheet
func (e *ExcelManagerImpl) calculateGainsRows(gains []tax.INRGains) (int, int, int) {
	if len(gains) == 0 {
		return 0, 0, 0
	}
	totalsRow := (len(gains) + 1) + 2
	stcgRow := totalsRow + 2
	ltcgRow := stcgRow + 1
	return totalsRow, stcgRow, ltcgRow
}

// calculateTotalsRow computes the TOTALS row position for a data slice
func (e *ExcelManagerImpl) calculateTotalsRow(data any) int {
	switch v := data.(type) {
	case []tax.INRDividend:
		if len(v) == 0 {
			return 0
		}
		return (len(v) + 1) + 2
	case []tax.INRInterest:
		if len(v) == 0 {
			return 0
		}
		return (len(v) + 1) + 2
	}
	return 0
}

// writeSummaryHeader writes the title row (Row 1: "SUMMARY" - bold)
func (e *ExcelManagerImpl) writeSummaryHeader(f *excelize.File, sheetName string) error {
	return e.writeTotalsLabel(f, sheetName, 1, "SUMMARY")
}

// applyBoldStyle applies bold font style to an entire row
func (e *ExcelManagerImpl) applyBoldStyle(f *excelize.File, sheetName string, rowNum int) error {
	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return fmt.Errorf("failed to create bold style: %w", err)
	}
	if err := f.SetRowStyle(sheetName, rowNum, rowNum, style); err != nil {
		return fmt.Errorf("failed to set bold row style: %w", err)
	}
	return nil
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
	headers := []any{
		"Amount (USD)", "Tax (USD)", "Net (USD)",
		"Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := util.WriteRow(f, sheetName, headerRow, headers); err != nil {
		return fmt.Errorf("failed to write row %d in sheet %s: %w", headerRow, sheetName, err)
	}

	// Apply bold style to headers
	if err := e.applyBoldStyle(f, sheetName, headerRow); err != nil {
		return err
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

	if err := util.WriteFormulaRange(f, sheetName, valuesRow, formulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", valuesRow, sheetName, err)
	}
	return nil
}

// writeGainsSection writes the Gains section with Short Term and Long Term subsections
// startRow: where to start writing (typically row 3)
// gainsSTCGRow, gainsLTCGRow: row numbers in Gains sheet for STCG/LTCG breakdowns (0 if no data)
// gainsTotalsRow: row number in Gains sheet for TOTALS row (0 if no data)
func (e *ExcelManagerImpl) writeGainsSection(f *excelize.File, sheetName string, startRow, gainsSTCGRow, gainsLTCGRow, gainsTotalsRow int) error {
	const ltcgSectionRowOffset = 5

	// Write Short Term section
	if err := e.writeGainsCategorySection(f, sheetName, startRow, "Short Term", gainsSTCGRow, gainsTotalsRow); err != nil {
		return err
	}

	// Write Long Term section
	return e.writeGainsCategorySection(f, sheetName, startRow+ltcgSectionRowOffset, "Long Term", gainsLTCGRow, gainsTotalsRow)
}

// writeGainsCategorySection writes either Short Term or Long Term subsection
func (e *ExcelManagerImpl) writeGainsCategorySection(f *excelize.File, sheetName string, startRow int, category string, gainsRow, totalsRow int) error {
	headers := []any{"PNL (USD)", "Commission (USD)", "PNL (INR)"}

	// Row startRow: Category header (bold section header)
	if err := e.writeTotalsLabel(f, sheetName, startRow, category); err != nil {
		return err
	}

	// Row startRow+1: Column headers (bold)
	headerRow := startRow + 1
	if err := util.WriteRow(f, sheetName, headerRow, headers); err != nil {
		return fmt.Errorf("failed to write row %d in sheet %s: %w", headerRow, sheetName, err)
	}
	if err := e.applyBoldStyle(f, sheetName, headerRow); err != nil {
		return err
	}

	// Row startRow+2: Formulas
	valuesRow := startRow + 2
	formulas := map[int]string{
		1: fmt.Sprintf("=Gains!E%d", gainsRow),  // A: PNL USD
		2: fmt.Sprintf("=Gains!F%d", totalsRow), // B: Commission from TOTALS
		3: fmt.Sprintf("=Gains!J%d", gainsRow),  // C: PNL INR
	}
	if err := util.WriteFormulaRange(f, sheetName, valuesRow, formulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", valuesRow, sheetName, err)
	}
	return nil
}

// writeInterestSection writes the Interest Income section with cross-referenced formulas
// startRow: where to start writing section header
// interestTotalsRow: row number in Interest sheet where TOTALS row exists (0 if no data)
func (e *ExcelManagerImpl) writeInterestSection(f *excelize.File, sheetName string, startRow, interestTotalsRow int) error {
	// Row startRow: "Interest Income" (bold section header)
	if err := e.writeTotalsLabel(f, sheetName, startRow, "Interest Income"); err != nil {
		return err
	}

	// Row startRow+1: Column headers (USD first, then INR) - bold
	headerRow := startRow + 1
	headers := []any{
		"Amount (USD)", "Tax (USD)", "Net (USD)",
		"Amount (INR)", "Tax (INR)", "Net (INR)",
	}
	if err := util.WriteRow(f, sheetName, headerRow, headers); err != nil {
		return fmt.Errorf("failed to write row %d in sheet %s: %w", headerRow, sheetName, err)
	}
	if err := e.applyBoldStyle(f, sheetName, headerRow); err != nil {
		return err
	}

	// Row startRow+2: Formulas referencing Interest TOTALS row
	valuesRow := startRow + 2
	formulas := map[int]string{
		1: fmt.Sprintf("=Interest!C%d", interestTotalsRow), // A: Amount USD
		2: fmt.Sprintf("=Interest!D%d", interestTotalsRow), // B: Tax USD
		3: fmt.Sprintf("=Interest!E%d", interestTotalsRow), // C: Net USD
		4: fmt.Sprintf("=Interest!H%d", interestTotalsRow), // D: Amount INR
		5: fmt.Sprintf("=Interest!I%d", interestTotalsRow), // E: Tax INR
		6: fmt.Sprintf("=Interest!J%d", interestTotalsRow), // F: Net INR
	}

	if err := util.WriteFormulaRange(f, sheetName, valuesRow, formulas); err != nil {
		return fmt.Errorf("failed to write formula range at row %d in sheet %s: %w", valuesRow, sheetName, err)
	}
	return nil
}
