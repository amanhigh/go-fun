package manager_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	clientMocks "github.com/amanhigh/go-fun/components/kohan/clients/mocks"
	manager "github.com/amanhigh/go-fun/components/kohan/manager"
	repoMocks "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testCSV = `DATE,TT BUY,TT SELL
2024-01-23 Wednesday,82.50,83.50
2024-01-24 Thursday,82.75,83.75
`

var _ = Describe("SBIManager", func() {
	var (
		mockClient   *clientMocks.SBIClient
		mockExchange *repoMocks.ExchangeRepository
		sbiManager   *manager.SBIManagerImpl
		testDir      string
		ctx          = context.Background()
	)

	BeforeEach(func() {
		mockClient = clientMocks.NewSBIClient(GinkgoT())
		mockExchange = repoMocks.NewExchangeRepository(GinkgoT())

		var err error
		testDir, err = os.MkdirTemp("", "sbi-test-*")
		Expect(err).NotTo(HaveOccurred())

		filePath := filepath.Join(testDir, tax.SBI_RATES_FILENAME)
		sbiManager = manager.NewSBIManager(mockClient, filePath, mockExchange)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("DownloadRates", func() {
		It("should download and save rates successfully", func() {
			mockClient.EXPECT().FetchExchangeRates(ctx).Return(testCSV, nil)

			err := sbiManager.DownloadRates(ctx)
			Expect(err).To(BeNil())

			// Verify file content
			content, readErr := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})

		It("should skip download if file already exists", func() {
			// Create test file first
			err := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			// Expect no client calls since file exists
			// mockClient.EXPECT().FetchExchangeRates(ctx) should not be called

			err = sbiManager.DownloadRates(ctx)
			Expect(err).To(BeNil())

			// Verify file content remains unchanged
			content, readErr := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})

		It("should handle client error", func() {
			expectedErr := common.NewHttpError("Failed to fetch exchange rates", http.StatusInternalServerError)
			mockClient.EXPECT().FetchExchangeRates(ctx).Return("", expectedErr)

			err := sbiManager.DownloadRates(ctx)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Context("GetTTBuyRate", func() {
		It("should get rate for date using repository", func() {
			testDate := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
			expectedRate := 82.50

			mockExchange.EXPECT().GetAllRecords(ctx).Return([]tax.SbiRate{
				{Date: "2024-01-23 Wednesday", TTBuy: expectedRate, TTSell: 83.50},
				{Date: "2024-01-24 Thursday", TTBuy: 82.75, TTSell: 83.75},
			}, nil)

			rate, err := sbiManager.GetTTBuyRate(ctx, testDate)
			Expect(err).To(BeNil())
			Expect(rate).To(Equal(expectedRate))
		})
	})

	Context("When exact date not found", func() {
		var (
			requestedDate = time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC)
			closestDate   = time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
			expectedRate  = 82.50
		)

		BeforeEach(func() {
			mockExchange.EXPECT().
				GetAllRecords(ctx).
				Return([]tax.SbiRate{
					{Date: "2024-01-23 Wednesday", TTBuy: expectedRate, TTSell: 83.50},
					{Date: "2024-01-22 Tuesday", TTBuy: 82.25, TTSell: 83.25},
				}, nil)
		})

		It("should return closest previous date with ClosestDateError", func() {
			rate, err := sbiManager.GetTTBuyRate(ctx, requestedDate)
			Expect(rate).To(Equal(expectedRate))

			// Verify ClosestDateError details
			closestErr, ok := err.(tax.ClosestDateError)
			Expect(ok).To(BeTrue())
			Expect(closestErr.Code()).To(Equal(http.StatusOK))
			Expect(closestErr.GetRequestedDate()).To(Equal(requestedDate))
			Expect(closestErr.GetClosestDate()).To(Equal(closestDate))
		})
	})
})
