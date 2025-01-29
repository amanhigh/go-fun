package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("DividendManager", func() {
	var (
		ctx                 = context.Background()
		mockExchangeManager *mocks.ExchangeManager
		dividendManager     manager.DividendManager

		// Common test data
		ticker   = "AAPL"
		date     = "2024-01-15"
		amount   = 100.00
		taxValue = 10.00
		net      = 90.00
	)

	BeforeEach(func() {
		mockExchangeManager = mocks.NewExchangeManager(GinkgoT())
		dividendManager = manager.NewDividendManager(mockExchangeManager)
	})

	Context("Basic Dividend Processing", func() {
		var (
			dividends []tax.Dividend
		)

		BeforeEach(func() {
			dividends = []tax.Dividend{
				tax.Dividend{
					Symbol: ticker,
					Date:   date,
					Amount: amount,
					Tax:    taxValue,
					Net:    net,
				},
			}

			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.Anything).
				Return(nil)
		})

		It("should process dividend correctly", func() {
			inrDividends, err := dividendManager.ProcessDividends(ctx, dividends)

			Expect(err).To(BeNil())
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
			dividends = []tax.Dividend{tax.Dividend{
				Symbol: ticker,
				Date:   date,
				Amount: amount,
				Tax:    taxValue,
				Net:    net,
			}}

			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.Anything).
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

			// Verify that exchangeables passed contain correct dividend amounts
			mockExchangeManager.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
					Expect(exchangeables).To(HaveLen(2))
					Expect(exchangeables[0].GetUSDAmount()).To(Equal(100.00))
					Expect(exchangeables[1].GetUSDAmount()).To(Equal(200.00))
				}).
				Return(nil)
		})

		It("should process multiple dividends correctly", func() {
			inrDividends, err := dividendManager.ProcessDividends(ctx, dividends)

			Expect(err).To(BeNil())
			Expect(inrDividends).To(HaveLen(2))

			// Verify first dividend
			Expect(inrDividends[0].Symbol).To(Equal("AAPL"))
			Expect(inrDividends[0].Amount).To(Equal(100.00))

			// Verify second dividend
			Expect(inrDividends[1].Symbol).To(Equal("MSFT"))
			Expect(inrDividends[1].Amount).To(Equal(200.00))
		})
	})
})
