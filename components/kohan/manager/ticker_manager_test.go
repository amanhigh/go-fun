package manager_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients/mocks"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TickerManager", func() {
	var (
		mockClient    *mocks.StockDataClient
		tickerManager *manager.TickerManagerImpl
		testDir       string
		ctx           = context.Background()
		err           common.HttpError
	)

	BeforeEach(func() {
		mockClient = mocks.NewStockDataClient(GinkgoT())

		var err error
		testDir, err = os.MkdirTemp("", "ticker-test-*")
		Expect(err).NotTo(HaveOccurred())

		tickerManager = manager.NewTickerManager(mockClient, testDir)
	})

	AfterEach(func() {
		mockClient.AssertExpectations(GinkgoT())
		os.RemoveAll(testDir)
	})

	Context("DownloadTicker", func() {
		var (
			ticker    = "TEST"
			stockData tax.StockData
			filePath  string
		)

		BeforeEach(func() {
			stockData = tax.StockData{
				Prices: map[string]float64{
					"2024-01-23": 100.00,
				},
				Splits: []tax.YahooSplit{},
			}

			filePath = filepath.Join(testDir, "TEST.json")
		})

		It("should download and save ticker data successfully", func() {
			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Return(stockData, nil)

			// Test download
			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).ToNot(HaveOccurred())

			// Verify file exists and content
			fileContent, err := os.ReadFile(filePath)
			Expect(err).ToNot(HaveOccurred())

			var savedData tax.StockData
			err = json.Unmarshal(fileContent, &savedData)
			Expect(err).ToNot(HaveOccurred())
			Expect(savedData).To(Equal(stockData))
		})

		It("should skip download if file exists", func() {
			data, err := json.Marshal(stockData)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(filePath, data, util.APPEND_PERM)
			Expect(err).ToNot(HaveOccurred())

			// Call download without Mock Expectations
			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle API errors", func() {
			expectedErr := common.NewHttpError("API Error", http.StatusInternalServerError)

			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Return(tax.StockData{}, expectedErr)

			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Context("GetPrice", func() {
		var (
			ticker    = "TEST"
			stockData tax.StockData
		)

		BeforeEach(func() {
			stockData = tax.StockData{
				Prices: map[string]float64{
					"2024-01-15": 100.00,
					"2024-01-16": 101.00,
					"2024-01-17": 102.00,
				},
				Splits: []tax.YahooSplit{},
			}

			// Save test data
			filePath := filepath.Join(testDir, ticker+".json")
			data, _ := json.Marshal(stockData)
			err := os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return exact date price", func() {
			date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			price, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).ToNot(HaveOccurred())
			Expect(price).To(Equal(100.00))
		})

		It("should return closest previous date price", func() {
			// Request price for Jan 18 (not in data)
			date := time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC)
			price, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).ToNot(HaveOccurred())
			// Should return Jan 17 price
			Expect(price).To(Equal(102.00))
		})

		It("should use cache for subsequent requests", func() {
			date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			// First request - loads from file
			price1, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).ToNot(HaveOccurred())

			// Modify file to verify cache is used
			stockData.Prices["2024-01-15"] = 999.00
			data, _ := json.Marshal(stockData)
			filePath := filepath.Join(testDir, ticker+".json")
			writeErr := os.WriteFile(filePath, data, 0600)
			Expect(writeErr).ToNot(HaveOccurred())

			// Second request - should use cache
			price2, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).ToNot(HaveOccurred())
			Expect(price2).To(Equal(price1)) // Should return cached value
		})

		It("should handle missing data errors", func() {
			date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			_, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No price data found"))
			Expect(err.Code()).To(Equal(http.StatusNotFound))
		})
	})

	Context("GetDailyPrices", func() {
		var (
			ticker    = "AAPL"
			year      = 2023
			stockData tax.StockData
			filePath  string
		)

		BeforeEach(func() {
			// Setup test data with prices from multiple years
			stockData = tax.StockData{
				Prices: map[string]float64{
					// 2022 prices (should be filtered out)
					"2022-12-30": 150.00,
					"2022-12-31": 151.00,
					// 2023 prices (should be included)
					"2023-01-15": 152.00,
					"2023-02-14": 153.00,
					"2023-03-31": 154.00,
					"2023-06-30": 155.00,
					"2023-07-15": 156.00,
					"2023-12-31": 157.00,
					// 2024 prices (should be filtered out)
					"2024-01-15": 158.00,
					"2024-12-31": 159.00,
				},
				Splits: []tax.YahooSplit{},
			}

			// Save test data to file
			filePath = filepath.Join(testDir, ticker+".json")
			data, _ := json.Marshal(stockData)
			err := os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return all prices for requested year", func() {
			prices, err := tickerManager.GetDailyPrices(ctx, ticker, year)

			Expect(err).ToNot(HaveOccurred())
			Expect(prices).NotTo(BeNil())

			// Verify 2023 prices + previous year-end price are included (for backfill support)
			Expect(prices).To(HaveLen(7))

			// Verify specific dates and prices
			Expect(prices["2022-12-31"]).To(Equal(151.00)) // Previous year-end for backfill
			Expect(prices["2023-01-15"]).To(Equal(152.00))
			Expect(prices["2023-02-14"]).To(Equal(153.00))
			Expect(prices["2023-03-31"]).To(Equal(154.00))
			Expect(prices["2023-06-30"]).To(Equal(155.00))
			Expect(prices["2023-07-15"]).To(Equal(156.00))
			Expect(prices["2023-12-31"]).To(Equal(157.00))

			// Verify other years' prices are NOT included
			Expect(prices).NotTo(HaveKey("2022-12-30"))
			Expect(prices).NotTo(HaveKey("2024-01-15"))
			Expect(prices).NotTo(HaveKey("2024-12-31"))
		})

		It("should return error when no prices found for year", func() {
			prices, err := tickerManager.GetDailyPrices(ctx, ticker, 2025)

			Expect(err).To(HaveOccurred())
			Expect(prices).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("no price data found"))
			Expect(err.Error()).To(ContainSubstring("2025"))
			Expect(err.Code()).To(Equal(http.StatusNotFound))
		})

		It("should handle file not found error", func() {
			// Mock the client to return error when trying to fetch missing data
			mockClient.EXPECT().
				FetchDailyPrices(ctx, "NONEXISTENT").
				Return(tax.StockData{}, common.NewHttpError("File not found", http.StatusNotFound))

			prices, err := tickerManager.GetDailyPrices(ctx, "NONEXISTENT", year)

			Expect(err).To(HaveOccurred())
			Expect(prices).To(BeNil())
		})

		It("should return prices for single entry year", func() {
			singleYearData := tax.StockData{
				Prices: map[string]float64{
					"2023-06-15": 120.00,
				},
				Splits: []tax.YahooSplit{},
			}

			singleYearPath := filepath.Join(testDir, "SINGLE.json")
			data, _ := json.Marshal(singleYearData)
			err := os.WriteFile(singleYearPath, data, 0600)
			Expect(err).ToNot(HaveOccurred())

			prices, err := tickerManager.GetDailyPrices(ctx, "SINGLE", year)

			Expect(err).ToNot(HaveOccurred())
			Expect(prices).To(HaveLen(1))
			Expect(prices["2023-06-15"]).To(Equal(120.00))
		})

		It("should handle sparse dates correctly", func() {
			sparseData := tax.StockData{
				Prices: map[string]float64{
					"2023-01-15": 100.00,
					// Gap in dates
					"2023-06-30": 110.00,
					// Another gap
					"2023-12-31": 120.00,
				},
				Splits: []tax.YahooSplit{},
			}

			sparsePath := filepath.Join(testDir, "SPARSE.json")
			data, _ := json.Marshal(sparseData)
			err := os.WriteFile(sparsePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())

			prices, err := tickerManager.GetDailyPrices(ctx, "SPARSE", year)

			Expect(err).ToNot(HaveOccurred())
			Expect(prices).To(HaveLen(3))
			Expect(prices["2023-01-15"]).To(Equal(100.00))
			Expect(prices["2023-06-30"]).To(Equal(110.00))
			Expect(prices["2023-12-31"]).To(Equal(120.00))
		})

		It("should cache data for subsequent calls", func() {
			// First call - loads from file
			prices1, err1 := tickerManager.GetDailyPrices(ctx, ticker, year)
			Expect(err1).ToNot(HaveOccurred())

			// Modify file to verify cache is used
			modifiedData := tax.StockData{
				Prices: map[string]float64{
					"2023-01-15": 999.00, // Different value
				},
			}
			data, _ := json.Marshal(modifiedData)
			err := os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())

			// Second call - should use cache
			prices2, err2 := tickerManager.GetDailyPrices(ctx, ticker, year)
			Expect(err2).ToNot(HaveOccurred())

			// Should return cached values (original length of 7 including previous year-end)
			Expect(prices2).To(HaveLen(7))
			Expect(prices1).To(Equal(prices2))
		})

		It("should use consistent date format YYYY-MM-DD", func() {
			prices, err := tickerManager.GetDailyPrices(ctx, ticker, year)

			Expect(err).ToNot(HaveOccurred())

			// Verify all keys are in YYYY-MM-DD format
			for dateStr := range prices {
				parts := strings.Split(dateStr, "-")
				Expect(parts).To(HaveLen(3))
				// Year part should be 4 digits
				Expect(parts[0]).To(HaveLen(4))
				// Month and day should be 2 digits
				Expect(parts[1]).To(HaveLen(2))
				Expect(parts[2]).To(HaveLen(2))
			}
		})

		It("should handle year boundary correctly", func() {
			boundaryData := tax.StockData{
				Prices: map[string]float64{
					"2022-12-31": 100.00,
					"2023-01-01": 101.00,
					"2023-01-02": 102.00,
					"2023-12-30": 110.00,
					"2023-12-31": 111.00,
					"2024-01-01": 112.00,
				},
				Splits: []tax.YahooSplit{},
			}

			boundaryPath := filepath.Join(testDir, "BOUNDARY.json")
			data, _ := json.Marshal(boundaryData)
			err := os.WriteFile(boundaryPath, data, 0600)
			Expect(err).ToNot(HaveOccurred())

			prices, err := tickerManager.GetDailyPrices(ctx, "BOUNDARY", year)

			Expect(err).ToNot(HaveOccurred())
			// Should include 2023 dates + previous year-end for backfill support
			Expect(prices).To(HaveLen(5))
			Expect(prices).To(HaveKey("2022-12-31")) // Previous year-end for backfill
			Expect(prices).To(HaveKey("2023-01-01"))
			Expect(prices).To(HaveKey("2023-01-02"))
			Expect(prices).To(HaveKey("2023-12-30"))
			Expect(prices).To(HaveKey("2023-12-31"))
			// Should NOT include 2024 prices
			Expect(prices).NotTo(HaveKey("2024-01-01"))
		})

		It("should return different data for different years", func() {
			data2023, err2023 := tickerManager.GetDailyPrices(ctx, ticker, 2023)
			Expect(err2023).ToNot(HaveOccurred())

			data2022, err2022 := tickerManager.GetDailyPrices(ctx, ticker, 2022)
			Expect(err2022).ToNot(HaveOccurred())

			// 2023 should have 7 entries (6 from 2023 + 2022-12-31 for backfill)
			Expect(data2023).To(HaveLen(7))
			// 2022 should have 2 entries (2022-12-30, 2022-12-31; no 2021-12-31 in test data)
			Expect(data2022).To(HaveLen(2))

			// Data should be different
			Expect(data2023).NotTo(Equal(data2022))
		})
	})

	Context("GetSplits", func() {
		var (
			ticker    = "SPLIT"
			stockData tax.StockData
			from, to  time.Time
			splits    []tax.YahooSplit
			getErr    common.HttpError
		)

		BeforeEach(func() {
			// Base test data with splits on 2024-01-15 (2:1), 2024-03-01 (3:1), 2024-05-01 (1:4 reverse)
			stockData = tax.StockData{
				Prices: map[string]float64{
					"2024-01-10": 100.00,
				},
				Splits: []tax.YahooSplit{
					{Date: 1705276800, Numerator: 2, Denominator: 1}, // 2024-01-15
					{Date: 1709251200, Numerator: 3, Denominator: 1}, // 2024-03-01
					{Date: 1714521600, Numerator: 1, Denominator: 4}, // 2024-05-01 (reverse)
				},
			}
			filePath := filepath.Join(testDir, ticker+".json")
			data, err := json.Marshal(stockData)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			splits, getErr = tickerManager.GetSplits(ctx, ticker, from, to)
		})

		Context("with overlapping date range", func() {
			BeforeEach(func() {
				from = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
			})

			It("should return splits within date range inclusively", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(1))
				Expect(splits[0].Date).To(Equal(int64(1709251200))) // 2024-03-01
				Expect(splits[0].Numerator).To(Equal(3.0))
				Expect(splits[0].Denominator).To(Equal(1.0))
			})
		})

		Context("with no splits in stock data", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{"2024-01-10": 100.00},
					Splits: []tax.YahooSplit{},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				from = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
			})

			It("should return empty non-nil slice when no splits exist", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).ToNot(BeNil())
				Expect(splits).To(BeEmpty())
			})
		})
		Context("with inverted date range", func() {
			BeforeEach(func() {
				from = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			})

			It("should reject inverted date range", func() {
				Expect(getErr).To(HaveOccurred())
				Expect(getErr.Code()).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with malformed split ratio", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{"2024-01-10": 100.00},
					Splits: []tax.YahooSplit{
						{Date: 1705276800, Numerator: 0, Denominator: 1}, // zero numerator
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				from = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
			})

			It("should reject malformed split ratios", func() {
				Expect(getErr).To(HaveOccurred())
				Expect(getErr.Code()).To(Equal(http.StatusBadRequest))
				Expect(getErr.Error()).To(ContainSubstring(ticker))
			})
		})

		Context("defensive copy behavior", func() {
			var (
				secondSplits []tax.YahooSplit
				secondErr    common.HttpError
			)

			BeforeEach(func() {
				from = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
			})

			JustBeforeEach(func() {
				// Second call, runs after parent JustBeforeEach (first call)
				secondSplits, secondErr = tickerManager.GetSplits(ctx, ticker, from, to)

				// Mutate first-call result to test defensive copy independence
				if len(splits) > 0 {
					splits[0] = tax.YahooSplit{}
				}
			})

			It("should return a defensive copy", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(secondErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(3))
				Expect(secondSplits).To(HaveLen(3))

				// First result is zero-valued after mutation
				Expect(splits[0].Numerator).To(BeZero())
				Expect(splits[0].Denominator).To(BeZero())

				// Second result retains original values (independent defensive copy)
				Expect(secondSplits[0].Numerator).To(Equal(2.0))
				Expect(secondSplits[0].Denominator).To(Equal(1.0))

				// They now differ
				Expect(splits[0]).ToNot(Equal(secondSplits[0]))
			})
		})

		Context("with from boundary split", func() {
			BeforeEach(func() {
				from = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Same as first split
				to = time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
			})

			It("should include split exactly on from boundary", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(1))
				Expect(splits[0].Date).To(Equal(int64(1705276800))) // 2024-01-15
			})
		})

		Context("with to boundary split", func() {
			BeforeEach(func() {
				from = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC) // Same as second split
			})

			It("should include split exactly on to boundary", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(1))
				Expect(splits[0].Date).To(Equal(int64(1709251200))) // 2024-03-01
			})
		})

		Context("with unsorted splits", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{"2024-01-10": 100.00},
					Splits: []tax.YahooSplit{
						{Date: 1714521600, Numerator: 1, Denominator: 4}, // 2024-05-01 (reverse)
						{Date: 1705276800, Numerator: 2, Denominator: 1}, // 2024-01-15
						{Date: 1709251200, Numerator: 3, Denominator: 1}, // 2024-03-01
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				from = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
			})

			It("should return splits in chronological order", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(3))
				Expect(splits[0].Date).To(Equal(int64(1705276800))) // 2024-01-15
				Expect(splits[1].Date).To(Equal(int64(1709251200))) // 2024-03-01
				Expect(splits[2].Date).To(Equal(int64(1714521600))) // 2024-05-01
			})
		})

		Context("with intraday split timestamp on boundary", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{"2024-01-10": 100.00},
					Splits: []tax.YahooSplit{
						{Date: 1705276800, Numerator: 2, Denominator: 1}, // 2024-01-15 00:00 UTC
						{Date: 1709251300, Numerator: 3, Denominator: 1}, // 2024-03-01 00:01:40 UTC
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				from = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
				to = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
			})

			It("should include intraday split matching on calendar date boundary", func() {
				Expect(getErr).ToNot(HaveOccurred())
				Expect(splits).To(HaveLen(1))
				Expect(splits[0].Date).To(Equal(int64(1709251300)))
				Expect(splits[0].Numerator).To(Equal(3.0))
				Expect(splits[0].Denominator).To(Equal(1.0))
			})
		})
	})

	Context("GetPrice with split adjustments", func() {
		var (
			ticker    = "SPLIT_ADJ"
			stockData tax.StockData
			queryDate time.Time
			price     float64
			priceErr  common.HttpError
		)

		BeforeEach(func() {
			// StockData with a 2:1 split on 2024-03-01 and prices before/on the split
			stockData = tax.StockData{
				Prices: map[string]float64{
					"2024-01-10": 100.00, // Pre-split price
					"2024-03-01": 50.00,  // On split date (post-split trading price)
				},
				Splits: []tax.YahooSplit{
					{Date: 1709251200, Numerator: 2, Denominator: 1}, // 2024-03-01 2:1 split
				},
			}
			filePath := filepath.Join(testDir, ticker+".json")
			data, err := json.Marshal(stockData)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("querying pre-split date", func() {
			BeforeEach(func() {
				queryDate = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
			})

			JustBeforeEach(func() {
				price, priceErr = tickerManager.GetPrice(ctx, ticker, queryDate)
			})

			It("should reconstruct pre-split price by future split factor", func() {
				Expect(priceErr).ToNot(HaveOccurred())
				// Pre-split cached price 100 * 2/1 (future split) = 200
				Expect(price).To(Equal(200.00))
			})
		})

		Context("querying on split date", func() {
			BeforeEach(func() {
				queryDate = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
			})

			JustBeforeEach(func() {
				price, priceErr = tickerManager.GetPrice(ctx, ticker, queryDate)
			})

			It("should return unchanged price on the split date", func() {
				Expect(priceErr).ToNot(HaveOccurred())
				// On split date, no adjustment: cached 50.00 stays 50.00
				Expect(price).To(Equal(50.00))
			})
		})

		Context("with malformed split data", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{"2024-01-10": 100.00},
					Splits: []tax.YahooSplit{
						{Date: 1709251200, Numerator: 0, Denominator: 1},
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				queryDate = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
			})

			JustBeforeEach(func() {
				price, priceErr = tickerManager.GetPrice(ctx, ticker, queryDate)
			})

			It("should fail with BadRequest on malformed split data", func() {
				Expect(priceErr).To(HaveOccurred())
				Expect(priceErr.Code()).To(Equal(http.StatusBadRequest))
				Expect(priceErr.Error()).To(ContainSubstring(ticker))
			})
		})

		Context("querying with intraday split timestamp", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{
						"2024-01-10": 100.00,
						"2024-03-01": 80.00,
					},
					Splits: []tax.YahooSplit{
						{Date: 1709251300, Numerator: 2, Denominator: 1}, // 2024-03-01 00:01:40 UTC
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
				queryDate = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
			})

			JustBeforeEach(func() {
				price, priceErr = tickerManager.GetPrice(ctx, ticker, queryDate)
			})

			It("should not adjust price when split has intraday timestamp on same calendar date", func() {
				Expect(priceErr).ToNot(HaveOccurred())
				Expect(price).To(Equal(80.00))
			})
		})
	})

	Context("GetDailyPrices with split adjustments", func() {
		var (
			ticker    = "SPLIT_ADJ_DLY"
			stockData tax.StockData
			year      = 2024
			prices    map[string]float64
			pricesErr common.HttpError
		)

		BeforeEach(func() {
			stockData = tax.StockData{
				Prices: map[string]float64{
					"2024-01-10": 100.00, // Pre-split price
					"2024-03-01": 50.00,  // On split date
				},
				Splits: []tax.YahooSplit{
					{Date: 1709251200, Numerator: 2, Denominator: 1}, // 2024-03-01 2:1 split
				},
			}
			filePath := filepath.Join(testDir, ticker+".json")
			data, err := json.Marshal(stockData)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(filePath, data, 0600)
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			prices, pricesErr = tickerManager.GetDailyPrices(ctx, ticker, year)
		})

		It("should adjust pre-split prices by future split factors", func() {
			Expect(pricesErr).ToNot(HaveOccurred())
			Expect(prices).To(HaveLen(2))
			// Pre-split price should be multiplied by cumulative future split factor
			Expect(prices["2024-01-10"]).To(Equal(200.00)) // 100 * 2/1
		})

		It("should not adjust prices on the split date", func() {
			Expect(pricesErr).ToNot(HaveOccurred())
			Expect(prices["2024-03-01"]).To(Equal(50.00))
		})

		Context("with missing price on split date", func() {
			BeforeEach(func() {
				stockData = tax.StockData{
					Prices: map[string]float64{
						"2024-01-10": 100.00, // Pre-split price only; no exact price on split date
					},
					Splits: []tax.YahooSplit{
						{Date: 1709251200, Numerator: 2, Denominator: 1}, // 2024-03-01 2:1 split
					},
				}
				filePath := filepath.Join(testDir, ticker+".json")
				data, err := json.Marshal(stockData)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(filePath, data, 0600)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return error when no exact cached price exists on split date", func() {
				Expect(pricesErr).To(HaveOccurred())
				Expect(pricesErr.Code()).To(Equal(http.StatusNotFound))
				Expect(pricesErr.Error()).To(ContainSubstring(ticker))
				Expect(pricesErr.Error()).To(ContainSubstring("2024-03-01"))
			})
		})
	})
})
