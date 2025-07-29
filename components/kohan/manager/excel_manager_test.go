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
	"github.com/xuri/excelize/v2"
)

var _ = Describe("ExcelManagerImpl", func() {
	var (
		ctx         context.Context
		baseTempDir string
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		baseTempDir, err = os.MkdirTemp(os.TempDir(), "excel_manager_test_run_*")
		Expect(err).ToNot(HaveOccurred())
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

	Describe("GenerateTaxSummaryExcel", func() {
		Context("when generating the 'Gains' sheet with data", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
				sampleSummary      tax.Summary
				sheetName          = "Gains"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "gains_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_with_data.xlsx")

				excelManager = manager.NewExcelManager(tempOutputFilePath)

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
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
				Expect(err).ToNot(HaveOccurred())
				Expect(tempOutputFilePath).Should(BeARegularFile())
			})

			It("should contain a 'Gains' sheet with the correct headers in order", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
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
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, errGetRows := f.GetRows(sheetName)
				Expect(errGetRows).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1+len(sampleSummary.INRGains)), "Number of rows should be headers + data records")

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
				pnlINR1, _ := getCellFloat(f, sheetName, "J2")
				Expect(pnlINR1).To(BeNumerically("~", gain1.INRValue(), 0.001))

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
				pnlINR2, _ := getCellFloat(f, sheetName, "J3")
				Expect(pnlINR2).To(BeNumerically("~", gain2.INRValue(), 0.001))

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
				pnlINR3, _ := getCellFloat(f, sheetName, "J4")
				Expect(pnlINR3).To(BeNumerically("~", gain3.INRValue(), 0.001))
			})
		})

		Context("when INRGains slice is empty", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
				emptySummary       tax.Summary
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_gains_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_empty_gains.xlsx")

				excelManager = manager.NewExcelManager(tempOutputFilePath)
				emptySummary = tax.Summary{INRGains: []tax.INRGains{}}
			})

			It("should create the 'Gains' sheet with only headers", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, emptySummary)
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
				excelManager       manager.ExcelManager
				tempOutputFilePath string
				sampleSummary      tax.Summary
				sheetName          = "Dividends"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "dividends_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_with_dividends.xlsx")

				excelManager = manager.NewExcelManager(tempOutputFilePath)

				div1TTDate, _ := time.Parse(time.DateOnly, "2023-04-05")
				div1 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "AAPL", Date: "2023-04-10", Amount: 50.25},
					TTDate:   div1TTDate, TTRate: 82.10,
				}
				div2TTDate, _ := time.Parse(time.DateOnly, "2023-05-12")
				div2 := tax.INRDividend{
					Dividend: tax.Dividend{Symbol: "GOOG", Date: "2023-05-15", Amount: 75.50},
					TTDate:   div2TTDate, TTRate: 82.50,
				}
				sampleSummary.INRDividends = []tax.INRDividend{div1, div2}
			})

			It("should create the 'Dividends' sheet with correct headers and data", func() {
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1 + len(sampleSummary.INRDividends)))

				// Verify Headers
				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "TTDate", "TTRate", "Amount (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))

				// Verify first data row for accuracy (basic check)
				div1 := sampleSummary.INRDividends[0]
				Expect(rows[1][0]).To(Equal(div1.Symbol))
				Expect(rows[1][1]).To(Equal(div1.Date))
				amountUSD1, _ := getCellFloat(f, sheetName, "C2")
				Expect(amountUSD1).To(BeNumerically("~", div1.Amount, 0.001))
				Expect(rows[1][3]).To(Equal(div1.TTDate.Format(time.DateOnly)))
				rate1, _ := getCellFloat(f, sheetName, "E2")
				Expect(rate1).To(BeNumerically("~", div1.TTRate, 0.001))
				amountINR1, _ := getCellFloat(f, sheetName, "F2")
				Expect(amountINR1).To(BeNumerically("~", div1.INRValue(), 0.001))
			})
		})

		Context("when INRDividends slice is empty", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
			)
			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_dividends_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_empty_dividends.xlsx")
				excelManager = manager.NewExcelManager(tempOutputFilePath)
			})

			It("should create the 'Dividends' sheet with only headers", func() {
				emptySummary := tax.Summary{INRDividends: []tax.INRDividend{}}
				err := excelManager.GenerateTaxSummaryExcel(ctx, emptySummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				sheetName := "Dividends"
				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))

				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "TTDate", "TTRate", "Amount (INR)",
				}
				Expect(rows[0]).To(Equal(expectedHeaders))
			})
		})

		Context("when generating the 'Valuations' sheet with data", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
				sampleSummary      tax.Summary
				sheetName          = "Valuations"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "valuations_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_with_valuations.xlsx")

				excelManager = manager.NewExcelManager(tempOutputFilePath)

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
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(2)) // Header + 1 Data Row

				// Verify Headers
				expectedHeaders := []string{
					"Symbol",
					"Date (First)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (Peak)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
					"Date (YearEnd)", "Qty", "Price", "ValUSD", "TTDate", "TTRate", "ValINR",
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
				valUSD, _ := getCellFloat(f, sheetName, "E2")
				Expect(valUSD).To(Equal(posFirst.USDValue()))
				Expect(rows[1][5]).To(Equal(posFirst.TTDate.Format(time.DateOnly)))
				rate, _ := getCellFloat(f, sheetName, "G2")
				Expect(rate).To(Equal(posFirst.TTRate))
				valINR, _ := getCellFloat(f, sheetName, "H2")
				Expect(valINR).To(Equal(posFirst.INRValue()))

				// Peak Position
				posPeak := val1.PeakPosition
				Expect(rows[1][8]).To(Equal(posPeak.Date.Format(time.DateOnly)))
				qty, _ = getCellFloat(f, sheetName, "J2")
				Expect(qty).To(Equal(posPeak.Quantity))
				price, _ = getCellFloat(f, sheetName, "K2")
				Expect(price).To(Equal(posPeak.USDPrice))
				valUSD, _ = getCellFloat(f, sheetName, "L2")
				Expect(valUSD).To(Equal(posPeak.USDValue()))
				Expect(rows[1][12]).To(Equal(posPeak.TTDate.Format(time.DateOnly)))
				rate, _ = getCellFloat(f, sheetName, "N2")
				Expect(rate).To(Equal(posPeak.TTRate))
				valINR, _ = getCellFloat(f, sheetName, "O2")
				Expect(valINR).To(Equal(posPeak.INRValue()))

				// Year End Position
				posYearEnd := val1.YearEndPosition
				Expect(rows[1][15]).To(Equal(posYearEnd.Date.Format(time.DateOnly)))
				qty, _ = getCellFloat(f, sheetName, "Q2")
				Expect(qty).To(Equal(posYearEnd.Quantity))
				price, _ = getCellFloat(f, sheetName, "R2")
				Expect(price).To(Equal(posYearEnd.USDPrice))
				valUSD, _ = getCellFloat(f, sheetName, "S2")
				Expect(valUSD).To(Equal(posYearEnd.USDValue()))
				Expect(rows[1][19]).To(Equal(posYearEnd.TTDate.Format(time.DateOnly)))
				rate, _ = getCellFloat(f, sheetName, "U2")
				Expect(rate).To(Equal(posYearEnd.TTRate))
				valINR, _ = getCellFloat(f, sheetName, "V2")
				Expect(valINR).To(BeNumerically("~", posYearEnd.INRValue(), 0.001))
			})
		})

		Context("when generating the 'Interest' sheet with data", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
				sampleSummary      tax.Summary
				sheetName          = "Interest"
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "interest_data_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "summary_with_interest.xlsx")

				excelManager = manager.NewExcelManager(tempOutputFilePath)

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
				err := excelManager.GenerateTaxSummaryExcel(ctx, sampleSummary)
				Expect(err).ToNot(HaveOccurred())

				f, err := excelize.OpenFile(tempOutputFilePath)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()

				rows, err := f.GetRows(sheetName)
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(2)) // Header + 1 Data Row

				// Verify Headers
				expectedHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
					"TTDate", "TTRate", "Amount (INR)",
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
				inrValue, _ := getCellFloat(f, sheetName, "H2")
				Expect(inrValue).To(BeNumerically("~", interest1.INRValue(), 0.001))
			})
		})

		Context("when the tax summary is completely empty", func() {
			var (
				excelManager       manager.ExcelManager
				tempOutputFilePath string
			)

			BeforeEach(func() {
				contextTempDir, err := os.MkdirTemp(baseTempDir, "empty_summary_test_*")
				Expect(err).ToNot(HaveOccurred())
				tempOutputFilePath = filepath.Join(contextTempDir, "empty_summary.xlsx")
				excelManager = manager.NewExcelManager(tempOutputFilePath)
			})

			It("should create a valid Excel file with all sheets containing only headers", func() {
				emptySummary := tax.Summary{}
				err := excelManager.GenerateTaxSummaryExcel(ctx, emptySummary)
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
					"Symbol", "Date", "Amount (USD)", "TTDate", "TTRate", "Amount (INR)",
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
				}
				Expect(rows[0]).To(Equal(expectedValuationsHeaders))

				// Check Interest Sheet
				rows, err = f.GetRows("Interest")
				Expect(err).ToNot(HaveOccurred())
				Expect(rows).To(HaveLen(1))
				expectedInterestHeaders := []string{
					"Symbol", "Date", "Amount (USD)", "Tax (USD)", "Net (USD)",
					"TTDate", "TTRate", "Amount (INR)",
				}
				Expect(rows[0]).To(Equal(expectedInterestHeaders))
			})
		})

		Context("regarding file system operations", func() {
			It("should create parent directories for the output file if they do not exist", func() {
				nestedDirPath := filepath.Join(baseTempDir, "reports_test", "fy_temp")
				specificOutputFilePath := filepath.Join(nestedDirPath, "final_summary.xlsx")

				excelManager := manager.NewExcelManager(specificOutputFilePath)

				err := excelManager.GenerateTaxSummaryExcel(ctx, tax.Summary{})
				Expect(err).ToNot(HaveOccurred())

				Expect(nestedDirPath).Should(BeADirectory())
				Expect(specificOutputFilePath).Should(BeARegularFile())
			})

			It("should return an error if saving the file to an invalid target fails", func() {
				invalidOutputFilePath := baseTempDir
				excelManager := manager.NewExcelManager(invalidOutputFilePath)

				err := excelManager.GenerateTaxSummaryExcel(ctx, tax.Summary{})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
