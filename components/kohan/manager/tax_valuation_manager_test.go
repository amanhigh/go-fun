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

var _ = Describe("Tax Valuation Manager", func() {
	var (
		mockSBIManager      *mocks.SBIManager
		taxValuationManager manager.TaxValuationManager
		ctx                 = context.Background()
	)

	BeforeEach(func() {
		mockSBIManager = mocks.NewSBIManager(GinkgoT())
		taxValuationManager = manager.NewTaxValuationManager(mockSBIManager)
	})

	Context("Single Valuation", func() {
		var (
			valuation    tax.Valuation
			taxValuation []tax.TaxValuation
			err          common.HttpError

			// Test data
			ticker      = "AAPL"
			firstDate   = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			peakDate    = time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
			yearEndDate = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

			// Exchange rates
			firstRate   = 82.50
			peakRate    = 83.25
			yearEndRate = 84.00
		)

		BeforeEach(func() {
			// Setup test valuation
			valuation = tax.Valuation{
				Ticker: ticker,
				FirstPosition: tax.Position{
					Date:     firstDate,
					Quantity: 100,
					USDPrice: 150,
				},
				PeakPosition: tax.Position{
					Date:     peakDate,
					Quantity: 150,
					USDPrice: 160,
				},
				YearEndPosition: tax.Position{
					Date:     yearEndDate,
					Quantity: 120,
					USDPrice: 155,
				},
			}

		})

		Context("Full Valuation", func() {
			BeforeEach(func() {
				// Setup mock expectations
				mockSBIManager.EXPECT().
					GetTTBuyRate(firstDate).
					Return(firstRate, nil)
				mockSBIManager.EXPECT().
					GetTTBuyRate(peakDate).
					Return(peakRate, nil)
				mockSBIManager.EXPECT().
					GetTTBuyRate(yearEndDate).
					Return(yearEndRate, nil)
			})

			It("should process valuation successfully", func() {
				taxValuation, err = taxValuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
				Expect(err).To(BeNil())
				Expect(taxValuation).To(HaveLen(1))

				result := taxValuation[0]
				Expect(result.Ticker).To(Equal(ticker))

				// Verify first position
				Expect(result.FirstPosition.TTDate).To(Equal(firstDate))
				Expect(result.FirstPosition.TTBuyRate).To(Equal(firstRate))
				Expect(result.FirstPosition.USDValue()).To(Equal(valuation.FirstPosition.USDValue()))
				Expect(result.FirstPosition.INRValue()).To(Equal(valuation.FirstPosition.USDValue() * firstRate))

				// Verify peak position
				Expect(result.PeakPosition.TTDate).To(Equal(peakDate))
				Expect(result.PeakPosition.TTBuyRate).To(Equal(peakRate))
				Expect(result.PeakPosition.USDValue()).To(Equal(valuation.PeakPosition.USDValue()))
				Expect(result.PeakPosition.INRValue()).To(Equal(valuation.PeakPosition.USDValue() * peakRate))

				// Verify year end position
				Expect(result.YearEndPosition.TTDate).To(Equal(yearEndDate))
				Expect(result.YearEndPosition.TTBuyRate).To(Equal(yearEndRate))
				Expect(result.YearEndPosition.USDValue()).To(Equal(valuation.YearEndPosition.USDValue()))
				Expect(result.YearEndPosition.INRValue()).To(Equal(valuation.YearEndPosition.USDValue() * yearEndRate))
			})
		})

		Context("Empty Position", func() {
			BeforeEach(func() {
				// Set year end position to empty
				valuation.YearEndPosition = tax.Position{}

				// Update mock expectations (only first and peak needed)
				mockSBIManager.EXPECT().
					GetTTBuyRate(firstDate).
					Return(firstRate, nil)
				mockSBIManager.EXPECT().
					GetTTBuyRate(peakDate).
					Return(peakRate, nil)
			})

			It("should handle empty position", func() {
				taxValuation, err = taxValuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
				Expect(err).To(BeNil())
				Expect(taxValuation).To(HaveLen(1))

				result := taxValuation[0]
				// Verify year end position is empty
				Expect(result.YearEndPosition.Quantity).To(Equal(0.0))
				Expect(result.YearEndPosition.USDValue()).To(Equal(0.0))
				Expect(result.YearEndPosition.INRValue()).To(Equal(0.0))
			})
		})
	})

	Context("Error Cases", func() {
		var (
			valuation tax.Valuation
			err       common.HttpError
		)

		BeforeEach(func() {
			valuation = tax.Valuation{
				Ticker: "AAPL",
				FirstPosition: tax.Position{
					Date:     time.Now(),
					Quantity: 100,
					USDPrice: 150,
				},
			}
		})

		It("should handle missing exchange rate", func() {
			mockSBIManager.EXPECT().
				GetTTBuyRate(valuation.FirstPosition.Date).
				Return(0.0, common.ErrNotFound)

			_, err = taxValuationManager.ProcessValuations(ctx, []tax.Valuation{valuation})
			Expect(err).To(Equal(common.ErrNotFound))
		})
	})

	Context("Multiple Valuations", func() {
		var (
			valuations []tax.Valuation
			date       = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			rate       = 82.50
		)

		BeforeEach(func() {
			// Create two simple valuations
			valuations = []tax.Valuation{
				{
					Ticker: "AAPL",
					FirstPosition: tax.Position{
						Date:     date,
						Quantity: 100,
						USDPrice: 150,
					},
				},
				{
					Ticker: "MSFT",
					FirstPosition: tax.Position{
						Date:     date,
						Quantity: 50,
						USDPrice: 200,
					},
				},
			}

			// Same date should reuse rate
			mockSBIManager.EXPECT().
				GetTTBuyRate(date).
				Return(rate, nil).
				Times(2)
		})

		It("should process multiple valuations", func() {
			taxValuations, err := taxValuationManager.ProcessValuations(ctx, valuations)
			Expect(err).To(BeNil())
			Expect(taxValuations).To(HaveLen(2))

			// Verify both valuations processed
			Expect(taxValuations[0].Ticker).To(Equal("AAPL"))
			Expect(taxValuations[1].Ticker).To(Equal("MSFT"))

			// Verify rate applied correctly to both
			Expect(taxValuations[0].FirstPosition.TTBuyRate).To(Equal(rate))
			Expect(taxValuations[1].FirstPosition.TTBuyRate).To(Equal(rate))
		})
	})
})
