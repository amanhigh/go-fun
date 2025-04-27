package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	mockRepo "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("CapitalGainManager", func() {
	var (
		ctx                 = context.Background()
		mockExchangeManager *mocks.ExchangeManager
		mockFYManager       *mocks.FinancialYearManager[tax.Gains]
		mockGainsRepo       *mockRepo.GainsRepository
		gainManager         manager.CapitalGainManager

		// Common test data
		ticker   = "AAPL"
		sellDate = "2024-01-15"
	)
	// FIXME: Add unit tests for edge date conditions (e.g., buy/sell on April 1st, March 31st) for financial year boundaries.
	// FIXME: Enhance unit tests to verify correct processing and aggregation details when multiple valid gains exist within the same financial year.

	BeforeEach(func() {
		mockExchangeManager = mocks.NewExchangeManager(GinkgoT())
		mockFYManager = mocks.NewFinancialYearManager[tax.Gains](GinkgoT())
		mockGainsRepo = mockRepo.NewGainsRepository(GinkgoT())
		gainManager = manager.NewCapitalGainManager(mockExchangeManager, mockGainsRepo, mockFYManager)
	})

	Context("Basic Gains Processing", func() {
		var (
			gains []tax.Gains
			pnl   = 1000.00
		)

		BeforeEach(func() {
			gains = []tax.Gains{
				{
					Symbol:   ticker,
					SellDate: sellDate,
					PNL:      pnl,
				},
			}

			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.Anything).
				Return(nil)
		})

		It("should process gain and keep original Gain Values", func() {
			taxGains, err := gainManager.ProcessTaxGains(ctx, gains)
			Expect(err).ToNot(HaveOccurred())
			Expect(taxGains).To(HaveLen(1))

			result := taxGains[0]
			Expect(result.Gains).To(Equal(gains[0]))
		})
	})

	Context("Exchange Rate Error", func() {
		var gains []tax.Gains

		BeforeEach(func() {
			gains = []tax.Gains{{
				Symbol:   ticker,
				SellDate: sellDate,
				PNL:      1000.00,
			}}

			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.Anything).
				Return(common.ErrNotFound)
		})

		It("should handle missing exchange rate", func() {
			_, err := gainManager.ProcessTaxGains(ctx, gains)
			Expect(err).To(Equal(common.ErrNotFound))
		})
	})

	Context("Multiple Gains", func() {
		var (
			gains []tax.Gains
		)

		BeforeEach(func() {
			gains = []tax.Gains{
				{
					Symbol:   "AAPL",
					SellDate: "2024-01-15",
					PNL:      1000.00,
				},
				{
					Symbol:   "MSFT",
					SellDate: "2024-01-16",
					PNL:      2000.00,
				},
			}

			// Verify that exchangeables passed contain correct gain amounts
			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
					Expect(exchangeables).To(HaveLen(2))
					Expect(exchangeables[0].GetUSDAmount()).To(Equal(1000.00))
					Expect(exchangeables[1].GetUSDAmount()).To(Equal(2000.00))
				}).
				Return(nil)
		})

		It("should process multiple gains correctly", func() {
			taxGains, err := gainManager.ProcessTaxGains(ctx, gains)

			Expect(err).ToNot(HaveOccurred())
			Expect(taxGains).To(HaveLen(2))

			// Verify first gain
			Expect(taxGains[0].Gains.Symbol).To(Equal("AAPL"))
			Expect(taxGains[0].Gains.PNL).To(Equal(1000.00))

			// Verify second gain
			Expect(taxGains[1].Gains.Symbol).To(Equal("MSFT"))
			Expect(taxGains[1].Gains.PNL).To(Equal(2000.00))
		})
	})

	// BUG: Put Upper Methods in Function Context
	Context("GetGainsForYear", func() {
		var (
			testYear = 2024
			// BUG: Create Constructor
			allGains = []tax.Gains{
				{
					Symbol:   "AAPL",
					SellDate: "2024-04-15",
					PNL:      1000.00,
				},
				{
					Symbol:   "GOOGL",
					SellDate: "2024-05-20",
					PNL:      2000.00,
				},
			}
			filteredGains = []tax.Gains{allGains[0]} // Assume only first gain matches FY
		)

		Context("when successful", func() {
			BeforeEach(func() {
				// Setup repository mock to return test gains
				mockGainsRepo.EXPECT().
					GetAllRecords(ctx).
					Return(allGains, nil)

				// Setup FY manager to filter gains
				mockFYManager.EXPECT().
					FilterRecordsByFY(ctx, allGains, testYear).
					Return(filteredGains, nil)
			})

			It("should return filtered gains for the year", func() {
				gains, err := gainManager.GetGainsForYear(ctx, testYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(gains).To(Equal(filteredGains))
			})
		})

		Context("when repository fails", func() {
			BeforeEach(func() {
				mockGainsRepo.EXPECT().
					GetAllRecords(ctx).
					Return(nil, common.ErrInternalServerError)
			})

			It("should return repository error", func() {
				_, err := gainManager.GetGainsForYear(ctx, testYear)

				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})

		Context("when filtering fails", func() {
			BeforeEach(func() {
				mockGainsRepo.EXPECT().
					GetAllRecords(ctx).
					Return(allGains, nil)

				mockFYManager.EXPECT().
					FilterRecordsByFY(ctx, allGains, testYear).
					Return(nil, common.ErrInternalServerError)
			})

			It("should return filtering error", func() {
				_, err := gainManager.GetGainsForYear(ctx, testYear)

				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})
	})
})
