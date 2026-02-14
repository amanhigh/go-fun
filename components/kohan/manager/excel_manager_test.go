//nolint:dupl // Test files often have similar setup patterns
package manager_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		baseTempDir, err = os.MkdirTemp(os.TempDir(), "excel_manager_test_run_*")
		Expect(err).ToNot(HaveOccurred())
		tempOutputFilePath = filepath.Join(baseTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
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

	Describe("GenerateTaxSummaryExcel", func() {
		Context("when generating the 'Gains' sheet with data", func() {
			var (
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Gains"
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "gains_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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

				sheetFound := false
				for _, name := range f.GetSheetList() {
					if name == sheetName {
						sheetFound = true
						break
					}
				}
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

				// Verify gain1 (row 2 in Excel, index 1 in `rows` slice)
				gain1 := sampleSummary.INRGains[0]
				Expect(rows[1][0]).To(Equal(gain1.Symbol))
				Expect(rows[1][1]).To(Equal(gain1.BuyDate))
				Expect(rows[1][2]).To(Equal(gain1.SellDate))
				qty1, _ := getCellFloat(f, sheetName, "D2")
				Expect(qty1).To(BeNumerically("~", gain1.Quantity, 0.001))
				pnlUSD1, _ := getCellFloat(f, sheetName, "E2")
				Expect(pnlUSD1).To(BeNumerically("~", gain1.PNL, 0.001))
				comm1, _ := getCellFloat(f, sheetName, "F2")
				Expect(comm1).To(BeNumerically("~", gain1.Commission, 0.001))
				Expect(rows[1][6]).To(Equal(gain1.Type))
				Expect(rows[1][7]).To(Equal(gain1.TTDate.Format(time.DateOnly)))
				rate1, _ := getCellFloat(f, sheetName, "I2")
				Expect(rate1).To(BeNumerically("~", gain1.TTRate, 0.001))
				// Column J has formula =E2*I2 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J2", "=E2*I2", gain1.PNL*gain1.TTRate)

				// Verify gain2 (row 3 in Excel, index 2 in `rows` slice)
				gain2 := sampleSummary.INRGains[1]
				Expect(rows[2][0]).To(Equal(gain2.Symbol))
				Expect(rows[2][1]).To(Equal(gain2.BuyDate))
				Expect(rows[2][2]).To(Equal(gain2.SellDate))
				qty2, _ := getCellFloat(f, sheetName, "D3")
				Expect(qty2).To(BeNumerically("~", gain2.Quantity, 0.001))
				pnlUSD2, _ := getCellFloat(f, sheetName, "E3")
				Expect(pnlUSD2).To(BeNumerically("~", gain2.PNL, 0.001))
				comm2, _ := getCellFloat(f, sheetName, "F3")
				Expect(comm2).To(BeNumerically("~", gain2.Commission, 0.001))
				Expect(rows[2][6]).To(Equal(gain2.Type))
				Expect(rows[2][7]).To(Equal(gain2.TTDate.Format(time.DateOnly)))
				rate2, _ := getCellFloat(f, sheetName, "I3")
				Expect(rate2).To(BeNumerically("~", gain2.TTRate, 0.001))
				// Column J has formula =E3*I3 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J3", "=E3*I3", gain2.PNL*gain2.TTRate)

				// Verify gain3WithZeroTTDate (row 4 in Excel, index 3 in `rows` slice)
				gain3 := sampleSummary.INRGains[2]
				Expect(rows[3][0]).To(Equal(gain3.Symbol))
				Expect(rows[3][1]).To(Equal(gain3.BuyDate))
				Expect(rows[3][2]).To(Equal(gain3.SellDate))
				qty3, _ := getCellFloat(f, sheetName, "D4")
				Expect(qty3).To(BeNumerically("~", gain3.Quantity, 0.001))
				pnlUSD3, _ := getCellFloat(f, sheetName, "E4")
				Expect(pnlUSD3).To(BeNumerically("~", gain3.PNL, 0.001))
				comm3, _ := getCellFloat(f, sheetName, "F4")
				Expect(comm3).To(BeNumerically("~", gain3.Commission, 0.001))
				Expect(rows[3][6]).To(Equal(gain3.Type))
				Expect(rows[3][7]).To(Equal(""), "TTDate for zero time should be an empty string")
				rate3, _ := getCellFloat(f, sheetName, "I4")
				Expect(rate3).To(BeNumerically("~", gain3.TTRate, 0.001))
				// Column J has formula =E4*I4 and calculates PNL (INR)
				expectFormulaCell(f, sheetName, "J4", "=E4*I4", gain3.PNL*gain3.TTRate)
			})
		})

		Context("when INRGains slice is empty", func() {
			var (
				excelManager manager.ExcelManager
				emptySummary tax.Summary
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_gains_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(contextTempDir)
				emptySummary = tax.Summary{INRGains: []tax.INRGains{}}
			})

			It("should create the 'Gains' sheet with only headers", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())
				Expect(tempOutputFilePath).Should(BeARegularFile())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetName := "Gains"
				sheetFound := false
				for _, name := range f.GetSheetList() {
					if name == sheetName {
						sheetFound = true
						break
					}
				}
				Expect(sheetFound).To(BeTrue(), "Sheet 'Gains' should exist")

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1), "Sheet should only contain the header row")

				expectedHeaders := []string{
					"Symbol", "BuyDate", "SellDate", "Quantity", "PNL (USD)",
					"Commission (USD)", "Type", "TTDate", "TTRate", "PNL (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when generating the 'Dividends' sheet with data", func() {
			var (
				excelManager  manager.ExcelManager
				sampleSummary tax.Summary
				sheetName     = "Dividends"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "dividends_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(contextTempDir)

				div1TTDate, _ := time.Parse(time.DateOnly, "2023-04-05")
				div1 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "AAPL", Date: "2023-04-10", Amount: 50.25, Tax: 7.54, Net: 42.71},
					TTDate:   div1TTDate, TTRate: 82.10,
				}
				div2TTDate, _ := time.Parse(time.DateOnly, "2023-05-12")
				div2 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "GOOG", Date: "2023-05-15", Amount: 75.50, Tax: 11.33, Net: 64.17},
					TTDate:   div2TTDate, TTRate: 82.50,
				}
				sampleSummary.INRDividends = []tax.INRDividend{div1, div2}
			})

			It("should create the 'Dividends' sheet with correct headers and data", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				// Rows include: headers + data records + empty row + TOTALS row
				expectedRowCount := 1 + len(sampleSummary.INRDividends) + 2
				Expect(rows).To(HaveLen(expectedRowCount))

				// Verify Headers
				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)", "TTDate", "TTRate",
					"Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))

				// Verify first data row for accuracy (basic check)
				div1 := sampleSummary.INRDividends[0]
				Expect(rows[1][0]).To(Equal(div1.Symbol))
				Expect(rows[1][1]).To(Equal(div1.Date))
				amountUSD1, _ := getCellFloat(f, sheetName, "C2")
				Expect(amountUSD1).To(BeNumerically("~", div1.Amount, 0.001))
				taxUSD1, _ := getCellFloat(f, sheetName, "D2")
				Expect(taxUSD1).To(BeNumerically("~", div1.Tax, 0.001))
				netUSD1, _ := getCellFloat(f, sheetName, "E2")
				Expect(netUSD1).To(BeNumerically("~", div1.Net, 0.001))
				Expect(rows[1][5]).To(Equal(div1.TTDate.Format(time.DateOnly)))
				rate1, _ := getCellFloat(f, sheetName, "G2")
				Expect(rate1).To(BeNumerically("~", div1.TTRate, 0.001))
				// Verify formulas for INR columns
				expectFormulaCell(f, sheetName, "H2", "=C2*G2", div1.Amount*div1.TTRate) // Amount (INR)
				expectFormulaCell(f, sheetName, "I2", "=D2*G2", div1.Tax*div1.TTRate)    // Tax (INR)
				expectFormulaCell(f, sheetName, "J2", "=E2*G2", div1.Net*div1.TTRate)    // Net (INR)
			})

			It("should write TOTALS row with correct formulas and calculated values", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

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

		Context("when INRDividends slice is empty", func() {
			var (
				excelManager manager.ExcelManager
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_dividends_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
				excelManager = manager.NewExcelManager(contextTempDir)
			})

			It("should create the 'Dividends' sheet with only headers", func() {
				emptySummary := tax.Summary{INRDividends: []tax.INRDividend{}}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetName := "Dividends"
				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))

				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)", "TTDate", "TTRate",
					"Amount (INR)", "Tax (INR)", "Net (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when generating the 'Valuations' sheet with data", func() {
			var (
				excelManager  manager.ExcelManager
				sampleSummary tax.Summary
				sheetName     = "Valuations"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "valuations_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(contextTempDir)

				// Define dates
				firstDate, _ := time.Parse(time.DateOnly, "2022-01-10")
				firstTTDate, _ := time.Parse(time.DateOnly, "2022-01-11")
				peakDate, _ := time.Parse(time.DateOnly, "2022-11-25")
				yearEndDate, _ := time.Parse(time.DateOnly, "2023-03-31")

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
				defer f.Close()

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

				// Calculate expected value using lo.SumBy (1 total)
				totalAmountPaidINR := lo.SumBy(sampleSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.AmountPaid
				})

				// Verify TOTALS row position
				lastDataRow := len(sampleSummary.INRValuations) + 1 // Row 2 (1 data row)
				totalsRow := lastDataRow + 2                        // Row 4 (skip empty row 3)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

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

				// Calculate expected value using lo.SumBy (2 valuations)
				totalAmountPaidINR := lo.SumBy(nonZeroSummary.INRValuations, func(v tax.INRValuation) float64 {
					return v.AmountPaid
				})

				// Verify TOTALS row position
				lastDataRow := len(nonZeroSummary.INRValuations) + 1 // Row 3 (2 data rows)
				totalsRow := lastDataRow + 2                         // Row 5 (skip empty row 4)

				// Verify TOTALS label
				totalsLabel, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", totalsRow))
				Expect(err).ToNot(HaveOccurred())
				Expect(totalsLabel).To(Equal("TOTALS"))

				// Verify AmountPaid (INR) column W with non-zero total
				// Expected: 5432.10 + 3210.50 = 8642.60
				expectFormulaCell(f, sheetName, fmt.Sprintf("W%d", totalsRow),
					fmt.Sprintf("=SUM(W2:W%d)", lastDataRow), totalAmountPaidINR)
			})
		})

		Context("when INRValuations slice is empty", func() {
			var (
				excelManager manager.ExcelManager
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_valuations_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
				excelManager = manager.NewExcelManager(contextTempDir)
			})

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
				excelManager  manager.ExcelManager
				sampleSummary tax.Summary
				sheetName     = "Interest"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "interest_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(contextTempDir)

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
			var (
				excelManager manager.ExcelManager
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_interest_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
				excelManager = manager.NewExcelManager(contextTempDir)
			})

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

		Context("when the tax summary is completely empty", func() {
			var (
				excelManager manager.ExcelManager
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_summary_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
				excelManager = manager.NewExcelManager(contextTempDir)
			})

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
			})

			It("should create a file with exactly 4 sheets and no 'Sheet1'", func() {
				emptySummary := tax.Summary{}
				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheets := f.GetSheetList()
				Expect(sheets).To(HaveLen(5), "There should be exactly 5 sheets")
				Expect(sheets).To(ConsistOf("Summary", "Gains", "Dividends", "Valuations", "Interest"))
			})
		})

		Context("regarding file system operations", func() {
			It("should create parent directories for the output file if they do not exist", func() {
				nestedDirPath := filepath.Join(baseTempDir, "reports_test", "fy_temp")
				specificOutputFilePath := filepath.Join(nestedDirPath, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager := manager.NewExcelManager(nestedDirPath)

				err := excelManager.GenerateTaxSummaryExcel(ctx, testYear, tax.Summary{})
				Expect(err).ToNot(HaveOccurred())

				Expect(nestedDirPath).Should(BeADirectory())
				Expect(specificOutputFilePath).Should(BeARegularFile())
			})

			It("should return an error if saving the file to an invalid target fails", func() {
				// Use a directory path that will make SaveAs fail - create a directory with the same name as the expected file
				invalidDir := filepath.Join(baseTempDir, "invalid_test")
				invalidFile := filepath.Join(invalidDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				// Create a directory where the file should be created
				err := os.MkdirAll(invalidFile, 0755)
				Expect(err).ToNot(HaveOccurred())

				excelManager := manager.NewExcelManager(invalidDir)

				err = excelManager.GenerateTaxSummaryExcel(ctx, testYear, tax.Summary{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when generating Summary sheet with all sections", func() {
			var (
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "summary_all_sections_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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
				Expect(sheets).To(HaveLen(5)) // Summary, Gains, Dividends, Valuations, Interest
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
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "summary_only_gains_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "summary_only_dividends_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "summary_only_interest_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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
				excelManager  manager.ExcelManager
				tempOutputDir string
				sampleSummary tax.Summary
				sheetName     = "Summary"
				f             *excelize.File
			)

			BeforeEach(func() {
				var err error
				tempOutputDir, err = os.MkdirTemp(baseTempDir, "summary_empty_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(tempOutputDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))

				excelManager = manager.NewExcelManager(tempOutputDir)

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
