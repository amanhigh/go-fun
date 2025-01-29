package manager_test

import (
	"context"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("TaxValuationManager", func() {
	var (
		ctx              context.Context
		mockExchange     *mocks.ExchangeManager
		valuationManager manager.TaxValuationManager
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockExchange = mocks.NewExchangeManager(GinkgoT())
		valuationManager = manager.NewTaxValuationManager(mockExchange)
	})

	Context("Position Processing", func() {
		var (
			valuation tax.Valuation
		)

		BeforeEach(func() {
			valuation = tax.Valuation{
				Ticker: "AAPL",
			}
		})

		It("should process positions", func() {
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

			mockExchange.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Return(nil)

			result, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(1))

			taxVal := result[0]
			Expect(taxVal.Ticker).To(Equal(valuation.Ticker))
			Expect(taxVal.FirstPosition.Position).To(Equal(valuation.FirstPosition))
			Expect(taxVal.PeakPosition.Position).To(Equal(valuation.PeakPosition))
			Expect(taxVal.YearEndPosition.Position).To(Equal(valuation.YearEndPosition))
		})

		It("should return error if exchange fails", func() {
			errTest := common.NewHttpError("test error", http.StatusInternalServerError)
			mockExchange.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Return(errTest)

			_, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
			Expect(err).To(Equal(errTest))
		})
	})

	Context("Batch Processing", func() {
		var (
			valuations []tax.Valuation
			testDate   time.Time
		)

		BeforeEach(func() {
			testDate = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			// Create multiple valuations
			valuations = []tax.Valuation{
				{
					Ticker: "AAPL",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 100,
						USDPrice: 150,
					},
				},
				{
					Ticker: "MSFT",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 50,
						USDPrice: 200,
					},
				},
			}

			mockExchange.EXPECT().
				Exchange(ctx, mock.Anything).
				Times(2).
				Return(nil)
		})

		It("should process multiple valuations", func() {
			result, err := valuationManager.ProcessValuations(ctx, valuations)

			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(2))
			Expect(result[0].Ticker).To(Equal("AAPL"))
			Expect(result[1].Ticker).To(Equal("MSFT"))
		})
	})
})
