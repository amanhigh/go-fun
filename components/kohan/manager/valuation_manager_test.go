package manager_test

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	repoMocks "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to assert valuation positions
var assertValuationPositions = func(valuation tax.Valuation, expectedFirst tax.Position, expectedPeak tax.Position, expectedYearEnd tax.Position) {
	Expect(valuation.FirstPosition.Date).To(Equal(expectedFirst.Date))
	Expect(valuation.FirstPosition.Quantity).To(Equal(expectedFirst.Quantity))
	Expect(valuation.FirstPosition.USDPrice).To(Equal(expectedFirst.USDPrice))

	Expect(valuation.PeakPosition.Date).To(Equal(expectedPeak.Date))
	Expect(valuation.PeakPosition.Quantity).To(Equal(expectedPeak.Quantity))
	Expect(valuation.PeakPosition.USDPrice).To(Equal(expectedPeak.USDPrice))

	Expect(valuation.YearEndPosition.Date).To(Equal(expectedYearEnd.Date))
	Expect(valuation.YearEndPosition.Quantity).To(Equal(expectedYearEnd.Quantity))
	Expect(valuation.YearEndPosition.USDPrice).To(Equal(expectedYearEnd.USDPrice))
}

var _ = Describe("ValuationManager", func() {
	var (
		ctx                 = context.Background()
		mockTickerManager   *mocks.TickerManager
		mockAccountManager  *mocks.AccountManager
		mockTradeRepository *repoMocks.TradeRepository
		mockFyManager       *mocks.FinancialYearManager[tax.Trade]
		mockSBIManager      *mocks.SBIManager
		valuationManager    manager.ValuationManager

		// Common variables
		AAPL         = "AAPL"
		MSFT         = "MSFT"
		year         = 2024
		yearEndDate  = time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		yearEndPrice = 150.00
	)

	BeforeEach(func() {
		mockTickerManager = mocks.NewTickerManager(GinkgoT())
		mockAccountManager = mocks.NewAccountManager(GinkgoT())
		mockTradeRepository = repoMocks.NewTradeRepository(GinkgoT())
		mockFyManager = mocks.NewFinancialYearManager[tax.Trade](GinkgoT())
		mockSBIManager = mocks.NewSBIManager(GinkgoT())
		valuationManager = manager.NewValuationManager(mockTickerManager, mockAccountManager, mockTradeRepository, mockFyManager, mockSBIManager)
	})

	Context("Analyse Valuation", func() {

		Context("Fresh Start", func() {
			BeforeEach(func() {
				// All tests under Fresh Start expect no last year position
				mockAccountManager.EXPECT().
					GetRecord(ctx, AAPL, year-1).
					Return(tax.Account{}, common.ErrNotFound)
			})

			Context("Basic Position Tracking", func() {
				Context("Single Buy and Hold", func() {
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						// Daily prices for peak calculation (match trade execution price)
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
						}
						// Daily rates for INR calculation
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.5,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should compute correct positions", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
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
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
							tax.NewTrade(AAPL, "2024-02-15", "SELL", 10, 120),
						}

						// Daily prices for peak calculation
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-02-15": 120.0,
						}
						// Daily rates for INR calculation
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.5,
							"2024-02-15": 83.0,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						// NOTE: No GetPrice mock needed - position fully exits (quantity = 0),
						// so determineYearEndPosition doesn't call GetPrice (see line 249 in valuation_manager.go)
					})

					It("should compute positions with zero year-end", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position from buy
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
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
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 5, 100),
							tax.NewTrade(AAPL, "2024-02-15", "BUY", 5, 110),
							tax.NewTrade(AAPL, "2024-03-15", "BUY", 5, 120),
						}

						// Daily prices for peak calculation
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-02-15": 110.0,
							"2024-03-15": 120.0,
						}
						// Daily rates for INR calculation (increasing rates ensure peak stays on last trade)
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.5,
							"2024-02-15": 83.0,
							"2024-03-15": 83.5,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track increasing position correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position from first buy
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
						Expect(valuation.FirstPosition.Quantity).To(Equal(trades[0].Quantity))
						Expect(valuation.FirstPosition.USDPrice).To(Equal(trades[0].USDPrice))

						// Peak position after all buys
						date, getDateErr = trades[2].GetDate() // Reuse getDateErr
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.PeakPosition.Date).To(Equal(date))
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
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 5, 100), // $500
							tax.NewTrade(AAPL, "2024-02-15", "BUY", 10, 80), // $800  - Buying dip
							tax.NewTrade(AAPL, "2024-03-15", "BUY", 5, 90),  // $450  - Recovery buy
						}

						// Daily prices for peak calculation
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-02-15": 80.0,
							"2024-03-15": 90.0,
						}
						// Daily rates for INR calculation (increasing rates ensure peak stays on last trade)
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.5,
							"2024-02-15": 83.0,
							"2024-03-15": 83.5,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track averaged position correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position from initial buy
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
						Expect(valuation.FirstPosition.Quantity).To(Equal(5.0))
						Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))
						Expect(valuation.FirstPosition.USDValue()).To(Equal(500.0))

						// Peak position at final state
						date, getDateErr = trades[2].GetDate() // Reuse getDateErr
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.PeakPosition.Date).To(Equal(date))
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
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 5, 100),
							tax.NewTrade(AAPL, "2024-12-31", "BUY", 5, 120), // Year end trade
						}

						// Daily prices for peak calculation (match trade execution prices)
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-12-31": 120.0, // Year-end trade price
						}
						// Daily rates for INR calculation (increasing rates ensure Dec 31 peak stays)
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.5,
							"2024-12-31": 84.0, // Higher rate ensures Dec 31 wins
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should handle year end trades correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						firstPosDate, _ := trades[0].GetDate()
						peakPosDate, _ := trades[1].GetDate()

						assertValuationPositions(valuation, tax.Position{
							Date:     firstPosDate,
							Quantity: 5.0,
							USDPrice: 100.0,
						}, tax.Position{
							Date:     peakPosDate,
							Quantity: 10.0,
							USDPrice: 120.0,
						}, tax.Position{
							Date:     yearEndDate,
							Quantity: 10.0,
							USDPrice: yearEndPrice,
						})
					})
				})

				Context("Multiple Position Peaks", func() {
					var trades []tax.Trade
					BeforeEach(func() {
						// FIXME: #A Multiple Peaks with Same Value (Take Second higher TBBR Rate) or Throw Error.
						// Test validates that when multiple position peaks exist with SAME quantity,
						// the second peak wins because it has HIGHER USD price AND HIGHER TBBR rate
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),  // Initial 10 shares
							tax.NewTrade(AAPL, "2024-02-15", "BUY", 5, 110),   // Peak 1: 15 shares @ $110
							tax.NewTrade(AAPL, "2024-03-15", "SELL", 5, 120),  // Down to 10 shares
							tax.NewTrade(AAPL, "2024-04-15", "BUY", 5, 115),   // Peak 2: 15 shares @ $115 ← Expected peak (same qty, higher price+rate)
							tax.NewTrade(AAPL, "2024-05-15", "SELL", 12, 125), // Down to 3 shares
						}

						// Daily prices for peak calculation
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-02-15": 110.0,
							"2024-03-15": 120.0,
							"2024-04-15": 115.0,
							"2024-05-15": 125.0,
						}

						// Daily rates for peak calculation
						// Strategy: Both Feb 15 and Apr 15 have 15 shares (SAME quantity)
						// Difference: Apr 15 has higher USD price ($115 vs $110) AND higher TBBR rate
						// Apr 15 wins when using > comparison because: 15×115×85.0 > 15×110×82.0
						//
						// INR Calculations:
						// Jan 15: 10 × 100 × 80.0 = 80,000 INR
						// Feb 15: 15 × 110 × 82.0 = 132,300 INR (first peak candidate, same qty as Apr 15)
						// Mar 15: 10 × 120 × 82.5 = 99,000 INR (position reduces from 15)
						// Apr 15: 15 × 115 × 85.0 = 147,375 INR ← PEAK (second peak, same qty but higher price+rate)
						// May 15: 3 × 125 × 85.5 = 31,912.5 INR (position reduces from 15)
						mergedDailyRates := map[string]float64{
							"2024-01-15": 80.0,
							"2024-02-15": 82.0,
							"2024-03-15": 82.5,
							"2024-04-15": 85.0, // Higher TBBR rate for Apr 15 (second peak advantage)
							"2024-05-15": 85.5,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(mergedDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should identify highest peak through multiple cycles", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						firstPosDate, _ := trades[0].GetDate()
						peakPosDate, _ := trades[3].GetDate()

						assertValuationPositions(valuation, tax.Position{
							Date:     firstPosDate,
							Quantity: 10.0,
							USDPrice: 100.0,
						}, tax.Position{
							Date:     peakPosDate,
							Quantity: 15.0,
							USDPrice: 115.0,
						}, tax.Position{
							Date:     yearEndDate,
							Quantity: 3.0,
							USDPrice: yearEndPrice,
						})
					})
				})
			})

			Context("Peak Should be INR Based (Qty × Price × Rate)", func() {
				// Test suite validating that peak is determined by INR value (Qty × USD_Price × SBI_Rate),
				// NOT just by quantity. All scenarios use constant quantity to isolate price/rate variations.
				// This ensures compliance with Tax.md Line 124 daily peak calculation requirement.

				Context("Scenario 1: Same Qty, Lower USD Price WINS due to Higher Rate", func() {
					// Prove that exchange rate can dominate over USD price
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						// Same quantity (10 shares) held on both dates
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0, // Lower price
							"2024-06-15": 120.0, // Higher price (but won't be peak)
						}

						// Rate determines peak: higher rate on lower price date
						aaplDailyRates := map[string]float64{
							"2024-01-15": 85.0, // Higher rate on lower price date
							"2024-06-15": 70.0, // Lower rate on higher price date
						}

						// INR Calculations:
						// Jan 15: 10 × $100 × 85.0 = 85,000 INR ← PEAK
						// Jun 15: 10 × $120 × 70.0 = 84,000 INR (higher price but lower INR)

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should identify peak on lower USD price date due to higher exchange rate", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// Peak should be Jan 15 (lower USD price but higher rate)
						peakDate, _ := trades[0].GetDate()
						Expect(valuation.PeakPosition.Date).To(Equal(peakDate))
						Expect(valuation.PeakPosition.Quantity).To(Equal(10.0))
						Expect(valuation.PeakPosition.USDPrice).To(Equal(100.0))

						// Verify INR values: Jan 15 > Jun 15
						jan15INR := 10.0 * 100.0 * 85.0 // = 85,000 INR
						jun15INR := 10.0 * 120.0 * 70.0 // = 84,000 INR
						Expect(jan15INR).To(BeNumerically(">", jun15INR))

						// Verify peak is indeed the maximum INR value
						actualPeakINR := valuation.PeakPosition.Quantity * valuation.PeakPosition.USDPrice * 85.0
						Expect(actualPeakINR).To(Equal(jan15INR))
						Expect(actualPeakINR).To(BeNumerically(">", jun15INR))
					})
				})

				Context("Scenario 2: Same Qty, Same Price, Rate Determines Peak", func() {
					// Pure rate variation test: isolates rate as the only changing factor
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						// Same quantity and same price on both dates
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-06-15": 100.0, // Same price as Jan 15
						}

						// Only rate varies: higher rate later
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.0,
							"2024-06-15": 84.0, // Higher rate
						}

						// INR Calculations:
						// Jan 15: 10 × $100 × 82.0 = 82,000 INR
						// Jun 15: 10 × $100 × 84.0 = 84,000 INR ← PEAK (only rate differs)

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should identify peak on higher rate date when price and quantity are same", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// Peak should be Jun 15 (higher rate)
						Expect(valuation.PeakPosition.Date).To(Equal(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)))
						Expect(valuation.PeakPosition.Quantity).To(Equal(10.0))
						Expect(valuation.PeakPosition.USDPrice).To(Equal(100.0))

						// Verify INR values: Jun 15 > Jan 15
						jan15INR := 10.0 * 100.0 * 82.0 // = 82,000 INR
						jun15INR := 10.0 * 100.0 * 84.0 // = 84,000 INR
						Expect(jun15INR).To(BeNumerically(">", jan15INR))

						// Verify only rate changed
						actualPeakINR := valuation.PeakPosition.Quantity * valuation.PeakPosition.USDPrice * 84.0
						Expect(actualPeakINR).To(Equal(jun15INR))
						Expect(actualPeakINR).To(BeNumerically(">", jan15INR))
					})
				})

				Context("Scenario 3: Same Qty, Price-Rate Tradeoff (Multi-Factor Optimization)", func() {
					// Both price and rate vary across 3 dates: tests multi-factor optimization
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						// Quantity constant (10 shares) across all dates, but price varies
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0, // Medium price
							"2024-06-15": 120.0, // High price ← Peak here due to price × rate product
							"2024-09-15": 90.0,  // Low price
						}

						// Rate varies inversely to price to create optimization scenario
						aaplDailyRates := map[string]float64{
							"2024-01-15": 82.0, // Medium rate
							"2024-06-15": 75.0, // Low rate (but high price compensates)
							"2024-09-15": 92.0, // High rate (but low price doesn't compensate)
						}

						// INR Calculations:
						// Jan 15: 10 × $100 × 82.0 = 82,000 INR
						// Jun 15: 10 × $120 × 75.0 = 90,000 INR ← PEAK (high price × rate product wins)
						// Sep 15: 10 × $90 × 92.0 = 82,800 INR (high rate can't compensate for low price)

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should identify peak where Price × Rate product is optimal", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// Peak should be Jun 15 (high price dominates despite low rate)
						Expect(valuation.PeakPosition.Date).To(Equal(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)))
						Expect(valuation.PeakPosition.Quantity).To(Equal(10.0))
						Expect(valuation.PeakPosition.USDPrice).To(Equal(120.0))

						// Verify INR calculations for all three dates
						jan15INR := 10.0 * 100.0 * 82.0 // = 82,000 INR
						jun15INR := 10.0 * 120.0 * 75.0 // = 90,000 INR ← PEAK
						sep15INR := 10.0 * 90.0 * 92.0  // = 82,800 INR

						// Verify Jun 15 is the maximum
						Expect(jun15INR).To(BeNumerically(">", jan15INR))
						Expect(jun15INR).To(BeNumerically(">", sep15INR))

						// Verify peak INR matches Jun 15 calculation
						actualPeakINR := valuation.PeakPosition.Quantity * valuation.PeakPosition.USDPrice * 75.0
						Expect(actualPeakINR).To(Equal(jun15INR))
					})
				})

				Context("Scenario 4: Same Qty, Lower Price & Rate LOSES (Negative Test)", func() {
					// Validate that when both price and rate are lower, peak doesn't shift
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						// Quantity constant (10 shares), but price lower on later date
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-06-15": 95.0, // Lower price
						}

						// Rate also lower on the later date
						aaplDailyRates := map[string]float64{
							"2024-01-15": 83.0,
							"2024-06-15": 80.0, // Lower rate
						}

						// INR Calculations:
						// Jan 15: 10 × $100 × 83.0 = 83,000 INR ← PEAK
						// Jun 15: 10 × $95 × 80.0 = 76,000 INR (both factors lower)

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should not shift peak when both price and rate are lower", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// Peak should remain Jan 15 (both price and rate lower on Jun 15)
						peakDate, _ := trades[0].GetDate()
						Expect(valuation.PeakPosition.Date).To(Equal(peakDate))
						Expect(valuation.PeakPosition.Quantity).To(Equal(10.0))
						Expect(valuation.PeakPosition.USDPrice).To(Equal(100.0))

						// Verify INR values: Jan 15 > Jun 15
						jan15INR := 10.0 * 100.0 * 83.0 // = 83,000 INR
						jun15INR := 10.0 * 95.0 * 80.0  // = 76,000 INR

						// Verify Jan 15 remains peak
						Expect(jan15INR).To(BeNumerically(">", jun15INR))
						actualPeakINR := valuation.PeakPosition.Quantity * valuation.PeakPosition.USDPrice * 83.0
						Expect(actualPeakINR).To(Equal(jan15INR))
						Expect(actualPeakINR).To(BeNumerically(">", jun15INR))
					})
				})
			})

			Context("Position Reduction", func() {
				Context("Partial Position Selling", func() {
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100), // Initial 10 shares
							tax.NewTrade(AAPL, "2024-02-15", "SELL", 3, 110), // Sell 3 shares
							tax.NewTrade(AAPL, "2024-03-15", "SELL", 4, 120), // Sell 4 shares
						}

						// Daily prices for peak calculation
						aaplDailyPrices := map[string]float64{
							"2024-01-15": 100.0,
							"2024-02-15": 110.0,
							"2024-03-15": 120.0,
						}

						// Daily rates for peak calculation (decreasing to ensure Jan 15 is peak)
						// Jan 15: 10 × 100 × 83.0 = 83,000 INR ← PEAK
						// Feb 15: 7 × 110 × 82.5 = 63,525 INR
						// Mar 15: 3 × 120 × 82.0 = 29,520 INR
						aaplDailyRates := map[string]float64{
							"2024-01-15": 83.0,
							"2024-02-15": 82.5,
							"2024-03-15": 82.0,
						}

						mockTickerManager.EXPECT().
							GetDailyPrices(ctx, AAPL, year).
							Return(aaplDailyPrices, nil)
						mockSBIManager.EXPECT().
							GetDailyRates(ctx, year).
							Return(aaplDailyRates, nil)
						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track partial sells correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position from initial buy
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
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
			})

		})

		Context("First Trade is Sell on Fresh Start", func() {
			It("should return error", func() {
				trades := []tax.Trade{
					tax.NewTrade(AAPL, "2024-01-15", "SELL", 10, 100),
				}

				// Daily prices and rates needed (calculateDailyPeak called before processTrades validation)
				aaplDailyPrices := map[string]float64{
					"2024-01-15": 100.0,
				}
				aaplDailyRates := map[string]float64{
					"2024-01-15": 82.5,
				}

				mockAccountManager.EXPECT().GetRecord(ctx, AAPL, year-1).Return(tax.Account{}, common.ErrNotFound)
				mockTickerManager.EXPECT().
					GetDailyPrices(ctx, AAPL, year).
					Return(aaplDailyPrices, nil)
				mockSBIManager.EXPECT().
					GetDailyRates(ctx, year).
					Return(aaplDailyRates, nil)

				_, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring("first trade can't be sell on fresh start"))
			})
		})

		Context("Title Case Trade Types (Real DriveWealth Data)", func() {
			It("should handle 'Buy'/'Sell' trade types case-insensitively", func() {
				// Real DriveWealth data has "Buy"/"Sell" (title case), not "BUY"/"SELL"
				trades := []tax.Trade{
					{Symbol: AAPL, Date: "2024-01-15", Type: "Buy", Quantity: 10, USDPrice: 100.0},
					{Symbol: AAPL, Date: "2024-02-15", Type: "Buy", Quantity: 5, USDPrice: 110.0},
					{Symbol: AAPL, Date: "2024-03-15", Type: "Sell", Quantity: 3, USDPrice: 120.0},
				}

				// Daily prices for peak calculation
				aaplDailyPrices := map[string]float64{
					"2024-01-15": 100.0,
					"2024-02-15": 110.0,
					"2024-03-15": 120.0,
				}

				// Daily rates for peak calculation (increasing to ensure Feb 15 is peak)
				// Jan 15: 10 × 100 × 82.0 = 82,000 INR
				// Feb 15: 15 × 110 × 83.0 = 137,250 INR ← PEAK
				// Mar 15: 12 × 120 × 82.5 = 118,800 INR
				aaplDailyRates := map[string]float64{
					"2024-01-15": 82.0,
					"2024-02-15": 83.0,
					"2024-03-15": 82.5,
				}

				mockAccountManager.EXPECT().
					GetRecord(ctx, AAPL, year-1).
					Return(tax.Account{}, common.ErrNotFound) // Fresh start

				mockTickerManager.EXPECT().
					GetDailyPrices(ctx, AAPL, year).
					Return(aaplDailyPrices, nil)
				mockSBIManager.EXPECT().
					GetDailyRates(ctx, year).
					Return(aaplDailyRates, nil)
				mockTickerManager.EXPECT().
					GetPrice(ctx, AAPL, yearEndDate).
					Return(yearEndPrice, nil)

				valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
				Expect(err).ToNot(HaveOccurred())

				// Verify correct quantity calculations despite title case
				// 10 (buy) + 5 (buy) - 3 (sell) = 12 shares at year end
				Expect(valuation.YearEndPosition.Quantity).To(Equal(12.0))

				// Verify peak was tracked correctly (after second buy)
				Expect(valuation.PeakPosition.Quantity).To(Equal(15.0))

				// Verify first position set correctly
				Expect(valuation.FirstPosition.Quantity).To(Equal(10.0))

				// Verify all positions have valid dates (critical for ProcessValuations -> Exchange)
				Expect(valuation.FirstPosition.Date).ToNot(BeZero())
				Expect(valuation.PeakPosition.Date).ToNot(BeZero())
				Expect(valuation.YearEndPosition.Date).ToNot(BeZero())
				Expect(valuation.YearEndPosition.Date).To(Equal(yearEndDate))
			})
		})

		Context("Error Cases", func() {
			Context("Empty Trades", func() {
				It("should return error for empty trades", func() {
					trades := []tax.Trade{}
					mockAccountManager.EXPECT().GetRecord(ctx, AAPL, year-1).Return(tax.Account{}, common.ErrNotFound).Maybe()
					_, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(AAPL))
					Expect(err.Error()).To(ContainSubstring("no trades or carry-over position provided"))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Multiple Ticker Trades", func() {
				var trades []tax.Trade

				BeforeEach(func() {
					trades = []tax.Trade{
						tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						tax.NewTrade("MSFT", "2024-02-15", "BUY", 5, 200), // Different ticker
					}
				})

				It("should return error for mixed ticker trades", func() {
					// We pass the first trade's ticker as the expected one.
					// The validation logic should then find the mismatch.
					_, err := valuationManager.AnalyzeValuation(ctx, trades[0].Symbol, trades, year)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("trade symbol mismatch"))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Year End Price Error", func() {
				var trades []tax.Trade

				BeforeEach(func() {
					trades = []tax.Trade{
						tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
					}

					// Daily prices and rates for peak calculation (succeeds)
					aaplDailyPrices := map[string]float64{
						"2024-01-15": 100.0,
					}
					aaplDailyRates := map[string]float64{
						"2024-01-15": 82.5,
					}

					mockAccountManager.EXPECT().
						GetRecord(ctx, AAPL, year-1).
						Return(tax.Account{}, common.ErrNotFound)
					mockTickerManager.EXPECT().
						GetDailyPrices(ctx, AAPL, year).
						Return(aaplDailyPrices, nil)
					mockSBIManager.EXPECT().
						GetDailyRates(ctx, year).
						Return(aaplDailyRates, nil)
					mockTickerManager.EXPECT().
						GetPrice(ctx, AAPL, yearEndDate).
						Return(0.0, common.ErrNotFound)
				})

				It("should fail when ticker price fetch fails", func() {
					_, err := valuationManager.AnalyzeValuation(ctx, AAPL, trades, year)
					Expect(err).To(HaveOccurred())                                           // Should fail fast
					Expect(err.Code()).To(Equal(http.StatusInternalServerError))             // Server error expected
					Expect(err.Error()).To(ContainSubstring("failed to get year end price")) // Specific error message
				})
			})
		})
	})

	Context("GetYearlyValuationsUSD", func() {
		var (
			// Define sample trades for multiple tickers
			tradeAAPL1 = tax.NewTrade(AAPL, "2024-01-10", "BUY", 10, 100) // Date: 2024-01-10
			tradeAAPL2 = tax.NewTrade(AAPL, "2024-05-15", "SELL", 5, 120) // Date: 2024-05-15
			tradeMSFT1 = tax.NewTrade(MSFT, "2024-02-20", "BUY", 20, 200) // Date: 2024-02-20
			allTrades  []tax.Trade
		)

		BeforeEach(func() {
			// Reset trades for each test
			allTrades = []tax.Trade{tradeAAPL1, tradeAAPL2, tradeMSFT1}
		})

		It("should process multiple tickers successfully", func() {
			// Mock Repo: GetAllRecords returns combined trades
			mockTradeRepository.EXPECT().GetAllRecords(ctx).Return(allTrades, nil).Once()

			// Mock FY Manager: FilterUS returns the same trades (assuming all are in the year for this test)
			mockFyManager.EXPECT().FilterUS(ctx, allTrades, year).Return(allTrades, nil).Once()

			// Mock AccountManager: GetAllRecords for carry-over (none in this test)
			mockAccountManager.EXPECT().GetAllRecords(ctx, year-1).Return([]tax.Account{}, common.ErrNotFound).Once()

			// Daily prices for AAPL (two trades: buy on Jan 10, sell on May 15)
			aaplDailyPrices := map[string]float64{
				"2024-01-10": 100.0,
				"2024-05-15": 120.0,
			}

			// Daily prices for MSFT (one trade: buy on Feb 20)
			msftDailyPrices := map[string]float64{
				"2024-02-20": 200.0,
			}

			// Merged daily rates (Q4: Option B - overlapping rates when dates conflict)
			// Strategy: Ensure each ticker's peak date has highest rate for that ticker
			// AAPL: Peak on Jan 10 (10 shares), Jan 10 rate must be > Feb 20 rate to avoid backfill issue
			// MSFT: Peak on Feb 20 (30 shares), Feb 20 rate highest
			mergedDailyRates := map[string]float64{
				"2024-01-10": 83.5, // Highest for AAPL (Jan 10: 10×100×83.5 = 83,500 INR > Feb 20: 10×100×82.5)
				"2024-02-20": 82.5, // Lower for Feb 20 to ensure AAPL Jan 10 stays peak
				"2024-05-15": 82.0,
			}

			// Mock dependencies needed by AnalyzeValuation for AAPL
			mockAccountManager.EXPECT().GetRecord(ctx, AAPL, year-1).Return(tax.Account{}, common.ErrNotFound).Once() // Fresh start for AAPL
			mockTickerManager.EXPECT().GetDailyPrices(ctx, AAPL, year).Return(aaplDailyPrices, nil).Once()
			mockTickerManager.EXPECT().GetPrice(ctx, AAPL, yearEndDate).Return(150.0, nil).Once() // AAPL year end price

			// Mock dependencies needed by AnalyzeValuation for MSFT
			mockAccountManager.EXPECT().GetRecord(ctx, "MSFT", year-1).Return(tax.Account{Quantity: 10, MarketValue: 1800}, nil).Once() // Start MSFT with 10 shares @ 180
			mockTickerManager.EXPECT().GetDailyPrices(ctx, "MSFT", year).Return(msftDailyPrices, nil).Once()
			mockTickerManager.EXPECT().GetPrice(ctx, "MSFT", yearEndDate).Return(210.0, nil).Once() // MSFT year end price

			// Merged SBI rates (Q1: Answer B - called once per ticker, so expect 2 calls total)
			mockSBIManager.EXPECT().GetDailyRates(ctx, year).Return(mergedDailyRates, nil).Times(2)

			// Call the target method
			valuations, err := valuationManager.GetYearlyValuationsUSD(ctx, year)

			// Assertions
			Expect(err).ToNot(HaveOccurred())
			Expect(valuations).To(HaveLen(2))

			// Sort valuations by ticker to ensure consistent order
			sort.Slice(valuations, func(i, j int) bool {
				return valuations[i].Ticker < valuations[j].Ticker
			})

			var aaplVal = valuations[0]
			var msftVal = valuations[1]

			// Assert AAPL Valuation (based on AnalyzeValuation logic)
			Expect(aaplVal.Ticker).To(Equal(AAPL))
			aaplFirstDate, getDateErr := tradeAAPL1.GetDate()
			Expect(getDateErr).NotTo(HaveOccurred())
			Expect(aaplVal.FirstPosition.Quantity).To(Equal(10.0))                 // From tradeAAPL1
			Expect(aaplVal.FirstPosition.Date).To(Equal(aaplFirstDate))            // Date of first buy
			Expect(aaplVal.PeakPosition.Quantity).To(Equal(10.0))                  // Peak was initial buy
			Expect(aaplVal.PeakPosition.Date).To(Equal(aaplFirstDate))             // Date peak reached
			Expect(aaplVal.YearEndPosition.Quantity).To(BeNumerically("~", 5.0))   // 10 - 5
			Expect(aaplVal.YearEndPosition.Date).To(Equal(yearEndDate))            // Dec 31st
			Expect(aaplVal.YearEndPosition.USDPrice).To(BeNumerically("~", 150.0)) // Mocked year end price

			// Assert MSFT Valuation (based on AnalyzeValuation logic)
			Expect(msftVal.Ticker).To(Equal("MSFT"))
			Expect(msftVal.FirstPosition.Quantity).To(Equal(10.0))                                        // From starting position
			Expect(msftVal.FirstPosition.Date).To(Equal(time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC))) // Date of start pos (Dec 31st of previous year)
			Expect(msftVal.PeakPosition.Quantity).To(Equal(30.0))                                         // 10 start + 20 buy
			msftPeakDate, getDateErr := tradeMSFT1.GetDate()                                              // Use getDateErr here too
			Expect(getDateErr).NotTo(HaveOccurred())
			Expect(msftVal.PeakPosition.Date).To(Equal(msftPeakDate))              // Date peak reached
			Expect(msftVal.YearEndPosition.Quantity).To(BeNumerically("~", 30.0))  // Final quantity
			Expect(msftVal.YearEndPosition.Date).To(Equal(yearEndDate))            // Dec 31st
			Expect(msftVal.YearEndPosition.USDPrice).To(BeNumerically("~", 210.0)) // Mocked year end price
		})

		It("should return empty slice if no trades found", func() {
			// Mock Repo: GetAllRecords returns empty list initially
			initialTrades := []tax.Trade{}
			mockTradeRepository.EXPECT().GetAllRecords(ctx).Return(initialTrades, nil).Once()

			// Mock FY Manager: FilterUS is called even with empty initial trades
			// It should return an empty slice and no error in this scenario before the length check.
			mockFyManager.EXPECT().FilterUS(ctx, initialTrades, year).Return([]tax.Trade{}, nil).Once()

			// Mock AccountManager: GetAllRecords for carry-over (none in this test)
			mockAccountManager.EXPECT().GetAllRecords(ctx, year-1).Return([]tax.Account{}, common.ErrNotFound).Once()

			valuations, err := valuationManager.GetYearlyValuationsUSD(ctx, year)

			// Assert that the specific StatusNotFound error is returned
			Expect(err).To(HaveOccurred())
			Expect(err.Code()).To(Equal(http.StatusNotFound))
			Expect(err.Error()).To(ContainSubstring("no trades or carry-over positions found for year"))
			Expect(valuations).To(BeEmpty()) // Expect empty valuations slice
		})

		It("should return error from GetAllRecords", func() {
			// Mock Repo: GetAllRecords returns a generic error
			expectedErr := common.NewServerError(errors.New("repo failed"))
			mockTradeRepository.EXPECT().GetAllRecords(ctx).Return(nil, expectedErr).Once()

			// No need to mock FilterUS as it won't be called if GetAllRecords fails

			_, err := valuationManager.GetYearlyValuationsUSD(ctx, year)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedErr))
		})

		It("should fail when any ticker price fetch fails", func() {
			// Mock Repo: GetAllRecords returns trades for two tickers
			mockTradeRepository.EXPECT().GetAllRecords(ctx).Return(allTrades, nil).Once()

			// Mock FY Manager: FilterUS returns the same trades
			mockFyManager.EXPECT().FilterUS(ctx, allTrades, year).Return(allTrades, nil).Once()

			// Mock AccountManager: GetAllRecords for carry-over (none in this test)
			mockAccountManager.EXPECT().GetAllRecords(ctx, year-1).Return([]tax.Account{}, common.ErrNotFound).Once()

			// Daily prices for peak calculation (AAPL only - MSFT never reached)
			aaplDailyPrices := map[string]float64{
				"2024-01-10": 100.0,
				"2024-05-15": 120.0,
			}

			// Merged daily rates for both tickers (used by AAPL peak calc before failure)
			mergedDailyRates := map[string]float64{
				"2024-01-10": 82.0,
				"2024-02-20": 83.0, // MSFT trade date (never reached)
				"2024-05-15": 82.5, // AAPL trade date
			}

			// Mock dependencies for AAPL (price fetch fails)
			mockAccountManager.EXPECT().GetRecord(ctx, AAPL, year-1).Return(tax.Account{}, common.ErrNotFound).Once()
			mockTickerManager.EXPECT().GetDailyPrices(ctx, AAPL, year).Return(aaplDailyPrices, nil).Once()
			mockSBIManager.EXPECT().GetDailyRates(ctx, year).Return(mergedDailyRates, nil).Once() // Called once for AAPL
			expectedErr := common.NewServerError(errors.New("price fetch failed"))
			mockTickerManager.EXPECT().GetPrice(ctx, AAPL, yearEndDate).Return(0.0, expectedErr).Once()

			// MSFT dependencies should NOT be called since AAPL fails first (fail-fast behavior)
			// No MSFT mocks needed - they should never be reached

			_, err := valuationManager.GetYearlyValuationsUSD(ctx, year)
			// Assertions: Should fail fast when any ticker has price fetch error
			Expect(err).To(HaveOccurred())
			Expect(err.Code()).To(Equal(http.StatusInternalServerError))
			Expect(err.Error()).To(ContainSubstring("failed to get year end price"))
		})

		Context("Carry-Over Without Trades", func() {
			// TDD Bug Fix: Ticker with carry-over from previous year but NO trades in current year
			// Should still appear in valuations (SIVR/TLT missing from 2023 accounts bug)
			var (
				tickerWithTrades    string
				tickerWithoutTrades string
				prevYearAccount     tax.Account
				tradeInYear         tax.Trade
			)

			BeforeEach(func() {
				tickerWithTrades = "AAPL"
				tickerWithoutTrades = "MSFT"
				prevYearAccount = tax.Account{Symbol: tickerWithoutTrades, Quantity: 50, Cost: 10000, MarketValue: 10000}
				tradeInYear = tax.NewTrade(tickerWithTrades, "2024-06-15", tax.TRADE_TYPE_BUY, 10, 150)

				// Mock: Only AAPL has trades in target year
				mockTradeRepository.EXPECT().GetAllRecords(ctx).Return([]tax.Trade{tradeInYear}, nil).Once()
				mockFyManager.EXPECT().FilterUS(ctx, []tax.Trade{tradeInYear}, year).Return([]tax.Trade{tradeInYear}, nil).Once()

				// Mock: MSFT has carry-over from previous year (no trades in current year)
				mockAccountManager.EXPECT().GetAllRecords(ctx, year-1).Return([]tax.Account{prevYearAccount}, nil).Once()

				// Daily prices for AAPL (fresh start with one trade on Jun 15)
				aaplDailyPrices := map[string]float64{
					"2024-06-15": 150.0,
				}

				// Daily prices for MSFT (empty - no trades in current year, Q3: Yes, empty maps work)
				msftDailyPrices := map[string]float64{}

				// Merged daily rates (both tickers)
				mergedDailyRates := map[string]float64{
					"2024-06-15": 82.5, // AAPL trade
				}

				// Mock: AAPL dependencies (fresh start in current year)
				mockAccountManager.EXPECT().GetRecord(ctx, tickerWithTrades, year-1).Return(tax.Account{}, common.ErrNotFound).Once()
				mockTickerManager.EXPECT().GetDailyPrices(ctx, tickerWithTrades, year).Return(aaplDailyPrices, nil).Once()
				mockTickerManager.EXPECT().GetPrice(ctx, tickerWithTrades, yearEndDate).Return(160.0, nil).Once()

				// Mock: MSFT dependencies (carry-over, no trades)
				mockAccountManager.EXPECT().GetRecord(ctx, tickerWithoutTrades, year-1).Return(prevYearAccount, nil).Once()
				mockTickerManager.EXPECT().GetDailyPrices(ctx, tickerWithoutTrades, year).Return(msftDailyPrices, nil).Once()
				mockTickerManager.EXPECT().GetPrice(ctx, tickerWithoutTrades, yearEndDate).Return(210.0, nil).Once()

				// Merged SBI rates (Q1: Answer B - called once per ticker)
				mockSBIManager.EXPECT().GetDailyRates(ctx, year).Return(mergedDailyRates, nil).Times(2)
			})

			It("should include ticker with carry-over but no trades", func() {
				// Execute
				valuations, err := valuationManager.GetYearlyValuationsUSD(ctx, year)

				// Assert: Should include BOTH tickers (AAPL with trades + MSFT carry-over only)
				Expect(err).ToNot(HaveOccurred())
				Expect(valuations).To(HaveLen(2), "Should include ticker with trades AND ticker with carry-over only")

				// Find MSFT (carry-over without trades)
				var msftVal *tax.Valuation
				for i := range valuations {
					if valuations[i].Ticker == tickerWithoutTrades {
						msftVal = &valuations[i]
						break
					}
				}
				Expect(msftVal).ToNot(BeNil(), "MSFT should be included despite no trades")

				// Assert FirstPosition: Should be carried from previous year
				prevYearEnd := time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC)
				Expect(msftVal.FirstPosition.Date).To(Equal(prevYearEnd), "FirstPosition date from carry-over (Dec 31 of previous year)")
				Expect(msftVal.FirstPosition.Quantity).To(Equal(50.0), "FirstPosition quantity from carry-over")
				Expect(msftVal.FirstPosition.USDPrice).To(Equal(200.0), "FirstPosition price from carry-over (10000/50)")

				// Assert PeakPosition: Should remain same as FirstPosition (no trades to change it)
				Expect(msftVal.PeakPosition.Date).To(Equal(prevYearEnd), "PeakPosition date unchanged (no trades)")
				Expect(msftVal.PeakPosition.Quantity).To(Equal(50.0), "PeakPosition quantity unchanged (no trades)")
				Expect(msftVal.PeakPosition.USDPrice).To(Equal(200.0), "PeakPosition price unchanged (no trades)")

				// Assert YearEndPosition: Should have current year date with year-end price
				Expect(msftVal.YearEndPosition.Date).To(Equal(yearEndDate), "YearEndPosition date is current year Dec 31")
				Expect(msftVal.YearEndPosition.Quantity).To(Equal(50.0), "YearEndPosition quantity unchanged (no trades)")
				Expect(msftVal.YearEndPosition.USDPrice).To(Equal(210.0), "YearEndPosition price from year-end lookup")
			})
		})
	})

	Context("With Carry-Over Position", func() {
		var (
			carryOverAccount tax.Account
			tradesInYear     []tax.Trade
		)

		BeforeEach(func() {
			carryOverAccount = tax.Account{
				Symbol:      AAPL,
				Quantity:    50,
				MarketValue: 8000.00, // Implies $160.00 per share
			}
		})

		Context("With Trades in Year", func() {
			var (
				testYear     = 2023
				yearEndDate  = time.Date(testYear, 12, 31, 0, 0, 0, 0, time.UTC)
				yearEndPrice = 181.00
			)
			BeforeEach(func() {
				// Daily prices for the three trades in year
				// Q2: Strategy C - Let backfill logic handle opening position naturally
				aaplDailyPrices := map[string]float64{
					"2023-03-15": 150.00,
					"2023-06-01": 165.00, // Peak date (highest quantity after two buys)
					"2023-10-20": 170.00,
				}

				// Daily rates for INR calculation
				// Strategy: Increasing rates to emphasize Jun 01 peak with higher rate
				aaplDailyRates := map[string]float64{
					"2023-03-15": 82.0,
					"2023-06-01": 83.0, // Higher rate on peak date
					"2023-10-20": 82.5,
				}

				mockAccountManager.EXPECT().
					GetRecord(ctx, AAPL, testYear-1).
					Return(carryOverAccount, nil).Once()

				mockTickerManager.EXPECT().
					GetDailyPrices(ctx, AAPL, testYear).
					Return(aaplDailyPrices, nil).Once()

				mockSBIManager.EXPECT().
					GetDailyRates(ctx, testYear).
					Return(aaplDailyRates, nil).Once()

				mockTickerManager.EXPECT().
					GetPrice(ctx, AAPL, yearEndDate).
					Return(yearEndPrice, nil).Once()

				tradesInYear = []tax.Trade{
					tax.NewTrade(AAPL, "2023-03-15", "BUY", 20, 150.00),  // 50 (start) + 20 = 70 shares
					tax.NewTrade(AAPL, "2023-06-01", "BUY", 30, 165.00),  // 70 + 30 = 100 shares (New Peak)
					tax.NewTrade(AAPL, "2023-10-20", "SELL", 10, 170.00), // 100 - 10 = 90 shares
				}
			})

			It("should correctly set FirstPosition based on carry-over and track subsequent trades to Peak and YearEnd", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, tradesInYear, testYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(valuation.Ticker).To(Equal(AAPL))

				// 1. Assert FirstPosition (Opening balance for the period)
				expectedFirstPosDate := time.Date(testYear-1, 12, 31, 0, 0, 0, 0, time.UTC)
				Expect(valuation.FirstPosition.Date.Format(time.DateOnly)).To(Equal(expectedFirstPosDate.Format(time.DateOnly)))
				Expect(valuation.FirstPosition.Quantity).To(Equal(50.0))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(160.00))

				// 2. Assert PeakPosition
				// Start: 50. After Buy1 (20 shares @ $150 on Mar 15): 70 shares.
				// After Buy2 (30 shares @ $165 on Jun 01): 100 shares. This is the new peak quantity.
				expectedPeakPosDate, _ := time.Parse(time.DateOnly, "2023-06-01")
				Expect(valuation.PeakPosition.Date.Format(time.DateOnly)).To(Equal(expectedPeakPosDate.Format(time.DateOnly)))
				Expect(valuation.PeakPosition.Quantity).To(Equal(100.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(165.00))

				// 3. Assert YearEndPosition
				// After Sell1 (10 shares): 100 - 10 = 90 shares remaining.
				Expect(valuation.YearEndPosition.Date.Format(time.DateOnly)).To(Equal(yearEndDate.Format(time.DateOnly)))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(90.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})

		Context("With No Trades in Year", func() {
			var (
				testYear     = 2023
				yearEndDate  = time.Date(testYear, 12, 31, 0, 0, 0, 0, time.UTC)
				yearEndPrice = 185.00 // A different price for clarity
			)

			BeforeEach(func() {
				carryOverAccount.Quantity = 75
				carryOverAccount.MarketValue = 12000.00 // Implies $160.00 per share

				// Daily prices map is empty (no trades in year)
				// Q3: Yes, empty maps work - getClosestPrice returns 0, day skipped
				aaplDailyPrices := map[string]float64{}

				// Daily rates map is empty (no trades in year)
				aaplDailyRates := map[string]float64{}

				mockAccountManager.EXPECT().
					GetRecord(ctx, AAPL, testYear-1).
					Return(carryOverAccount, nil).Once()

				mockTickerManager.EXPECT().
					GetDailyPrices(ctx, AAPL, testYear).
					Return(aaplDailyPrices, nil).Once()

				mockSBIManager.EXPECT().
					GetDailyRates(ctx, testYear).
					Return(aaplDailyRates, nil).Once()

				mockTickerManager.EXPECT().
					GetPrice(ctx, AAPL, yearEndDate).
					Return(yearEndPrice, nil).Once()

				tradesInYear = []tax.Trade{} // Explicitly empty
			})

			It("should correctly set Valuations based on carry-over when no trades occur", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, tradesInYear, testYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(valuation.Ticker).To(Equal(AAPL))

				// 1. Assert FirstPosition (Opening balance for the period)
				expectedFirstPosDate := time.Date(testYear-1, 12, 31, 0, 0, 0, 0, time.UTC)
				Expect(valuation.FirstPosition.Date.Format(time.DateOnly)).To(Equal(expectedFirstPosDate.Format(time.DateOnly)))
				Expect(valuation.FirstPosition.Quantity).To(Equal(75.0))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(160.00)) // From carryOverAccount MarketValue/Quantity

				// 2. Assert PeakPosition (Should be the same as FirstPosition as no trades increased quantity)
				Expect(valuation.PeakPosition.Date.Format(time.DateOnly)).To(Equal(expectedFirstPosDate.Format(time.DateOnly)))
				Expect(valuation.PeakPosition.Quantity).To(Equal(75.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(160.00))

				// 3. Assert YearEndPosition (Quantity remains the same, price is year-end price)
				Expect(valuation.YearEndPosition.Date.Format(time.DateOnly)).To(Equal(yearEndDate.Format(time.DateOnly)))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(75.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(yearEndPrice))
			})
		})

		Context("With Only Sell Trades", func() {
			var (
				MSFT        = "MSFT"
				testYear    = 2024
				yearEndDate = time.Date(testYear, 12, 31, 0, 0, 0, 0, time.UTC)
			)

			BeforeEach(func() {
				carryOverAccount = tax.Account{
					Symbol:      MSFT,
					Quantity:    50,
					MarketValue: 10000.00, // Implies $200.00 per share
				}

				// Daily prices for the sell trade
				msftDailyPrices := map[string]float64{
					"2024-02-15": 210.00,
				}

				// Daily rates for peak calculation
				msftDailyRates := map[string]float64{
					"2024-02-15": 82.5,
				}

				mockAccountManager.EXPECT().
					GetRecord(ctx, MSFT, testYear-1).
					Return(carryOverAccount, nil).Once()

				mockTickerManager.EXPECT().
					GetDailyPrices(ctx, MSFT, testYear).
					Return(msftDailyPrices, nil).Once()

				mockSBIManager.EXPECT().
					GetDailyRates(ctx, testYear).
					Return(msftDailyRates, nil).Once()

				// NOTE: No GetPrice mock needed - position fully exits (quantity = 0 after sell)
				// Code skips GetPrice call when quantity = 0 (see determineYearEndPosition line 249)

				tradesInYear = []tax.Trade{
					tax.NewTrade(MSFT, "2024-02-15", "SELL", 50, 210.00),
				}
			})

			It("should correctly calculate valuation", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, MSFT, tradesInYear, testYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(valuation.Ticker).To(Equal(MSFT))

				// 1. Assert FirstPosition (Opening balance for the period)
				expectedFirstPosDate := time.Date(testYear-1, 12, 31, 0, 0, 0, 0, time.UTC)
				Expect(valuation.FirstPosition.Date.Format(time.DateOnly)).To(Equal(expectedFirstPosDate.Format(time.DateOnly)))
				Expect(valuation.FirstPosition.Quantity).To(Equal(50.0))
				Expect(valuation.FirstPosition.USDPrice).To(Equal(200.00))

				// 2. Assert PeakPosition (Should be the same as FirstPosition as no trades increased quantity)
				Expect(valuation.PeakPosition.Date.Format(time.DateOnly)).To(Equal(expectedFirstPosDate.Format(time.DateOnly)))
				Expect(valuation.PeakPosition.Quantity).To(Equal(50.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(200.00))

				// 3. Assert YearEndPosition (Quantity is zero after the sell)
				Expect(valuation.YearEndPosition.Date.Format(time.DateOnly)).To(Equal(yearEndDate.Format(time.DateOnly)))
				Expect(valuation.YearEndPosition.Quantity).To(Equal(0.0))
				Expect(valuation.YearEndPosition.USDPrice).To(Equal(0.0)) // Price is irrelevant for zero quantity
			})
		})

		Context("Peak INR Calculation with Sparse Prices", func() {
			// Scenario: AAPL 2023 - Backfill Price from Previous Year
			// Opening: 50 shares @ $160 (2022-12-31 carry-over)
			// Trades: Mar 15 BUY 20, Jul 10 BUY 30 (100 max), Oct 20 SELL 15
			// Sparse prices: Only Nov 10 ($175), Dec 31 ($181) - no early year prices!
			//
			// Key Positions (Qty | Backfilled Price | TT Rate | INR Value):
			// Jan 1-Mar 14:  50 qty | $160 (backfilled 2022-12-31) | 82.00 | ₹656,000
			// Mar 15-Jul 9:  70 qty | $160 (backfilled 2022-12-31) | 82.00-82.50 | ₹915,200-₹918,400
			// Jul 10:       100 qty | $160 (backfilled 2022-12-31) | 82.50 | ₹1,320,000
			// Jul 11-Aug 30: 100 qty | $160 (backfilled 2022-12-31) | 82.50 (backfilled) | ₹1,320,000
			// Aug 31:       100 qty | $160 (backfilled 2022-12-31) | 82.55 | ₹1,320,800 ← PEAK (highest INR value)
			// Sep 1-Oct 19: 100 qty | $160 (backfilled 2022-12-31) | 82.55 (backfilled) | ₹1,320,800
			// Oct 20-Nov 9:  85 qty | $160 (backfilled 2022-12-31) | 82.55 (backfilled) | ₹1,120,780
			// Nov 10:        85 qty | $175 (price exists)          | 82.55 (backfilled) | ₹1,224,436
			// Nov 15:        85 qty | $175 (backfilled Nov 10)     | 83.20 | ₹1,236,700
			// Dec 31:        85 qty | $181 (price exists)          | 82.00 | ₹1,260,170

			var (
				testYear     = 2023
				yearEndDate  = time.Date(testYear, 12, 31, 0, 0, 0, 0, time.UTC)
				yearEndPrice = 181.00
			)

			BeforeEach(func() {
				carryOverAccount.Quantity = 50
				carryOverAccount.MarketValue = 8000.00 // $160 per share

				tradesInYear = []tax.Trade{
					tax.NewTrade(AAPL, "2023-03-15", "BUY", 20, 150.00),
					tax.NewTrade(AAPL, "2023-07-10", "BUY", 30, 165.00),
					tax.NewTrade(AAPL, "2023-10-20", "SELL", 15, 170.00),
				}

				// Sparse prices (only Nov 10, Dec 31) + previous year-end for backfill
				aaplDailyPrices := map[string]float64{
					"2022-12-31": 160.00, // Previous year-end for backfill support
					"2023-11-10": 175.00,
					"2023-12-31": 181.00,
				}

				aaplDailyRates := map[string]float64{
					"2023-03-15": 82.00,
					"2023-07-10": 82.50,
					"2023-08-31": 82.55,
					"2023-11-15": 83.20,
					"2023-12-31": 82.00,
				}

				mockAccountManager.EXPECT().GetRecord(ctx, AAPL, testYear-1).Return(carryOverAccount, nil).Once()
				mockTickerManager.EXPECT().GetDailyPrices(ctx, AAPL, testYear).Return(aaplDailyPrices, nil).Once()
				mockSBIManager.EXPECT().GetDailyRates(ctx, testYear).Return(aaplDailyRates, nil).Once()
				mockTickerManager.EXPECT().GetPrice(ctx, AAPL, yearEndDate).Return(yearEndPrice, nil).Once()
			})

			It("should identify Aug 31 as peak based on INR calculation", func() {
				valuation, err := valuationManager.AnalyzeValuation(ctx, AAPL, tradesInYear, testYear)
				Expect(err).ToNot(HaveOccurred())

				Expect(valuation.PeakPosition.Date).To(Equal(time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC)))
				Expect(valuation.PeakPosition.Quantity).To(Equal(100.0))
				Expect(valuation.PeakPosition.USDPrice).To(Equal(160.00))
				Expect(valuation.PeakPosition.Quantity * valuation.PeakPosition.USDPrice * 82.55).To(Equal(1320800.0))
			})
		})
	})
})
