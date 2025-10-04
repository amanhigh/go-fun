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

var _ = Describe("BrokerageParserBase", func() {
	var (
		tempTestDir      string
		base             manager.BrokerageParserBase
		taxConfig        config.TaxConfig
		mockGainsManager *mocks.GainsComputationManager
		ctx              = context.Background()
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "brokerage_base_test_*")
		Expect(err).ToNot(HaveOccurred())

		taxConfig = config.TaxConfig{
			TradesPath:       filepath.Join(tempTestDir, "trades.csv"),
			DividendFilePath: filepath.Join(tempTestDir, "dividends.csv"),
			GainsFilePath:    filepath.Join(tempTestDir, "gains.csv"),
			InterestFilePath: filepath.Join(tempTestDir, "interest.csv"),
		}

		mockGainsManager = mocks.NewGainsComputationManager(GinkgoT())
		base = manager.NewBrokerageParserBase(taxConfig, mockGainsManager)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("GenerateCsv", func() {
		Context("with complete brokerage info", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				info = tax.BrokerageInfo{
					Interests: []tax.Interest{
						{Symbol: "CASH", Date: "2024-01-15", Amount: 10.50, Tax: 0, Net: 10.50},
					},
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
						{Symbol: "AAPL", Date: "2024-02-10", Type: "SELL", Quantity: 10, USDPrice: 160.0, USDValue: 1600, Commission: 1.0},
					},
					Dividends: []tax.Dividend{
						{Symbol: "MSFT", Date: "2024-01-20", Amount: 50.0, Tax: 7.5, Net: 42.5},
					},
				}

				expectedGains := []tax.Gains{
					{Symbol: "AAPL", BuyDate: "2024-01-10", SellDate: "2024-02-10", Type: "STCG", Quantity: 10, PNL: 98, Commission: 2},
				}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(expectedGains, nil)

				err := base.GenerateCsv(ctx, info)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create interest file with correct data", func() {
				data, err := os.ReadFile(taxConfig.InterestFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("CASH,2024-01-15,10.5,0,10.5"))
			})

			It("should create trades file with correct data", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,BUY"))
				Expect(string(data)).To(ContainSubstring("AAPL,2024-02-10,SELL"))
			})

			It("should create dividends file with correct data", func() {
				data, err := os.ReadFile(taxConfig.DividendFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("MSFT,2024-01-20,50,7.5,42.5"))
			})

			It("should create gains file with correct data", func() {
				data, err := os.ReadFile(taxConfig.GainsFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("AAPL,2024-01-10,2024-02-10"))
			})
		})

		Context("with empty interest data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				info = tax.BrokerageInfo{
					Interests: []tax.Interest{},
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
					Dividends: []tax.Dividend{},
				}

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(expectedGains, nil)

				err := base.GenerateCsv(ctx, info)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create empty interest file with headers only", func() {
				data, err := os.ReadFile(taxConfig.InterestFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("Symbol"))
			})

			It("should create trades file", func() {
				_, err := os.Stat(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with empty trades data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				info = tax.BrokerageInfo{
					Interests: []tax.Interest{},
					Trades:    []tax.Trade{},
					Dividends: []tax.Dividend{},
				}

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(expectedGains, nil)

				err := base.GenerateCsv(ctx, info)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create empty trades file with headers only", func() {
				data, err := os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("Symbol"))
			})
		})

		Context("with empty dividends data", func() {
			var info tax.BrokerageInfo

			BeforeEach(func() {
				info = tax.BrokerageInfo{
					Interests: []tax.Interest{},
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "BUY", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
					Dividends: []tax.Dividend{},
				}

				expectedGains := []tax.Gains{}
				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(expectedGains, nil)

				err := base.GenerateCsv(ctx, info)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create empty dividends file with headers only", func() {
				data, err := os.ReadFile(taxConfig.DividendFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("Symbol"))
			})
		})

		Context("when gains computation fails", func() {
			var info tax.BrokerageInfo
			var csvErr error

			BeforeEach(func() {
				info = tax.BrokerageInfo{
					Trades: []tax.Trade{
						{Symbol: "AAPL", Date: "2024-01-10", Type: "SELL", Quantity: 10, USDPrice: 150.0, USDValue: 1500, Commission: 1.0},
					},
				}

				mockGainsManager.EXPECT().ComputeGainsFromTrades(ctx, info.Trades).Return(nil, mockError("no buy positions"))

				csvErr = base.GenerateCsv(ctx, info)
			})

			It("should return error from gains manager", func() {
				Expect(csvErr).To(HaveOccurred())
			})
		})
	})

	Describe("MatchDividendWithTax", func() {
		var (
			taxMap   map[string]map[string]float64
			dividend *tax.Dividend
		)

		BeforeEach(func() {
			taxMap = map[string]map[string]float64{
				"MSFT": {
					"2024-01-20": 7.5,
					"2024-02-20": 8.0,
				},
				"AAPL": {
					"2024-01-15": 5.0,
				},
			}
		})

		Context("when tax exists for dividend", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 50.0,
				}
				base.MatchDividendWithTax(dividend, taxMap)
			})

			It("should match tax and calculate net", func() {
				Expect(dividend.Tax).To(Equal(7.5))
				Expect(dividend.Net).To(Equal(42.5))
			})

			It("should remove tax from pool", func() {
				_, exists := taxMap["MSFT"]["2024-01-20"]
				Expect(exists).To(BeFalse())
			})
		})

		Context("when no tax exists for dividend", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "GOOGL",
					Date:   "2024-01-20",
					Amount: 100.0,
				}
				base.MatchDividendWithTax(dividend, taxMap)
			})

			It("should set tax to 0 and net equals amount", func() {
				Expect(dividend.Tax).To(Equal(0.0))
				Expect(dividend.Net).To(Equal(100.0))
			})
		})

		Context("when symbol exists but date doesn't match", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-03-20",
					Amount: 50.0,
				}
				base.MatchDividendWithTax(dividend, taxMap)
			})

			It("should set tax to 0", func() {
				Expect(dividend.Tax).To(Equal(0.0))
				Expect(dividend.Net).To(Equal(50.0))
			})
		})

		Context("when multiple dividends share same tax pool", func() {
			var dividend2 *tax.Dividend

			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 30.0,
				}
				dividend2 = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 20.0,
				}
				base.MatchDividendWithTax(dividend, taxMap)
				base.MatchDividendWithTax(dividend2, taxMap)
			})

			It("should match first dividend and remove tax from pool", func() {
				Expect(dividend.Tax).To(Equal(7.5))
				Expect(dividend.Net).To(Equal(22.5))
			})

			It("should not match second dividend (tax already consumed)", func() {
				Expect(dividend2.Tax).To(Equal(0.0))
				Expect(dividend2.Net).To(Equal(20.0))
			})
		})
	})
})
