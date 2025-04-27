package core_test

import (
	"context"
	"path/filepath"
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
			Expect(aaplGain.PNL).To(BeNumerically("~", 1000.00))
			Expect(aaplGain.Type).To(Equal("STCG")) // Holding < 730 days
			Expect(aaplGain.BuyDate).To(Equal("2024-01-15"))
			Expect(aaplGain.SellDate).To(Equal("2024-01-17"))
			Expect(aaplGain.TTRate).To(BeNumerically("~", 82.90)) // Rate for Jan 17
			Expect(aaplGain.INRValue()).To(BeNumerically("~", 1000.00*82.90))
			Expect(aaplGain.TTDate.Format(time.DateOnly)).To(Equal("2024-01-17"))

			// --- Assertions for MSFT (Expected at index 1 after sort) ---
			msftGain := summary.INRGains[1]
			Expect(msftGain.Symbol).To(Equal("MSFT"))
			Expect(msftGain.PNL).To(BeNumerically("~", 500.00))
			Expect(msftGain.Type).To(Equal("LTCG")) // Holding > 730 days
			Expect(msftGain.BuyDate).To(Equal("2022-01-10"))
			Expect(msftGain.SellDate).To(Equal("2024-02-15"))
			Expect(msftGain.TTRate).To(BeNumerically("~", 83.00)) // Rate for Feb 15 (from test data)
			Expect(msftGain.INRValue()).To(BeNumerically("~", 500.00*83.00))
			Expect(msftGain.TTDate.Format(time.DateOnly)).To(Equal("2024-02-15"))
		})
	})

	// FUTURE CONTEXT: Placeholder for Dividends
	/*
		Context("Dividend Calculation (INRDividends)", func() {
		    // ... tests for dividends ...
		})
	*/

	// FUTURE CONTEXT: Placeholder for Interest
	/*
		Context("Interest Calculation (INRInterest)", func() {
		    // ... tests for interest ...
		})
	*/

	// FUTURE CONTEXT: Placeholder for Valuations
	/*
		Context("Valuation Calculation (INRValuation)", func() {
		    // ... tests for valuations ...
		})
	*/
})
