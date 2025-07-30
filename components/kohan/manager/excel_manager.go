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

//go:generate mockery --name ExcelManager
type ExcelManager interface {
	GenerateTaxSummaryExcel(ctx context.Context, summary tax.Summary) error
}

type ExcelManagerImpl struct {
	outputFilePath string // Full exact path, injected via constructor
}

func NewExcelManager(outputFilePath string) ExcelManager {
	return &ExcelManagerImpl{
		outputFilePath: outputFilePath,
	}
}

// In ExcelManagerImpl

func (e *ExcelManagerImpl) GenerateTaxSummaryExcel(ctx context.Context, summary tax.Summary) error {
	if err := e.ensureDirectoryExists(ctx); err != nil {
		return err
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error closing Excel file")
		}
	}()

	// --- Gains Sheet (MVP) ---
	if err := e.writeGainsSheet(ctx, f, summary.INRGains); err != nil {
		return err // Error already logged in helper
	}

	// --- Future Sheets (Example structure) ---
	if err := e.writeDividendsSheet(ctx, f, summary.INRDividends); err != nil {
		return err
	}

	// --- Valuations Sheet ---
	if err := e.writeValuationsSheet(ctx, f, summary.INRValuations); err != nil {
		return err
	}

	// --- Interest Sheet ---
	if err := e.writeInterestSheet(ctx, f, summary.INRInterest); err != nil {
		return err
	}

	// Delete the default "Sheet1"
	f.DeleteSheet("Sheet1")

	if err := f.SaveAs(e.outputFilePath); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("path", e.outputFilePath).Msg("Failed to save Excel file")
		return fmt.Errorf("failed to save excel file to %s: %w", e.outputFilePath, err)
	}

	log.Ctx(ctx).Info().Str("path", e.outputFilePath).Msg("Excel summary generated successfully.")
	return nil
}

// ensureDirectoryExists creates the directory for the output file if it doesn't exist.
func (e *ExcelManagerImpl) ensureDirectoryExists(ctx context.Context) error {
	dir := filepath.Dir(e.outputFilePath)
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
		rowData := []interface{}{
			gainRecord.Symbol,
			gainRecord.BuyDate,  // Assuming string
			gainRecord.SellDate, // Assuming string
			gainRecord.Quantity,
			gainRecord.PNL, // USD PNL
			gainRecord.Commission,
			gainRecord.Type,
			e.formatDateForExcel(gainRecord.TTDate), // Format date
			gainRecord.TTRate,
			gainRecord.INRValue(), // Calculated PNL (INR)
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
			return err
		}
	}

	return nil
}

// writeDividendsSheet handles the creation and population of the "Dividends" sheet.
// It assumes tax.INRDividend has fields: Symbol, Date, Amount, TTDate, TTRate and a method INRValue().
func (e *ExcelManagerImpl) writeDividendsSheet(ctx context.Context, f *excelize.File, dividends []tax.INRDividend) error {
	sheetName := "Dividends"
	headers := []string{
		"Symbol", "Date", "Amount (USD)", "TTDate", "TTRate", "Amount (INR)",
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
			e.formatDateForExcel(dividendRecord.TTDate),
			dividendRecord.TTRate,
			dividendRecord.INRValue(),
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
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
	}

	return nil
}

func (e *ExcelManagerImpl) getValuationHeaders() []string {
	return []string{
		"Symbol",
		"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
		"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
		"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
	}
}

func (e *ExcelManagerImpl) buildValuationRow(valuation tax.INRValuation) []interface{} {
	rowData := []interface{}{valuation.Ticker}
	rowData = append(rowData, e.getPositionRowData(&valuation.FirstPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.PeakPosition)...)
	rowData = append(rowData, e.getPositionRowData(&valuation.YearEndPosition)...)
	return rowData
}

func (e *ExcelManagerImpl) getPositionRowData(pos *tax.INRPosition) []interface{} {
	return []interface{}{
		e.formatDateForExcel(pos.Date),
		pos.Quantity,
		pos.USDPrice,
		pos.USDValue(),
		e.formatDateForExcel(pos.TTDate),
		pos.TTRate,
		pos.INRValue(),
	}
}

// writeInterestSheet handles the creation and population of the "Interest" sheet.
func (e *ExcelManagerImpl) writeInterestSheet(ctx context.Context, f *excelize.File, interest []tax.INRInterest) error {
	sheetName := "Interest"
	headers := []string{
		"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
		"TTDate", "TTRate", "Amount (INR)",
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
			interestRecord.INRValue(),
		}
		if err := e.writeRow(f, sheetName, rowNum, rowData); err != nil {
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
