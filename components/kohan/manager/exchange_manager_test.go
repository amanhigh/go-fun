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
				GetTTBuyRate(testDate).
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
					GetTTBuyRate(dates[i]).
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
				GetTTBuyRate(testDate).
				Return(0.0, common.ErrNotFound)

			err := exchangeMgr.Exchange(ctx, exchangeables)
			Expect(err).To(Equal(common.ErrNotFound))
		})

		It("should handle empty exchangeables list", func() {
			err := exchangeMgr.Exchange(ctx, nil)
			Expect(err).To(BeNil())
		})
	})
})
