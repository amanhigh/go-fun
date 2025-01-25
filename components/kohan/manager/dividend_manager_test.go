package manager

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testCSV = `Security,Dividend Date,Dividend Per Share (USD),Dividend Received (USD),Dividend Tax (USD),Net Dividend (USD)
VGK,2023-06-26,1.099,197.82,-49.46,148.36
VGK,2023-09-22,0.2606,46.91,-11.73,35.18`

var (
	mockSBIManager  *mocks.SBIManager
	dividendManager DividendManager
	testDir         string
	ctx             = context.Background()
)

func timeFromStr(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

var _ = Describe("DividendManager", func() {
	BeforeEach(func() {
		var err error
		mockSBIManager = mocks.NewSBIManager(GinkgoT())
		testDir, err = os.MkdirTemp("", "dividend-test-*")
		Expect(err).NotTo(HaveOccurred())

		dividendManager = NewDividendManager(mockSBIManager, testDir, "dividends.csv")
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("GetDividendTransactions", func() {
		BeforeEach(func() {
			err := os.WriteFile(filepath.Join(testDir, "dividends.csv"), []byte(testCSV), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			mockSBIManager.EXPECT().
				GetTTBuyRate(timeFromStr("2023-06-26")).
				Return(82.6773, nil)
			mockSBIManager.EXPECT().
				GetTTBuyRate(timeFromStr("2023-09-22")).
				Return(82.6784, nil)
		})

		It("should process dividends correctly", func() {
			transactions, err := dividendManager.GetDividendTransactions(ctx)
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			first := transactions[0]
			Expect(first.Security).To(Equal("VGK"))
			Expect(first.DividendDate).To(Equal("2023-06-26"))
			Expect(first.DividendPerShare).To(Equal(1.099))
			Expect(first.DividendReceived).To(Equal(197.82))
			Expect(first.DividendTax).To(Equal(-49.46))
			Expect(first.NetDividend).To(Equal(148.36))
			Expect(first.USDINRRate).To(Equal(82.6773))
			Expect(first.NetDividendINR).To(BeNumerically("~", 12266.0, 0.1))
			Expect(first.DividendTaxINR).To(BeNumerically("~", -4089.22, 0.1))

			second := transactions[1]
			Expect(second.Security).To(Equal("VGK"))
			Expect(second.DividendDate).To(Equal("2023-09-22"))
			Expect(second.DividendPerShare).To(Equal(0.2606))
			Expect(second.DividendReceived).To(Equal(46.91))
			Expect(second.DividendTax).To(Equal(-11.73))
			Expect(second.NetDividend).To(Equal(35.18))
			Expect(second.USDINRRate).To(Equal(82.6784))
			Expect(second.NetDividendINR).To(BeNumerically("~", 2908.63, 0.1))
			Expect(second.DividendTaxINR).To(BeNumerically("~", -969.82, 0.1))
		})
	})

	Context("Error Handling", func() {
		It("should handle missing file", func() {
			transactions, err := dividendManager.GetDividendTransactions(ctx)
			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should handle invalid CSV", func() {
			err := os.WriteFile(filepath.Join(testDir, "dividends.csv"),
				[]byte("invalid,csv"), util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			transactions, err := dividendManager.GetDividendTransactions(ctx)
			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should handle SBI manager rate error", func() {
			err := os.WriteFile(filepath.Join(testDir, "dividends.csv"), []byte(testCSV),
				util.DEFAULT_PERM)
			Expect(err).To(BeNil())

			mockSBIManager.EXPECT().
				GetTTBuyRate(timeFromStr("2023-06-26")).
				Return(0.0, common.ErrNotFound)

			transactions, err := dividendManager.GetDividendTransactions(ctx)
			Expect(err).To(Equal(common.ErrNotFound))
			Expect(transactions).To(BeNil())
		})
	})
})
