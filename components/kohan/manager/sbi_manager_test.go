package manager_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients/mocks"
	manager "github.com/amanhigh/go-fun/components/kohan/manager"
	repoMock "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
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
		mockClient   *mocks.SBIClient
		mockExchange *repoMock.ExchangeRepository
		sbiManager   *manager.SBIManagerImpl
		testDir      string
		ctx          = context.Background()
	)

	BeforeEach(func() {
		mockClient = mocks.NewSBIClient(GinkgoT())
		mockExchange = repoMock.NewExchangeRepository(GinkgoT())

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
			content, err := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(err).To(BeNil())
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
			content, err := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})

		It("should handle client error", func() {
			expectedErr := common.NewHttpError("Failed to fetch exchange rates", http.StatusInternalServerError)
			mockClient.EXPECT().FetchExchangeRates(ctx).Return("", expectedErr)

			err := sbiManager.DownloadRates(ctx)
			Expect(err).To(Equal(expectedErr))
		})

		It("should create directory if not exists", func() {
			// Create nested test directory
			nestedDir := filepath.Join(testDir, "nested", "dir")
			sbiManager = manager.NewSBIManager(mockClient, nestedDir, mockExchange)

			mockClient.EXPECT().FetchExchangeRates(ctx).Return(testCSV, nil)

			err := sbiManager.DownloadRates(ctx)
			Expect(err).To(BeNil())

			// Verify directory was created with file
			content, err := os.ReadFile(filepath.Join(nestedDir, tax.SBI_RATES_FILENAME))
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})
	})

	Context("GetTTBuyRate", func() {
		It("should get rate for date using repository", func() {
			testDate := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
			expectedRate := 82.50

			mockExchange.EXPECT().
				GetRecordsForDate(ctx, testDate).
				Return([]tax.SbiRate{{TTBuy: expectedRate}}, nil)

			rate, err := sbiManager.GetTTBuyRate(testDate)
			Expect(err).To(BeNil())
			Expect(rate).To(Equal(expectedRate))
		})
	})
})
