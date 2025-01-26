package manager_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

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
		os.RemoveAll(testDir)
	})

	Context("DownloadTicker", func() {
		It("should download and save ticker data successfully", func() {
			// Mock return data
			stockData := tax.VantageStockData{
				MetaData: tax.MetaData{Symbol: "TEST"},
				TimeSeries: map[string]tax.DayPrice{
					"2024-01-23": {Close: "100.00"},
				},
			}
			mockClient.EXPECT().
				FetchDailyPrices(ctx, stockData.MetaData.Symbol).
				Return(stockData, nil)

			// Test download
			err = tickerManager.DownloadTicker(ctx, stockData.MetaData.Symbol)
			Expect(err).To(BeNil())

			// Verify file exists and content
			filePath := filepath.Join(testDir, "TEST.json")
			fileContent, err := os.ReadFile(filePath)
			Expect(err).To(BeNil())

			var savedData tax.VantageStockData
			err = json.Unmarshal(fileContent, &savedData)
			Expect(err).To(BeNil())
			Expect(savedData).To(Equal(stockData))
		})

		It("should skip download if file exists", func() {
			ticker := "TEST"
			filePath := filepath.Join(testDir, ticker+".json")

			// Create file
			err := os.WriteFile(filePath, []byte("{}"), 0644)
			Expect(err).To(BeNil())

			// Mock should not be called
			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Times(0)

			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(BeNil())
		})

		It("should handle API errors", func() {
			ticker := "TEST"
			expectedErr := common.NewHttpError("API Error", http.StatusInternalServerError)

			mockClient.EXPECT().
				FetchDailyPrices(ctx, ticker).
				Return(tax.VantageStockData{}, expectedErr)

			err = tickerManager.DownloadTicker(ctx, ticker)
			Expect(err).To(Equal(expectedErr))
		})
	})
})
