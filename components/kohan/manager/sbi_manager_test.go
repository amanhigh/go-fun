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
		mockClient *mocks.SBIClient
		sbiManager *manager.SBIManagerImpl
		testDir    string
		ctx        = context.Background()
		err        common.HttpError
	)

	BeforeEach(func() {
		mockClient = mocks.NewSBIClient(GinkgoT())

		var err error
		testDir, err = os.MkdirTemp("", "sbi-test-*")
		Expect(err).NotTo(HaveOccurred())

		sbiManager = manager.NewSBIManager(mockClient, testDir)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("DownloadRates", func() {
		It("should download and save rates successfully", func() {
			mockClient.EXPECT().FetchExchangeRates(ctx).Return(testCSV, nil)

			err = sbiManager.DownloadRates(ctx)
			Expect(err).To(BeNil())

			// Verify file content
			content, err := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})

		// FIXME: #B Should not downolad Rates if exists

		It("should handle client error", func() {
			expectedErr := common.NewHttpError("Failed to fetch exchange rates", http.StatusInternalServerError)
			mockClient.EXPECT().FetchExchangeRates(ctx).Return("", expectedErr)

			err = sbiManager.DownloadRates(ctx)
			Expect(err).To(Equal(expectedErr))
		})

		It("should create directory if not exists", func() {
			// Create nested test directory
			nestedDir := filepath.Join(testDir, "nested", "dir")
			sbiManager = manager.NewSBIManager(mockClient, nestedDir)

			mockClient.EXPECT().FetchExchangeRates(ctx).Return(testCSV, nil)

			err = sbiManager.DownloadRates(ctx)
			Expect(err).To(BeNil())

			// Verify directory was created with file
			content, err := os.ReadFile(filepath.Join(nestedDir, tax.SBI_RATES_FILENAME))
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal(testCSV))
		})
	})

	Context("GetTTBuyRate", func() {
		var testDate time.Time

		BeforeEach(func() {
			var err error
			testDate, err = time.Parse("2006-01-02", "2024-01-23")
			Expect(err).To(BeNil())
		})

		It("should return error when file not found", func() {
			_, err = sbiManager.GetTTBuyRate(testDate)
			Expect(err.Error()).To(Equal("SBI rates file not found"))
			Expect(err.Code()).To(Equal(http.StatusNotFound))
		})

		It("should return rate for matching date", func() {
			// Create test file
			err := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			rate, err := sbiManager.GetTTBuyRate(testDate)
			Expect(err).To(BeNil())
			Expect(rate).To(Equal(82.50))
		})

		It("should return not found for missing date", func() {
			// Create test file
			err := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			missingDate := testDate.AddDate(0, 0, -1)
			_, err = sbiManager.GetTTBuyRate(missingDate)
			Expect(err).To(Equal(common.ErrNotFound))
		})

		PIt("should handle invalid CSV file", func() {
			// FIXME: Test is failing after moving to gocsv
			
			// Create invalid CSV file
			writeErr := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte("invalid,csv"), util.DEFAULT_PERM)
			Expect(writeErr).To(BeNil())

			_, err = sbiManager.GetTTBuyRate(testDate)
			Expect(err).NotTo(BeNil())
			Expect(err.Code()).To(Equal(http.StatusInternalServerError))
		})

		It("should handle empty file", func() {
			// Create empty file
			writeErr := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(""), util.DEFAULT_PERM)
			Expect(writeErr).To(BeNil())

			_, err = sbiManager.GetTTBuyRate(testDate)
			Expect(err).NotTo(BeNil())
			Expect(err.Code()).To(Equal(http.StatusInternalServerError))
		})

		It("should cache rates after first call", func() {
			// Create test file
			err := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			// First call should load cache
			rate1, err := sbiManager.GetTTBuyRate(testDate)
			Expect(err).To(BeNil())
			Expect(rate1).To(Equal(82.50))

			// Second call should use cache
			rate2, err := sbiManager.GetTTBuyRate(testDate)
			Expect(err).To(BeNil())
			Expect(rate2).To(Equal(82.50))

			// Verify cache miss
			missingDate := testDate.AddDate(0, 0, -1)
			_, err = sbiManager.GetTTBuyRate(missingDate)
			Expect(err).To(Equal(common.ErrNotFound))
		})
	})
})
