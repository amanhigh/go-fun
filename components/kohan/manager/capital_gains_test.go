package manager_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	adiFirstBuyDate  = "2024-01-04"
	snpsFirstBuyDate = "2024-02-21"
)

var _ = Describe("CapitalGainsManager", func() {
	// TODO: #B Not Working Test
	var (
		ctx               context.Context
		mockTickerManager *mocks.TickerManager
		cgManager         manager.CapitalGainsManager
		testDir           string
		statementFile     string
		year              = 2024
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockTickerManager = mocks.NewTickerManager(GinkgoT())

		var err error
		testDir, err = os.MkdirTemp("", "capital-gains-test-*")
		Expect(err).NotTo(HaveOccurred())

		statementFile = "broker_statement.csv"
		cgManager = manager.NewCapitalGainsManager(mockTickerManager, testDir, statementFile)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("Success Cases", func() {
		const (
			adiSellDate  = "2024-01-19"
			snpsSellDate = "2024-02-22"
		)

		BeforeEach(func() {
			// Create test CSV
			testData := `Security,Quantity Sold,Date Acquired,Buying Price (USD),Date Sold,Selling Price (USD),Proceeds (USD),Cost Basis (USD),Gains/Losses (USD)
ADI,2,2024-01-04,182.08,2024-01-19,194.56,389.12,364.16,24.96
SNPS,1,2024-02-21,531.03,2024-02-22,589.44,589.44,531.03,58.41`
			err := os.WriteFile(filepath.Join(testDir, statementFile), []byte(testData), 0644)
			Expect(err).NotTo(HaveOccurred())

			// Setup price mocks for ADI
			setupADIPriceMocks(ctx, mockTickerManager)
			// Setup price mocks for SNPS
			setupSNPSPriceMocks(ctx, mockTickerManager)
		})

		PIt("should process all tickers correctly", func() {
			result, err := cgManager.AnalysePositions(ctx, year)
			Expect(err).NotTo(HaveOccurred())

			// Verify ADI Analysis
			adiAnalysis := result["ADI"]
			Expect(adiAnalysis.FirstPosition.Quantity).To(Equal(2.0))
			Expect(adiAnalysis.FirstPosition.USDPrice).To(Equal(182.08))
			Expect(adiAnalysis.PeakPosition.USDPrice).To(Equal(194.56))
			Expect(adiAnalysis.FirstPosition.Date).To(Equal(parseDateMust("2024-01-04")))
			Expect(adiAnalysis.PeakPosition.Date).To(Equal(parseDateMust("2024-01-19")))

			// Verify SNPS Analysis
			snpsAnalysis := result["SNPS"]
			Expect(snpsAnalysis.FirstPosition.Quantity).To(Equal(1.0))
			Expect(snpsAnalysis.FirstPosition.USDPrice).To(Equal(531.03))
			Expect(snpsAnalysis.PeakPosition.USDPrice).To(Equal(589.44))
			Expect(snpsAnalysis.PeakPosition.Date).To(Equal(parseDateMust("2024-02-22")))
		})
	})

	Context("Error Cases", func() {
		It("should handle missing file", func() {
			result, err := cgManager.AnalysePositions(ctx, year)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to open broker statement"))
			Expect(result).To(BeNil())
		})

		It("should handle malformed CSV", func() {
			// Write invalid CSV
			testData := "Invalid,CSV,Format\nWithout,Proper,Headers"
			err := os.WriteFile(filepath.Join(testDir, statementFile), []byte(testData), 0644)
			Expect(err).NotTo(HaveOccurred())

			result, err := cgManager.AnalysePositions(ctx, year)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid date format"))
			Expect(result).To(BeNil())
		})

		It("should handle price fetch errors", func() {
			// Write valid CSV
			testData := `Security,Quantity Sold,Date Acquired,Buying Price (USD),Date Sold,Selling Price (USD),Proceeds (USD),Cost Basis (USD),Gains/Losses (USD)
ADI,2,2024-01-04,182.08,2024-01-19,194.56,389.12,364.16,24.96`
			err := os.WriteFile(filepath.Join(testDir, statementFile), []byte(testData), 0644)
			Expect(err).NotTo(HaveOccurred())

			// Setup mock to return error
			mockTickerManager.EXPECT().
				GetPrice(ctx, "ADI", parseDateMust(adiFirstBuyDate)).
				Return(0.0, fmt.Errorf("price fetch error"))

			result, err := cgManager.AnalysePositions(ctx, year)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get price"))
			Expect(result).To(BeNil())
		})
	})
})

func setupADIPriceMocks(ctx context.Context, mockTicker *mocks.TickerManager) {
	// First buy date
	mockTicker.EXPECT().
		GetPrice(ctx, "ADI", parseDateMust(adiFirstBuyDate)).
		Return(182.08, nil)

	// Sell/Peak date
	mockTicker.EXPECT().
		GetPrice(ctx, "ADI", parseDateMust("2024-01-19")).
		Return(194.56, nil)

	// Year end
	mockTicker.EXPECT().
		GetPrice(ctx, "ADI", parseDateMust("2024-12-31")).
		Return(190.00, nil)
}

func setupSNPSPriceMocks(ctx context.Context, mockTicker *mocks.TickerManager) {
	// First buy date
	mockTicker.EXPECT().
		GetPrice(ctx, "SNPS", parseDateMust(snpsFirstBuyDate)).
		Return(531.03, nil)

	// Sell/Peak date
	mockTicker.EXPECT().
		GetPrice(ctx, "SNPS", parseDateMust("2024-02-22")).
		Return(589.44, nil)

	// Year end
	mockTicker.EXPECT().
		GetPrice(ctx, "SNPS", parseDateMust("2024-12-31")).
		Return(550.00, nil)
}

func parseDateMust(date string) time.Time {
	t, err := time.Parse(common.DateOnly, date)
	if err != nil {
		panic(err)
	}
	return t
}
