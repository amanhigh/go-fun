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
		ctx                  context.Context
		mockExchange         *mocks.ExchangeManager
		mockValuationManager *mocks.ValuationManager
		valuationManager     manager.TaxValuationManager
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockExchange = mocks.NewExchangeManager(GinkgoT())
		mockValuationManager = mocks.NewValuationManager(GinkgoT())
		// Pass the new mock to the constructor
		valuationManager = manager.NewTaxValuationManager(mockExchange, mockValuationManager)
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

			result, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation}, []tax.INRDividend{})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(1))

			taxVal := result[0]
			Expect(taxVal.Ticker).To(Equal(valuation.Ticker))
			Expect(taxVal.FirstPosition.Position).To(Equal(valuation.FirstPosition))
			Expect(taxVal.PeakPosition.Position).To(Equal(valuation.PeakPosition))
			Expect(taxVal.YearEndPosition.Position).To(Equal(valuation.YearEndPosition))
			Expect(taxVal.AmountPaid).To(Equal(0.0))
		})

		It("should return error if exchange fails", func() {
			errTest := common.NewHttpError("test error", http.StatusInternalServerError)
			mockExchange.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Return(errTest)

			_, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation}, []tax.INRDividend{})
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
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
					Expect(exchangeables).To(HaveLen(2 * 3))
					Expect(exchangeables[0].GetUSDAmount()).To(Equal(15000.00))
					Expect(exchangeables[3].GetUSDAmount()).To(Equal(10000.00))
				}).
				Return(nil)
		})

		It("should process multiple valuations", func() {
			result, err := valuationManager.ProcessValuations(ctx, valuations, []tax.INRDividend{})

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(2))
			Expect(result[0].Ticker).To(Equal("AAPL"))
			Expect(result[0].AmountPaid).To(Equal(0.0))
			Expect(result[1].Ticker).To(Equal("MSFT"))
			Expect(result[1].AmountPaid).To(Equal(0.0))
		})
	})

	Context("AmountPaid Calculation", func() {
		var (
			valuations []tax.Valuation
			dividends  []tax.INRDividend
			testDate   time.Time
		)

		BeforeEach(func() {
			testDate = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			valuations = []tax.Valuation{
				{
					Ticker: "IEF",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 100,
						USDPrice: 100,
					},
				},
				{
					Ticker: "IVV",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 50,
						USDPrice: 200,
					},
				},
				{
					Ticker: "TLT",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 25,
						USDPrice: 150,
					},
				},
			}

			mockExchange.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Return(nil)
		})

		It("should calculate AmountPaid when dividends are provided", func() {
			dividends = []tax.INRDividend{
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-01-10", Amount: 10.0},
					TTRate:   82.5,
				},
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-02-10", Amount: 15.0},
					TTRate:   83.0,
				},
				{
					Dividend: tax.Dividend{Symbol: "IVV", Date: "2024-03-10", Amount: 20.0},
					TTRate:   82.0,
				},
			}

			result, err := valuationManager.ProcessValuations(ctx, valuations, dividends)

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(3))

			// IEF: (10 * 82.5) + (15 * 83.0) = 825 + 1245 = 2070
			Expect(result[0].Ticker).To(Equal("IEF"))
			Expect(result[0].AmountPaid).To(Equal(2070.0))

			// IVV: (20 * 82.0) = 1640
			Expect(result[1].Ticker).To(Equal("IVV"))
			Expect(result[1].AmountPaid).To(Equal(1640.0))

			// TLT: no dividends = 0
			Expect(result[2].Ticker).To(Equal("TLT"))
			Expect(result[2].AmountPaid).To(Equal(0.0))
		})

		It("should set AmountPaid to 0 when no dividends provided", func() {
			result, err := valuationManager.ProcessValuations(ctx, valuations, []tax.INRDividend{})

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(3))

			Expect(result[0].AmountPaid).To(Equal(0.0))
			Expect(result[1].AmountPaid).To(Equal(0.0))
			Expect(result[2].AmountPaid).To(Equal(0.0))
		})

		It("should handle multiple dividends for same ticker", func() {
			dividends = []tax.INRDividend{
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-01-10", Amount: 5.0},
					TTRate:   80.0,
				},
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-02-10", Amount: 10.0},
					TTRate:   81.0,
				},
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-03-10", Amount: 8.0},
					TTRate:   82.0,
				},
			}

			result, err := valuationManager.ProcessValuations(ctx, valuations, dividends)

			Expect(err).ToNot(HaveOccurred())

			// IEF: (5 * 80) + (10 * 81) + (8 * 82) = 400 + 810 + 656 = 1866
			Expect(result[0].Ticker).To(Equal("IEF"))
			Expect(result[0].AmountPaid).To(Equal(1866.0))
		})

		It("should handle dividends for tickers not in valuations", func() {
			dividends = []tax.INRDividend{
				{
					Dividend: tax.Dividend{Symbol: "IEF", Date: "2024-01-10", Amount: 10.0},
					TTRate:   82.0,
				},
				{
					Dividend: tax.Dividend{Symbol: "UNKNOWN", Date: "2024-02-10", Amount: 100.0},
					TTRate:   83.0,
				},
			}

			result, err := valuationManager.ProcessValuations(ctx, valuations, dividends)

			Expect(err).ToNot(HaveOccurred())

			// IEF gets its dividend
			Expect(result[0].Ticker).To(Equal("IEF"))
			Expect(result[0].AmountPaid).To(Equal(820.0))

			// IVV and TLT get 0
			Expect(result[1].AmountPaid).To(Equal(0.0))
			Expect(result[2].AmountPaid).To(Equal(0.0))

			// UNKNOWN ticker dividend is ignored (not in valuations)
		})

		Context("Zero Quantity Positions", func() {
			var (
				testDate  time.Time
				yearEnd   time.Time
				valuation tax.Valuation
			)

			BeforeEach(func() {
				testDate = time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
				yearEnd = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

				// Simulate ADI: Bought Jan 4, Sold Jan 20 (fully liquidated)
				valuation = tax.Valuation{
					Ticker: "ADI",
					FirstPosition: tax.Position{
						Date:     testDate,
						Quantity: 2,
						USDPrice: 181.90,
					},
					PeakPosition: tax.Position{
						Date:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
						Quantity: 2,
						USDPrice: 194.75,
					},
					YearEndPosition: tax.Position{
						Date:     yearEnd,
						Quantity: 0, // ⚠️ FULLY LIQUIDATED
						USDPrice: 0,
					},
				}
			})

			It("should pass all positions to exchange including zero-quantity year-end", func() {
				// Exchange should receive ALL 3 positions
				// Exchange manager should skip zero-value positions internally
				mockExchange.EXPECT().
					Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
					Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
						// All 3 positions should be passed to Exchange
						Expect(exchangeables).To(HaveLen(3),
							"All positions should be passed to Exchange (including zero-qty)")

						// Verify positions
						Expect(exchangeables[0].GetUSDAmount()).To(BeNumerically(">", 0), "First position should have value")
						Expect(exchangeables[1].GetUSDAmount()).To(BeNumerically(">", 0), "Peak position should have value")
						Expect(exchangeables[2].GetUSDAmount()).To(Equal(0.0), "YearEnd position should have zero value")
					}).
					Return(nil)

				result, err := valuationManager.ProcessValuations(ctx, []tax.Valuation{valuation}, nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(1))

				adi := result[0]
				Expect(adi.Ticker).To(Equal("ADI"))

				// ✅ All three positions should be present in output
				Expect(adi.FirstPosition.Quantity).To(Equal(2.0), "FirstPosition preserved")
				Expect(adi.FirstPosition.Date).To(Equal(testDate))

				Expect(adi.PeakPosition.Quantity).To(Equal(2.0), "PeakPosition preserved")

				// ✅ YearEndPosition preserved with zero values (not exchanged)
				Expect(adi.YearEndPosition.Quantity).To(Equal(0.0),
					"YearEnd position preserved in output (audit trail)")
				Expect(adi.YearEndPosition.Date).To(Equal(yearEnd),
					"YearEnd date preserved")
				Expect(adi.YearEndPosition.TTRate).To(Equal(0.0),
					"Zero-quantity position not exchanged (no TTRate)")
				Expect(adi.YearEndPosition.INRValue()).To(Equal(0.0),
					"Zero-quantity position has zero INR value")
			})

		})
	})
})
