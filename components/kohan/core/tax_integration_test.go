package core_test

import (
	"context"
	"path/filepath"
	"sort" // Add import
	"time"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tax Integration", Label("it"), func() {
	var (
		ctx        context.Context
		taxManager manager.TaxManager
		testYear   = 2023
	)

	BeforeEach(func() {
		ctx = context.Background()
		testDataBasePath := filepath.Join("..", "testdata", "tax")

		// Configure KohanConfig with TaxConfig pointing to test data files
		kohanConfig := config.KohanConfig{
			Tax: config.TaxConfig{
				// DownloadsDir is separate, points to base testdata path for this test
				DownloadsDir: testDataBasePath,
				// File Paths using constants and joined with base path
				BrokerStatementPath: filepath.Join(testDataBasePath, tax.TRADES_FILENAME),
				DividendFilePath:    filepath.Join(testDataBasePath, tax.DIVIDENDS_FILENAME),
				SBIFilePath:         filepath.Join(testDataBasePath, tax.SBI_RATES_FILENAME),
				AccountFilePath:     filepath.Join(testDataBasePath, tax.ACCOUNTS_FILENAME),
				GainsFilePath:       filepath.Join(testDataBasePath, tax.GAINS_FILENAME),
				InterestFilePath:    filepath.Join(testDataBasePath, tax.INTEREST_FILENAME),
			},
		}

		// Setup the global injector with test configuration
		core.SetupKohanInjector(kohanConfig)

		// Retrieve the TaxManager instance
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
			Expect(summary.INRGains).To(HaveLen(2)) // Expecting AAPL (STCG) and MSFT (LTCG)

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

			// --- Assertions for MSFT (Expected at index 1 after sort) ---
			msftGain := summary.INRGains[1]
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
			// Based on testdata: AAPL Jan 15, MSFT Feb 20 (added), AAPL Mar 15 fall in this FY.
			// AAPL Apr 15 should be filtered out.
			Expect(summary.INRDividends).To(HaveLen(3)) // Expecting 3 dividends after filtering

			// Sort results by date to ensure consistent order for assertions
			sort.Slice(summary.INRDividends, func(i, j int) bool {
				dateI, err := summary.INRDividends[i].GetDate()
				Expect(err).NotTo(HaveOccurred())
				dateJ, err := summary.INRDividends[j].GetDate()
				Expect(err).NotTo(HaveOccurred())
				return dateI.Before(dateJ)
			})

			// --- Assertions for Jan 15 Dividend (AAPL) - Full Detail ---
			janDividend := summary.INRDividends[0]
			Expect(janDividend.Symbol).To(Equal("AAPL"))
			Expect(janDividend.Date).To(Equal("2024-01-15"))
			Expect(janDividend.Amount).To(Equal(115.00))
			Expect(janDividend.Tax).To(Equal(17.25))
			Expect(janDividend.Net).To(Equal(97.75))
			Expect(janDividend.TTRate).To(Equal(82.50)) // Rate for Jan 15 from sbi_rates.csv
			Expect(janDividend.TTDate.Format(time.DateOnly)).To(Equal("2024-01-15"))
			Expect(janDividend.INRValue()).To(Equal(9487.50)) // 115.00 * 82.50

			// --- Assertions for Feb 20 Dividend (MSFT) - Key Details ---
			febDividend := summary.INRDividends[1]
			Expect(febDividend.Symbol).To(Equal("MSFT"))
			Expect(febDividend.Amount).To(Equal(50.00)) // Check Amount for MSFT
			Expect(febDividend.TTRate).To(Equal(83.05)) // Assumed rate for Feb 20
			Expect(febDividend.TTDate.Format(time.DateOnly)).To(Equal("2024-02-20"))
			Expect(febDividend.INRValue()).To(Equal(4152.50)) // 50.00 * 83.05

			// --- Assertions for Mar 15 Dividend (AAPL) - Key Details ---
			marDividend := summary.INRDividends[2]
			Expect(marDividend.Symbol).To(Equal("AAPL"))
			Expect(marDividend.Amount).To(Equal(100.00)) // Check Amount for AAPL Mar
			Expect(marDividend.TTRate).To(Equal(83.10))  // Rate for Mar 15 from sbi_rates.csv
			Expect(marDividend.TTDate.Format(time.DateOnly)).To(Equal("2024-03-15"))
			Expect(marDividend.INRValue()).To(Equal(8310.00)) // 100.00 * 83.10
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
			Expect(decInterest.TTRate).To(Equal(82.00)) // Assumed rate for Dec 31
			Expect(decInterest.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(decInterest.INRValue()).To(Equal(1640.00)) // 20.00 * 82.00

			// --- Assertions for Jan 10 Interest (AAPL) - Key Details ---
			janInterest := summary.INRInterest[1]
			Expect(janInterest.Symbol).To(Equal("AAPL"))
			Expect(janInterest.Amount).To(Equal(5.50))  // Check Amount for AAPL
			Expect(janInterest.TTRate).To(Equal(82.40)) // Assumed rate for Jan 10
			Expect(janInterest.TTDate.Format(time.DateOnly)).To(Equal("2024-01-10"))
			Expect(janInterest.Tax).To(Equal(1.10))
			Expect(janInterest.Net).To(Equal(4.40))
			Expect(janInterest.INRValue()).To(Equal(453.20)) // 5.50 * 82.40
		})
	})

	Context("Valuation Calculation (INRValuation)", func() {
		It("should calculate valuations correctly for carry-over and fresh-start tickers", func() {
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary.INRValuations).ToNot(BeNil())
			Expect(summary.INRValuations).To(HaveLen(2))

			// Sort by Ticker for consistent assertion order
			sort.Slice(summary.INRValuations, func(i, j int) bool {
				return summary.INRValuations[i].Ticker < summary.INRValuations[j].Ticker
			})

			aaplVal := summary.INRValuations[0]
			msftVal := summary.INRValuations[1]

			// Assert AAPL (Carry-over with new trades for 2023)
			Expect(aaplVal.Ticker).To(Equal("AAPL"))

			// FirstPosition for AAPL (opening balance for 2023 period, from Dec 31, 2022 accounts.csv)
			Expect(aaplVal.FirstPosition.Quantity).To(Equal(50.0))
			Expect(aaplVal.FirstPosition.USDPrice).To(Equal(160.00))
			Expect(aaplVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2023-01-01"))
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

			// Assert MSFT (Fresh Start)
			Expect(msftVal.Ticker).To(Equal("MSFT"))
			// First Position (MSFT)
			Expect(msftVal.FirstPosition.Quantity).To(Equal(20.0))
			Expect(msftVal.FirstPosition.USDPrice).To(Equal(205.00))
			Expect(msftVal.FirstPosition.Date.Format(time.DateOnly)).To(Equal("2023-05-01"))
			Expect(msftVal.FirstPosition.TTRate).To(Equal(82.00))
			Expect(msftVal.FirstPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-05-01"))
			// Peak Position (MSFT)
			Expect(msftVal.PeakPosition.Quantity).To(Equal(50.0))
			Expect(msftVal.PeakPosition.USDPrice).To(Equal(215.00))
			Expect(msftVal.PeakPosition.Date.Format(time.DateOnly)).To(Equal("2023-09-01"))
			Expect(msftVal.PeakPosition.TTRate).To(Equal(82.55))
			Expect(msftVal.PeakPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-08-31"))
			// Year End Position (MSFT)
			Expect(msftVal.YearEndPosition.Quantity).To(Equal(50.0))
			Expect(msftVal.YearEndPosition.USDPrice).To(Equal(221.00))
			Expect(msftVal.YearEndPosition.Date.Format(time.DateOnly)).To(Equal("2023-12-31"))
			Expect(msftVal.YearEndPosition.TTRate).To(Equal(82.00))
			Expect(msftVal.YearEndPosition.TTDate.Format(time.DateOnly)).To(Equal("2023-12-31"))
		})
	})
})
