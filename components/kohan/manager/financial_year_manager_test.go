package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FinancialYearManager", func() {
	var (
		ctx       context.Context
		fyManager manager.FinancialYearManager[tax.Interest]
	)

	BeforeEach(func() {
		ctx = context.Background()
		fyManager = manager.NewFinancialYearManager[tax.Interest]()
	})

	Context("FilterIndia", func() {
		var (
			year        = 2024 // Testing for FY 2024-25
			testRecords []tax.Interest
		)

		BeforeEach(func() {
			testRecords = []tax.Interest{
				{Symbol: "TEST", Date: "2024-04-01", Amount: 100}, // Start of FY
				{Symbol: "TEST", Date: "2024-08-15", Amount: 200}, // Mid FY
				{Symbol: "TEST", Date: "2025-03-31", Amount: 300}, // End of FY
				{Symbol: "TEST", Date: "2024-03-31", Amount: 400}, // Previous FY
				{Symbol: "TEST", Date: "2025-04-01", Amount: 500}, // Next FY
			}
		})

		It("should filter records for correct financial year", func() {
			filtered, err := fyManager.FilterIndia(ctx, testRecords, year)
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(HaveLen(3))
			Expect(filtered[0].Amount).To(Equal(100.0))
			Expect(filtered[1].Amount).To(Equal(200.0))
			Expect(filtered[2].Amount).To(Equal(300.0))
		})

		It("should handle empty record list", func() {
			filtered, err := fyManager.FilterIndia(ctx, []tax.Interest{}, year)
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(BeEmpty())
		})
	})

	Context("FilterUS", func() {
		var (
			year        = 2024 // Testing for FY 2024 (Jan-Dec)
			testRecords []tax.Interest
		)

		BeforeEach(func() {
			testRecords = []tax.Interest{
				{Symbol: "AAPL", Date: "2023-12-31", Amount: 400}, // Previous FY
				{Symbol: "AAPL", Date: "2024-01-01", Amount: 100}, // Start of FY
				{Symbol: "AAPL", Date: "2024-06-15", Amount: 200}, // Mid FY
				{Symbol: "AAPL", Date: "2024-12-31", Amount: 300}, // End of FY
				{Symbol: "AAPL", Date: "2025-01-01", Amount: 500}, // Next FY
			}
		})

		It("should filter records for correct financial year", func() {
			filtered, err := fyManager.FilterUS(ctx, testRecords, year)
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(HaveLen(3))
			Expect(filtered[0].Amount).To(Equal(100.0))
			Expect(filtered[1].Amount).To(Equal(200.0))
			Expect(filtered[2].Amount).To(Equal(300.0))
		})

		It("should handle empty record list", func() {
			filtered, err := fyManager.FilterUS(ctx, []tax.Interest{}, year)
			Expect(err).ToNot(HaveOccurred())
			Expect(filtered).To(BeEmpty())
		})

		It("should handle invalid date format", func() {
			testRecords = []tax.Interest{
				{Symbol: "AAPL", Date: "invalid-date", Amount: 100}, // Invalid date
			}
			filtered, err := fyManager.FilterUS(ctx, testRecords, year)
			Expect(err).To(HaveOccurred())
			Expect(filtered).To(BeEmpty())
		})
	})
})
