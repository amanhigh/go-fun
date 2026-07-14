package manager_test

import (
	"context"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mock "github.com/stretchr/testify/mock"
)

// SplitManager — Yahoo event-based normalization tests.

var _ = Describe("SplitManager", func() {
	const (
		VO      = "VO"
		AAPL    = "AAPL"
		BAD     = "BAD"
		BADDATE = "BADDATE"
		COMP    = "COMP"
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

	// ===================================================================
	// NormalizeTrades — Yahoo event-based normalization
	// ===================================================================
	Describe("NormalizeTrades", func() {

		// -------- 1. forward split: pre-split trade --------
		Context("forward split: pre-split trade", func() {
			var (
				input  []tax.Trade
				output []tax.Trade
				err    common.HttpError
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: VO, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000, Commission: 1.0},
				}
				// 4:1 split on 2024-09-01 — trade is before the event
				mockTicker.EXPECT().
					GetSplits(ctx, VO, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 1},
					}, nil).Once()
				output, err = splitManager.NormalizeTrades(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should multiply pre-split quantity by 4", func() {
				Expect(output[0].Quantity).To(Equal(40.0))
			})
			It("should divide pre-split price by 4", func() {
				Expect(output[0].USDPrice).To(Equal(25.0))
			})
			It("should preserve USDValue", func() {
				Expect(output[0].USDValue).To(Equal(1000.0))
			})
			It("should preserve commission", func() {
				Expect(output[0].Commission).To(Equal(1.0))
			})
			It("should preserve symbol, date and type", func() {
				Expect(output[0].Symbol).To(Equal(VO))
				Expect(output[0].Date).To(Equal("2024-06-01"))
				Expect(output[0].Type).To(Equal("BUY"))
			})
		})

		// -------- 2. no split events --------
		Context("no split events", func() {
			var (
				input  []tax.Trade
				output []tax.Trade
				err    common.HttpError
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: AAPL, Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 150.00, USDValue: 15000, Commission: 10.0},
				}
				mockTicker.EXPECT().
					GetSplits(ctx, AAPL, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{}, nil).Once()
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

		// -------- 3. trade on event date (unchanged) --------
		Context("trade on event date", func() {
			var (
				input  []tax.Trade
				output []tax.Trade
				err    common.HttpError
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: VO, Date: "2024-09-01", Type: "BUY", Quantity: 10, USDPrice: 25.00, USDValue: 250, Commission: 0.50},
				}
				// 4:1 split on same day — trade is on/after event date → unchanged
				mockTicker.EXPECT().
					GetSplits(ctx, VO, time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 1},
					}, nil).Once()
				output, err = splitManager.NormalizeTrades(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should leave the trade unchanged", func() {
				Expect(output[0].Quantity).To(Equal(10.0))
				Expect(output[0].USDPrice).To(Equal(25.0))
				Expect(output[0].USDValue).To(Equal(250.0))
				Expect(output[0].Commission).To(Equal(0.50))
			})
		})

		// -------- 4. compound forward and reverse event --------
		Context("compound forward and reverse event", func() {
			var (
				input  []tax.Trade
				output []tax.Trade
				err    common.HttpError
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: COMP, Date: "2024-03-01", Type: "BUY", Quantity: 100, USDPrice: 200.00, USDValue: 20000, Commission: 2.0},
				}
				// 4:1 forward then 1:4 reverse — cumulative ratio = 4 * 0.25 = 1
				mockTicker.EXPECT().
					GetSplits(ctx, COMP, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 1},
						{Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 1, Denominator: 4},
					}, nil).Once()
				output, err = splitManager.NormalizeTrades(ctx, input)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should net to unchanged when forward and reverse compound to unity", func() {
				Expect(output[0].Quantity).To(Equal(100.0))
				Expect(output[0].USDPrice).To(Equal(200.0))
				Expect(output[0].USDValue).To(Equal(20000.0))
				Expect(output[0].Commission).To(Equal(2.0))
			})
		})

		// -------- 5. FIFO gain across forward split --------
		Context("FIFO gain across forward split", func() {
			var (
				trades   []tax.Trade
				results  []tax.Trade
				gains    []tax.Gains
				gainsErr common.HttpError
			)

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: VO, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000},
					{Symbol: VO, Date: "2024-10-01", Type: "SELL", Quantity: 15, USDPrice: 30.00, USDValue: 450},
				}
				// 4:1 split between buy and sell dates
				mockTicker.EXPECT().
					GetSplits(ctx, VO, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 1},
					}, nil).Once()
				var err common.HttpError
				results, err = splitManager.NormalizeTrades(ctx, trades)
				Expect(err).ToNot(HaveOccurred())

				gainsManager := manager.NewGainsComputationManager()
				gains, gainsErr = gainsManager.ComputeGainsFromTrades(ctx, results)
				Expect(gainsErr).ToNot(HaveOccurred())
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
			It("should produce a FIFO gain of $75 when passed through GainsComputationManager", func() {
				Expect(gains).To(HaveLen(1))
				Expect(gains[0].Quantity).To(Equal(15.0))
				Expect(gains[0].BuyDate).To(Equal("2024-06-01"))
				Expect(gains[0].SellDate).To(Equal("2024-10-01"))
				// PNL = 15 x ($30 - $25) = $75
				Expect(gains[0].PNL).To(Equal(75.0))
				Expect(gains[0].Commission).To(Equal(0.0))
			})
		})

		// -------- 6. malformed split event data --------
		Context("malformed split event data", func() {
			var err common.HttpError

			BeforeEach(func() {
				trades := []tax.Trade{
					{Symbol: BAD, Date: "2024-01-15", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000},
				}
				// GetSplits succeeds but returns a split with invalid ratio
				// (Denominator=0) — requires SplitManager to validate.
				mockTicker.EXPECT().
					GetSplits(ctx, BAD, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 0},
					}, nil).Once()
				_, err = splitManager.NormalizeTrades(ctx, trades)
			})

			It("should fail with BadRequest containing ticker context", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring(BAD))
			})
		})

		// -------- 9. invalid later trade date --------
		Context("invalid later trade date", func() {
			var err common.HttpError

			BeforeEach(func() {
				trades := []tax.Trade{
					{Symbol: BADDATE, Date: "2024-01-15", Type: "BUY", Quantity: 10, USDPrice: 100.00, USDValue: 1000},
					{Symbol: BADDATE, Date: "not-a-date", Type: "BUY", Quantity: 5, USDPrice: 50.00, USDValue: 250},
				}
				mockTicker.EXPECT().
					GetSplits(ctx, BADDATE, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), mock.Anything).
					Return([]tax.YahooSplit{}, nil).Once()
				_, err = splitManager.NormalizeTrades(ctx, trades)
			})

			It("should fail with BadRequest containing ticker context and the invalid date string", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring(BADDATE))
				Expect(err.Error()).To(ContainSubstring("not-a-date"))
			})
		})
	})

})
