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
)

// SplitManager — Yahoo event-based normalization tests.

var _ = Describe("SplitManager", func() {
	const (
		VO      = "VO"
		AAPL    = "AAPL"
		BAD     = "BAD"
		BADDATE = "BADDATE"
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

		// -------- 1. Happy path: two tickers, five trades --------
		Context("happy path", func() {
			var (
				input         []tax.Trade
				expectedInput []tax.Trade // deep copy for immutability check
				output        []tax.Trade
				err           common.HttpError
			)

			BeforeEach(func() {
				input = []tax.Trade{
					{Symbol: VO, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100, USDValue: 1000, Commission: 1.0},
					{Symbol: AAPL, Date: "2024-03-01", Type: "BUY", Quantity: 100, USDPrice: 150, USDValue: 15000, Commission: 10.0},
					{Symbol: VO, Date: "2024-09-01", Type: "BUY", Quantity: 10, USDPrice: 30, USDValue: 300, Commission: 0.5},
					{Symbol: AAPL, Date: "2024-10-15", Type: "SELL", Quantity: 100, USDPrice: 200, USDValue: 20000, Commission: 5.0},
					{Symbol: VO, Date: "2025-06-01", Type: "SELL", Quantity: 5, USDPrice: 60, USDValue: 300},
				}
				expectedInput = append([]tax.Trade(nil), input...)

				// VO: 3:2 split on 2024-09-01, 4:3 split on 2025-01-15
				// Exact GetSplits range: 2024-06-01 through 2025-06-01
				mockTicker.EXPECT().
					GetSplits(ctx, VO, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 3, Denominator: 2},
						{Date: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 3},
					}, nil).Once()

				// AAPL: no split events — trade untouched
				// Exact GetSplits range: 2024-03-01 through 2024-10-15
				mockTicker.EXPECT().
					GetSplits(ctx, AAPL, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)).
					Return([]tax.YahooSplit{}, nil).Once()

				output, err = splitManager.NormalizeTrades(ctx, input)
			})

			It("should normalize trades across all scenarios preserving input order and invariant fields", func() {
				Expect(err).ToNot(HaveOccurred())

				// 0. Original input unchanged (immutability guarantee)
				Expect(input).To(Equal(expectedInput))

				// 1. Original input order preserved
				Expect(output[0].Symbol).To(Equal(VO))
				Expect(output[1].Symbol).To(Equal(AAPL))
				Expect(output[2].Symbol).To(Equal(VO))
				Expect(output[3].Symbol).To(Equal(AAPL))
				Expect(output[4].Symbol).To(Equal(VO))

				// 2. VO pre-split (2024-06-01): cumulative factor 3:2 × 4:3 = 2.0
				Expect(output[0].Quantity).To(Equal(20.0))
				Expect(output[0].USDPrice).To(Equal(50.0))
				Expect(output[0].USDValue).To(Equal(1000.0))
				Expect(output[0].Commission).To(Equal(1.0))
				Expect(output[0].Date).To(Equal("2024-06-01"))
				Expect(output[0].Type).To(Equal("BUY"))

				// 3. AAPL (2024-03-01): no splits — unchanged
				Expect(output[1].Quantity).To(Equal(100.0))
				Expect(output[1].USDPrice).To(Equal(150.0))
				Expect(output[1].USDValue).To(Equal(15000.0))
				Expect(output[1].Commission).To(Equal(10.0))
				Expect(output[1].Date).To(Equal("2024-03-01"))
				Expect(output[1].Type).To(Equal("BUY"))

				// 4. VO same-day (2024-09-01): 3:2 split same day does not apply,
				//    only 4:3 split on 2025-01-15 applies — factor = 4/3
				Expect(output[2].Quantity).To(BeNumerically("~", 10.0*4.0/3.0, 1e-9))
				Expect(output[2].USDPrice).To(Equal(22.5))
				Expect(output[2].USDValue).To(Equal(300.0))
				Expect(output[2].Commission).To(Equal(0.5))
				Expect(output[2].Date).To(Equal("2024-09-01"))
				Expect(output[2].Type).To(Equal("BUY"))

				// 5. AAPL (2024-10-15): no splits — unchanged
				Expect(output[3].Quantity).To(Equal(100.0))
				Expect(output[3].USDPrice).To(Equal(200.0))
				Expect(output[3].USDValue).To(Equal(20000.0))
				Expect(output[3].Commission).To(Equal(5.0))
				Expect(output[3].Date).To(Equal("2024-10-15"))
				Expect(output[3].Type).To(Equal("SELL"))

				// 6. VO post-final (2025-06-01): after both splits — unchanged
				Expect(output[4].Quantity).To(Equal(5.0))
				Expect(output[4].USDPrice).To(Equal(60.0))
				Expect(output[4].USDValue).To(Equal(300.0))
				Expect(output[4].Commission).To(Equal(0.0))
				Expect(output[4].Date).To(Equal("2025-06-01"))
				Expect(output[4].Type).To(Equal("SELL"))
			})
		})

		// -------- 2. Split provider error --------
		Context("split provider error", func() {
			var (
				result []tax.Trade
				err    common.HttpError
			)

			BeforeEach(func() {
				// GetSplits returns a provider-level error
				mockTicker.EXPECT().
					GetSplits(ctx, VO, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)).
					Return(nil, common.NewHttpError("provider unavailable", http.StatusServiceUnavailable)).Once()
				result, err = splitManager.NormalizeTrades(ctx, []tax.Trade{
					{Symbol: VO, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100, USDValue: 1000},
				})
			})

			It("should return nil result and propagate the provider error unchanged", func() {
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusServiceUnavailable))
				Expect(err.Error()).To(Equal("provider unavailable"))
			})
		})

		// -------- 3. Malformed split event data --------
		Context("malformed split event data", func() {
			var err common.HttpError

			BeforeEach(func() {
				// GetSplits succeeds but returns a split with invalid ratio
				// (Denominator=0) — SplitManager must validate.
				mockTicker.EXPECT().
					GetSplits(ctx, BAD, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)).
					Return([]tax.YahooSplit{
						{Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC).Unix(), Numerator: 4, Denominator: 0},
					}, nil).Once()
				_, err = splitManager.NormalizeTrades(ctx, []tax.Trade{
					{Symbol: BAD, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100, USDValue: 1000},
				})
			})

			It("should fail with BadRequest containing ticker context", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring(BAD))
			})
		})

		// -------- 4. Invalid trade date --------
		Context("invalid trade date", func() {
			var err common.HttpError

			BeforeEach(func() {
				// One valid date and one unparseable date for the same ticker.
				// No GetSplits expectation — even though current production
				// fetches splits before validating individual trade dates.
				trades := []tax.Trade{
					{Symbol: BADDATE, Date: "2024-06-01", Type: "BUY", Quantity: 10, USDPrice: 100, USDValue: 1000},
					{Symbol: BADDATE, Date: "not-a-date", Type: "BUY", Quantity: 5, USDPrice: 50, USDValue: 250},
				}
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
