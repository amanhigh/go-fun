package manager_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitManager", func() {
	const (
		VO   = "VO"
		AAPL = "AAPL"
		BAD  = "BAD"
	)
	var (
		ctx          = context.Background()
		mockTicker   *mocks.TickerManager
		splitManager manager.SplitManager
	)

	BeforeEach(func() {
		mockTicker = mocks.NewTickerManager(GinkgoT())
		splitManager = manager.NewSplitManager(mockTicker)
	})

	AfterEach(func() {
		mockTicker.AssertExpectations(GinkgoT())
	})

	// -----------------------------------------------------------------------
	// NormalizeTrades
	// -----------------------------------------------------------------------
	Describe("NormalizeTrades", func() {
		Context("VO-like 4:1 inference", func() {
			var (
				input      []tax.Trade
				output     []tax.Trade
				tradeDt    = time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
				yahooClose = 72.555
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: VO, Date: "2025-12-31", Type: "BUY", Quantity: 9, USDPrice: 292.11, USDValue: 2628.99, Commission: 0.33765725},
				}
				mockTicker.EXPECT().
					GetPrice(ctx, VO, tradeDt).
					Return(yahooClose, nil).Once()
				var err common.HttpError
				output, err = splitManager.NormalizeTrades(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should normalize quantity by factor 4", func() {
				Expect(output[0].Quantity).To(Equal(36.0))
			})
			It("should normalize price by factor 1/4", func() {
				Expect(output[0].USDPrice).To(Equal(73.03))
			})
			It("should preserve USDValue", func() {
				Expect(output[0].USDValue).To(Equal(2628.99))
			})
			It("should preserve commission", func() {
				Expect(output[0].Commission).To(Equal(0.33765725))
			})
			It("should preserve symbol, date and type", func() {
				Expect(output[0].Symbol).To(Equal(VO))
				Expect(output[0].Date).To(Equal("2025-12-31"))
				Expect(output[0].Type).To(Equal("BUY"))
			})
		})

		Context("no split", func() {
			var (
				input   []tax.Trade
				output  []tax.Trade
				tradeDt = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: AAPL, Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 150.00, USDValue: 15000, Commission: 10.0},
				}
				mockTicker.EXPECT().
					GetPrice(ctx, AAPL, tradeDt).
					Return(150.50, nil).Once()
				var err common.HttpError
				output, err = splitManager.NormalizeTrades(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an economically and structurally unchanged trade", func() {
				Expect(output[0].Quantity).To(Equal(100.0))
				Expect(output[0].USDPrice).To(Equal(150.00))
				Expect(output[0].USDValue).To(Equal(15000.0))
				Expect(output[0].Commission).To(Equal(10.0))
			})
		})

		Context("ambiguous large mismatch", func() {
			var err common.HttpError

			BeforeEach(func() {
				trades := []tax.Trade{
					{Symbol: BAD, Date: "2024-01-15", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000},
				}
				mockTicker.EXPECT().
					GetPrice(ctx, BAD, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)).
					Return(31.0, nil).Once()
				_, err = splitManager.NormalizeTrades(ctx, trades)
			})

			It("should fail rather than select an unsupported factor", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported split ratio"))
			})
		})

		Context("pre-split buy and post-split sale", func() {
			var (
				trades  []tax.Trade
				results []tax.Trade
			)

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: VO, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000},
					{Symbol: VO, Date: "2024-07-01", Type: "SELL", Quantity: 15, USDPrice: 30.00, USDValue: 450},
				}
				// Pre-split buy: Yahoo close $25 → factor 4
				mockTicker.EXPECT().
					GetPrice(ctx, VO, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)).
					Return(25.0, nil).Once()
				// Post-split sell: Yahoo close $30.50 → ratio ≈ 0.98 → factor 1
				mockTicker.EXPECT().
					GetPrice(ctx, VO, time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)).
					Return(30.50, nil).Once()
				var err common.HttpError
				results, err = splitManager.NormalizeTrades(ctx, trades)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should normalize pre-split buy to 40 @ $25", func() {
				Expect(results[0].Quantity).To(Equal(40.0))
				Expect(results[0].USDPrice).To(Equal(25.0))
				Expect(results[0].USDValue).To(Equal(1000.0))
			})
			It("should leave post-split sell unchanged", func() {
				Expect(results[1].Quantity).To(Equal(15.0))
				Expect(results[1].USDPrice).To(Equal(30.0))
				Expect(results[1].USDValue).To(Equal(450.0))
			})
		})
	})

	// -----------------------------------------------------------------------
	// NormalizeAccount
	// -----------------------------------------------------------------------
	Describe("NormalizeAccount", func() {
		Context("legacy broker-basis account", func() {
			var (
				input    tax.Account
				output   tax.Account
				originDt = time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
			)

			BeforeEach(func() {
				input = tax.Account{
					Symbol:      VO,
					Quantity:    9,
					MarketValue: 653,
					OriginDate:  "2025-12-31",
					OriginQty:   9,
					OriginPrice: 292.11,
				}
				mockTicker.EXPECT().
					GetPrice(ctx, VO, originDt).
					Return(72.555, nil).Once()
				var err common.HttpError
				output, err = splitManager.NormalizeAccount(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should normalize quantity by factor 4", func() {
				Expect(output.Quantity).To(Equal(36.0))
			})
			It("should scale MarketValue by factor 4", func() {
				Expect(output.MarketValue).To(Equal(653.0 * 4))
			})
			It("should normalize origin qty by factor 4", func() {
				Expect(output.OriginQty).To(Equal(36.0))
			})
			It("should normalize origin price by factor 1/4", func() {
				Expect(output.OriginPrice).To(Equal(73.03))
			})
		})

		Context("already normalized account", func() {
			var (
				input  tax.Account
				output tax.Account
			)

			BeforeEach(func() {
				input = tax.Account{
					Symbol:      VO,
					Quantity:    36,
					MarketValue: 2611.98,
					OriginDate:  "2025-12-31",
					OriginQty:   36,
					OriginPrice: 73.0275,
				}
				// Origin price ≈ Yahoo close → factor stays 1
				mockTicker.EXPECT().
					GetPrice(ctx, VO, time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)).
					Return(72.555, nil).Once()
				var err common.HttpError
				output, err = splitManager.NormalizeAccount(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should remain unchanged when normalized again", func() {
				Expect(output.Quantity).To(Equal(36.0))
				Expect(output.MarketValue).To(Equal(2611.98))
				Expect(output.OriginQty).To(Equal(36.0))
				// OriginPrice was already 73.0275 on input — factor ≈ 1, so unchanged
				Expect(output.OriginPrice).To(Equal(73.0275))
			})
		})
	})
})
