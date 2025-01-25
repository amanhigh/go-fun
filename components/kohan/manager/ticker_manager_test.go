package manager_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/clients/mocks"
	manager "github.com/amanhigh/go-fun/components/kohan/manager"
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

	// Initial basic test case
	Context("DownloadTicker", func() {
		It("should download and save ticker data successfully", func() {
			// Mock return data
			stockData := tax.StockData{
				MetaData: tax.MetaData{Symbol: "TEST"},
			}
			mockClient.EXPECT().FetchDailyPrices(ctx, stockData.MetaData.Symbol).Return(stockData, nil)

			// Test download
			err = tickerManager.DownloadTicker(ctx, stockData.MetaData.Symbol)
			Expect(err).To(BeNil())

			// Verify file exists
			_, err := os.Stat(filepath.Join(testDir, "TEST.json"))
			Expect(err).To(BeNil())
		})
	})

	// FIXME: #B Add Test for Remaining Functions.
})
