package manager_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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
		mockClient    *mocks.AlphaClient
		tickerManager *manager.TickerManagerImpl
		testDir       string
		ctx           = context.Background()
		err           common.HttpError
	)

	BeforeEach(func() {
		mockClient = mocks.NewAlphaClient(GinkgoT())

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
			stockData tax.VantageStockData
			filePath  string
		)

		BeforeEach(func() {
			stockData = tax.VantageStockData{
				MetaData: tax.MetaData{Symbol: ticker},
				TimeSeries: map[string]tax.DayPrice{
					"2024-01-23": {Close: "100.00"},
				},
			}

			filePath = filepath.Join(testDir, "TEST.json")
		})

		It("should download and save ticker data successfully", func() {
			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Return(stockData, nil)

			// Test download
			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(BeNil())

			// Verify file exists and content
			fileContent, err := os.ReadFile(filePath)
			Expect(err).To(BeNil())

			var savedData tax.VantageStockData
			err = json.Unmarshal(fileContent, &savedData)
			Expect(err).To(BeNil())
			Expect(savedData).To(Equal(stockData))
		})

		It("should skip download if file exists", func() {
			data, err := json.Marshal(stockData)
			Expect(err).To(BeNil())
			err = os.WriteFile(filePath, data, util.APPEND_PERM)
			Expect(err).To(BeNil())

			// Call download without Mock Expectations
			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(BeNil())
		})

		It("should handle API errors", func() {
			expectedErr := common.NewHttpError("API Error", http.StatusInternalServerError)

			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Return(tax.VantageStockData{}, expectedErr)

			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Context("ValueTicker", func() {
		var (
			ticker    = "TEST"
			year      = 2024
			stockData tax.VantageStockData
		)

		BeforeEach(func() {
			// Create test data with known peak and year-end prices
			stockData = tax.VantageStockData{
				MetaData: tax.MetaData{Symbol: ticker},
				TimeSeries: map[string]tax.DayPrice{
					"2024-03-15": {Close: "150.00"}, // Peak price
					"2024-02-01": {Close: "120.00"},
					"2024-12-31": {Close: "130.00"}, // Year end
				},
			}

			// Save test data to file
			filePath := filepath.Join(testDir, ticker+".json")
			data, _ := json.Marshal(stockData)
			err := os.WriteFile(filePath, data, 0644)
			Expect(err).To(BeNil())
		})

		It("should correctly analyze yearly data", func() {
			valuation, err := tickerManager.ValueTicker(ctx, ticker, year)
			Expect(err).To(BeNil())

			Expect(valuation.Ticker).To(Equal(ticker))
			Expect(valuation.PeakDate).To(Equal("2024-03-15"))
			Expect(valuation.PeakPrice).To(Equal(150.00))
			Expect(valuation.YearEndDate).To(Equal("2024-12-31"))
			Expect(valuation.YearEndPrice).To(Equal(130.00))
		})

		It("should handle missing data", func() {
			// Test with non-existent ticker
			_, err := tickerManager.ValueTicker(ctx, "INVALID", year)
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("GetPrice", func() {
		var (
			ticker    = "TEST"
			stockData tax.VantageStockData
		)

		BeforeEach(func() {
			stockData = tax.VantageStockData{
				MetaData: tax.MetaData{Symbol: ticker},
				TimeSeries: map[string]tax.DayPrice{
					"2024-01-15": {Close: "100.00"},
					"2024-01-16": {Close: "101.00"},
					"2024-01-17": {Close: "102.00"},
				},
			}

			// Save test data
			filePath := filepath.Join(testDir, ticker+".json")
			data, _ := json.Marshal(stockData)
			err := os.WriteFile(filePath, data, 0644)
			Expect(err).To(BeNil())
		})

		It("should return exact date price", func() {
			date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			price, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(BeNil())
			Expect(price).To(Equal(100.00))
		})

		It("should return closest previous date price", func() {
			// Request price for Jan 18 (not in data)
			date := time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC)
			price, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(BeNil())
			// Should return Jan 17 price
			Expect(price).To(Equal(102.00))
		})

		It("should use cache for subsequent requests", func() {
			date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			// First request - loads from file
			price1, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(BeNil())

			// Modify file to verify cache is used
			stockData.TimeSeries["2024-01-15"] = tax.DayPrice{Close: "999.00"}
			data, _ := json.Marshal(stockData)
			filePath := filepath.Join(testDir, ticker+".json")
			writeErr := os.WriteFile(filePath, data, 0644)
			Expect(writeErr).To(BeNil())

			// Second request - should use cache
			price2, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(BeNil())
			Expect(price2).To(Equal(price1)) // Should return cached value
		})

		It("should handle missing data errors", func() {
			date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			_, err := tickerManager.GetPrice(ctx, ticker, date)
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(ContainSubstring("No price data found"))
			Expect(err.Code()).To(Equal(http.StatusNotFound))
		})
	})
})
