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

		Context("Averaging Down Position", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 5, 100), // $500
					tax.NewTrade(ticker, parseDate("2024-02-15"), "BUY", 10, 80), // $800  - Buying dip
					tax.NewTrade(ticker, parseDate("2024-03-15"), "BUY", 5, 90),  // $450  - Recovery buy
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
				Expect(valuation.FirstPosition.USDValue()).To(Equal(500.0))

				// Peak position at final state
				Expect(valuation.PeakPosition.Date).To(Equal(trades[2].Date))
				Expect(valuation.PeakPosition.Quantity).To(Equal(20.0))     // Total shares: 5 + 10 + 5
				Expect(valuation.PeakPosition.USDPrice).To(Equal(90.0))     // Last trade price
				Expect(valuation.PeakPosition.USDValue()).To(Equal(1800.0)) // 20 * 90

				// Year end position
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(20.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
				Expect(valuation.YearEndPosition.USDValue()).To(Equal(20.0 * yearEndPrice))
			})
		})
	})

	Context("Complex Scenarios", func() {
		Context("Year End Trading", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 5, 100),
					tax.NewTrade(ticker, parseDate("2024-12-31"), "BUY", 5, 120), // Year end trade
				}

				mockTickerManager.EXPECT().
					GetPrice(ctx, ticker, yearEndDate).
					Return(yearEndPrice, nil)
			})

			It("should handle year end trades correctly", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
				Expect(err).To(BeNil())

				// First position from initial buy
				Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
				Expect(valuation.FirstPosition.Quantity).To(Equal(5.0))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))

				// Peak position should be final position
				Expect(valuation.PeakPosition.Date).To(Equal(trades[1].Date))
				Expect(valuation.PeakPosition.Quantity).To(Equal(10.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(120.0))

				// Year end position
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(10.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})

		Context("Multiple Position Peaks", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				// HACK: Multiple Peaks with Same Value (Take Second higher TBBR Rate) or Throw Error.
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),  // Initial 10
					tax.NewTrade(ticker, parseDate("2024-02-15"), "BUY", 5, 110),   // Peak 1: 15 shares
					tax.NewTrade(ticker, parseDate("2024-03-15"), "SELL", 8, 120),  // Down to 7
					tax.NewTrade(ticker, parseDate("2024-04-15"), "BUY", 10, 115),  // Peak 2: 17 shares
					tax.NewTrade(ticker, parseDate("2024-05-15"), "SELL", 12, 125), // Down to 5
				}

				mockTickerManager.EXPECT().
					GetPrice(ctx, ticker, yearEndDate).
					Return(yearEndPrice, nil)
			})

			It("should identify highest peak through multiple cycles", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
				Expect(err).To(BeNil())

				// First position
				Expect(valuation.FirstPosition.Date).To(Equal(trades[0].Date))
				Expect(valuation.FirstPosition.Quantity).To(Equal(10.0))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))

				// Peak should be second peak with 17 shares
				Expect(valuation.PeakPosition.Date).To(Equal(trades[3].Date))
				Expect(valuation.PeakPosition.Quantity).To(Equal(17.0))  // 7 + 10 shares
				Expect(valuation.PeakPosition.USDPrice).To(Equal(115.0)) // Price at peak

				// Year end position shows final holdings
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(5.0)) // Final position after all trades
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})
	})

	Context("Position Reduction", func() {
		Context("Partial Position Selling", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100), // Initial 10 shares
					tax.NewTrade(ticker, parseDate("2024-02-15"), "SELL", 3, 110), // Sell 3 shares
					tax.NewTrade(ticker, parseDate("2024-03-15"), "SELL", 4, 120), // Sell 4 shares
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
				Expect(valuation.FirstPosition.USDValue()).To(Equal(1000.0))

				// Peak position should be initial position
				Expect(valuation.PeakPosition).To(Equal(valuation.FirstPosition))

				// Year end position (3 shares remaining)
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(3.0)) // 10 - 3 - 4 shares
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
				Expect(valuation.YearEndPosition.USDValue()).To(Equal(3.0 * yearEndPrice))
			})
		})

		Context("Error Cases", func() {
			Context("Empty Trades", func() {
				It("should return error for empty trades", func() {
					trades := []tax.Trade{}
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(Not(BeNil()))
					Expect(err.Error()).To(ContainSubstring("no trades provided"))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Multiple Ticker Trades", func() {
				var trades []tax.Trade

				BeforeEach(func() {
					trades = []tax.Trade{
						tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),
						tax.NewTrade("MSFT", parseDate("2024-02-15"), "BUY", 5, 200), // Different ticker
					}
				})

				It("should return error for mixed ticker trades", func() {
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(Not(BeNil()))
					Expect(err.Error()).To(ContainSubstring("multiple tickers found"))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Year End Price Error", func() {
				var trades []tax.Trade

				BeforeEach(func() {
					trades = []tax.Trade{
						tax.NewTrade(ticker, parseDate("2024-01-15"), "BUY", 10, 100),
					}

					mockTickerManager.EXPECT().
						GetPrice(ctx, ticker, yearEndDate).
						Return(0.0, common.ErrNotFound)
				})

				It("should handle ticker price fetch error", func() {
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(Not(BeNil()))
					Expect(err.Error()).To(ContainSubstring("failed to get year end price"))
					Expect(err.Code()).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})
