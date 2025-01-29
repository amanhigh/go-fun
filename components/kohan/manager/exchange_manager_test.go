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

			Expect(err).To(BeNil())
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

			Expect(err).To(BeNil())
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
			Expect(err).To(BeNil())
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

				Expect(err).To(BeNil())
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

				Expect(err).To(BeNil())
				Expect(position.TTRate).To(Equal(expectedRate))
				Expect(position.TTDate).To(Equal(requestedDate)) // Verify requested date was set
			})
		})
	})
})
