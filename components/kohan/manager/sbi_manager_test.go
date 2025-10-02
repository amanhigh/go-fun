package manager_test

import (
	"context"
	"errors"
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
2024-01-23,82.50,83.50
2024-01-24,82.75,83.75
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
			Expect(err).ToNot(HaveOccurred())

			// Verify file content
			content, readErr := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(readErr).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal(testCSV))
		})

		It("should skip download if file already exists", func() {
			// Create test file first
			err := os.WriteFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).ToNot(HaveOccurred())

			// Expect no client calls since file exists
			// mockClient.EXPECT().FetchExchangeRates(ctx) should not be called

			err = sbiManager.DownloadRates(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Verify file content remains unchanged
			content, readErr := os.ReadFile(filepath.Join(testDir, tax.SBI_RATES_FILENAME))
			Expect(readErr).ToNot(HaveOccurred())
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
			expectedRecord := tax.SbiRate{Date: "2024-01-23", TTBuy: expectedRate, TTSell: 83.50}

			// Mock GetRecordsForTicker for the exact date
			mockExchange.EXPECT().GetRecordsForTicker(ctx, "2024-01-23").Return([]tax.SbiRate{expectedRecord}, nil)

			rate, err := sbiManager.GetTTBuyRate(ctx, testDate)
			Expect(err).ToNot(HaveOccurred())
			Expect(rate).To(Equal(expectedRate))
		})

		Context("When exact date not found", func() {
			var (
				requestedDate = time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC)
				closestDate   = time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
				expectedRate  = 82.50
			)

			BeforeEach(func() {
				// First mock the exact date lookup to return nothing
				mockExchange.EXPECT().GetRecordsForTicker(ctx, "2024-01-24").Return([]tax.SbiRate{}, nil)

				// Mock GetAllRecords for the fallback mechanism
				closestRecord := tax.SbiRate{Date: "2024-01-23", TTBuy: expectedRate, TTSell: 83.50}
				mockExchange.EXPECT().GetAllRecords(ctx).Return([]tax.SbiRate{closestRecord}, nil)
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
				Expect(closestErr.Error()).To(ContainSubstring("exact rate not found for 2024-01-24"))
				Expect(closestErr.Error()).To(ContainSubstring("using closest available date 2024-01-23"))
			})
		})

		Context("When weekend/holiday with zero TTBuy rate", func() {
			It("should skip zero rates and use previous non-zero rate", func() {
				// Scenario: Request rate for Monday (2022-05-02) after weekend
				// Weekend Saturday has TTBuy=0 (no forex trading)
				// Should skip Saturday and use Friday's rate
				requestedDate := time.Date(2022, 5, 2, 0, 0, 0, 0, time.UTC) // Monday
				fridayDate := time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC)   // Friday
				expectedRate := 76.00

				// Mock exact date lookup (Monday not found - holiday)
				mockExchange.EXPECT().GetRecordsForTicker(ctx, "2022-05-02").Return([]tax.SbiRate{}, nil)

				// Mock GetAllRecords with Friday (valid) and Saturday (zero) rates
				allRates := []tax.SbiRate{
					{Date: "2022-04-29", TTBuy: 76.00, TTSell: 77.00}, // Friday - valid
					{Date: "2022-04-30", TTBuy: 0, TTSell: 0},         // Saturday - weekend (zero)
				}
				mockExchange.EXPECT().GetAllRecords(ctx).Return(allRates, nil)

				rate, err := sbiManager.GetTTBuyRate(ctx, requestedDate)

				// Should skip Saturday (zero) and use Friday
				Expect(rate).To(Equal(expectedRate))

				// Verify ClosestDateError with Friday's date (not Saturday)
				closestErr, ok := err.(tax.ClosestDateError)
				Expect(ok).To(BeTrue())
				Expect(closestErr.GetRequestedDate()).To(Equal(requestedDate))
				Expect(closestErr.GetClosestDate()).To(Equal(fridayDate))
			})

			It("should return error when all previous rates are zero", func() {
				requestedDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

				// Mock exact date lookup
				mockExchange.EXPECT().GetRecordsForTicker(ctx, "2024-01-10").Return([]tax.SbiRate{}, nil)

				// All rates are zero (all weekends)
				allRates := []tax.SbiRate{
					{Date: "2024-01-07", TTBuy: 0, TTSell: 0},
					{Date: "2024-01-06", TTBuy: 0, TTSell: 0},
				}
				mockExchange.EXPECT().GetAllRecords(ctx).Return(allRates, nil)

				rate, err := sbiManager.GetTTBuyRate(ctx, requestedDate)

				// Should return error when no valid rate found
				Expect(rate).To(Equal(0.0))
				_, ok := err.(tax.RateNotFoundError)
				Expect(ok).To(BeTrue())
			})
		})
	})

	Context("GetLastMonthEndRate", func() {
		var allRates []tax.SbiRate

		BeforeEach(func() {
			allRates = []tax.SbiRate{
				{Date: "2024-01-10", TTBuy: 82.40, TTSell: 83.40},
				{Date: "2024-01-15", TTBuy: 82.50, TTSell: 83.50},
				{Date: "2024-01-31", TTBuy: 82.75, TTSell: 83.75},
				{Date: "2024-02-05", TTBuy: 83.00, TTSell: 84.00},
			}
		})

		It("should return last rate in preceding month for dividend/interest date", func() {
			mockExchange.EXPECT().GetAllRecords(ctx).Return(allRates, nil).Once()

			inputDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Rate).To(Equal(82.75))
			Expect(result.ActualDate).To(Equal(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)))
		})

		It("should return last available rate when month-end date missing", func() {
			ratesWithoutMonthEnd := []tax.SbiRate{
				{Date: "2024-01-10", TTBuy: 82.40, TTSell: 83.40},
				{Date: "2024-01-15", TTBuy: 82.50, TTSell: 83.50},
				{Date: "2024-02-05", TTBuy: 83.00, TTSell: 84.00},
			}

			mockExchange.EXPECT().GetAllRecords(ctx).Return(ratesWithoutMonthEnd, nil).Once()

			inputDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Rate).To(Equal(82.50))
			Expect(result.ActualDate).To(Equal(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)))
		})

		It("should cache results for subsequent calls", func() {
			mockExchange.EXPECT().GetAllRecords(ctx).Return(allRates, nil).Once()

			inputDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			result1, err1 := sbiManager.GetLastMonthEndRate(ctx, inputDate)
			Expect(err1).ToNot(HaveOccurred())

			result2, err2 := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err2).ToNot(HaveOccurred())
			Expect(result2.Rate).To(Equal(result1.Rate))
			Expect(result2.ActualDate).To(Equal(result1.ActualDate))
		})

		It("should return error when no rates in preceding month", func() {
			emptyMonthRates := []tax.SbiRate{
				{Date: "2024-02-05", TTBuy: 83.00, TTSell: 84.00},
			}

			mockExchange.EXPECT().GetAllRecords(ctx).Return(emptyMonthRates, nil).Once()

			inputDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).To(HaveOccurred())
			Expect(result.Rate).To(Equal(0.0))
			Expect(result.ActualDate.IsZero()).To(BeTrue())

			_, ok := err.(tax.RateNotFoundError)
			Expect(ok).To(BeTrue())
		})

		It("should handle different months independently", func() {
			mockExchange.EXPECT().GetAllRecords(ctx).Return(allRates, nil).Once()

			febDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			janResult, err1 := sbiManager.GetLastMonthEndRate(ctx, febDate)
			Expect(err1).ToNot(HaveOccurred())
			Expect(janResult.Rate).To(Equal(82.75))
			Expect(janResult.ActualDate.Month()).To(Equal(time.January))

			marDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
			febResult, err2 := sbiManager.GetLastMonthEndRate(ctx, marDate)
			Expect(err2).ToNot(HaveOccurred())
			Expect(febResult.Rate).To(Equal(83.00))
			Expect(febResult.ActualDate.Month()).To(Equal(time.February))
		})

		It("should propagate repository errors", func() {
			expectedErr := common.NewServerError(errors.New("repository error"))
			mockExchange.EXPECT().GetAllRecords(ctx).Return(nil, expectedErr).Once()

			inputDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).To(HaveOccurred())
			Expect(result.Rate).To(Equal(0.0))
			Expect(result.ActualDate.IsZero()).To(BeTrue())
		})

		It("should handle leap year February correctly", func() {
			leapYearRates := []tax.SbiRate{
				{Date: "2024-02-28", TTBuy: 83.00, TTSell: 84.00},
				{Date: "2024-02-29", TTBuy: 83.05, TTSell: 84.05},
			}

			mockExchange.EXPECT().GetAllRecords(ctx).Return(leapYearRates, nil).Once()

			inputDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Rate).To(Equal(83.05))
			Expect(result.ActualDate.Day()).To(Equal(29))
		})

		It("should calculate preceding month correctly for January", func() {
			decemberRates := []tax.SbiRate{
				{Date: "2023-12-28", TTBuy: 82.00, TTSell: 83.00},
				{Date: "2023-12-31", TTBuy: 82.10, TTSell: 83.10},
			}

			mockExchange.EXPECT().GetAllRecords(ctx).Return(decemberRates, nil).Once()

			inputDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			result, err := sbiManager.GetLastMonthEndRate(ctx, inputDate)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Rate).To(Equal(82.10))
			Expect(result.ActualDate).To(Equal(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)))
		})
	})
})
