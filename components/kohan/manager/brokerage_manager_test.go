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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		dwPath           string
		ibPath           string
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "brokerage_manager_test_*")
		Expect(err).ToNot(HaveOccurred())

		dwPath = filepath.Join(tempTestDir, "drivewealth.xlsx")
		ibPath = filepath.Join(tempTestDir, "ib_realized.csv")

		taxConfig = config.TaxConfig{
			TradesPath:       filepath.Join(tempTestDir, "trades.csv"),
			DividendFilePath: filepath.Join(tempTestDir, "dividends.csv"),
			GainsFilePath:    filepath.Join(tempTestDir, "gains.csv"),
			InterestFilePath: filepath.Join(tempTestDir, "interest.csv"),
			DriveWealthPath:  dwPath,
			IBPath:           ibPath,
		}

		mockGainsManager = mocks.NewGainsComputationManager(GinkgoT())
		mockDWBroker = mocks.NewBroker(GinkgoT())
		mockIBBroker = mocks.NewBroker(GinkgoT())

		brokerageManager = manager.NewBrokerageManager(mockDWBroker, mockIBBroker, mockGainsManager, taxConfig)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("ParseAndGenerate", func() {
		Context("with both DriveWealth and IB files present", func() {
			BeforeEach(func() {
				err := os.WriteFile(dwPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(ibPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

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

				err = brokerageManager.ParseAndGenerate(ctx)
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

		//nolint:dupl
		Context("with only DriveWealth file present", func() {
			BeforeEach(func() {
				err := os.WriteFile(dwPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

				dwInfo := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
				}

				mockDWBroker.EXPECT().Parse().Return(dwInfo, nil)

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, dwInfo.Trades).Return(expectedGains, nil)

				err = brokerageManager.ParseAndGenerate(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should only parse DriveWealth and skip IB", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,BUY"))
			})
		})

		Context("with only IB file present", func() {
			BeforeEach(func() {
				err := os.WriteFile(ibPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

				ibInfo := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "GOOGL", Date: "2024-02-10", Type: "BUY", Quantity: 5, USDPrice: 120.0, USDValue: 600, Commission: 0.5},
					},
				}

				mockIBBroker.EXPECT().Parse().Return(ibInfo, nil)

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, ibInfo.Trades).Return(expectedGains, nil)

				err = brokerageManager.ParseAndGenerate(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should only parse IB and skip DriveWealth", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("GOOGL,2024-02-10,BUY"))
			})
		})

		Context("with no broker files present", func() {
			It("should return error", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no broker files found"))
			})
		})

		Context("when DriveWealth parsing fails", func() {
			BeforeEach(func() {
				err := os.WriteFile(dwPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

				mockDWBroker.EXPECT().Parse().Return(tax.BrokerageInfo{}, mockError("parse error"))
			})

			It("should return parse error", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("parse failed"))
			})
		})

		Context("when gains computation fails", func() {
			BeforeEach(func() {
				err := os.WriteFile(dwPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

				dwInfo := tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "SELL", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
				}

				mockDWBroker.EXPECT().Parse().Return(dwInfo, nil)
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, dwInfo.Trades).Return(nil, mockError("no buy positions"))
			})

			It("should return gains computation error", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with empty data from all brokers", func() {
			BeforeEach(func() {
				err := os.WriteFile(dwPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())
				err = os.WriteFile(ibPath, []byte("dummy"), 0600)
				Expect(err).ToNot(HaveOccurred())

				emptyInfo := tax.BrokerageInfo{}
				mockDWBroker.EXPECT().Parse().Return(emptyInfo, nil)
				mockIBBroker.EXPECT().Parse().Return(emptyInfo, nil)
			})

			It("should return error for empty data", func() {
				err := brokerageManager.ParseAndGenerate(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no broker files found or all files empty"))
			})
		})
	})
})
