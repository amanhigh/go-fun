package util

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// WriteRow writes a slice of any data to a specific row in the given sheet.
func WriteRow(f *excelize.File, sheetName string, rowNum int, data []any) error {
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

// WriteFormulaCell writes a formula to a specific cell.
func WriteFormulaCell(f *excelize.File, sheetName string, rowNum, colNum int, formula string) error {
	cellName, err := excelize.CoordinatesToCellName(colNum, rowNum)
	if err != nil {
		return fmt.Errorf("failed to get cell name for col %d, row %d: %w", colNum, rowNum, err)
	}
	if err := f.SetCellFormula(sheetName, cellName, formula); err != nil {
		return fmt.Errorf("failed to set formula at %s: %w", cellName, err)
	}
	return nil
}

// WriteFormulaRange writes multiple formulas for a single row.
// formulas is a map of columnIndex -> formulaString.
func WriteFormulaRange(f *excelize.File, sheetName string, rowNum int, formulas map[int]string) error {
	for colNum, formula := range formulas {
		if err := WriteFormulaCell(f, sheetName, rowNum, colNum, formula); err != nil {
			return err
		}
	}
	return nil
}

// ApplyAutoFilter applies an AutoFilter from row 1 through the given lastColumn and lastDataRow.
// lastColumn is the final column letter (e.g., "C", "J"). The range is "A1:{lastColumn}{lastDataRow}".
func ApplyAutoFilter(f *excelize.File, sheetName, lastColumn string, lastDataRow int) error {
	if lastColumn == "" || lastDataRow < 1 {
		return fmt.Errorf("invalid AutoFilter range: lastColumn=%q, lastDataRow=%d", lastColumn, lastDataRow)
	}

	rangeRef := fmt.Sprintf("A1:%s%d", lastColumn, lastDataRow)
	if err := f.AutoFilter(sheetName, rangeRef, []excelize.AutoFilterOptions{}); err != nil {
		return fmt.Errorf("failed to apply AutoFilter: %w", err)
	}
	return nil
}
