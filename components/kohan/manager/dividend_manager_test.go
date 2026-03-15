//nolint:dupl // False positives: Similar test patterns for dividend/interest processing
package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	mockrepo "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("DividendManager", func() {
	var (
		ctx                      = context.Background()
		mockExchangeManager      *mocks.ExchangeManager
		mockFinancialYearManager *mocks.FinancialYearManager[tax.Dividend]
		mockDividendRepository   *mockrepo.DividendRepository
		dividendManager          manager.DividendManager

		// Common test data
		ticker   = "AAPL"
		date     = "2024-01-15"
		amount   = 100.00
		taxValue = 10.00
		net      = 90.00
	)

	BeforeEach(func() {
		mockExchangeManager = mocks.NewExchangeManager(GinkgoT())
		mockFinancialYearManager = mocks.NewFinancialYearManager[tax.Dividend](GinkgoT())
		mockDividendRepository = mockrepo.NewDividendRepository(GinkgoT())
		dividendManager = manager.NewDividendManager(mockExchangeManager, mockFinancialYearManager, mockDividendRepository)
	})

	Context("Basic Dividend Processing", func() {
		var (
			dividends []tax.Dividend
		)

		BeforeEach(func() {
			dividends = []tax.Dividend{
				{
					Symbol: ticker,
					Date:   date,
					Amount: amount,
					Tax:    taxValue,
					Net:    net,
				},
			}

			mockExchangeManager.EXPECT().
				ExchangeWithPrecedingMonth(ctx, mock.Anything).
				Return(nil)
		})

		It("should process dividend correctly", func() {
			inrDividends, err := dividendManager.ProcessDividends(ctx, dividends)

			Expect(err).ToNot(HaveOccurred())
			Expect(inrDividends).To(HaveLen(1))

			result := inrDividends[0]
			Expect(result.Symbol).To(Equal(ticker))
			Expect(result.Date).To(Equal(date))
			Expect(result.Amount).To(Equal(amount))
			Expect(result.Tax).To(Equal(taxValue))
			Expect(result.Net).To(Equal(net))
		})
	})

	Context("Exchange Rate Error", func() {
		var dividends []tax.Dividend

		BeforeEach(func() {
			dividends = []tax.Dividend{{
				Symbol: ticker,
				Date:   date,
				Amount: amount,
				Tax:    taxValue,
				Net:    net,
			}}

			mockExchangeManager.EXPECT().
				ExchangeWithPrecedingMonth(ctx, mock.Anything).
				Return(common.ErrNotFound)
		})

		It("should handle exchange rate error", func() {
			_, err := dividendManager.ProcessDividends(ctx, dividends)
			Expect(err).To(Equal(common.ErrNotFound))
		})
	})

	Context("Multiple Dividends", func() {
		var (
			dividends []tax.Dividend
		)

		BeforeEach(func() {
			dividends = []tax.Dividend{
				{
					Symbol: "AAPL",
					Date:   "2024-01-15",
					Amount: 100.00,
					Tax:    10.00,
					Net:    90.00,
				},
				{
					Symbol: "MSFT",
					Date:   "2024-01-16",
					Amount: 200.00,
					Tax:    20.00,
					Net:    180.00,
				},
			}

			mockExchangeManager.EXPECT().
				ExchangeWithPrecedingMonth(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
					Expect(exchangeables).To(HaveLen(2))
					Expect(exchangeables[0].GetUSDAmount()).To(Equal(100.00))
					Expect(exchangeables[1].GetUSDAmount()).To(Equal(200.00))
				}).
				Return(nil)
		})

		It("should process multiple dividends correctly", func() {
			inrDividends, err := dividendManager.ProcessDividends(ctx, dividends)

			Expect(err).ToNot(HaveOccurred())
			Expect(inrDividends).To(HaveLen(2))

			// Verify first dividend
			Expect(inrDividends[0].Symbol).To(Equal("AAPL"))
			Expect(inrDividends[0].Amount).To(Equal(100.00))

			// Verify second dividend
			Expect(inrDividends[1].Symbol).To(Equal("MSFT"))
			Expect(inrDividends[1].Amount).To(Equal(200.00))
		})
	})

	Context("GetDividendsForUSYear", func() {
		var (
			allDividends      []tax.Dividend
			filteredDividends []tax.Dividend
			year              int
		)

		BeforeEach(func() {
			year = 2022
			allDividends = []tax.Dividend{
				{Symbol: "IEF", Date: "2021-12-15", Amount: 50.00},  // Before 2022
				{Symbol: "IEF", Date: "2022-01-15", Amount: 100.00}, // In 2022
				{Symbol: "IVV", Date: "2022-06-20", Amount: 150.00}, // In 2022
				{Symbol: "TLT", Date: "2022-12-31", Amount: 200.00}, // In 2022
				{Symbol: "VGK", Date: "2023-01-05", Amount: 75.00},  // After 2022
			}

			filteredDividends = []tax.Dividend{
				{Symbol: "IEF", Date: "2022-01-15", Amount: 100.00},
				{Symbol: "IVV", Date: "2022-06-20", Amount: 150.00},
				{Symbol: "TLT", Date: "2022-12-31", Amount: 200.00},
			}

			mockDividendRepository.EXPECT().
				GetAllRecords(ctx).
				Return(allDividends, nil)

			mockFinancialYearManager.EXPECT().
				FilterUS(ctx, allDividends, year).
				Return(filteredDividends, nil)
		})

		It("should filter dividends by US calendar year", func() {
			result, err := dividendManager.GetDividendsForUSYear(ctx, year)

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(3))
			Expect(result[0].Symbol).To(Equal("IEF"))
			Expect(result[0].Date).To(Equal("2022-01-15"))
			Expect(result[1].Symbol).To(Equal("IVV"))
			Expect(result[2].Symbol).To(Equal("TLT"))
		})

		Context("when repository fails", func() {
			It("should return error", func() {
				mockDividendRepository.ExpectedCalls = nil
				mockFinancialYearManager.ExpectedCalls = nil

				mockDividendRepository.EXPECT().
					GetAllRecords(ctx).
					Return(nil, common.ErrNotFound)

				_, err := dividendManager.GetDividendsForUSYear(ctx, year)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})
	})
})
