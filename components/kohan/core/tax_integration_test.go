package core_test

import (
	"context"
	"fmt"
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

	Context("Tax Summary Calculation", func() {
		It("should calculate tax summary correctly for the given year", func() {
			summary, err := taxManager.GetTaxSummary(ctx, testYear)
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).ToNot(BeNil())

			// Debug output
			fmt.Printf("Tax Summary: %+v\n", summary)
			fmt.Printf("INRGains: %+v\n", summary.INRGains)

			Expect(summary.INRGains).To(HaveLen(1))

			gain := summary.INRGains[0]
			Expect(gain.Symbol).To(Equal("AAPL"))
			Expect(gain.PNL).To(BeNumerically("~", 1000.00))
			Expect(gain.Type).To(Equal("STCG"))
			Expect(gain.BuyDate).To(Equal("2024-01-15"))
			Expect(gain.SellDate).To(Equal("2024-01-17"))

			Expect(gain.TTRate).To(BeNumerically("~", 82.90))
			Expect(gain.INRValue()).To(BeNumerically("~", 1000.00*82.90))
			Expect(gain.TTDate.Format(time.DateOnly)).To(Equal("2024-01-17"))
		})
	})
})
