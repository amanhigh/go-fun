package manager_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValuationManager", func() {
	var (
		ctx               = context.Background()
		mockTickerManager *mocks.TickerManager
		valuationManager  manager.ValuationManager

		// Common variables
		ticker       = "AAPL"
		year         = 2024
		yearEndDate  = time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		yearEndPrice = 150.00
	)

	BeforeEach(func() {
		mockTickerManager = mocks.NewTickerManager(GinkgoT())
		valuationManager = manager.NewValuationManager(mockTickerManager)
	})

	// Helper function to parse date string to time.Time
	var parseDate = func(date string) time.Time {
		t, err := time.Parse(common.DateOnly, date)
		if err != nil {
			panic(err)
		}
		return t
	}

	Context("Basic Position Tracking", func() {
		Context("Single Buy and Hold", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),
				}

				mockTickerManager.EXPECT().
					GetPrice(ctx, ticker, yearEndDate).
					Return(yearEndPrice, nil)
			})

			It("should compute correct positions", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
				Expect(err).To(BeNil())

				// First position
				Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
				Expect(valuation.FirstPosition.Quantity).To(Equal(trades[0].Quantity))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(trades[0].USDPrice))

				// Peak matches first position for single buy
				Expect(valuation.PeakPosition).To(Equal(valuation.FirstPosition))

				// Year end position
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(trades[0].Quantity))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})

		Context("Complete Position Exit", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),
					tax.NewTrade(ticker, parseDate("2024-02-15"), "SELL", 10, 120),
				}
			})

			It("should compute positions with zero year-end", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
				Expect(err).To(BeNil())

				// First position from buy
				Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
				Expect(valuation.FirstPosition.Quantity).To(Equal(trades[0].Quantity))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(trades[0].USDPrice))

				// Peak position matches first buy
				Expect(valuation.PeakPosition).To(Equal(valuation.FirstPosition))

				// Empty year end position (zero quantity)
				Expect(valuation.YearEndPosition.Quantity).To(Equal(0.0))
			})
		})
	})

	Context("Position Building", func() {
		Context("Gradual Position Increase", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 5, 100),
					tax.NewTrade(ticker, parseDate("2024-02-15"), "BUY", 5, 110),
					tax.NewTrade(ticker, parseDate("2024-03-15"), "BUY", 5, 120),
				}

				mockTickerManager.EXPECT().
					GetPrice(ctx, ticker, yearEndDate).
					Return(yearEndPrice, nil)
			})

			It("should track increasing position correctly", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
				Expect(err).To(BeNil())

				// First position from first buy
				Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
				Expect(valuation.FirstPosition.Quantity).To(Equal(trades[0].Quantity))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(trades[0].USDPrice))

				// Peak position after all buys
				Expect(valuation.PeakPosition.Date).To(Equal(trades[2].Date))
				Expect(valuation.PeakPosition.Quantity).To(Equal(15.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(trades[2].USDPrice))

				// Year end position
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(15.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})
	})
})

Context("Complex Position Building", func() {
    Context("Averaging Down Position", func() {
        var trades []tax.Trade

        BeforeEach(func() {
            trades = []tax.Trade{
                tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 5, 100),   // $500
                tax.NewTrade(ticker, parseDate("2024-02-15"), "BUY", 10, 80),   // $800  - Buying dip
                tax.NewTrade(ticker, parseDate("2024-03-15"), "BUY", 5, 90),    // $450  - Recovery buy
            }
            

            mockTickerManager.EXPECT().
                GetPrice(ctx, ticker, yearEndDate).
                Return(yearEndPrice, nil)
        })

        It("should track averaged position correctly", func() {
            valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
            Expect(err).To(BeNil())

            // First position from initial buy
            Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
            Expect(valuation.FirstPosition.Quantity).To(Equal(5.0))
            Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))
            Expect(valuation.FirstPosition.USDValue).To(Equal(500.0))

            // Peak position at final state
            Expect(valuation.PeakPosition.Date).To(Equal(trades[2].Date))
            Expect(valuation.PeakPosition.Quantity).To(Equal(20.0))     // Total shares: 5 + 10 + 5
            Expect(valuation.PeakPosition.USDPrice).To(Equal(90.0))     // Last trade price
            Expect(valuation.PeakPosition.USDValue).To(Equal(1800.0))   // 20 * 90

            // Year end position
            Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
            Expect(valuation.YearEndPosition.Quantity).To(Equal(20.0))
            Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
            Expect(valuation.YearEndPosition.USDValue).To(Equal(20.0 * yearEndPrice))
        })
    })
})

Context("Position Reduction", func() {
    Context("Partial Position Selling", func() {
        var trades []tax.Trade

        BeforeEach(func() {
            trades = []tax.Trade{
                tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),    // Initial 10 shares
                tax.NewTrade(ticker, parseDate("2024-02-15"), "SELL", 3, 110),    // Sell 3 shares
                tax.NewTrade(ticker, parseDate("2024-03-15"), "SELL", 4, 120),    // Sell 4 shares
            }

            mockTickerManager.EXPECT().
                GetPrice(ctx, ticker, yearEndDate).
                Return(yearEndPrice, nil)
        })

        It("should track partial sells correctly", func() {
            valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
            Expect(err).To(BeNil())

            // First position from initial buy
            Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
            Expect(valuation.FirstPosition.Quantity).To(Equal(10.0))
            Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))
            Expect(valuation.FirstPosition.USDValue).To(Equal(1000.0))

            // Peak position should be initial position
            Expect(valuation.PeakPosition).To(Equal(valuation.FirstPosition))

            // Year end position (3 shares remaining)
            Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
            Expect(valuation.YearEndPosition.Quantity).To(Equal(3.0))    // 10 - 3 - 4 shares
            Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
            Expect(valuation.YearEndPosition.USDValue).To(Equal(3.0 * yearEndPrice))
        })
    })
})

