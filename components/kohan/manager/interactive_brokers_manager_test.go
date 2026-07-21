package manager_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TODO: #C (2025-10-18) Investigate flaky Interactive Brokers CSV edge-case specs.
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
		var info tax.BrokerageInfo

		BeforeEach(func() {
			csvContent := `Statement,Header,Field Name,Field Value
Statement,Data,BrokerName,Interactive Brokers LLC
Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,MPC,"2024-10-31, 09:30:00",8,146.21,146.20,-1169.68,-0.36024125,1170.04024125,0,O
Trades,Data,Order,Stocks,USD,MPC,"2024-12-17, 09:31:09",-8,136.85,136.86,1094.8,-0.38419669,-1170.04024,-75.624437,C;IM;P
Trades,Data,Order,Stocks,USD,SIVR,"2024-09-04, 10:13:06",1,26.9,26.91,-26.9,-0.271397715,27.171397715,0,O
Trades,Data,Order,Stocks,USD,SIVR,"2024-09-13, 11:15:37",-1,29.122,29.12,29.122,-0.292661638,-27.171398,1.65794,C
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Withholding Tax,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share - US Tax,-1.82,
Interest,Header,Currency,Date,Description,Amount
Interest,Data,USD,2024-12-15,USD Credit Interest for Nov-2024,2.50
Interest,Data,USD,2025-01-05,USD Credit Interest for Dec-2024,1.75
Interest,Data,Total,,,4.25`

			err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
			Expect(err).ToNot(HaveOccurred())

			ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
			var parseErr error
			info, parseErr = ibManager.Parse(testYear)
			Expect(parseErr).ToNot(HaveOccurred())
		})

		Context("period coverage", func() {
			It("should have correct CoverageThrough", func() {
				Expect(info.CoverageThrough).To(Equal(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)))
			})
		})

		Context("when parsing trades", func() {
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
			It("should extract dividend entries and match withholding tax", func() {
				Expect(info.Dividends).To(HaveLen(1))

				Expect(info.Dividends[0].Symbol).To(Equal("MPC"))
				Expect(info.Dividends[0].Date).To(Equal("2024-12-10"))
				Expect(info.Dividends[0].Amount).To(Equal(7.28))
				Expect(info.Dividends[0].Tax).To(Equal(1.82))
				Expect(info.Dividends[0].Net).To(Equal(5.46))
			})
		})

		Context("when parsing interest", func() {
			It("should extract USD interest entries and skip Total rows", func() {
				Expect(info.Interests).To(HaveLen(2))

				Expect(info.Interests[0].Symbol).To(Equal("CASH"))
				Expect(info.Interests[0].Date).To(Equal("2024-12-15"))
				Expect(info.Interests[0].Amount).To(Equal(2.50))
				Expect(info.Interests[0].Tax).To(Equal(0.0))
				Expect(info.Interests[0].Net).To(Equal(2.50))

				Expect(info.Interests[1].Symbol).To(Equal("CASH"))
				Expect(info.Interests[1].Date).To(Equal("2025-01-05"))
				Expect(info.Interests[1].Amount).To(Equal(1.75))
				Expect(info.Interests[1].Tax).To(Equal(0.0))
				Expect(info.Interests[1].Net).To(Equal(1.75))
			})
		})
	})

	Context("with invalid CSV files", func() {
		Context("when CSV file is missing", func() {
			BeforeEach(func() {
				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
			})

			It("should return an error about missing file", func() {
				_, err := ibManager.Parse(testYear)
				Expect(err).To(MatchError("failed to open CSV file"))
			})
		})

		Context("when CSV has malformed data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,INVALID,,BADQTY,1,1,0,0,0,0,O`

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

		Context("when CSV has no Period metadata", func() {
			BeforeEach(func() {
				csvContent := `Some,Random,Data`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
			})

			It("should return an error about missing period", func() {
				_, err := ibManager.Parse(testYear)
				Expect(err).To(MatchError("period metadata not found"))
			})
		})

		Context("when CSV has malformed Period format", func() {
			BeforeEach(func() {
				csvContent := `Statement,Data,Period,InvalidPeriodFormat`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
			})

			It("should return an error about invalid period format", func() {
				_, err := ibManager.Parse(testYear)
				Expect(err).To(MatchError("invalid period format: InvalidPeriodFormat"))
			})
		})
	})

	Context("with edge case CSV data", func() {
		Context("when dividend has no matching tax", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				// Use basePath so Parse(year) resolves basePath_YYYY.csv
				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
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

		Context("when CSV has SubTotal and Total rows", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,MPC,"2024-10-31, 09:30:00",8,146.21,146.20,-1169.68,-0.36024125,1170.04024125,0,O
Trades,SubTotal,,Stocks,USD,MPC,,0,,,-74.88,-0.74443794,0.00000125,-75.624437,
Trades,Total,,Stocks,USD,,,,,,-129.958,-2.097372693,0.000001215,-132.055372,
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend,7.28
Dividends,Data,Total,,,7.28`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				// Use basePath so Parse(year) resolves basePath_YYYY.csv
				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should skip them and only parse Data rows", func() {
				Expect(info.Trades).To(HaveLen(1))
				Expect(info.Dividends).To(HaveLen(1))
			})
		})

		Context("when CSV has only headers with no data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				csvContent := `Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Dividends,Header,Currency,Date,Description,Amount
Interest,Header,Currency,Date,Description,Amount`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				// Use basePath so Parse(year) resolves basePath_YYYY.csv
				ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
				var parseErr error
				info, parseErr = ibManager.Parse(testYear)
				Expect(parseErr).ToNot(HaveOccurred())
			})

			It("should return empty arrays", func() {
				Expect(info.Trades).To(BeEmpty())
				Expect(info.Dividends).To(BeEmpty())
				Expect(info.Interests).To(BeEmpty())
			})
		})
	})

	Context("multi-file merge", func() {
		var info tax.BrokerageInfo

		writeCSV := func(year int, content string) {
			err := os.WriteFile(
				fmt.Sprintf("%s_%d.csv", basePath, year),
				[]byte(content), 0600)
			Expect(err).ToNot(HaveOccurred())
		}

		BeforeEach(func() {
			// File 1 (2024): AAPL trade, MPC dividend with $1.82 withholding, USD Interest $2.50
			writeCSV(2024, `Statement,Data,Period,"January 1, 2024 - December 31, 2024"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,AAPL,"2024-10-31, 09:30:00",8,146.21,146.20,-1169.68,-0.36024125,1170.04024125,0,O
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Withholding Tax,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share - US Tax,-1.82,
Interest,Header,Currency,Date,Description,Amount
Interest,Data,USD,2024-12-15,USD Credit Interest for Nov-2024,2.50`)

			// File 2 (2025): MSFT trade, MPC dividend with $2.50 withholding, USD Interest $1.75
			writeCSV(2025, `Statement,Data,Period,"January 1, 2025 - December 31, 2025"
Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,C. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,MSFT,"2025-06-15, 10:30:00",5,350.00,349.50,-1750.00,-0.50,1750.50,0,O
Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Withholding Tax,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share - US Tax,-2.50,
Interest,Header,Currency,Date,Description,Amount
Interest,Data,USD,2025-01-05,USD Credit Interest for Dec-2024,1.75`)

			ibManager = manager.NewInteractiveBrokersManagerImpl(basePath)
			var err error
			info, err = ibManager.Parse(2024)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should aggregate both trades", func() {
			Expect(info.Trades).To(HaveLen(2))
			Expect(info.Trades[0].Symbol).To(Equal("AAPL"))
			Expect(info.Trades[1].Symbol).To(Equal("MSFT"))
		})

		It("should aggregate both dividends with per-file withholding-tax isolation", func() {
			Expect(info.Dividends).To(HaveLen(2))
			var mpcTaxes []float64
			for _, d := range info.Dividends {
				Expect(d.Symbol).To(Equal("MPC"))
				mpcTaxes = append(mpcTaxes, d.Tax)
			}
			Expect(mpcTaxes).To(ConsistOf(1.82, 2.50))
		})

		It("should aggregate both interests", func() {
			Expect(info.Interests).To(HaveLen(2))
			var amounts []float64
			for _, i := range info.Interests {
				amounts = append(amounts, i.Amount)
			}
			Expect(amounts).To(ConsistOf(2.50, 1.75))
		})

		It("should set CoverageThrough to the latest period end", func() {
			Expect(info.CoverageThrough).To(Equal(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)))
		})
	})
})
