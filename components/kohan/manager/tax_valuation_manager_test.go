package manager_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("TaxValuationManager", func() {
	var (
		ctx              context.Context
		mockExchange     *mocks.ExchangeManager
		valuationManager manager.TaxValuationManager
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockExchange = mocks.NewExchangeManager(GinkgoT())
		valuationManager = manager.NewTaxValuationManager(mockExchange)

		Context("Position Processing", func() {
			var (
				valuation tax.Valuation
			)

			BeforeEach(func() {
				valuation = tax.Valuation{
					Ticker: "AAPL",
				}
			})

			It("should skip empty positions", func() {
				// Empty positions
				valuation.FirstPosition = tax.Position{Quantity: 0}
				valuation.PeakPosition = tax.Position{Quantity: 0}
				valuation.YearEndPosition = tax.Position{Quantity: 0}

				result, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
				Expect(err).To(BeNil())
				Expect(result).To(HaveLen(1))
			})

			It("should process only non-empty positions", func() {
				date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

				// Mixed positions
				valuation.FirstPosition = tax.Position{
					Date:     date,
					Quantity: 100,
					USDPrice: 150,
				}
				valuation.PeakPosition = tax.Position{Quantity: 0}
				valuation.YearEndPosition = tax.Position{
					Date:     date,
					Quantity: 200,
					USDPrice: 160,
				}

				// Exchange should be called only for non-empty positions
				mockExchange.EXPECT().
					Exchange(ctx, valuation).
					Return(nil)

				result, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{})
				Expect(err).To(BeNil())
				Expect(result).To(HaveLen(1))
			})
		})
	})
})
