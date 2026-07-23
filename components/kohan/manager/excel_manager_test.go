//nolint:dupl // False positives: Similar Excel test patterns for dividends/interest sheets
package manager_test

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("ExcelManagerImpl", func() {
	var (
		ctx                context.Context
		baseTempDir        string
		testYear           = 2023
		tempOutputFilePath string
		excelManager       manager.ExcelManager
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		baseTempDir, err = os.MkdirTemp(os.TempDir(), "excel_manager_test_run_*")
		Expect(err).ToNot(HaveOccurred())
		tempOutputFilePath = filepath.Join(baseTempDir, fmt.Sprintf("%d_Tax_Summary.xlsx", testYear))
		excelManager = manager.NewExcelManager(baseTempDir)
	})

	AfterEach(func() {
		if baseTempDir != "" {
			err := os.RemoveAll(baseTempDir)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	getCellFloat := func(f *excelize.File, sheetName, axis string) (float64, error) {
		val, err := f.GetCellValue(sheetName, axis)
		if err != nil {
			return 0, err
		}
		if val == "" {
			return 0, nil
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse float from cell %s value '%s': %w", axis, val, err)
		}
		return floatVal, nil
	}

	// mustParseDate parses a fixture date using time.DateOnly, failing the test immediately on error.
	mustParseDate := func(dateStr string) time.Time {
		t, err := time.Parse(time.DateOnly, dateStr)
		Expect(err).ToNot(HaveOccurred(), "mustParseDate(%q)", dateStr)
		return t
	}

	// expectFormulaCell verifies both the formula string and calculated value
	expectFormulaCell := func(f *excelize.File, sheetName, axis, expectedFormula string, expectedValue float64) {
		// Verify formula string
		formula, err := f.GetCellFormula(sheetName, axis)
		Expect(err).ToNot(HaveOccurred())
		Expect(formula).To(Equal(expectedFormula))

		// Verify calculated value (with 2 decimal places precision)
		calculatedStr, err := f.CalcCellValue(sheetName, axis)
		Expect(err).ToNot(HaveOccurred())
		calculatedVal, err := strconv.ParseFloat(calculatedStr, 64)
		Expect(err).ToNot(HaveOccurred())
		Expect(calculatedVal).To(BeNumerically("~", expectedValue, 0.01))
	}

	// expectCrossSheetFormula validates a cross-sheet reference formula and its value
	// sourceSheet: Sheet containing the formula (e.g., "Summary")
	// sourceCell: Cell with the formula (e.g., "A5")
	// targetSheet: Sheet being referenced (e.g., "Dividends")
	// targetCell: Cell being referenced (e.g., "C5")
	// Verifies: formula string matches "=targetSheet!targetCell" and calculated values match
	expectCrossSheetFormula := func(f *excelize.File, sourceSheet, sourceCell, targetSheet, targetCell string) {
		// Get formula from source cell
		formula, err := f.GetCellFormula(sourceSheet, sourceCell)
		Expect(err).ToNot(HaveOccurred())

		// Expected formula format: =targetSheet!targetCell
		expectedFormula := fmt.Sprintf("=%s!%s", targetSheet, targetCell)
		Expect(formula).To(Equal(expectedFormula))

		// Get calculated value from source sheet
		sourceCalcStr, err := f.CalcCellValue(sourceSheet, sourceCell)
		Expect(err).ToNot(HaveOccurred())
		sourceVal, err := strconv.ParseFloat(sourceCalcStr, 64)
		Expect(err).ToNot(HaveOccurred())

		// Get calculated value from target sheet
		targetCalcStr, err := f.CalcCellValue(targetSheet, targetCell)
		Expect(err).ToNot(HaveOccurred())
		targetVal, err := strconv.ParseFloat(targetCalcStr, 64)
		Expect(err).ToNot(HaveOccurred())

		// Verify they match
		Expect(sourceVal).To(BeNumerically("~", targetVal, 0.01))
	}

	// readAutoFilterRanges inspects the raw workbook XML to extract autoFilter
	// ranges for each detail sheet. It maps sheet names to files via workbook.xml
	// relationships to avoid fragile assumptions about sheet XML numbering.
	readAutoFilterRanges := func(xlsxPath string) map[string]string {
		zr, zErr := zip.OpenReader(xlsxPath)
		Expect(zErr).ToNot(HaveOccurred())
		defer func() { Expect(zr.Close()).To(Succeed()) }()

		readZipEntry := func(name string) []byte {
			for _, zf := range zr.File {
				if zf.Name == name {
					rc, oErr := zf.Open()
					Expect(oErr).ToNot(HaveOccurred())
					defer func() { Expect(rc.Close()).To(Succeed()) }()
					data, rErr := io.ReadAll(rc)
					Expect(rErr).ToNot(HaveOccurred())
					return data
				}
			}
			Fail(fmt.Sprintf("zip entry not found: %s", name))
			return nil
		}

		// Parse workbook.xml to get sheet name ↔ rId mapping
		wbXML := string(readZipEntry("xl/workbook.xml"))
		sheetRe := regexp.MustCompile(`<sheet[^>]*name="([^"]+)"[^>]*r:id="([^"]+)"`)
		sheetMatches := sheetRe.FindAllStringSubmatch(wbXML, -1)

		// Parse workbook.xml.rels to get rId ↔ target (path) mapping
		relsXML := string(readZipEntry("xl/_rels/workbook.xml.rels"))
		relRe := regexp.MustCompile(`<Relationship[^>]*Id="([^"]+)"[^>]*Target="([^"]+)"`)
		relMatches := relRe.FindAllStringSubmatch(relsXML, -1)

		rIDToTarget := make(map[string]string)
		for _, m := range relMatches {
			target := m[2]
			// Normalize target path: remove leading "/" and ensure "xl/" prefix
			target = strings.TrimPrefix(target, "/")
			if !strings.HasPrefix(target, "xl/") {
				target = "xl/" + target
			}
			rIDToTarget[m[1]] = target
		}

		sheetToTarget := make(map[string]string)
		for _, m := range sheetMatches {
			if target, ok := rIDToTarget[m[2]]; ok {
				sheetToTarget[m[1]] = target
			}
		}

		// Read each worksheet XML and extract autoFilter ref
		afRe := regexp.MustCompile(`<autoFilter[^>]*ref="([^"]+)"`)
		sheetRanges := make(map[string]string)
		for sName, target := range sheetToTarget {
			if sName == "Summary" {
				continue
			}
			wsXML := string(readZipEntry(target))
			if m := afRe.FindStringSubmatch(wsXML); len(m) > 1 {
				// Normalize ref: ensure absolute (A1:J3 → $A$1:$J$3),
				// preserve already-absolute refs ($A$1:$J$3) without $$ doubling.
				parts := strings.Split(m[1], ":")
				Expect(parts).To(HaveLen(2), "autoFilter ref should have two parts")
				makeAbsolute := func(ref string) string {
					if strings.Contains(ref, "$") {
						return ref
					}
					return "$" + ref
				}
				sheetRanges[sName] = makeAbsolute(parts[0]) + ":" + makeAbsolute(parts[1])
			}
		}
		return sheetRanges
	}

	Describe("GenerateTaxSummaryExcel", func() {
		Context("when generating the 'Gains' sheet with data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Gains"
			)

			BeforeEach(func() {
				gain1TTDate := mustParseDate("2023-01-15")
				gain1 := tax.INRGains{
					Gains:  tax.Gains{Symbol: "AAPL", BuyDate: "2022-10-01", SellDate: "2023-01-20", Quantity: 10.5, PNL: 100.75, Commission: 5.25, Type: "STCG"},
					TTDate: gain1TTDate, TTRate: 82.50,
				}
				gain2TTDate := mustParseDate("2023-02-10")
				gain2 := tax.INRGains{
					Gains:  tax.Gains{Symbol: "MSFT", BuyDate: "2020-05-01", SellDate: "2023-02-15", Quantity: 5, PNL: -50.20, Commission: 0, Type: "LTCG"},
					TTDate: gain2TTDate, TTRate: 83.10,
				}
				gain3WithZeroTTDate := tax.INRGains{
					Gains:  tax.Gains{Symbol: "GOOG", BuyDate: "2021-01-01", SellDate: "2023-03-10", Quantity: 20, PNL: 200.00, Commission: 1.50, Type: "LTCG"},
					TTDate: time.Time{}, TTRate: 81.75,
				}
				sampleSummary.INRGains = []tax.INRGains{gain1, gain2, gain3WithZeroTTDate}
			})

			It("should create the Excel file successfully at the specified path", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())
				Expect(tempOutputFilePath).Should(BeARegularFile())
			})

			It("should contain a 'Gains' sheet with the correct headers in order", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetFound := slices.Contains(f.GetSheetList(), sheetName)
				Expect(sheetFound).To(BeTrue(), "Sheet 'Gains' should exist")

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).ToNot(BeEmpty(), "Sheet should have at least a header row")

				expectedHeaders := []string{
					"Symbol", "BuyDate", "SellDate", "Quantity", "PNL (USD)",
					"Commission (USD)", "Type", "TTDate", "TTRate", "PNL (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders), "Headers in 'Gains' sheet are incorrect")
			})

			It("should write all INRGains records accurately to the sheet", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, errGetRows := f.GetRows(sheetName)
				Expect(errGetRows).ToNot(HaveOccurred())
				// Rows include: headers + data records + empty row + TOTALS row + empty row + STCG row + LTCG row
				expectedRowCount := 1 + len(sampleSummary.INRGains) + 5
				Expect(rows).To(HaveLen(expectedRowCount), "Number of rows should be headers + data records + totals section")

				// Sorted by SellDate descending: GOOG (2023-03-10), MSFT (2023-02-15), AAPL (2023-01-20)
				// Verify GOOG (row 2 in Excel, index 1 in `rows` slice) — latest sell first
				gain3 := sampleSummary.INRGains[2]
				Expect(rows[1][0]).To(Equal(gain3.Symbol))
				Expect(rows[1][1]).To(Equal(gain3.BuyDate))
				Expect(rows[1][2]).To(Equal(gain3.SellDate))
				qty1, err := getCellFloat(f, sheetName, "D2")
				Expect(err).ToNot(HaveOccurred())
				Expect(qty1).To(BeNumerically("~", gain3.Quantity, 0.001))
				pnlUSD1, err := getCellFloat(f, sheetName, "E2")
				Expect(err).ToNot(HaveOccurred())
				Expect(pnlUSD1).To(BeNumerically("~", gain3.PNL, 0.001))
				comm1, err := getCellFloat(f, sheetName, "F2")
				Expect(err).ToNot(HaveOccurred())
				Expect(comm1).To(BeNumerically("~", gain3.Commission, 0.001))
				Expect(rows[1][6]).To(Equal(gain3.Type))
				Expect(rows[1][7]).To(Equal(""), "TTDate for zero time should be an empty string")
				rate1, err := getCellFloat(f, sheetName, "I2")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate1).To(BeNumerically("~", gain3.TTRate, 0.001))
				// Column J has formula =E2*I2 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J2", "=E2*I2", gain3.PNL*gain3.TTRate)

				// Verify MSFT (row 3 in Excel, index 2 in `rows` slice)
				gain2 := sampleSummary.INRGains[1]
				Expect(rows[2][0]).To(Equal(gain2.Symbol))
				Expect(rows[2][1]).To(Equal(gain2.BuyDate))
				Expect(rows[2][2]).To(Equal(gain2.SellDate))
				qty2, err := getCellFloat(f, sheetName, "D3")
				Expect(err).ToNot(HaveOccurred())
				Expect(qty2).To(BeNumerically("~", gain2.Quantity, 0.001))
				pnlUSD2, err := getCellFloat(f, sheetName, "E3")
				Expect(err).ToNot(HaveOccurred())
				Expect(pnlUSD2).To(BeNumerically("~", gain2.PNL, 0.001))
				comm2, err := getCellFloat(f, sheetName, "F3")
				Expect(err).ToNot(HaveOccurred())
				Expect(comm2).To(BeNumerically("~", gain2.Commission, 0.001))
				Expect(rows[2][6]).To(Equal(gain2.Type))
				Expect(rows[2][7]).To(Equal(gain2.TTDate.Format(time.DateOnly)))
				rate2, err := getCellFloat(f, sheetName, "I3")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate2).To(BeNumerically("~", gain2.TTRate, 0.001))
				// Column J has formula =E3*I3 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J3", "=E3*I3", gain2.PNL*gain2.TTRate)

				// Verify AAPL (row 4 in Excel, index 3 in `rows` slice) — earliest sell last
				gain1 := sampleSummary.INRGains[0]
				Expect(rows[3][0]).To(Equal(gain1.Symbol))
				Expect(rows[3][1]).To(Equal(gain1.BuyDate))
				Expect(rows[3][2]).To(Equal(gain1.SellDate))
				qty3, err := getCellFloat(f, sheetName, "D4")
				Expect(err).ToNot(HaveOccurred())
				Expect(qty3).To(BeNumerically("~", gain1.Quantity, 0.001))
				pnlUSD3, err := getCellFloat(f, sheetName, "E4")
				Expect(err).ToNot(HaveOccurred())
				Expect(pnlUSD3).To(BeNumerically("~", gain1.PNL, 0.001))
				comm3, err := getCellFloat(f, sheetName, "F4")
				Expect(err).ToNot(HaveOccurred())
				Expect(comm3).To(BeNumerically("~", gain1.Commission, 0.001))
				Expect(rows[3][6]).To(Equal(gain1.Type))
				Expect(rows[3][7]).To(Equal(gain1.TTDate.Format(time.DateOnly)))
				rate3, err := getCellFloat(f, sheetName, "I4")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate3).To(BeNumerically("~", gain1.TTRate, 0.001))
				// Column J has formula =E4*I4 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J4", "=E4*I4", gain1.PNL*gain1.TTRate)
			})

			It("should write TOTALS, STCG, and LTCG rows with correct labels and formulas", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer func() { Expect(f.Close()).To(Succeed()) }()

				lastDataRow := len(sampleSummary.INRGains) + 1 // Row 4 (3 data rows: rows 2-4)
				totalsRow := lastDataRow + 2                   // Row 6 (skip empty row 5)
				stcgRow := totalsRow + 2                       // Row 8 (skip empty row 7)
				ltcgRow := stcgRow + 1                         // Row 9

				// Compute expected totals from the fixture data
				// gain1 (AAPL, STCG): PNL=100.75, Commission=5.25, TTRate=82.50
				// gain2 (MSFT, LTCG): PNL=-50.20, Commission=0, TTRate=83.10
				// gain3 (GOOG, LTCG): PNL=200.00, Commission=1.50, TTRate=81.75
				gains := sampleSummary.INRGains
				totalPNLUSD := lo.SumBy(gains, func(g tax.INRGains) float64 { return g.PNL })
				totalCommissionUSD := lo.SumBy(gains, func(g tax.INRGains) float64 { return g.Commission })
				totalPNLINR := lo.SumBy(gains, func(g tax.INRGains) float64 { return g.PNL * g.TTRate })
				stcgPNLUSD := gains[0].PNL // AAPL is the only STCG
				stcgPNLINR := gains[0].PNL * gains[0].TTRate
				ltcgPNLUSD := gains[1].PNL + gains[2].PNL
				ltcgPNLINR := gains[1].PNL*gains[1].TTRate + gains[2].PNL*gains[2].TTRate

				// Verify TOTALS label and formulas
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))
				expectFormulaCell(f, sheetName, fmt.Sprintf("E%d", totalsRow),
					fmt.Sprintf("=SUM(E2:E%d)", lastDataRow), totalPNLUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("F%d", totalsRow),
					fmt.Sprintf("=SUM(F2:F%d)", lastDataRow), totalCommissionUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("J%d", totalsRow),
					fmt.Sprintf("=SUM(J2:J%d)", lastDataRow), totalPNLINR)

				// Verify STCG label and SUMIF formulas
				stcgLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", stcgRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(stcgLabel).To(Equal("STCG"))
				expectFormulaCell(f, sheetName, fmt.Sprintf("E%d", stcgRow),
					fmt.Sprintf("=SUMIF(G2:G%d,\"STCG\",E2:E%d)", lastDataRow, lastDataRow), stcgPNLUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("J%d", stcgRow),
					fmt.Sprintf("=SUMIF(G2:G%d,\"STCG\",J2:J%d)", lastDataRow, lastDataRow), stcgPNLINR)

				// Verify LTCG label and SUMIF formulas
				ltcgLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", ltcgRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(ltcgLabel).To(Equal("LTCG"))
				expectFormulaCell(f, sheetName, fmt.Sprintf("E%d", ltcgRow),
					fmt.Sprintf("=SUMIF(G2:G%d,\"LTCG\",E2:E%d)", lastDataRow, lastDataRow), ltcgPNLUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("J%d", ltcgRow),
					fmt.Sprintf("=SUMIF(G2:G%d,\"LTCG\",J2:J%d)", lastDataRow, lastDataRow), ltcgPNLINR)
			})

			It("should set custom column widths for Gains sheet", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer func() { Expect(f.Close()).To(Succeed()) }()

				widthA, err := f.GetColWidth(sheetName, "A")
				Expect(err).ToNot(HaveOccurred())
				Expect(widthA).To(BeNumerically("==", 8.0), "Column A (Symbol) should be width 8")

				widthF, err := f.GetColWidth(sheetName, "F")
				Expect(err).ToNot(HaveOccurred())
				Expect(widthF).To(BeNumerically("==", 16.0), "Column F (Commission) should be width 16")
			})
		})

		Context("when generating the 'Dividends' sheet with data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Dividends"
				f             *excelize.File
			)

			BeforeEach(func() {
				div1TTDate := mustParseDate("2023-04-05")
				div1 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "AAPL", Date: "2023-04-10", Amount: 50.25, Tax: 7.54, Net: 42.71},
					TTDate:   div1TTDate, TTRate: 82.10,
				}
				div2TTDate := mustParseDate("2023-05-12")
				div2 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "GOOG", Date: "2023-05-15", Amount: 75.50, Tax: 11.33, Net: 64.17},
					TTDate:   div2TTDate, TTRate: 82.50,
				}
				sampleSummary.INRDividends = []tax.INRDividend{div1, div2}

				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					Expect(f.Close()).To(Succeed())
				}
			})

			It("should create the 'Dividends' sheet with correct headers and data", func() {
				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				// Rows include: headers + data records + empty row + TOTALS row
				expectedRowCount := 1 + len(sampleSummary.INRDividends) + 2
				Expect(rows).To(HaveLen(expectedRowCount), "unexpected row count for Dividends sheet")

				// Verify Headers
				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)", "TTDate", "TTRate",
					"Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders), "Dividends headers incorrect")

				// Sorted by Date descending: GOOG (2023-05-15) → row 2, AAPL (2023-04-10) → row 3
				// Verify GOOG (row 2 in Excel, index 1 in `rows` slice) — latest date first
				div2 := sampleSummary.INRDividends[1]
				Expect(rows[1][0]).To(Equal(div2.Symbol))
				Expect(rows[1][1]).To(Equal(div2.Date))
				amount2, err := getCellFloat(f, sheetName, "C2")
				Expect(err).ToNot(HaveOccurred())
				Expect(amount2).To(BeNumerically("~", div2.Amount, 0.001))
				tax2, err := getCellFloat(f, sheetName, "D2")
				Expect(err).ToNot(HaveOccurred())
				Expect(tax2).To(BeNumerically("~", div2.Tax, 0.001))
				net2, err := getCellFloat(f, sheetName, "E2")
				Expect(err).ToNot(HaveOccurred())
				Expect(net2).To(BeNumerically("~", div2.Net, 0.001))
				Expect(rows[1][5]).To(Equal(div2.TTDate.Format(time.DateOnly)))
				rate2, err := getCellFloat(f, sheetName, "G2")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate2).To(BeNumerically("~", div2.TTRate, 0.001))
				// Verify INR formulas for GOOG
				expectFormulaCell(f, sheetName, "H2", "=C2*G2", div2.Amount*div2.TTRate) // Amount (INR)
				expectFormulaCell(f, sheetName, "I2", "=D2*G2", div2.Tax*div2.TTRate)    // Tax (INR)
				expectFormulaCell(f, sheetName, "J2", "=E2*G2", div2.Net*div2.TTRate)    // Net (INR)

				// Verify AAPL (row 3 in Excel, index 2 in `rows` slice) — earliest date last
				div1 := sampleSummary.INRDividends[0]
				Expect(rows[2][0]).To(Equal(div1.Symbol))
				Expect(rows[2][1]).To(Equal(div1.Date))
				amount1, err := getCellFloat(f, sheetName, "C3")
				Expect(err).ToNot(HaveOccurred())
				Expect(amount1).To(BeNumerically("~", div1.Amount, 0.001))
				tax1, err := getCellFloat(f, sheetName, "D3")
				Expect(err).ToNot(HaveOccurred())
				Expect(tax1).To(BeNumerically("~", div1.Tax, 0.001))
				net1, err := getCellFloat(f, sheetName, "E3")
				Expect(err).ToNot(HaveOccurred())
				Expect(net1).To(BeNumerically("~", div1.Net, 0.001))
				Expect(rows[2][5]).To(Equal(div1.TTDate.Format(time.DateOnly)))
				rate1, err := getCellFloat(f, sheetName, "G3")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate1).To(BeNumerically("~", div1.TTRate, 0.001))
				// Verify INR formulas for AAPL
				expectFormulaCell(f, sheetName, "H3", "=C3*G3", div1.Amount*div1.TTRate) // Amount (INR)
				expectFormulaCell(f, sheetName, "I3", "=D3*G3", div1.Tax*div1.TTRate)    // Tax (INR)
				expectFormulaCell(f, sheetName, "J3", "=E3*G3", div1.Net*div1.TTRate)    // Net (INR)
			})

			It("should write TOTALS row with correct formulas and calculated values", func() {
				// Calculate expected values using lo.SumBy
				totalAmountUSD := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Amount
				})
				totalTaxUSD := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Tax
				})
				totalNetUSD := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Net
				})
				totalAmountINR := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Amount * d.TTRate
				})
				totalTaxINR := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Tax * d.TTRate
				})
				totalNetINR := lo.SumBy(sampleSummary.INRDividends, func(d tax.INRDividend) float64 {
					return d.Net * d.TTRate
				})

				// Verify TOTALS row position
				lastDataRow := len(sampleSummary.INRDividends) + 1 // Row 3 (2 data rows)
				totalsRow := lastDataRow + 2                       // Row 5 (skip empty row 4)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

				// Verify USD columns (C, D, E)
				expectFormulaCell(f, sheetName, fmt.Sprintf("C%d", totalsRow),
					fmt.Sprintf("=SUM(C2:C%d)", lastDataRow), totalAmountUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("D%d", totalsRow),
					fmt.Sprintf("=SUM(D2:D%d)", lastDataRow), totalTaxUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("E%d", totalsRow),
					fmt.Sprintf("=SUM(E2:E%d)", lastDataRow), totalNetUSD)

				// Verify INR columns (H, I, J)
				expectFormulaCell(f, sheetName, fmt.Sprintf("H%d", totalsRow),
					fmt.Sprintf("=SUM(H2:H%d)", lastDataRow), totalAmountINR)
				expectFormulaCell(f, sheetName, fmt.Sprintf("I%d", totalsRow),
					fmt.Sprintf("=SUM(I2:I%d)", lastDataRow), totalTaxINR)
				expectFormulaCell(f, sheetName, fmt.Sprintf("J%d", totalsRow),
					fmt.Sprintf("=SUM(J2:J%d)", lastDataRow), totalNetINR)
			})
		})

		Context("when generating the 'Valuations' sheet with data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Valuations"
			)

			BeforeEach(func() {
				// Define dates
				firstDate := mustParseDate("2022-01-10")
				firstTTDate := mustParseDate("2022-01-11")
				peakDate := mustParseDate("2022-11-25")
				yearEndDate := mustParseDate("2023-03-31")

				// Create a full valuation object
				val1 := tax.INRValuation{
					Ticker: "TSLA",
					FirstPosition: tax.INRPosition{
						Position: tax.Position{Date: firstDate, Quantity: 10, USDPrice: 250.0},
						TTDate:   firstTTDate,
						TTRate:   80.5,
					},
					PeakPosition: tax.INRPosition{
						Position: tax.Position{Date: peakDate, Quantity: 15, USDPrice: 310.0},
						TTDate:   peakDate,
						TTRate:   81.90,
					},
					YearEndPosition: tax.INRPosition{
						Position: tax.Position{Date: yearEndDate, Quantity: 5, USDPrice: 207.46},
						TTDate:   yearEndDate,
						TTRate:   82.17,
					},
				}
				sampleSummary.INRValuations = []tax.INRValuation{val1}
			})

			It("should create the 'Valuations' sheet with correct headers and data for all positions", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer func() { Expect(f.Close()).To(Succeed()) }()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				// Rows include: headers + data records + empty row + TOTALS row
				expectedRowCount := 1 + len(sampleSummary.INRValuations) + 2
				Expect(rows).To(HaveLen(expectedRowCount))

				// Verify Headers
				expectedHeaders := []string{
					"Symbol",
					"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"AmountPaid (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))

				// Verify Data Row
				val1 := sampleSummary.INRValuations[0]

				// First Position
				posFirst := val1.FirstPosition
				Expect(rows[1][0]).To(Equal(val1.Ticker))
				Expect(rows[1][1]).To(Equal(posFirst.Date.Format(time.DateOnly)))
				qty, _ := getCellFloat(f, sheetName, "C2")
				Expect(qty).To(Equal(posFirst.Quantity))
				price, _ := getCellFloat(f, sheetName, "D2")
				Expect(price).To(Equal(posFirst.USDPrice))
				Expect(rows[1][5]).To(Equal(posFirst.TTDate.Format(time.DateOnly)))
				rate, _ := getCellFloat(f, sheetName, "G2")
				Expect(rate).To(Equal(posFirst.TTRate))
				// Verify First Position formulas
				expectFormulaCell(f, sheetName, "E2", "=C2*D2", posFirst.USDValue())
				expectFormulaCell(f, sheetName, "H2", "=E2*G2", posFirst.INRValue())

				// Peak Position
				posPeak := val1.PeakPosition
				Expect(rows[1][8]).To(Equal(posPeak.Date.Format(time.DateOnly)))
				qty, _ = getCellFloat(f, sheetName, "J2")
				Expect(qty).To(Equal(posPeak.Quantity))
				price, _ = getCellFloat(f, sheetName, "K2")
				Expect(price).To(Equal(posPeak.USDPrice))
				Expect(rows[1][12]).To(Equal(posPeak.TTDate.Format(time.DateOnly)))
				rate, _ = getCellFloat(f, sheetName, "N2")
				Expect(rate).To(Equal(posPeak.TTRate))
				// Verify Peak Position formulas
				expectFormulaCell(f, sheetName, "L2", "=J2*K2", posPeak.USDValue())
				expectFormulaCell(f, sheetName, "O2", "=L2*N2", posPeak.INRValue())

				// Year End Position
				posYearEnd := val1.YearEndPosition
				Expect(rows[1][15]).To(Equal(posYearEnd.Date.Format(time.DateOnly)))
				qty, _ = getCellFloat(f, sheetName, "Q2")
				Expect(qty).To(Equal(posYearEnd.Quantity))
				price, _ = getCellFloat(f, sheetName, "R2")
				Expect(price).To(Equal(posYearEnd.USDPrice))
				Expect(rows[1][19]).To(Equal(posYearEnd.TTDate.Format(time.DateOnly)))
				rate, _ = getCellFloat(f, sheetName, "U2")
				Expect(rate).To(Equal(posYearEnd.TTRate))
				// Verify YearEnd Position formulas
				expectFormulaCell(f, sheetName, "S2", "=Q2*R2", posYearEnd.USDValue())
				expectFormulaCell(f, sheetName, "V2", "=S2*U2", posYearEnd.INRValue())
			})

			It("should write TOTALS row for AmountPaid column", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Calculate expected totals using lo.SumBy (1 total)
				totalAmountPaidINR := lo.SumBy(sampleSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.AmountPaid
				})
				totalYearEndValINR := lo.SumBy(sampleSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.YearEndPosition.INRValue()
				})

				// Verify TOTALS row position
				lastDataRow := len(sampleSummary.INRValuations) + 1 // Row 2 (1 data row)
				totalsRow := lastDataRow + 2                        // Row 4 (skip empty row 3)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

				// Verify YearEnd ValINR (INR) column V
				expectFormulaCell(f, sheetName, fmt.Sprintf("V%d", totalsRow),
					fmt.Sprintf("=SUM(V2:V%d)", lastDataRow), totalYearEndValINR)

				// Verify AmountPaid (INR) column W
				expectFormulaCell(f, sheetName, fmt.Sprintf("W%d", totalsRow),
					fmt.Sprintf("=SUM(W2:W%d)", lastDataRow), totalAmountPaidINR)
			})

			It("should write TOTALS row for AmountPaid with non-zero value", func() {
				// Define dates for this test
				firstDate2, _ := time.Parse(time.DateOnly, "2022-02-15")
				firstTTDate2, _ := time.Parse(time.DateOnly, "2022-02-16")
				peakDate2, _ := time.Parse(time.DateOnly, "2022-12-20")
				yearEndDate2, _ := time.Parse(time.DateOnly, "2023-04-30")

				// Create new valuations with non-zero AmountPaid
				val2 := tax.INRValuation{
					Ticker: "AAPL",
					FirstPosition: tax.INRPosition{
						Position: tax.Position{Date: firstDate2, Quantity: 100, USDPrice: 150.0},
						TTDate:   firstTTDate2,
						TTRate:   82.0,
					},
					PeakPosition: tax.INRPosition{
						Position: tax.Position{Date: peakDate2, Quantity: 120, USDPrice: 180.0},
						TTDate:   peakDate2,
						TTRate:   82.5,
					},
					YearEndPosition: tax.INRPosition{
						Position: tax.Position{Date: yearEndDate2, Quantity: 110, USDPrice: 175.0},
						TTDate:   yearEndDate2,
						TTRate:   83.0,
					},
					AmountPaid: 5432.10, // Sum of gross dividends in INR
				}

				val3 := tax.INRValuation{
					Ticker: "MSFT",
					FirstPosition: tax.INRPosition{
						Position: tax.Position{Date: firstDate2, Quantity: 50, USDPrice: 300.0},
						TTDate:   firstTTDate2,
						TTRate:   82.0,
					},
					PeakPosition: tax.INRPosition{
						Position: tax.Position{Date: peakDate2, Quantity: 60, USDPrice: 350.0},
						TTDate:   peakDate2,
						TTRate:   82.5,
					},
					YearEndPosition: tax.INRPosition{
						Position: tax.Position{Date: yearEndDate2, Quantity: 55, USDPrice: 325.0},
						TTDate:   yearEndDate2,
						TTRate:   83.0,
					},
					AmountPaid: 3210.50, // Sum of gross dividends in INR
				}

				nonZeroSummary := tax.Summary{
					INRValuations: []tax.INRValuation{val2, val3},
				}

				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, nonZeroSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Calculate expected totals using lo.SumBy (2 valuations)
				totalAmountPaidINR := lo.SumBy(nonZeroSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.AmountPaid
				})
				totalYearEndValINR := lo.SumBy(nonZeroSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.YearEndPosition.INRValue()
				})

				// Verify TOTALS row position
				lastDataRow := len(nonZeroSummary.INRValuations) + 1 // Row 3 (2 data rows)
				totalsRow := lastDataRow + 2                         // Row 5 (skip empty row 4)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

				// Verify YearEnd ValINR (INR) column V with non-zero total
				expectFormulaCell(f, sheetName, fmt.Sprintf("V%d", totalsRow),
					fmt.Sprintf("=SUM(V2:V%d)", lastDataRow), totalYearEndValINR)

				// Verify AmountPaid (INR) column W with non-zero total
				// Expected: 5432.10 + 3210.50 = 8642.60
				expectFormulaCell(f, sheetName, fmt.Sprintf("W%d", totalsRow),
					fmt.Sprintf("=SUM(W2:W%d)", lastDataRow), totalAmountPaidINR)
			})
		})

		Context("when INRValuations slice is empty", func() {
			It("should create the 'Valuations' sheet with only headers", func() {
				emptySummary := tax.Summary{INRValuations: []tax.INRValuation{}}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetName := "Valuations"
				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1), "Sheet should only contain the header row")

				expectedHeaders := []string{
					"Symbol",
					"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"AmountPaid (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when generating the 'Interest' sheet with data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Interest"
			)

			BeforeEach(func() {
				// Define dates
				interestDate, _ := time.Parse(time.DateOnly, "2023-06-01")
				ttDate, _ := time.Parse(time.DateOnly, "2023-06-02")

				// Create a full interest object
				interest1 := tax.INRInterest{
					Interest: tax.Interest{
						Symbol: "US-TBILL",
						Date:   interestDate.Format(time.DateOnly),
						Amount: 100.0,
						Tax:    10.0,
						Net:    90.0,
					},
					TTDate: ttDate,
					TTRate: 82.5,
				}
				sampleSummary.INRInterest = []tax.INRInterest{interest1}
			})

			It("should create the 'Interest' sheet with correct headers and data", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				// Rows include: headers + data records + empty row + TOTALS row
				expectedRowCount := 1 + len(sampleSummary.INRInterest) + 2
				Expect(rows).To(HaveLen(expectedRowCount))

				// Verify Headers
				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
					"TTDate", "TTRate", "Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))

				// Verify Data Row
				interest1 := sampleSummary.INRInterest[0]
				Expect(rows[1][0]).To(Equal(interest1.Symbol))
				Expect(rows[1][1]).To(Equal(interest1.Date))
				amount, _ := getCellFloat(f, sheetName, "C2")
				Expect(amount).To(BeNumerically("==", interest1.Amount))
				taxVal, _ := getCellFloat(f, sheetName, "D2")
				Expect(taxVal).To(BeNumerically("==", interest1.Tax))
				net, _ := getCellFloat(f, sheetName, "E2")
				Expect(net).To(BeNumerically("==", interest1.Net))
				Expect(rows[1][5]).To(Equal(interest1.TTDate.Format(time.DateOnly)))
				rate, _ := getCellFloat(f, sheetName, "G2")
				Expect(rate).To(BeNumerically("==", interest1.TTRate))
				// Verify formulas for INR columns
				expectFormulaCell(f, sheetName, "H2", "=C2*G2", interest1.Amount*interest1.TTRate) // Amount (INR)
				expectFormulaCell(f, sheetName, "I2", "=D2*G2", interest1.Tax*interest1.TTRate)    // Tax (INR)
				expectFormulaCell(f, sheetName, "J2", "=E2*G2", interest1.Net*interest1.TTRate)    // Net (INR)
			})

			It("should write TOTALS row with correct formulas and calculated values", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Calculate expected values using lo.SumBy
				totalAmountUSD := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Amount
				})
				totalTaxUSD := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Tax
				})
				totalNetUSD := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Net
				})
				totalAmountINR := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Amount * i.TTRate
				})
				totalTaxINR := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Tax * i.TTRate
				})
				totalNetINR := lo.SumBy(sampleSummary.INRInterest, func(i tax.INRInterest) float64 {
					return i.Net * i.TTRate
				})

				// Verify TOTALS row position
				lastDataRow := len(sampleSummary.INRInterest) + 1 // Row 2 (1 data row)
				totalsRow := lastDataRow + 2                      // Row 4 (skip empty row 3)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

				// Verify USD columns (C, D, E)
				expectFormulaCell(f, sheetName, fmt.Sprintf("C%d", totalsRow),
					fmt.Sprintf("=SUM(C2:C%d)", lastDataRow), totalAmountUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("D%d", totalsRow),
					fmt.Sprintf("=SUM(D2:D%d)", lastDataRow), totalTaxUSD)
				expectFormulaCell(f, sheetName, fmt.Sprintf("E%d", totalsRow),
					fmt.Sprintf("=SUM(E2:E%d)", lastDataRow), totalNetUSD)

				// Verify INR columns (H, I, J)
				expectFormulaCell(f, sheetName, fmt.Sprintf("H%d", totalsRow),
					fmt.Sprintf("=SUM(H2:H%d)", lastDataRow), totalAmountINR)
				expectFormulaCell(f, sheetName, fmt.Sprintf("I%d", totalsRow),
					fmt.Sprintf("=SUM(I2:I%d)", lastDataRow), totalTaxINR)
				expectFormulaCell(f, sheetName, fmt.Sprintf("J%d", totalsRow),
					fmt.Sprintf("=SUM(J2:J%d)", lastDataRow), totalNetINR)
			})
		})

		Context("when INRInterest slice is empty", func() {
			It("should create the 'Interest' sheet with only headers", func() {
				emptySummary := tax.Summary{INRInterest: []tax.INRInterest{}}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetName := "Interest"
				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1), "Sheet should only contain the header row")

				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
					"TTDate", "TTRate", "Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when generating the 'TT Rates' sheet with data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "TT Rates"
			)

			BeforeEach(func() {
				marDate, _ := time.Parse(time.DateOnly, "2023-03-15")
				aprDate, _ := time.Parse(time.DateOnly, "2023-04-15")
				junDate, _ := time.Parse(time.DateOnly, "2023-06-15") // Thursday
				febDate, _ := time.Parse(time.DateOnly, "2024-02-20") // Tuesday

				sampleSummary = tax.Summary{
					Year: testYear,
					TTMonthEndRates: []tax.MonthEndRate{
						{Rate: 82.00, ActualDate: marDate, PDFFile: "-"},
						{Rate: 82.10, ActualDate: aprDate, PDFFile: "https://sbi.com/apr23.pdf"},
						{Rate: 82.30, ActualDate: junDate, PDFFile: "https://sbi.com/jun23.pdf"},
						{Rate: 83.05, ActualDate: febDate, PDFFile: "https://sbi.com/feb24.pdf"},
					},
				}
			})

			It("should create the 'TT Rates' sheet with correct headers", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetFound := slices.Contains(f.GetSheetList(), sheetName)
				Expect(sheetFound).To(BeTrue(), "Sheet 'TT Rates' should exist")

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())

				// Header row + 4 data rows
				Expect(rows).To(HaveLen(5), "Should have header + 4 data rows")

				expectedHeaders := []string{
					"Month", "Year", "TTDate", "TTRate", "PDF Link", "DayOfWeek",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})

			It("should write FY month labels in Apr→Mar order with correct year", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Row 2 (index 0) → APR, Year=2023
				month0, err := f.GetCellValue(sheetName, "A2")
				Expect(err).ToNot(HaveOccurred())
				Expect(month0).To(Equal("APR"))
				year0, err := f.GetCellValue(sheetName, "B2")
				Expect(err).ToNot(HaveOccurred())
				Expect(year0).To(Equal("2023"))

				// Row 3 (index 1) → MAY, Year=2023
				month1, err := f.GetCellValue(sheetName, "A3")
				Expect(err).ToNot(HaveOccurred())
				Expect(month1).To(Equal("MAY"))
				year1, err := f.GetCellValue(sheetName, "B3")
				Expect(err).ToNot(HaveOccurred())
				Expect(year1).To(Equal("2023"))

				// Row 4 (index 2) → JUN, Year=2023
				month2, err := f.GetCellValue(sheetName, "A4")
				Expect(err).ToNot(HaveOccurred())
				Expect(month2).To(Equal("JUN"))
				year2, err := f.GetCellValue(sheetName, "B4")
				Expect(err).ToNot(HaveOccurred())
				Expect(year2).To(Equal("2023"))

				// Row 5 (index 3) → JUL, Year=2023
				month3, err := f.GetCellValue(sheetName, "A5")
				Expect(err).ToNot(HaveOccurred())
				Expect(month3).To(Equal("JUL"))
				year3, err := f.GetCellValue(sheetName, "B5")
				Expect(err).ToNot(HaveOccurred())
				Expect(year3).To(Equal("2023"))
			})

			It("should write TTDate, TTRate, and DayOfWeek correctly", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Row 2: Index 0 (APR) - preceding March 2023
				ttDate2, err := f.GetCellValue(sheetName, "C2")
				Expect(err).ToNot(HaveOccurred())
				Expect(ttDate2).To(Equal("2023-03-15"))
				rate2, err := getCellFloat(f, sheetName, "D2")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate2).To(BeNumerically("~", 82.00, 0.001))
				dow2, err := f.GetCellValue(sheetName, "F2")
				Expect(err).ToNot(HaveOccurred())
				Expect(dow2).To(Equal("Wednesday"))

				// Row 3: Index 1 (MAY) - preceding April 2023
				ttDate3, err := f.GetCellValue(sheetName, "C3")
				Expect(err).ToNot(HaveOccurred())
				Expect(ttDate3).To(Equal("2023-04-15"))
				rate3, err := getCellFloat(f, sheetName, "D3")
				Expect(err).ToNot(HaveOccurred())
				Expect(rate3).To(BeNumerically("~", 82.10, 0.001))

				// Row 4: Index 2 (JUN) - Jun 15, 2023 = Thursday
				dow4, err := f.GetCellValue(sheetName, "F4")
				Expect(err).ToNot(HaveOccurred())
				Expect(dow4).To(Equal("Thursday"))

				// Row 5: Index 3 (JUL) - preceding February 2024 (year boundary test)
				ttDate5, err := f.GetCellValue(sheetName, "C5")
				Expect(err).ToNot(HaveOccurred())
				Expect(ttDate5).To(Equal("2024-02-20"))
				dow5, err := f.GetCellValue(sheetName, "F5")
				Expect(err).ToNot(HaveOccurred())
				Expect(dow5).To(Equal("Tuesday"))
			})

			It("should write clickable PDF hyperlink for URL PDFFile values", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Row 2 (APR): PDFFile = "-" → should NOT have hyperlink
				hasLink, target, err := f.GetCellHyperLink(sheetName, "E2")
				Expect(err).ToNot(HaveOccurred())
				Expect(hasLink).To(BeFalse(), "Row 2 should not have a hyperlink for '-' PDFFile")
				Expect(target).To(Equal(""))
				val2, err := f.GetCellValue(sheetName, "E2")
				Expect(err).ToNot(HaveOccurred())
				Expect(val2).To(Equal("-"))

				// Row 3 (MAY): PDFFile = "https://sbi.com/apr23.pdf" → should have hyperlink
				hasLink, target, err = f.GetCellHyperLink(sheetName, "E3")
				Expect(err).ToNot(HaveOccurred())
				Expect(hasLink).To(BeTrue(), "Row 3 should have a hyperlink for URL PDFFile")
				Expect(target).To(Equal("https://sbi.com/apr23.pdf"))
				val3, err := f.GetCellValue(sheetName, "E3")
				Expect(err).ToNot(HaveOccurred())
				Expect(val3).To(Equal("PDF"))

				// Row 4 (JUL): PDFFile = "https://sbi.com/jun23.pdf" → should have hyperlink
				hasLink, target, err = f.GetCellHyperLink(sheetName, "E4")
				Expect(err).ToNot(HaveOccurred())
				Expect(hasLink).To(BeTrue(), "Row 4 should have a hyperlink for URL PDFFile")
				Expect(target).To(Equal("https://sbi.com/jun23.pdf"))
				val4, err := f.GetCellValue(sheetName, "E4")
				Expect(err).ToNot(HaveOccurred())
				Expect(val4).To(Equal("PDF"))
			})
		})

		Context("when TTMonthEndRates slice is empty", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "TT Rates"
			)

			BeforeEach(func() {
				sampleSummary = tax.Summary{
					Year:            testYear,
					TTMonthEndRates: []tax.MonthEndRate{},
				}
			})

			It("should create the 'TT Rates' sheet with only headers", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1), "Sheet should only contain the header row")

				expectedHeaders := []string{
					"Month", "Year", "TTDate", "TTRate", "PDF Link", "DayOfWeek",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when the tax summary is completely empty", func() {
			It("should create a valid Excel file with all sheets containing only headers", func() {
				emptySummary := tax.Summary{}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())
				Expect(tempOutputFilePath).Should(BeARegularFile())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				// Check Gains Sheet
				rows, err := f.GetRows("Gains")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))
				expectedGainsHeaders := []string{
					"Symbol", "BuyDate", "SellDate", "Quantity", "PNL (USD)",
					"Commission (USD)", "Type", "TTDate", "TTRate", "PNL (INR)",
				}
				Expect(rows[0]).To(Equal(expectedGainsHeaders))

				// Check Dividends Sheet
				rows, err = f.GetRows("Dividends")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))
				expectedDividendsHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)", "TTDate", "TTRate",
					"Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedDividendsHeaders))

				// Check Valuations Sheet
				rows, err = f.GetRows("Valuations")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))
				expectedValuationsHeaders := []string{
					"Symbol",
					"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"AmountPaid (INR)",
				}
				Expect(rows[0]).To(Equal(expectedValuationsHeaders))

				// Check Interest Sheet
				rows, err = f.GetRows("Interest")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))
				expectedInterestHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
					"TTDate", "TTRate", "Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedInterestHeaders))

				// Check TT Rates Sheet
				rows, err = f.GetRows("TT Rates")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1), "TT Rates sheet should only contain the header row")
				expectedTTRatesHeaders := []string{
					"Month", "Year", "TTDate", "TTRate", "PDF Link", "DayOfWeek",
				}
				Expect(rows[0]).To(Equal(expectedTTRatesHeaders))
			})

			It("should create a file with exactly 6 sheets and no 'Sheet1'", func() {
				emptySummary := tax.Summary{}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheets := f.GetSheetList()
				Expect(sheets).To(HaveLen(6), "There should be exactly 6 sheets")
				Expect(sheets).To(ConsistOf("Summary", "Gains", "Dividends", "Valuations", "Interest", "TT Rates"))
			})
		})

		Context("regarding file system operations", func() {
			It("should create parent directories for the output file if they do not exist", func() {
				nestedDirPath := filepath.Join(baseTempDir, "reports_test", "fy_temp")
				specificOutputFilePath := filepath.Join(nestedDirPath, fmt.Sprintf("%d_Tax_Summary.xlsx", testYear))

				fileExcelManager := manager.NewExcelManager(nestedDirPath)

				err := fileExcelManager.GenerateTaxSummaryExcel(ctx, testYear, tax.Summary{})
				Expect(err).ToNot(HaveOccurred())

				Expect(nestedDirPath).Should(BeADirectory())
				Expect(specificOutputFilePath).Should(BeARegularFile())
			})

			It("should return an error if saving the file to an invalid target fails", func() {
				// Use a directory path that will make SaveAs fail - create a directory with the same name as the expected file
				invalidDir := filepath.Join(baseTempDir, "invalid_test")
				invalidFile := filepath.Join(invalidDir, fmt.Sprintf("%d_Tax_Summary.xlsx", testYear))

				// Create a directory where the file should be created
				err := os.MkdirAll(invalidFile, 0755)
				Expect(err).ToNot(HaveOccurred())

				fileExcelManager := manager.NewExcelManager(invalidDir)

				err = fileExcelManager.GenerateTaxSummaryExcel(ctx, testYear, tax.Summary{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when generating Summary sheet with all sections", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error

				// Create Gains data (3 records: 1 STCG, 2 LTCG)
				gain1TTDate, _ := time.Parse(time.DateOnly, "2023-01-15")
				gain1 := tax.INRGains{
					Gains:  tax.Gains{Symbol: "AAPL", BuyDate: "2022-10-01", SellDate: "2023-01-20", Quantity: 10.5, PNL: 100.75, Commission: 5.25, Type: "STCG"},
					TTDate: gain1TTDate, TTRate: 82.50,
				}
				gain2TTDate, _ := time.Parse(time.DateOnly, "2023-02-10")
				gain2 := tax.INRGains{
					Gains:  tax.Gains{Symbol: "MSFT", BuyDate: "2020-05-01", SellDate: "2023-02-15", Quantity: 5, PNL: -50.20, Commission: 0, Type: "LTCG"},
					TTDate: gain2TTDate, TTRate: 83.10,
				}
				gain3WithZeroTTDate := tax.INRGains{
					Gains:  tax.Gains{Symbol: "GOOG", BuyDate: "2021-01-01", SellDate: "2023-03-10", Quantity: 20, PNL: 200.00, Commission: 1.50, Type: "LTCG"},
					TTDate: time.Time{}, TTRate: 81.75,
				}

				// Create Dividends data (2 records)
				div1TTDate, _ := time.Parse(time.DateOnly, "2023-03-01")
				div1 := tax.INRDividend{
					Dividend: tax.Dividend{
						Symbol: "AAPL", Date: "2023-03-15",
						Amount: 50.25, Tax: 7.54, Net: 42.71,
					},
					TTDate: div1TTDate, TTRate: 82.10,
				}
				div2TTDate, _ := time.Parse(time.DateOnly, "2023-03-05")
				div2 := tax.INRDividend{
					Dividend: tax.Dividend{
						Symbol: "GOOG", Date: "2023-03-20",
						Amount: 75.50, Tax: 11.33, Net: 64.17,
					},
					TTDate: div2TTDate, TTRate: 82.50,
				}

				// Create Interest data (1 record)
				interestDate, _ := time.Parse(time.DateOnly, "2023-06-15")
				ttDate, _ := time.Parse(time.DateOnly, "2023-06-30")
				interest1 := tax.INRInterest{
					Interest: tax.Interest{
						Symbol: "US-TBILL", Date: interestDate.Format(time.DateOnly),
						Amount: 100.0, Tax: 10.0, Net: 90.0,
					},
					TTDate: ttDate, TTRate: 82.5,
				}

				// Build summary with ALL data types
				sampleSummary = tax.Summary{
					Year:          testYear,
					INRGains:      []tax.INRGains{gain1, gain2, gain3WithZeroTTDate},
					INRDividends:  []tax.INRDividend{div1, div2},
					INRInterest:   []tax.INRInterest{interest1},
					INRValuations: []tax.INRValuation{},
				}

				// Generate Excel
				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				// Open file for verification
				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					f.Close()
				}
			})

			It("should create Summary sheet as the first sheet", func() {
				sheets := f.GetSheetList()
				Expect(sheets).To(HaveLen(6)) // Summary, Gains, Dividends, Valuations, Interest, TT Rates
				Expect(sheets[0]).To(Equal("Summary"))
			})

			It("should write SUMMARY header in row 1", func() {
				header, err := f.GetCellValue(sheetName, "A1")
				Expect(err).ToNot(HaveOccurred())
				Expect(header).To(Equal("SUMMARY"))
			})

			Context("Short Term Gains section", func() {
				It("should write section header and column headers in rows 3-4", func() {
					// Section header (row 3)
					sectionHeader, err := f.GetCellValue(sheetName, "A3")
					Expect(err).ToNot(HaveOccurred())
					Expect(sectionHeader).To(Equal("Short Term"))

					// Column headers (row 4)
					headers := []struct{ cell, expected string }{
						{"A4", "PNL (USD)"},
						{"B4", "Commission (USD)"},
						{"C4", "PNL (INR)"},
					}
					for _, h := range headers {
						val, err := f.GetCellValue(sheetName, h.cell)
						Expect(err).ToNot(HaveOccurred())
						Expect(val).To(Equal(h.expected))
					}
				})

				It("should write formulas referencing Gains STCG row (row 5)", func() {
					// Gains sheet: 3 gains → lastDataRow=4, totalsRow=6, stcgRow=8
					expectCrossSheetFormula(f, "Summary", "A5", "Gains", "E8") // PNL USD
					expectCrossSheetFormula(f, "Summary", "B5", "Gains", "F6") // Commission from TOTALS
					expectCrossSheetFormula(f, "Summary", "C5", "Gains", "J8") // PNL INR
				})
			})

			Context("Long Term Gains section", func() {
				It("should write section header and column headers in rows 8-9", func() {
					// Section header (row 8)
					sectionHeader, err := f.GetCellValue(sheetName, "A8")
					Expect(err).ToNot(HaveOccurred())
					Expect(sectionHeader).To(Equal("Long Term"))

					// Column headers (row 9)
					headers := []struct{ cell, expected string }{
						{"A9", "PNL (USD)"},
						{"B9", "Commission (USD)"},
						{"C9", "PNL (INR)"},
					}
					for _, h := range headers {
						val, err := f.GetCellValue(sheetName, h.cell)
						Expect(err).ToNot(HaveOccurred())
						Expect(val).To(Equal(h.expected))
					}
				})

				It("should write formulas referencing Gains LTCG row (row 10)", func() {
					// Gains sheet: ltcgRow=9
					expectCrossSheetFormula(f, "Summary", "A10", "Gains", "E9") // PNL USD
					expectCrossSheetFormula(f, "Summary", "B10", "Gains", "F6") // Commission from TOTALS
					expectCrossSheetFormula(f, "Summary", "C10", "Gains", "J9") // PNL INR
				})
			})

			Context("Dividends section", func() {
				It("should write section header and column headers in rows 13-14", func() {
					// Section header (row 13)
					sectionHeader, err := f.GetCellValue(sheetName, "A13")
					Expect(err).ToNot(HaveOccurred())
					Expect(sectionHeader).To(Equal("Dividends"))

					// Column headers (row 14)
					headers := []struct{ cell, expected string }{
						{"A14", "Amount (USD)"},
						{"B14", "Tax (USD)"},
						{"C14", "Net (USD)"},
						{"D14", "Amount (INR)"},
						{"E14", "Tax (INR)"},
						{"F14", "Net (INR)"},
					}
					for _, h := range headers {
						val, err := f.GetCellValue(sheetName, h.cell)
						Expect(err).ToNot(HaveOccurred())
						Expect(val).To(Equal(h.expected))
					}
				})

				It("should write formulas referencing Dividends TOTALS row (row 15)", func() {
					// Dividends sheet: 2 dividends → lastDataRow=3, totalsRow=5
					expectCrossSheetFormula(f, "Summary", "A15", "Dividends", "C5") // Amount USD
					expectCrossSheetFormula(f, "Summary", "B15", "Dividends", "D5") // Tax USD
					expectCrossSheetFormula(f, "Summary", "C15", "Dividends", "E5") // Net USD
					expectCrossSheetFormula(f, "Summary", "D15", "Dividends", "H5") // Amount INR
					expectCrossSheetFormula(f, "Summary", "E15", "Dividends", "I5") // Tax INR
					expectCrossSheetFormula(f, "Summary", "F15", "Dividends", "J5") // Net INR
				})
			})

			Context("Interest section", func() {
				It("should write section header and column headers in rows 18-19", func() {
					// Section header (row 18)
					sectionHeader, err := f.GetCellValue(sheetName, "A18")
					Expect(err).ToNot(HaveOccurred())
					Expect(sectionHeader).To(Equal("Interest Income"))

					// Column headers (row 19)
					headers := []struct{ cell, expected string }{
						{"A19", "Amount (USD)"},
						{"B19", "Tax (USD)"},
						{"C19", "Net (USD)"},
						{"D19", "Amount (INR)"},
						{"E19", "Tax (INR)"},
						{"F19", "Net (INR)"},
					}
					for _, h := range headers {
						val, err := f.GetCellValue(sheetName, h.cell)
						Expect(err).ToNot(HaveOccurred())
						Expect(val).To(Equal(h.expected))
					}
				})

				It("should write formulas referencing Interest TOTALS row (row 20)", func() {
					// Interest sheet: 1 interest → lastDataRow=2, totalsRow=4
					expectCrossSheetFormula(f, "Summary", "A20", "Interest", "C4") // Amount USD
					expectCrossSheetFormula(f, "Summary", "B20", "Interest", "D4") // Tax USD
					expectCrossSheetFormula(f, "Summary", "C20", "Interest", "E4") // Net USD
					expectCrossSheetFormula(f, "Summary", "D20", "Interest", "H4") // Amount INR
					expectCrossSheetFormula(f, "Summary", "E20", "Interest", "I4") // Tax INR
					expectCrossSheetFormula(f, "Summary", "F20", "Interest", "J4") // Net INR
				})
			})
		})

		Context("when generating Summary sheet with only Gains data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error

				// Create Gains data only
				gain1TTDate, _ := time.Parse(time.DateOnly, "2023-01-15")
				gain1 := tax.INRGains{
					Gains:  tax.Gains{Symbol: "AAPL", BuyDate: "2022-10-01", SellDate: "2023-01-20", Quantity: 10.5, PNL: 100.75, Commission: 5.25, Type: "STCG"},
					TTDate: gain1TTDate, TTRate: 82.50,
				}

				sampleSummary = tax.Summary{
					Year:          testYear,
					INRGains:      []tax.INRGains{gain1},
					INRDividends:  []tax.INRDividend{},
					INRInterest:   []tax.INRInterest{},
					INRValuations: []tax.INRValuation{},
				}

				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					f.Close()
				}
			})

			It("should write Gains sections starting at row 3 when other sections are empty", func() {
				// Short Term starts at row 3 (no Dividends/Interest before it)
				sectionHeader, err := f.GetCellValue(sheetName, "A3")
				Expect(err).ToNot(HaveOccurred())
				Expect(sectionHeader).To(Equal("Short Term"))

				// Long Term starts at row 8
				sectionHeader, err = f.GetCellValue(sheetName, "A8")
				Expect(err).ToNot(HaveOccurred())
				Expect(sectionHeader).To(Equal("Long Term"))
			})
		})

		Context("when generating Summary sheet with only Dividends data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error

				// Create Dividends data only
				div1TTDate, _ := time.Parse(time.DateOnly, "2023-03-01")
				div1 := tax.INRDividend{
					Dividend: tax.Dividend{
						Symbol: "AAPL", Date: "2023-03-15",
						Amount: 50.25, Tax: 7.54, Net: 42.71,
					},
					TTDate: div1TTDate, TTRate: 82.10,
				}

				sampleSummary = tax.Summary{
					Year:          testYear,
					INRGains:      []tax.INRGains{},
					INRDividends:  []tax.INRDividend{div1},
					INRInterest:   []tax.INRInterest{},
					INRValuations: []tax.INRValuation{},
				}

				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					f.Close()
				}
			})

			It("should write Dividends section starting at row 3 when Gains/Interest are empty", func() {
				// Dividends starts at row 3 (no Gains before it)
				sectionHeader, err := f.GetCellValue(sheetName, "A3")
				Expect(err).ToNot(HaveOccurred())
				Expect(sectionHeader).To(Equal("Dividends"))
			})
		})

		Context("when generating Summary sheet with only Interest data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error

				// Create Interest data only
				interestDate, _ := time.Parse(time.DateOnly, "2023-06-15")
				ttDate, _ := time.Parse(time.DateOnly, "2023-06-30")
				interest1 := tax.INRInterest{
					Interest: tax.Interest{
						Symbol: "US-TBILL", Date: interestDate.Format(time.DateOnly),
						Amount: 100.0, Tax: 10.0, Net: 90.0,
					},
					TTDate: ttDate, TTRate: 82.5,
				}

				sampleSummary = tax.Summary{
					Year:          testYear,
					INRGains:      []tax.INRGains{},
					INRDividends:  []tax.INRDividend{},
					INRInterest:   []tax.INRInterest{interest1},
					INRValuations: []tax.INRValuation{},
				}

				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					f.Close()
				}
			})

			It("should write Interest section starting at row 3 when Gains/Dividends are empty", func() {
				// Interest starts at row 3 (no Gains/Dividends before it)
				sectionHeader, err := f.GetCellValue(sheetName, "A3")
				Expect(err).ToNot(HaveOccurred())
				Expect(sectionHeader).To(Equal("Interest Income"))
			})
		})

		Context("when generating Summary sheet with empty data", func() {
			var (
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error

				// All empty
				sampleSummary = tax.Summary{
					Year:          testYear,
					INRGains:      []tax.INRGains{},
					INRDividends:  []tax.INRDividend{},
					INRInterest:   []tax.INRInterest{},
					INRValuations: []tax.INRValuation{},
				}

				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err = excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				if f != nil {
					f.Close()
				}
			})

			It("should create Summary sheet with header only when all data is empty", func() {
				sheets := f.GetSheetList()
				Expect(sheets[0]).To(Equal("Summary"))

				// Should have SUMMARY header
				header, err := f.GetCellValue(sheetName, "A1")
				Expect(err).ToNot(HaveOccurred())
				Expect(header).To(Equal("SUMMARY"))

				// Should NOT have any section headers (row 3 should be empty)
				val, err := f.GetCellValue(sheetName, "A3")
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(BeEmpty())
			})
		})
	})
})
