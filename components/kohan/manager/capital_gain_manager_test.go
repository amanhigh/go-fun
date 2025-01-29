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

var _ = Describe("CapitalGainManager", func() {
	var (
		ctx                 = context.Background()
		mockExchangeManager *mocks.ExchangeManager
		gainManager         manager.CapitalGainManager

		// Common test data
		ticker   = "AAPL"
		sellDate = "2024-01-15"
	)

	BeforeEach(func() {
		mockExchangeManager = mocks.NewExchangeManager(GinkgoT())
		gainManager = manager.NewCapitalGainManager(mockExchangeManager)
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
			Expect(err).To(BeNil())
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

			Expect(err).To(BeNil())
			Expect(taxGains).To(HaveLen(2))

			// Verify first gain
			Expect(taxGains[0].Gains.Symbol).To(Equal("AAPL"))
			Expect(taxGains[0].Gains.PNL).To(Equal(1000.00))

			// Verify second gain
			Expect(taxGains[1].Gains.Symbol).To(Equal("MSFT"))
			Expect(taxGains[1].Gains.PNL).To(Equal(2000.00))
		})
	})
})
