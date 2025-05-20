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

var _ = Describe("ValuationManager", func() {
	var (
		ctx                 = context.Background()
		mockTickerManager   *mocks.TickerManager
		mockAccountManager  *mocks.AccountManager
		mockTradeRepository *repoMocks.TradeRepository
		mockFyManager       *mocks.FinancialYearManager[tax.Trade] // Added mock for FY Manager
		valuationManager    manager.ValuationManager

		// Common variables
		AAPL         = "AAPL"
		year         = 2024
		yearEndDate  = time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		yearEndPrice = 150.00
	)

	BeforeEach(func() {
		mockTickerManager = mocks.NewTickerManager(GinkgoT())
		mockAccountManager = mocks.NewAccountManager(GinkgoT())
		mockTradeRepository = repoMocks.NewTradeRepository(GinkgoT())
		mockFyManager = mocks.NewFinancialYearManager[tax.Trade](GinkgoT())                                                       // Initialize mock FY Manager
		valuationManager = manager.NewValuationManager(mockTickerManager, mockAccountManager, mockTradeRepository, mockFyManager) // Pass mockFyManager
	})

	Context("Analyse Valuation", func() {

		Context("Fresh Start", func() {
			BeforeEach(func() {
				// All tests under Fresh Start expect no last year position
				mockAccountManager.EXPECT().
					GetRecord(ctx, AAPL).
					Return(tax.Account{}, common.ErrNotFound)
			})

			Context("Basic Position Tracking", func() {
				Context("Single Buy and Hold", func() {
					var trades []tax.Trade

					BeforeEach(func() {
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
						}

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should compute correct positions", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
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
					})

					It("should compute positions with zero year-end", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
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

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track increasing position correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
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

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track averaged position correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
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

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should handle year end trades correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position from initial buy
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
						Expect(valuation.FirstPosition.Quantity).To(Equal(5.0))
						Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))

						// Peak position should be final position
						date, getDateErr = trades[1].GetDate() // Reuse getDateErr
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.PeakPosition.Date).To(Equal(date))
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
						// HACK: #C Multiple Peaks with Same Value (Take Second higher TBBR Rate) or Throw Error.
						trades = []tax.Trade{
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),  // Initial 10
							tax.NewTrade(AAPL, "2024-02-15", "BUY", 5, 110),   // Peak 1: 15 shares
							tax.NewTrade(AAPL, "2024-03-15", "SELL", 8, 120),  // Down to 7
							tax.NewTrade(AAPL, "2024-04-15", "BUY", 10, 115),  // Peak 2: 17 shares
							tax.NewTrade(AAPL, "2024-05-15", "SELL", 12, 125), // Down to 5
						}

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should identify highest peak through multiple cycles", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
						Expect(err).ToNot(HaveOccurred())

						// First position
						date, getDateErr := trades[0].GetDate()
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.FirstPosition.Date).To(Equal(date))
						Expect(valuation.FirstPosition.Quantity).To(Equal(10.0))
						Expect(valuation.FirstPosition.USDPrice).To(Equal(100.0))

						// Peak should be second peak with 17 shares
						date, getDateErr = trades[3].GetDate() // Reuse getDateErr
						Expect(getDateErr).NotTo(HaveOccurred())
						Expect(valuation.PeakPosition.Date).To(Equal(date))
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
							tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100), // Initial 10 shares
							tax.NewTrade(AAPL, "2024-02-15", "SELL", 3, 110), // Sell 3 shares
							tax.NewTrade(AAPL, "2024-03-15", "SELL", 4, 120), // Sell 4 shares
						}

						mockTickerManager.EXPECT().
							GetPrice(ctx, AAPL, yearEndDate).
							Return(yearEndPrice, nil)
					})

					It("should track partial sells correctly", func() {
						valuation, err := valuationManager.AnalyzeValuation(ctx, trades, year)
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

		Context("Error Cases", func() {
			Context("Empty Trades", func() {
				It("should return error for empty trades", func() {
					trades := []tax.Trade{}
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("no trades provided"))
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
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("multiple tickers found"))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Year End Price Error", func() {
				var trades []tax.Trade

				BeforeEach(func() {
					trades = []tax.Trade{
						tax.NewTrade(AAPL, "2024-01-15", "BUY", 10, 100),
					}

					mockTickerManager.EXPECT().
						GetPrice(ctx, AAPL, yearEndDate).
						Return(0.0, common.ErrNotFound)

					mockAccountManager.EXPECT().
						GetRecord(ctx, AAPL).
						Return(tax.Account{}, common.ErrNotFound)
				})

				It("should handle ticker price fetch error", func() {
					_, err := valuationManager.AnalyzeValuation(ctx, trades, year)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to get year end price"))
					Expect(err.Code()).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Context("GetYearlyValuationsUSD", func() {
		var (
			// Define sample trades for multiple tickers
			tradeAAPL1 = tax.NewTrade(AAPL, "2024-01-10", "BUY", 10, 100)   // Date: 2024-01-10
			tradeAAPL2 = tax.NewTrade(AAPL, "2024-05-15", "SELL", 5, 120)   // Date: 2024-05-15
			tradeMSFT1 = tax.NewTrade("MSFT", "2024-02-20", "BUY", 20, 200) // Date: 2024-02-20
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

			// Mock dependencies needed by AnalyzeValuation for AAPL
			mockAccountManager.EXPECT().GetRecord(ctx, AAPL).Return(tax.Account{}, common.ErrNotFound).Once() // Fresh start for AAPL
			mockTickerManager.EXPECT().GetPrice(ctx, AAPL, yearEndDate).Return(150.0, nil).Once()             // AAPL year end price

			// Mock dependencies needed by AnalyzeValuation for MSFT
			mockAccountManager.EXPECT().GetRecord(ctx, "MSFT").Return(tax.Account{Quantity: 10, MarketValue: 1800}, nil).Once() // Start MSFT with 10 shares @ 180
			mockTickerManager.EXPECT().GetPrice(ctx, "MSFT", yearEndDate).Return(210.0, nil).Once()                             // MSFT year end price

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
			Expect(msftVal.FirstPosition.Quantity).To(Equal(10.0))                                    // From starting position
			Expect(msftVal.FirstPosition.Date).To(Equal(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC))) // Date of start pos (Jan 1st of analysis year)
			Expect(msftVal.PeakPosition.Quantity).To(Equal(30.0))                                     // 10 start + 20 buy
			msftPeakDate, getDateErr := tradeMSFT1.GetDate()                                          // Use getDateErr here too
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

			valuations, err := valuationManager.GetYearlyValuationsUSD(ctx, year)

			// Assert that the specific StatusNotFound error is returned
			Expect(err).To(HaveOccurred())
			Expect(err.Code()).To(Equal(http.StatusNotFound))
			Expect(err.Error()).To(ContainSubstring("no trades found for year"))
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

		It("should fail fast if AnalyzeValuation returns error for one ticker", func() {
			// Mock Repo: GetAllRecords returns trades for two tickers
			mockTradeRepository.EXPECT().GetAllRecords(ctx).Return(allTrades, nil).Once()

			// Mock FY Manager: FilterUS returns the same trades
			mockFyManager.EXPECT().FilterUS(ctx, allTrades, year).Return(allTrades, nil).Once()

			// Mock dependencies for AAPL (enough to trigger AnalyzeValuation)
			mockAccountManager.EXPECT().GetRecord(ctx, AAPL).Return(tax.Account{}, common.ErrNotFound).Once()
			// Make TickerManager fail for AAPL
			expectedErr := common.NewServerError(errors.New("price fetch failed"))
			mockTickerManager.EXPECT().GetPrice(ctx, AAPL, yearEndDate).Return(0.0, expectedErr).Once()

			// We don't expect mocks for MSFT because it should fail fast on AAPL

			_, err := valuationManager.GetYearlyValuationsUSD(ctx, year)
			// Assertions
			Expect(err).To(HaveOccurred())
			// Check if the error is the one returned by GetPrice (or wrapped by AnalyzeValuation)
			Expect(err.Error()).To(ContainSubstring("failed to get year end price"))
			Expect(err.Code()).To(Equal(http.StatusInternalServerError)) // As wrapped by AnalyzeValuation
		})
	})

	Context("With Carry-Over Position", func() {
		var (
			carryOverAccount tax.Account
			tradesInYear     []tax.Trade
			testYear         = 2023
			yearEndDate      = time.Date(testYear, 12, 31, 0, 0, 0, 0, time.UTC)
			yearEndPrice     = 181.00
		)

		BeforeEach(func() {
			carryOverAccount = tax.Account{
				Symbol:      AAPL,
				Quantity:    50,
				MarketValue: 8000.00, // Implies $160.00 per share
			}

			mockAccountManager.EXPECT().
				GetRecord(ctx, AAPL).
				Return(carryOverAccount, nil).Once()

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
			valuation, err := valuationManager.AnalyzeValuation(ctx, tradesInYear, testYear)

			Expect(err).ToNot(HaveOccurred())
			Expect(valuation.Ticker).To(Equal(AAPL))

			// 1. Assert FirstPosition (Opening balance for the period)
			expectedFirstPosDate := time.Date(testYear, 1, 1, 0, 0, 0, 0, time.UTC)
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

	// FIXME: Carry Over with No Trades in Year.
})
