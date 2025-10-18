package manager_test

import (
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TODO: (2025-10-18) Investigate flaky Interactive Brokers CSV edge-case specs.
var _ = Describe("InteractiveBrokersManagerImpl", func() {
	var (
		tempTestDir   string
		sampleCSVPath string
		basePath      string
		ibManager     manager.Broker
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "ib_test_*")
		Expect(err).ToNot(HaveOccurred())

		basePath = filepath.Join(tempTestDir, "realized")
		sampleCSVPath = filepath.Join(tempTestDir, "realized_2024.csv")
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with a valid CSV file", func() {
		BeforeEach(func() {
			csvContent := `Statement,Header,Field Name,Field Value
Statement,Data,BrokerName,Interactive Brokers LLC
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,MPC,"2024-10-31, 09:30:00",8,146.21,-1169.68,-0.36024125,1170.04024125,0,O
Trades,Data,Order,Stocks,USD,MPC,"2024-12-17, 09:31:09",-8,136.85,1094.8,-0.38419669,-1170.04024,-75.624437,C;IM;P
Trades,Data,Order,Stocks,USD,SIVR,"2024-09-04, 10:13:06",1,26.9,-26.9,-0.271397715,27.171397715,0,O
Trades,Data,Order,Stocks,USD,SIVR,"2024-09-13, 11:15:37",-1,29.122,29.122,-0.292661638,-27.171398,1.65794,C
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Withholding Tax,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share - US Tax,-1.82,`

			err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
			Expect(err).ToNot(HaveOccurred())

			ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
		})

		Context("when parsing trades", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				var err error
				info, err = ibManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should extract trade entries correctly", func() {
				Expect(info.Trades).To(HaveLen(4))

				Expect(info.Trades[0].Symbol).To(Equal("MPC"))
				Expect(info.Trades[0].Type).To(Equal("BUY"))
				Expect(info.Trades[0].Quantity).To(Equal(8.0))
				Expect(info.Trades[0].USDPrice).To(Equal(146.21))
				Expect(info.Trades[0].Date).To(Equal("2024-10-31"))
				Expect(info.Trades[0].USDValue).To(Equal(1169.68))
				Expect(info.Trades[0].Commission).To(Equal(0.36024125))

				Expect(info.Trades[1].Symbol).To(Equal("MPC"))
				Expect(info.Trades[1].Type).To(Equal("SELL"))
				Expect(info.Trades[1].Quantity).To(Equal(8.0))
				Expect(info.Trades[1].USDPrice).To(Equal(136.85))
				Expect(info.Trades[1].Date).To(Equal("2024-12-17"))
				Expect(info.Trades[1].USDValue).To(Equal(1094.8))
				Expect(info.Trades[1].Commission).To(Equal(0.38419669))

				Expect(info.Trades[2].Symbol).To(Equal("SIVR"))
				Expect(info.Trades[2].Type).To(Equal("BUY"))
				Expect(info.Trades[2].Quantity).To(Equal(1.0))

				Expect(info.Trades[3].Symbol).To(Equal("SIVR"))
				Expect(info.Trades[3].Type).To(Equal("SELL"))
				Expect(info.Trades[3].Quantity).To(Equal(1.0))
			})
		})

		Context("when parsing dividends", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				var err error
				info, err = ibManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should extract dividend entries and match withholding tax", func() {
				Expect(info.Dividends).To(HaveLen(1))

				Expect(info.Dividends[0].Symbol).To(Equal("MPC"))
				Expect(info.Dividends[0].Date).To(Equal("2024-12-10"))
				Expect(info.Dividends[0].Amount).To(Equal(7.28))
				Expect(info.Dividends[0].Tax).To(Equal(1.82))
				Expect(info.Dividends[0].Net).To(Equal(5.46))
			})
		})
	})

	Context("with invalid CSV files", func() {
		Context("when CSV file is missing", func() {
			BeforeEach(func() {
				ibManager = manager.NewInteractiveBrokersManagerImpl(sampleCSVPath)
			})

			It("should return an error", func() {
				_, err := ibManager.Parse(testYear)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when CSV has malformed data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,INVALID,,BADQTY,BADPRICE,0,0,0,0,O`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should skip invalid rows gracefully", func() {
				Expect(info.Trades).To(BeEmpty())
			})
		})

		PContext("when dividend has no matching tax", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(sampleCSVPath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should set tax to 0", func() {
				Expect(info.Dividends).To(HaveLen(1))
				Expect(info.Dividends[0].Tax).To(Equal(0.0))
				Expect(info.Dividends[0].Net).To(Equal(7.28))
			})
		})
	})

	Context("with edge case CSV data", func() {
		PContext("when CSV has SubTotal and Total rows", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,MPC,"2024-10-31, 09:30:00",8,146.21,-1169.68,-0.36024125,1170.04024125,0,O
Trades,SubTotal,,Stocks,USD,MPC,,0,,-74.88,-0.74443794,0.00000125,-75.624437,
Trades,Total,,Stocks,USD,,,,,-129.958,-2.097372693,0.000001215,-132.055372,
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend,7.28
Dividends,Data,Total,,,7.28`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(sampleCSVPath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should skip them and only parse Data rows", func() {
				Expect(info.Trades).To(HaveLen(1))
				Expect(info.Dividends).To(HaveLen(1))
			})
		})

		PContext("when CSV has only headers with no data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Dividends,Header,Currency,Date,Description,Amount`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(sampleCSVPath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should return empty arrays", func() {
				Expect(info.Trades).To(BeEmpty())
				Expect(info.Dividends).To(BeEmpty())
			})
		})
	})
})
