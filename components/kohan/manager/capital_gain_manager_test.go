package manager_test

import (
	"context"
	"time"

	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CapitalGainManager", func() {
	var (
		ctx            = context.Background()
		mockSBIManager *mocks.SBIManager
		gainManager    manager.CapitalGainManager

		// Common test data
		ticker    = "AAPL"
		sellDate  = "2024-01-15"
		ttBuyRate = 82.50
	)

	BeforeEach(func() {
		mockSBIManager = mocks.NewSBIManager(GinkgoT())
		gainManager = manager.NewCapitalGainManager(mockSBIManager)
	})

	Context("Basic Gains Processing", func() {
		var (
			gains []tax.Gains
			pnl   = 1000.00
		)

		BeforeEach(func() {
			gains = []tax.Gains{
				{
					Symbol:   ticker,
					SellDate: sellDate,
					PNL:      pnl,
				},
			}

			parsedDate, _ := time.Parse(time.DateOnly, sellDate)
			mockSBIManager.EXPECT().
				GetTTBuyRate(parsedDate).
				Return(ttBuyRate, nil)
		})

		It("should process gain with correct INR values", func() {
			taxGains, err := gainManager.ProcessTaxGains(ctx, gains)
			Expect(err).To(BeNil())
			Expect(taxGains).To(HaveLen(1))

			result := taxGains[0]
			Expect(result.Symbol).To(Equal(ticker))
			Expect(result.TTRate).To(Equal(ttBuyRate))
			Expect(result.INRValue()).To(Equal(pnl * ttBuyRate))
		})

		// BUG: Test where exact Date is not found TTDate differs from sellDate.
	})

	Context("Error Cases", func() {
		var gains []tax.Gains
		Context("Invalid Date", func() {
			BeforeEach(func() {
				gains = []tax.Gains{{
					Symbol:   ticker,
					SellDate: "invalid-date",
					PNL:      1000.00,
				}}
			})

			It("should return error for invalid date format", func() {
				_, err := gainManager.ProcessTaxGains(ctx, gains)
				Expect(err).Should(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
			})
		})

		Context("Exchange Rate Error", func() {
			BeforeEach(func() {
				gains = []tax.Gains{{
					Symbol:   ticker,
					SellDate: sellDate,
					PNL:      1000.00,
				}}

				parsedDate, _ := time.Parse(time.DateOnly, sellDate)
				mockSBIManager.EXPECT().
					GetTTBuyRate(parsedDate).
					Return(0.0, common.ErrNotFound)
			})

			It("should handle missing exchange rate", func() {
				_, err := gainManager.ProcessTaxGains(ctx, gains)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})

		Context("Multiple Gains Processing", func() {
			var (
				gains []tax.Gains
				dates = []string{sellDate, "2024-01-16"}
			)

			BeforeEach(func() {
				gains = []tax.Gains{
					{Symbol: ticker, SellDate: dates[0], PNL: 1000.00},
					{Symbol: ticker, SellDate: dates[1], PNL: 2000.00},
				}

				// Setup expectations for both dates
				for _, dateStr := range dates {
					parsedDate, _ := time.Parse(time.DateOnly, dateStr)
					mockSBIManager.EXPECT().
						GetTTBuyRate(parsedDate).
						Return(ttBuyRate, nil)
				}
			})

			It("should process multiple gains correctly", func() {
				taxGains, err := gainManager.ProcessTaxGains(ctx, gains)
				Expect(err).To(BeNil())
				Expect(taxGains).To(HaveLen(2))

				// Verify each gain processed correctly
				for i, gain := range taxGains {
					Expect(gain.Symbol).To(Equal(ticker))
					Expect(gain.SellDate).To(Equal(dates[i]))
					Expect(gain.TTRate).To(Equal(ttBuyRate))
					Expect(gain.INRValue()).To(Equal(gain.PNL * ttBuyRate))
				}
			})
		})
	})
})
