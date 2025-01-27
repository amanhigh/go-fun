package repository_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testCSV = `Symbol,Date,Type,Quantity,Price,Value,Commission
AAPL,2024-01-15,BUY,100,150.50,15050.00,0.50
GOOGL,2024-01-15,BUY,50,2500.00,125000.00,1.00
AAPL,2024-01-16,SELL,50,155.00,7750.00,0.50`

var _ = Describe("TradeRepository", func() {
	var (
		tradeRepo repository.TradeRepository
		testDir   string
		ctx       = context.Background()
		err       error
	)

	BeforeEach(func() {
		// Create temp directory
		testDir, err = os.MkdirTemp("", "trade-test-*")
		Expect(err).NotTo(HaveOccurred())

		// Create test file
		tradeFile := filepath.Join(testDir, "trades.csv")
		err = os.WriteFile(tradeFile, []byte(testCSV), util.DEFAULT_PERM)
		Expect(err).To(BeNil())

		tradeRepo = repository.NewTradeRepository(tradeFile)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("Success Cases", func() {
		It("should read all trades", func() {
			trades, err := tradeRepo.GetAllRecords(ctx)
			Expect(err).To(BeNil())
			Expect(trades).To(HaveLen(3))

			// Verify first trade
			first := trades[0]
			Expect(first.Symbol).To(Equal("AAPL"))
			Expect(first.Type).To(Equal("BUY"))
			Expect(first.Quantity).To(Equal(100.0))
			Expect(first.USDPrice).To(Equal(150.50))
			Expect(first.USDValue).To(Equal(15050.00))
			Expect(first.Commission).To(Equal(0.50))
		})

		It("should get unique tickers", func() {
			tickers, err := tradeRepo.GetUniqueTickers(ctx)
			Expect(err).To(BeNil())
			Expect(tickers).To(ConsistOf("AAPL", "GOOGL"))
		})

		It("should filter by ticker", func() {
			trades, err := tradeRepo.GetRecordsForTicker(ctx, "AAPL")
			Expect(err).To(BeNil())
			Expect(trades).To(HaveLen(2))
			Expect(trades[0].Symbol).To(Equal("AAPL"))
			Expect(trades[1].Symbol).To(Equal("AAPL"))
		})
	})

	Context("Error Cases", func() {
		It("should handle missing file", func() {
			invalidRepo := repository.NewTradeRepository("invalid.csv")
			_, err := invalidRepo.GetAllRecords(ctx)
			Expect(err).To(Not(BeNil()))
		})

		It("should handle malformed CSV", func() {
			malformedFile := filepath.Join(testDir, "malformed.csv")
			err := os.WriteFile(malformedFile, []byte("invalid,csv"), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			invalidRepo := repository.NewTradeRepository(malformedFile)
			_, err = invalidRepo.GetAllRecords(ctx)
			Expect(err).To(Not(BeNil()))
		})

		It("should handle empty file", func() {
			emptyFile := filepath.Join(testDir, "empty.csv")
			err := os.WriteFile(emptyFile, []byte(""), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			emptyRepo := repository.NewTradeRepository(emptyFile)
			trades, err := emptyRepo.GetAllRecords(ctx)
			Expect(err).To(Not(BeNil()))
			Expect(trades).To(BeNil())
		})
	})
})
