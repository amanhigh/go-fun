package manager_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func mockError(message string) common.HttpError {
	return common.NewHttpError(message, http.StatusBadRequest)
}

var _ = Describe("BrokerageManager", func() {
	var (
		tempTestDir      string
		brokerageManager manager.BrokerageManager
		taxConfig        config.TaxConfig
		mockGainsManager *mocks.GainsComputationManager
		mockDWBroker     *mocks.Broker
		mockIBBroker     *mocks.Broker
		ctx              = context.Background()
		emptyInfo        tax.BrokerageInfo
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "brokerage_manager_test_*")
		Expect(err).ToNot(HaveOccurred())

		taxConfig = config.TaxConfig{
			TradesPath:       filepath.Join(tempTestDir, "trades.csv"),
			DividendFilePath: filepath.Join(tempTestDir, "dividends.csv"),
			GainsFilePath:    filepath.Join(tempTestDir, "gains.csv"),
			InterestFilePath: filepath.Join(tempTestDir, "interest.csv"),
		}

		mockGainsManager = mocks.NewGainsComputationManager(GinkgoT())
		mockDWBroker = mocks.NewBroker(GinkgoT())
		mockIBBroker = mocks.NewBroker(GinkgoT())

		mockDWBroker.EXPECT().GetName().Return("DriveWealth").Maybe()
		mockIBBroker.EXPECT().GetName().Return("InteractiveBrokers").Maybe()

		emptyInfo = tax.BrokerageInfo{}

		brokerageManager = manager.NewBrokerageManager(mockDWBroker, mockIBBroker, mockGainsManager, taxConfig)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("ParseAndGenerate", func() {
		Context("with both DriveWealth and IB returning data", func() {
			BeforeEach(func() {
				dwInfo := tax.BrokerageInfo{
					Interests: []tax.Interest{
						{Symbol: "CASH", Date: "2024-01-15", Amount: 10.50, Tax: 0, Net: 10.50},
					},
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
					Dividends: []tax.Dividend{
						{Symbol: "MSFT", Date: "2024-01-20", Amount: 50.0, Tax: 7.5, Net: 42.5},
					},
				}

				ibInfo := tax.BrokerageInfo{
					Interests: []tax.Interest{
						{Symbol: "CASH", Date: "2024-02-15", Amount: 5.25, Tax: 0, Net: 5.25},
					},
					Trades: []tax.Trade{
						{Symbol: "GOOGL", Date: "2024-02-10", Type: "BUY", Quantity: 5, USDPrice: 120.0, USDValue: 600, Commission: 0.5},
					},
					Dividends: []tax.Dividend{
						{Symbol: "TSLA", Date: "2024-02-20", Amount: 30.0, Tax: 4.5, Net: 25.5},
					},
				}

				mockDWBroker.EXPECT().Parse().Return(dwInfo, nil)
				mockIBBroker.EXPECT().Parse().Return(ibInfo, nil)

				mergedTrades := append(dwInfo.Trades, ibInfo.Trades...)
				expectedGains := []tax.Gains{
					{Symbol: "AAPL", BuyDate: "2024-01-10", SellDate: "2024-02-10", Type: "STCG", Quantity: 10, PNL: 98, Commission: 1.5},
				}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, mergedTrades).Return(expectedGains, nil)

				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should merge and write interest data from both brokers", func() {
				data, err := os.ReadFile(taxConfig.InterestFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("CASH,2024-01-15,10.5,0,10.5"))
				Expect(string(data)).To(ContainSubstring("CASH,2024-02-15,5.25,0,5.25"))
			})

			It("should merge and write trade data from both brokers", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,BUY"))
				Expect(string(data)).To(ContainSubstring("GOOGL,2024-02-10,BUY"))
			})

			It("should merge and write dividend data from both brokers", func() {
				data, err := os.ReadFile(taxConfig.DividendFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("MSFT,2024-01-20,50,7.5,42.5"))
				Expect(string(data)).To(ContainSubstring("TSLA,2024-02-20,30,4.5,25.5"))
			})

			It("should create gains file with merged trades", func() {
				data, err := os.ReadFile(taxConfig.GainsFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,2024-02-10"))
			})
		})

		Context("with one broker failing and one succeeding", func() {
			BeforeEach(func() {
				dwInfo := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
				}

				mockDWBroker.EXPECT().Parse().Return(dwInfo, nil)
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, mockError("file not found"))

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, dwInfo.Trades).Return(expectedGains, nil)

				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should only use data from successful broker", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,BUY"))
			})
		})

		Context("with all brokers failing", func() {
			BeforeEach(func() {
				mockDWBroker.EXPECT().Parse().Return(emptyInfo, mockError("file not found"))
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, mockError("file not found"))
			})

			It("should return error for no data", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no data found"))
			})
		})

		Context("with empty data from all brokers", func() {
			BeforeEach(func() {
				mockDWBroker.EXPECT().Parse().Return(emptyInfo, nil)
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, nil)
			})

			It("should return error for empty data", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no data found"))
			})
		})

		Context("when gains computation fails", func() {
			BeforeEach(func() {
				dwInfo := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "SELL", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
				}

				mockDWBroker.EXPECT().Parse().Return(dwInfo, nil)
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, nil)
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, dwInfo.Trades).Return(nil, mockError("no buy positions"))
			})

			It("should return gains computation error", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Trade Ordering", func() {
		var (
			unorderedTrades tax.BrokerageInfo
		)

		BeforeEach(func() {
			// Trades in wrong order: SELL before BUY on same date
			unorderedTrades = tax.BrokerageInfo{
				Trades: []tax.Trade{
					{Symbol: "AAPL", Date: "2024-08-05", Type: "SELL", Quantity: 6, USDPrice: 209, USDValue: 1254},
					{Symbol: "AAPL", Date: "2024-08-05", Type: "BUY", Quantity: 6, USDPrice: 199, USDValue: 1194},
					{Symbol: "MSFT", Date: "2024-08-01", Type: "BUY", Quantity: 10, USDPrice: 100, USDValue: 1000},
				},
			}
		})

		Context("when trades are in wrong order (SELL before BUY on same date)", func() {
			BeforeEach(func() {
				mockDWBroker.EXPECT().Parse().Return(unorderedTrades, nil)
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, nil)
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, mock.Anything).Return([]tax.Gains{}, nil)
			})

			It("should sort trades by date then type (BUY before SELL)", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Read generated trades file
				trades := readTradesCSV(taxConfig.TradesPath)
				Expect(trades).To(HaveLen(3))

				// Verify ordering: BUY should come before SELL on same date
				Expect(trades[0].Symbol).To(Equal("MSFT"))
				Expect(trades[0].Date).To(Equal("2024-08-01"))
				Expect(trades[0].Type).To(Equal("BUY"))

				Expect(trades[1].Symbol).To(Equal("AAPL"))
				Expect(trades[1].Date).To(Equal("2024-08-05"))
				Expect(trades[1].Type).To(Equal("BUY")) // BUY first

				Expect(trades[2].Symbol).To(Equal("AAPL"))
				Expect(trades[2].Date).To(Equal("2024-08-05"))
				Expect(trades[2].Type).To(Equal("SELL")) // SELL second
			})
		})
	})
})

func readTradesCSV(path string) []tax.Trade {
	file, err := os.Open(path)
	Expect(err).ToNot(HaveOccurred())
	defer file.Close()

	var trades []tax.Trade
	err = gocsv.UnmarshalFile(file, &trades)
	Expect(err).ToNot(HaveOccurred())

	return trades
}
