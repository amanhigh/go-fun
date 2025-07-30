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

var _ = Describe("ExchangeManager", func() {
	var (
		ctx         context.Context
		mockSBI     *mocks.SBIManager
		exchangeMgr manager.ExchangeManager
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockSBI = mocks.NewSBIManager(GinkgoT())
		exchangeMgr = manager.NewExchangeManager(mockSBI)
	})

	Context("Single Exchange", func() {
		var (
			exchangeables []tax.Exchangeable
			testDate      time.Time
			position      tax.INRPosition
		)

		BeforeEach(func() {
			testDate = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			position = tax.INRPosition{
				Position: tax.Position{
					Date:     testDate,
					Quantity: 100,
					USDPrice: 150,
				},
			}
			exchangeables = []tax.Exchangeable{&position}

			mockSBI.EXPECT().
				GetTTBuyRate(ctx, testDate).
				Return(82.0, nil)
		})

		It("should process single position exchange", func() {
			err := exchangeMgr.Exchange(ctx, exchangeables)

			Expect(err).ToNot(HaveOccurred())
			Expect(position.TTRate).To(Equal(82.0))
			Expect(position.TTDate).To(Equal(testDate))
		})
	})

	Context("Multiple Exchanges", func() {
		var (
			exchangeables []tax.Exchangeable
			positions     [3]tax.INRPosition
			dates         [3]time.Time
		)

		BeforeEach(func() {
			baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			for i := range dates {
				dates[i] = baseDate.AddDate(0, i, 0)
				positions[i] = tax.INRPosition{
					Position: tax.Position{
						Date:     dates[i],
						Quantity: 100 * float64(i+1),
						USDPrice: 150 + float64(i*10),
					},
				}
			}

			exchangeables = make([]tax.Exchangeable, 3)
			for i := range positions {
				exchangeables[i] = &positions[i]
				mockSBI.EXPECT().
					GetTTBuyRate(ctx, dates[i]).
					Return(82.0+float64(i), nil)
			}
		})

		It("should process multiple position exchanges", func() {
			err := exchangeMgr.Exchange(ctx, exchangeables)

			Expect(err).ToNot(HaveOccurred())
			for i, position := range positions {
				Expect(position.TTRate).To(Equal(82.0 + float64(i)))
				Expect(position.TTDate).To(Equal(dates[i]))
			}
		})
	})

	Context("Error Handling", func() {
		var (
			exchangeables []tax.Exchangeable
			position      tax.INRPosition
			testDate      time.Time
		)

		BeforeEach(func() {
			testDate = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			position = tax.INRPosition{
				Position: tax.Position{
					Date:     testDate,
					Quantity: 100,
					USDPrice: 150,
				},
			}
			exchangeables = []tax.Exchangeable{&position}
		})

		It("should handle missing exchange rate", func() {
			mockSBI.EXPECT().
				GetTTBuyRate(ctx, testDate).
				Return(0.0, common.ErrNotFound)

			err := exchangeMgr.Exchange(ctx, exchangeables)
			Expect(err).To(Equal(common.ErrNotFound))
		})

		It("should handle empty exchangeables list", func() {
			err := exchangeMgr.Exchange(ctx, nil)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("When closest date is used", func() {
			var (
				requestedDate = time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC)
				closestDate   = time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
				expectedRate  = 82.50
				exchangeables []tax.Exchangeable
				position      tax.INRPosition
			)

			BeforeEach(func() {
				// Setup test position
				position = tax.INRPosition{
					Position: tax.Position{
						Date:     requestedDate,
						Quantity: 100,
						USDPrice: 150,
					},
				}
				exchangeables = []tax.Exchangeable{&position}

				// Mock SBI manager to return ClosestDateError
				mockSBI.EXPECT().
					GetTTBuyRate(ctx, requestedDate).
					Return(expectedRate, tax.NewClosestDateError(requestedDate, closestDate))
			})

			It("should set closest date and rate", func() {
				err := exchangeMgr.Exchange(ctx, exchangeables)

				Expect(err).ToNot(HaveOccurred())
				Expect(position.TTRate).To(Equal(expectedRate))
				Expect(position.TTDate).To(Equal(closestDate)) // Verify closest date was set

				// Verify TTDate is different from requested date
				Expect(position.TTDate).ToNot(Equal(requestedDate))
			})
		})

		Context("When exact date is found", func() {
			var (
				requestedDate = time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
				expectedRate  = 82.50
				exchangeables []tax.Exchangeable
				position      tax.INRPosition
			)

			BeforeEach(func() {
				position = tax.INRPosition{
					Position: tax.Position{
						Date:     requestedDate,
						Quantity: 100,
						USDPrice: 150,
					},
				}
				exchangeables = []tax.Exchangeable{&position}

				mockSBI.EXPECT().
					GetTTBuyRate(ctx, requestedDate).
					Return(expectedRate, nil)
			})

			It("should set exact date and rate", func() {
				err := exchangeMgr.Exchange(ctx, exchangeables)

				Expect(err).ToNot(HaveOccurred())
				Expect(position.TTRate).To(Equal(expectedRate))
				Expect(position.TTDate).To(Equal(requestedDate)) // Verify requested date was set
			})
		})
	})

	Describe("ExchangeGains", func() {
		var (
			gains []tax.INRGains
		)

		Context("Successful Rate Fetch (Exact Date)", func() {
			BeforeEach(func() {
				gains = []tax.INRGains{
					{Gains: tax.Gains{Symbol: "AAPL", SellDate: "2023-04-15", PNL: 100}},
				}
				// Expected target date for SellDate "2023-04-15" is "2023-03-31"
				expectedTargetDate := time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC)
				mockSBI.EXPECT().GetTTBuyRate(ctx, expectedTargetDate).Return(82.50, nil).Once()
			})

			It("should set TTRate and TTDate correctly", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).ToNot(HaveOccurred())
				Expect(gains[0].TTRate).To(Equal(82.50))
				Expect(gains[0].TTDate.Format(time.DateOnly)).To(Equal("2023-03-31"))
			})
		})

		Context("Closest Date Scenario", func() {
			BeforeEach(func() {
				gains = []tax.INRGains{
					{Gains: tax.Gains{Symbol: "MSFT", SellDate: "2023-03-10", PNL: 200}},
				}
				// Expected target date for SellDate "2023-03-10" is "2023-02-28"
				requestedTargetDate := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
				closestDate := time.Date(2023, 2, 27, 0, 0, 0, 0, time.UTC)
				mockSBI.EXPECT().GetTTBuyRate(ctx, requestedTargetDate).Return(82.00, tax.NewClosestDateError(requestedTargetDate, closestDate)).Once()
			})

			It("should set TTRate to closest rate and TTDate to closest date", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).ToNot(HaveOccurred()) // ClosestDateError is not a processing failure for ExchangeGains itself
				Expect(gains[0].TTRate).To(Equal(82.00))
				Expect(gains[0].TTDate.Format(time.DateOnly)).To(Equal("2023-02-27"))
			})
		})

		Context("Error from SBIManager", func() {
			BeforeEach(func() {
				gains = []tax.INRGains{
					{Gains: tax.Gains{Symbol: "GOOG", SellDate: "2023-05-20", PNL: 150}},
				}
				// Expected target date for SellDate "2023-05-20" is "2023-04-30"
				expectedTargetDate := time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC)
				mockSBI.EXPECT().GetTTBuyRate(ctx, expectedTargetDate).Return(0.0, common.ErrInternalServerError).Once() // Using a predefined common.HttpError
			})

			It("should return the error from SBIManager", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})

		Context("Invalid SellDate in INRGains", func() {
			BeforeEach(func() {
				gains = []tax.INRGains{
					{Gains: tax.Gains{Symbol: "TSLA", SellDate: "invalid-date", PNL: 50}},
				}
				// No mock expectation for SBIManager as it shouldn't be called
			})

			It("should return an InvalidDateError", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).To(HaveOccurred())
				_, ok := err.(tax.InvalidDateError) // Check if the error is of the expected type
				Expect(ok).To(BeTrue(), "Error should be of type tax.InvalidDateError")
			})
		})

		Context("Multiple Gains Processing", func() {
			// Define dates clearly for readability
			var (
				sellDate1    = "2023-04-15" // Target Mar 31
				targetDate1  = time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC)
				rate1        = 82.50
				sellDate2    = "2023-03-10" // Target Feb 28
				targetDate2  = time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
				closestDate2 = time.Date(2023, 2, 27, 0, 0, 0, 0, time.UTC)
				rate2        = 82.00
				sellDate3    = "2023-01-05" // Target Dec 31 (prev year)
				targetDate3  = time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
				rate3        = 81.00
			)
			BeforeEach(func() {
				gains = []tax.INRGains{
					{Gains: tax.Gains{Symbol: "S1", SellDate: sellDate1, PNL: 10}},
					{Gains: tax.Gains{Symbol: "S2", SellDate: sellDate2, PNL: 20}},
					{Gains: tax.Gains{Symbol: "S3", SellDate: sellDate3, PNL: 30}},
				}
				mockSBI.EXPECT().GetTTBuyRate(ctx, targetDate1).Return(rate1, nil).Once()
				mockSBI.EXPECT().GetTTBuyRate(ctx, targetDate2).Return(rate2, tax.NewClosestDateError(targetDate2, closestDate2)).Once()
				mockSBI.EXPECT().GetTTBuyRate(ctx, targetDate3).Return(rate3, nil).Once()
			})

			It("should process all gains correctly", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).ToNot(HaveOccurred())

				// Gain 1
				Expect(gains[0].TTRate).To(Equal(rate1))
				Expect(gains[0].TTDate.Format(time.DateOnly)).To(Equal(targetDate1.Format(time.DateOnly)))

				// Gain 2 (Closest Date)
				Expect(gains[1].TTRate).To(Equal(rate2))
				Expect(gains[1].TTDate.Format(time.DateOnly)).To(Equal(closestDate2.Format(time.DateOnly)))

				// Gain 3
				Expect(gains[2].TTRate).To(Equal(rate3))
				Expect(gains[2].TTDate.Format(time.DateOnly)).To(Equal(targetDate3.Format(time.DateOnly)))
			})
		})

		Context("Empty Gains Slice", func() {
			BeforeEach(func() {
				gains = []tax.INRGains{}
				// No mock expectations as SBIManager should not be called.
			})

			It("should complete without error", func() {
				err := exchangeMgr.ExchangeGains(ctx, gains)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
