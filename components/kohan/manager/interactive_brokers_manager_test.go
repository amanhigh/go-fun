package manager_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InteractiveBrokersManager", func() {
	var (
		tempTestDir      string
		sampleCSVPath    string
		ibManager        manager.InteractiveBrokersManager
		taxConfig        config.TaxConfig
		mockGainsManager *mocks.GainsComputationManager
		ctx              = context.Background()
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "ib_test_*")
		Expect(err).ToNot(HaveOccurred())

		sampleCSVPath = filepath.Join(tempTestDir, "realized.csv")
		taxConfig = config.TaxConfig{
			TradesPath:       filepath.Join(tempTestDir, "trades.csv"),
			DividendFilePath: filepath.Join(tempTestDir, "dividends.csv"),
			GainsFilePath:    filepath.Join(tempTestDir, "gains.csv"),
			IBPath:           sampleCSVPath,
		}

		mockGainsManager = mocks.NewGainsComputationManager(GinkgoT())
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

			ibManager = manager.NewInteractiveBrokersManager(taxConfig, mockGainsManager)
		})

		Context("when parsing trades", func() {
			It("should extract trade entries correctly", func() {
				info, err := ibManager.Parse()
				Expect(err).ToNot(HaveOccurred())
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
				info, err := ibManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Dividends).To(HaveLen(1))

				Expect(info.Dividends[0].Symbol).To(Equal("MPC"))
				Expect(info.Dividends[0].Date).To(Equal("2024-12-10"))
				Expect(info.Dividends[0].Amount).To(Equal(7.28))
				Expect(info.Dividends[0].Tax).To(Equal(1.82))
				Expect(info.Dividends[0].Net).To(Equal(5.46))
			})
		})

		Context("when generating CSV", func() {
			It("should create valid csv files including gains.csv", func() {
				info := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "MPC", Date: "2024-10-31", Type: "BUY", Quantity: 8, USDPrice: 146.21, USDValue: 1169.68, Commission: 0.36024125},
						{Symbol: "MPC", Date: "2024-12-17", Type: "SELL", Quantity: 8, USDPrice: 136.85, USDValue: 1094.8, Commission: 0.38419669},
					},
					Dividends: []tax.Dividend{
						{Symbol: "MPC", Date: "2024-12-10", Amount: 7.28, Tax: 1.82, Net: 5.46},
					},
				}

				expectedGains := []tax.Gains{
					{Symbol: "MPC", BuyDate: "2024-10-31", SellDate: "2024-12-17", Type: "STCG", Quantity: 8, PNL: -75.62, Commission: 0.74},
				}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(expectedGains, nil)

				err := ibManager.GenerateCsv(ctx, info)
				Expect(err).ToNot(HaveOccurred())

				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("MPC,2024-10-31,BUY,8,146.21"))

				data, err = os.ReadFile(taxConfig.DividendFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("MPC,2024-12-10,7.28,1.82,5.46"))

				data, err = os.ReadFile(taxConfig.GainsFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("MPC,2024-10-31,2024-12-17"))
			})
		})
	})

	Context("with invalid CSV files", func() {
		Context("when CSV file is missing", func() {
			It("should return an error", func() {
				nonExistentManager := manager.NewInteractiveBrokersManager(taxConfig, mockGainsManager)
				_, err := nonExistentManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when CSV has malformed data", func() {
			It("should skip invalid rows gracefully", func() {
				csvContent := `Trades,Header,DataDiscriminator,Asset Category,Currency,Symbol,Date/Time,Quantity,T. Price,Proceeds,Comm/Fee,Basis,Realized P/L,Code
Trades,Data,Order,Stocks,USD,INVALID,,BADQTY,BADPRICE,0,0,0,0,O`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManager(taxConfig, mockGainsManager)
				info, err := ibManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Trades).To(BeEmpty())
			})
		})

		Context("when dividend has no matching tax", func() {
			It("should set tax to 0", func() {
				csvContent := `Dividends,Header,Currency,Date,Description,Amount
Dividends,Data,USD,2024-12-10,MPC(US56585A1025) Cash Dividend USD 0.91 per Share (Ordinary Dividend),7.28`

				err := os.WriteFile(sampleCSVPath, []byte(csvContent), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibManager = manager.NewInteractiveBrokersManager(taxConfig, mockGainsManager)
				info, err := ibManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Dividends).To(HaveLen(1))
				Expect(info.Dividends[0].Tax).To(Equal(0.0))
				Expect(info.Dividends[0].Net).To(Equal(7.28))
			})
		})
	})
})
