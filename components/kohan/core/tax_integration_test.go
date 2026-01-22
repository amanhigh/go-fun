package core_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort" // Add import
	"time"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("Tax Integration", Label("it"), func() {
	var (
		ctx         context.Context
		taxManager  manager.TaxManager
		testYear    = 2023
		kohanConfig config.KohanConfig
	)

	BeforeEach(func() {
		ctx = context.Background()
		testDataBasePath := filepath.Join("..", "testdata", "tax")

		kohanConfig = config.KohanConfig{
			Tax: config.TaxConfig{
				TaxDir: testDataBasePath,
				// Layer 1: Input - Raw broker statements
				DriveWealthBase: filepath.Join(testDataBasePath, "Input", "Brokerage", "vested"),
				IBKRBase:        filepath.Join(testDataBasePath, "Input", "Brokerage", "ibkr"),
				// Layer 2.5: Parsed - Generated from broker statements
				ParsedDir:        filepath.Join(testDataBasePath, "Input", "Parsed"),
				TradesPath:       filepath.Join(testDataBasePath, "Input", "Parsed", tax.TRADES_FILENAME),
				DividendFilePath: filepath.Join(testDataBasePath, "Input", "Parsed", tax.DIVIDENDS_FILENAME),
				InterestFilePath: filepath.Join(testDataBasePath, "Input", "Parsed", tax.INTEREST_FILENAME),
				// Layer 3: Reference data (tickers, exchange rates)
				TickerCacheDir: filepath.Join(testDataBasePath, "Data", "Tickers"),
				TTRateFilePath: filepath.Join(testDataBasePath, "Data", "Reference", tax.SBI_RATES_FILENAME),
				// Layer 4: Output - Computed and generated results
				GainsFilePath: filepath.Join(testDataBasePath, "Output", "Computed", tax.GAINS_FILENAME),
				AccountsDir:   filepath.Join(testDataBasePath, "Output", "YearEndBalance"),
				ReportsDir:    filepath.Join(testDataBasePath, "Output", "Reports"),
				ComputedDir:   filepath.Join(testDataBasePath, "Output", "Computed"),
			},
		}

		core.SetupKohanInjector(kohanConfig)
		var err error
		taxManager, err = core.GetKohanInterface().GetTaxManager()
		Expect(err).ToNot(HaveOccurred())
		Expect(taxManager).ToNot(BeNil())
	})

	// NEW CONTEXT: Basic check
	Context("Basic Summary Retrieval", func() {
		It("should retrieve the tax summary without errors", func() {
			summary, err := taxManager.GetTaxSummary(ctx, testYear)

			Expect(err).ToNot(HaveOccurred()) // Verify no error during retrieval
			Expect(summary).ToNot(BeNil())    // Verify the summary object itself is not nil
		})
	})

	Context("Capital Gains Calculation (INRGains)", func() {
		It("should calculate capital gains correctly for the given year", func() {
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).ToNot(BeNil())
			Expect(summary.INRGains).To(HaveLen(3)) // Expecting AAPL (STCG), ADI (STCG), and MSFT (LTCG)

			// --- Assertions for AAPL (Expected at index 0 after sort) ---
			aaplGain := summary.INRGains[0]
			Expect(aaplGain.Symbol).To(Equal("AAPL"))
			Expect(aaplGain.PNL).To(Equal(1000.00))
			Expect(aaplGain.Type).To(Equal("STCG")) // Holding < 730 days
			Expect(aaplGain.BuyDate).To(Equal("2024-01-15"))
			Expect(aaplGain.SellDate).To(Equal("2024-01-17"))
			Expect(aaplGain.TTRate).To(Equal(82.00))
			Expect(aaplGain.INRValue()).To(Equal(82000.00)) // 1000.00 * 82.00
			Expect(aaplGain.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))

			// --- Assertions for ADI (Expected at index 1) ---
			adiGain := summary.INRGains[1]
			Expect(adiGain.Symbol).To(Equal("ADI"))
			Expect(adiGain.PNL).To(Equal(25.70))
			Expect(adiGain.Type).To(Equal("STCG"))
			Expect(adiGain.BuyDate).To(Equal("2023-06-15"))
			Expect(adiGain.SellDate).To(Equal("2023-08-31"))
			Expect(adiGain.TTRate).To(Equal(82.50))
			Expect(adiGain.INRValue()).To(Equal(2120.25)) // 25.70 * 82.50
			Expect(adiGain.TTDate.Format(time.DateOnly)).To(Equal("2023-07-10"))

			// --- Assertions for MSFT (Expected at index 2) ---
			msftGain := summary.INRGains[2]
			Expect(msftGain.Symbol).To(Equal("MSFT"))
			Expect(msftGain.PNL).To(Equal(500.00))
			Expect(msftGain.Type).To(Equal("LTCG")) // Holding > 730 days
			Expect(msftGain.BuyDate).To(Equal("2022-01-10"))
			Expect(msftGain.SellDate).To(Equal("2024-02-15"))
			// Updated assertions for MSFT based on new logic (rate from 2024-01-17, as 2024-01-31 is missing)
			Expect(msftGain.TTRate).To(Equal(82.90))
			Expect(msftGain.INRValue()).To(Equal(41450.00)) // 500.00 * 82.90
			Expect(msftGain.TTDate.Format(time.DateOnly)).To(Equal("2024-01-17"))
		})
	})

	Context("Dividend Calculation (INRDividends)", func() {
		It("should calculate dividends correctly for multiple symbols, filtering by financial year", func() {
			// Retrieve the summary for the test year (FY 2023-24)
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).ToNot(BeNil())
			Expect(summary.INRDividends).ToNot(BeNil()) // Ensure the slice itself is not nil

			// --- Assertions for Dividends (FY 2023-24: 2023-04-01 to 2024-03-31) ---
			// Based on testdata: AAPL Jun 15, ADI Jul 10, AAPL Dec 30, AAPL Jan 15, MSFT Feb 20, AAPL Mar 15 fall in this FY.
			// AAPL Apr 15 should be filtered out.
			Expect(summary.INRDividends).To(HaveLen(6)) // Expecting 6 dividends after filtering

			// Sort results by date to ensure consistent order for assertions
			sort.Slice(summary.INRDividends, func(i, j int) bool {
				dateI, err := summary.INRDividends[i].GetDate()
				Expect(err).NotTo(HaveOccurred())
				dateJ, err := summary.INRDividends[j].GetDate()
				Expect(err).NotTo(HaveOccurred())
				return dateI.Before(dateJ)
			})

			// --- Assertions for Jun 15, 2023 Dividend (AAPL) - First after sorting ---
			junDividend := summary.INRDividends[0]
			Expect(junDividend.Symbol).To(Equal("AAPL"))
			Expect(junDividend.Date).To(Equal("2023-06-15"))
			Expect(junDividend.Amount).To(Equal(50.00))
			Expect(junDividend.Tax).To(Equal(7.50))
			Expect(junDividend.Net).To(Equal(42.50))
			Expect(junDividend.TTRate).To(Equal(82.00)) // May 2023 month-end
			Expect(junDividend.TTDate.Format(time.DateOnly)).To(Equal("2023-05-01"))
			Expect(junDividend.INRValue()).To(Equal(4100.00))

			// --- Assertions for Jul 10, 2023 Dividend (ADI) ---
			julDividend := summary.INRDividends[1]
			Expect(julDividend.Symbol).To(Equal("ADI"))
			Expect(julDividend.Date).To(Equal("2023-07-10"))
			Expect(julDividend.Amount).To(Equal(30.00))
			Expect(julDividend.Tax).To(Equal(4.50))
			Expect(julDividend.Net).To(Equal(25.50))
			Expect(julDividend.TTRate).To(Equal(82.10)) // Jun 2023 month-end
			Expect(julDividend.TTDate.Format(time.DateOnly)).To(Equal("2023-06-15"))
			Expect(julDividend.INRValue()).To(Equal(2463.00))

			// --- Assertions for Dec 30, 2023 Dividend (AAPL) ---
			decDividend := summary.INRDividends[2]
			Expect(decDividend.Symbol).To(Equal("AAPL"))
			Expect(decDividend.Date).To(Equal("2023-12-30"))
			Expect(decDividend.Amount).To(Equal(60.00))
			Expect(decDividend.Tax).To(Equal(9.00))
			Expect(decDividend.Net).To(Equal(51.00))
			Expect(decDividend.TTRate).To(Equal(83.20)) // Nov 2023 month-end
			Expect(decDividend.TTDate.Format(time.DateOnly)).To(Equal("2023-11-15"))
			Expect(decDividend.INRValue()).To(Equal(4992.00))

			// --- Assertions for Jan 15, 2024 Dividend (AAPL) - Full Detail ---
			janDividend := summary.INRDividends[3]
			Expect(janDividend.Symbol).To(Equal("AAPL"))
			Expect(janDividend.Date).To(Equal("2024-01-15"))
			Expect(janDividend.Amount).To(Equal(115.00))
			Expect(janDividend.Tax).To(Equal(17.25))
			Expect(janDividend.Net).To(Equal(97.75))
			Expect(janDividend.TTRate).To(Equal(82.00))
			Expect(janDividend.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(janDividend.INRValue()).To(Equal(9430.00))

			// --- Assertions for Feb 20, 2024 Dividend (MSFT) - Key Details ---
			febDividend := summary.INRDividends[4]
			Expect(febDividend.Symbol).To(Equal("MSFT"))
			Expect(febDividend.Amount).To(Equal(50.00))
			Expect(febDividend.Tax).To(Equal(7.50))
			Expect(febDividend.Net).To(Equal(42.50))
			Expect(febDividend.TTRate).To(Equal(82.90))
			Expect(febDividend.TTDate.Format(time.DateOnly)).To(Equal("2024-01-17"))
			Expect(febDividend.INRValue()).To(Equal(4145.00))

			// --- Assertions for Mar 15, 2024 Dividend (AAPL) - Key Details ---
			marDividend := summary.INRDividends[5]
			Expect(marDividend.Symbol).To(Equal("AAPL"))
			Expect(marDividend.Amount).To(Equal(100.00))
			Expect(marDividend.Tax).To(Equal(15.00))
			Expect(marDividend.Net).To(Equal(85.00))
			Expect(marDividend.TTRate).To(Equal(83.05))
			Expect(marDividend.TTDate.Format(time.DateOnly)).To(Equal("2024-02-20"))
			Expect(marDividend.INRValue()).To(Equal(8305.00))
		})
	})

	Context("Interest Calculation (INRInterest)", func() {
		It("should calculate interest correctly, filtering by financial year", func() {
			// Retrieve the summary for the test year (FY 2023-24)
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).ToNot(BeNil())
			Expect(summary.INRInterest).ToNot(BeNil()) // Ensure the slice itself is not nil

			// --- Assertions for Interest (FY 2023-24: 2023-04-01 to 2024-03-31) ---
			// Based on assumed testdata: MSFT Dec 31, AAPL Jan 10 fall in this FY.
			// AAPL May 10 should be filtered out.
			Expect(summary.INRInterest).To(HaveLen(2)) // Expecting 2 interest records after filtering

			// --- Assertions for Dec 31 Interest (MSFT) - Full Detail ---
			decInterest := summary.INRInterest[0]
			Expect(decInterest.Symbol).To(Equal("MSFT"))
			Expect(decInterest.Date).To(Equal("2023-12-31"))
			Expect(decInterest.Amount).To(Equal(20.00))
			Expect(decInterest.Tax).To(Equal(4.00))
			Expect(decInterest.Net).To(Equal(16.00))
			Expect(decInterest.TTRate).To(Equal(83.20)) // Nov 2023 month-end rate (preceding month)
			Expect(decInterest.TTDate.Format(time.DateOnly)).To(Equal("2023-11-15"))
			Expect(decInterest.INRValue()).To(Equal(1664.00)) // 20.00 * 83.20

			// --- Assertions for Jan 10 Interest (AAPL) - Key Details ---
			janInterest := summary.INRInterest[1]
			Expect(janInterest.Symbol).To(Equal("AAPL"))
			Expect(janInterest.Amount).To(Equal(5.50))  // Check Amount for AAPL
			Expect(janInterest.TTRate).To(Equal(82.00)) // Dec 2023 month-end rate (preceding month)
			Expect(janInterest.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(janInterest.Tax).To(Equal(1.10))
			Expect(janInterest.Net).To(Equal(4.40))
			Expect(janInterest.INRValue()).To(Equal(451.00)) // 5.50 * 82.00
		})
	})

	Context("Valuation Calculation (INRValuation)", func() {
		It("should calculate valuations correctly for carry-over and fresh-start tickers", func() {
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary.INRValuations).ToNot(BeNil())
			Expect(summary.INRValuations).To(HaveLen(4)) // AAPL, ADI, GOOGL, MSFT

			// Sort by Ticker for consistent assertion order
			sort.Slice(summary.INRValuations, func(i, j int) bool {
				return summary.INRValuations[i].Ticker < summary.INRValuations[j].Ticker
			})

			aaplVal := summary.INRValuations[0]
			googlVal := summary.INRValuations[2]
			msftVal := summary.INRValuations[3]

			// Assert AAPL (Carry-over with new trades for 2023)
			Expect(aaplVal.Ticker).To(Equal("AAPL"))

			// FirstPosition for AAPL (opening balance for 2023 period, from Dec 31, 2022 accounts.csv)
			Expect(aaplVal.FirstPosition.Quantity).To(Equal(50.0))
			Expect(aaplVal.FirstPosition.USDPrice).To(Equal(160.00))
			Expect(aaplVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2022-12-31"))
			Expect(aaplVal.FirstPosition.TTRate).To(Equal(81.50))
			Expect(aaplVal.FirstPosition.TTDate.Format(time.DateOnly)).To(Equal("2022-12-30"))

			// Peak Position for AAPL (achieved on Jul 10, 2023, after Buy1 and Buy2)
			// Opening: 50. Buy1 (Mar 15): +20 (Total 70). Buy2 (Jul 10): +30 (Total 100 - This is Peak Qty)
			Expect(aaplVal.PeakPosition.Quantity).To(Equal(100.0))
			Expect(aaplVal.PeakPosition.USDPrice).To(Equal(165.00)) // Price of the Buy2 trade on Jul 10
			Expect(aaplVal.PeakPosition.Date.Format(time.DateOnly)).To(Equal("2023-07-10"))
			Expect(aaplVal.PeakPosition.TTRate).To(Equal(82.50)) // Assumed rate for 2023-07-10
			Expect(aaplVal.PeakPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-07-10"))

			// Year End Position for AAPL (after Sell1 on Oct 20)
			// Peak Qty: 100. Sell1: -15. Year-End Qty: 85
			Expect(aaplVal.YearEndPosition.Quantity).To(Equal(85.0))
			Expect(aaplVal.YearEndPosition.USDPrice).To(Equal(181.00)) // From AAPL.json for 2023-12-31
			Expect(aaplVal.YearEndPosition.Date.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(aaplVal.YearEndPosition.TTRate).To(Equal(82.00)) // From sbi_rates.csv for 2023-12-31
			Expect(aaplVal.YearEndPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))

			// AmountPaid: AAPL has 2 dividends in US Year 2023 (Jun 15 + Dec 30)
			// Jun 15: $50 × 82.00 (May month-end) = ₹4,100
			// Dec 30: $60 × 83.20 (Nov month-end) = ₹4,992
			// Total: ₹9,092
			Expect(aaplVal.AmountPaid).To(Equal(9092.0))

			// Assert GOOGL (Carry-over without trades)
			Expect(googlVal.Ticker).To(Equal("GOOGL"))
			Expect(googlVal.FirstPosition.Quantity).To(Equal(25.0))
			Expect(googlVal.FirstPosition.USDPrice).To(Equal(200.00))
			Expect(googlVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2022-12-31"))
			Expect(googlVal.PeakPosition.Quantity).To(Equal(25.0))
			Expect(googlVal.PeakPosition.Date.Format(time.DateOnly)).To(Equal("2022-12-31"))
			Expect(googlVal.YearEndPosition.Quantity).To(Equal(25.0))
			Expect(googlVal.YearEndPosition.USDPrice).To(Equal(140.00))
			Expect(googlVal.YearEndPosition.Date.Format(time.DateOnly)).To(Equal("2023-12-31"))

			// AmountPaid: GOOGL has no dividends in US Year 2023
			Expect(googlVal.AmountPaid).To(Equal(0.0))

			// Assert MSFT (Fresh Start)
			Expect(msftVal.Ticker).To(Equal("MSFT"))
			// First Position (MSFT)
			Expect(msftVal.FirstPosition.Quantity).To(Equal(50.0))
			Expect(msftVal.FirstPosition.USDPrice).To(Equal(200.00))
			Expect(msftVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2022-12-31"))
			Expect(msftVal.FirstPosition.TTRate).To(Equal(81.50))
			Expect(msftVal.FirstPosition.TTDate.Format(time.DateOnly)).To(Equal("2022-12-30"))
			// Peak Position (MSFT)
			Expect(msftVal.PeakPosition.Quantity).To(Equal(100.0))
			Expect(msftVal.PeakPosition.USDPrice).To(Equal(215.00))
			Expect(msftVal.PeakPosition.Date.Format(time.DateOnly)).To(Equal("2023-09-01"))
			Expect(msftVal.PeakPosition.TTRate).To(Equal(82.55))
			Expect(msftVal.PeakPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-08-31"))
			// Year End Position (MSFT)
			Expect(msftVal.YearEndPosition.Quantity).To(Equal(100.0))
			Expect(msftVal.YearEndPosition.USDPrice).To(Equal(221.00))
			Expect(msftVal.YearEndPosition.Date.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(msftVal.YearEndPosition.TTRate).To(Equal(82.00))
			Expect(msftVal.YearEndPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))

			// AmountPaid: MSFT has no dividends in US Year 2023
			Expect(msftVal.AmountPaid).To(Equal(0.0))
		})
	})

	Context("Fail-Fast Ticker Download - TDD", func() {
		It("should fail fast when ticker data missing for positive positions", func() {
			// TDD Test: Validates proper fail-fast behavior for tax systems
			//
			// SCENARIO: 2022 has BUY trades (IEF=42 shares) but IEF.json missing
			// EXPECTED: System should fail with standardized error format
			// This ensures tax accuracy over convenience - no silent failures

			_, err := taxManager.GetTaxSummary(ctx, 2022)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get year end price for IEF"))
			Expect(err.Error()).To(ContainSubstring("failed to auto-download ticker"))
		})
	})

	Context("Excel File Generation", func() {
		It("should generate a valid Excel file with the correct sheets", func() {
			// NOTE: Output files (Layer 3) are cleaned up after each test for isolation.
			// This ensures tests are independent and don't accumulate state.
			// Get the tax summary
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())

			// Save the summary to Excel
			saveErr := taxManager.SaveTaxSummaryToExcel(ctx, testYear, summary)
			Expect(saveErr).ToNot(HaveOccurred())

			// Verify that the file was created in Output/Reports/
			filePath := filepath.Join(kohanConfig.Tax.ReportsDir, fmt.Sprintf("tax_summary_%d.xlsx", testYear))
			defer os.Remove(filePath) // Clean up test artifact
			Expect(filePath).Should(BeARegularFile())

			// Open the generated file to verify its integrity and sheets
			f, openErr := excelize.OpenFile(filePath)
			Expect(openErr).ToNot(HaveOccurred(), "Generated Excel file should be valid and readable")
			defer f.Close()

			// Check for the presence of all required sheets
			expectedSheets := []string{"Gains", "Dividends", "Valuations", "Interest"}
			for _, sheetName := range expectedSheets {
				_, sheetErr := f.GetRows(sheetName)
				Expect(sheetErr).ToNot(HaveOccurred(), "Sheet '%s' should exist and be readable", sheetName)
			}
		})
	})

	Context("Account CSV Generation", func() {
		It("should generate a CSV file with year-end account data", func() {
			// Get the tax summary, which triggers the CSV generation
			_, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())

			// Define the expected path for the generated CSV
			expectedCsvPath := filepath.Join(kohanConfig.Tax.AccountsDir, "accounts_2023.csv")
			defer os.Remove(expectedCsvPath)

			// Verify that the file was created
			Expect(expectedCsvPath).Should(BeARegularFile())
		})
	})

	Context("Zero Quantity Positions (Fully Liquidated) - ADI", func() {
		var (
			adiKohanConfig config.KohanConfig
		)

		BeforeEach(func() {
			testDataBasePath := filepath.Join("..", "testdata", "tax")

			adiKohanConfig = config.KohanConfig{
				Tax: config.TaxConfig{
					TaxDir: testDataBasePath,
					// Layer 1: Input - Raw broker statements
					DriveWealthBase: filepath.Join(testDataBasePath, "Input", "Brokerage", "vested"),
					IBKRBase:        filepath.Join(testDataBasePath, "Input", "Brokerage", "ibkr"),
					// Layer 2.5: Parsed - Generated from broker statements
					ParsedDir:        filepath.Join(testDataBasePath, "Input", "Parsed"),
					TradesPath:       filepath.Join(testDataBasePath, "Input", "Parsed", tax.TRADES_FILENAME),
					DividendFilePath: filepath.Join(testDataBasePath, "Input", "Parsed", tax.DIVIDENDS_FILENAME),
					InterestFilePath: filepath.Join(testDataBasePath, "Input", "Parsed", tax.INTEREST_FILENAME),
					// Layer 3: Reference data (tickers, exchange rates)
					TickerCacheDir: filepath.Join(testDataBasePath, "Data", "Tickers"),
					TTRateFilePath: filepath.Join(testDataBasePath, "Data", "Reference", tax.SBI_RATES_FILENAME),
					// Layer 4: Output - Computed and generated results
					GainsFilePath: filepath.Join(testDataBasePath, "Output", "Computed", tax.GAINS_FILENAME),
					AccountsDir:   filepath.Join(testDataBasePath, "Output", "YearEndBalance"),
					ReportsDir:    filepath.Join(testDataBasePath, "Output", "Reports"),
					ComputedDir:   filepath.Join(testDataBasePath, "Output", "Computed"),
				},
			}

			core.SetupKohanInjector(adiKohanConfig)
		})

		It("should handle fully liquidated positions with zero year-end quantity", func() {
			// Integration test for year 2024 where ADI is fully liquidated
			// ADI: Bought 2 shares (Jan 4) → Sold 2 shares (Jan 20) → Zero year-end
			// This tests the critical bug where exchange_manager tries to fetch
			// exchange rates for positions with USD amount = 0 (unnecessary computation)

			// Re-get the tax manager with the ADI-specific config
			adiTaxManager, err := core.GetKohanInterface().GetTaxManager()
			Expect(err).ToNot(HaveOccurred())
			Expect(adiTaxManager).ToNot(BeNil())

			summary, err := adiTaxManager.GetTaxSummary(ctx, 2024)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).ToNot(BeNil())

			// ============================================================
			// PART 1: Verify Capital Gains Processing (ADI trade worked)
			// ============================================================

			// ADI should have a capital gain entry (STCG)
			var adiGain *tax.INRGains
			for i := range summary.INRGains {
				if summary.INRGains[i].Symbol == "ADI" {
					adiGain = &summary.INRGains[i]
					break
				}
			}

			Expect(adiGain).ToNot(BeNil(), "ADI should have a capital gain entry")
			Expect(adiGain.Symbol).To(Equal("ADI"))

			// ADI P&L calculation:
			// Sell: 2 × $194.75 = $389.50
			// Buy:  2 × $181.90 = $363.80
			// Commission: $0.00 (both trades have 0 commission)
			// P&L: $389.50 - $363.80 = $25.70
			Expect(adiGain.PNL).To(Equal(25.70))
			Expect(adiGain.Type).To(Equal("STCG")) // Holding period < 730 days
			Expect(adiGain.BuyDate).To(Equal("2024-01-04"))
			Expect(adiGain.SellDate).To(Equal("2024-01-20"))

			// Exchange rate lookup for gains (sell date month-end precedent)
			// Jan 2024 sell → uses Jan 20 rate: 82.68 (TT Buy)
			Expect(adiGain.TTRate).To(Equal(82.68))
			Expect(adiGain.INRValue()).To(Equal(2125.08)) // 25.70 × 82.68 ≈ 2125.08

			// ============================================================
			// PART 2: Verify Valuation Processing (Zero-Quantity Handling)
			// ============================================================

			// Find ADI valuation entry
			var adiVal *tax.INRValuation
			for i := range summary.INRValuations {
				if summary.INRValuations[i].Ticker == "ADI" {
					adiVal = &summary.INRValuations[i]
					break
				}
			}

			// ADI should have a valuation entry (for audit trail and completeness)
			Expect(adiVal).ToNot(BeNil(), "ADI should have valuation entry despite zero year-end quantity")

			// FirstPosition: ADI has no carry-forward, so first position is from first trade
			Expect(adiVal.FirstPosition.Quantity).To(Equal(2.0))
			Expect(adiVal.FirstPosition.USDPrice).To(Equal(181.90))
			Expect(adiVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2024-01-04"))
			Expect(adiVal.FirstPosition.TTRate).To(Equal(82.84)) // Jan 4 rate

			// PeakPosition: Same as FirstPosition (no additional buys after first purchase)
			Expect(adiVal.PeakPosition.Quantity).To(Equal(2.0))
			Expect(adiVal.PeakPosition.USDPrice).To(Equal(181.90))
			Expect(adiVal.PeakPosition.Date.Format(time.DateOnly)).To(Equal("2024-01-04"))

			// YearEndPosition: ZERO quantity (fully liquidated)
			Expect(adiVal.YearEndPosition.Quantity).To(Equal(0.0))
			Expect(adiVal.YearEndPosition.USDPrice).To(Equal(230.06)) // Year-end price from ADI.json
			Expect(adiVal.YearEndPosition.Date.Format(time.DateOnly)).To(Equal("2024-12-31"))

			// ============================================================
			// CRITICAL ASSERTION: This is what the bug test is about!
			// ============================================================

			// USD Amount should be ZERO (Quantity × Price = 0 × $230.06 = $0)
			Expect(adiVal.YearEndPosition.GetUSDAmount()).To(Equal(0.0))

			// TTRate SHOULD BE ZERO (no exchange rate lookup for zero-value position)
			// BEFORE FIX: This will FAIL because system fetches 2024-12-31 rate (85.20)
			// AFTER FIX: This will PASS because system skips exchange (TTRate remains 0)
			Expect(adiVal.YearEndPosition.TTRate).To(Equal(0.0),
				"TTRate should be 0 for zero-quantity position (no exchange rate lookup needed)")

			// INR Value must be zero (0 × price × rate = 0 regardless of rate)
			expectedINRValue := adiVal.YearEndPosition.Quantity *
				adiVal.YearEndPosition.USDPrice *
				adiVal.YearEndPosition.TTRate
			Expect(expectedINRValue).To(Equal(0.0))

			// AmountPaid: ADI has no dividends in 2024
			Expect(adiVal.AmountPaid).To(Equal(0.0))
		})
	})
})
