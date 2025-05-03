package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	repomock "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("InterestManager", func() {
	var (
		ctx                  = context.Background()
		mockExchange         *mocks.ExchangeManager
		mockFinancialYearMgr *mocks.FinancialYearManager[tax.Interest]
		mockInterestRepo     *repomock.InterestRepository
		interestManager      manager.InterestManager

		// Common test data
		ticker   = "AAPL"
		date     = "2024-01-15"
		amount   = 100.00
		taxValue = 10.00
		net      = 90.00
	)

	BeforeEach(func() {
		mockExchange = mocks.NewExchangeManager(GinkgoT())
		mockFinancialYearMgr = mocks.NewFinancialYearManager[tax.Interest](GinkgoT())
		mockInterestRepo = repomock.NewInterestRepository(GinkgoT())
		interestManager = manager.NewInterestManager(
			mockExchange,
			mockFinancialYearMgr,
			mockInterestRepo,
		)
	})

	Context("Basic Interest Processing", func() {
		var (
			interests []tax.Interest
		)

		BeforeEach(func() {
			interests = []tax.Interest{
				{
					Symbol: ticker,
					Date:   date,
					Amount: amount,
					Tax:    taxValue,
					Net:    net,
				},
			}

			mockExchange.EXPECT().
				Exchange(ctx, mock.Anything).
				Return(nil)
		})

		It("should process interest correctly", func() {
			inrInterests, err := interestManager.ProcessInterest(ctx, interests)

			Expect(err).ToNot(HaveOccurred())
			Expect(inrInterests).To(HaveLen(1))

			result := inrInterests[0]
			Expect(result.Interest).To(Equal(interests[0]))
		})
	})

	Context("Exchange Rate Error", func() {
		var interests []tax.Interest

		BeforeEach(func() {
			interests = []tax.Interest{{
				Symbol: ticker,
				Date:   date,
				Amount: amount,
				Tax:    taxValue,
				Net:    net,
			}}

			mockExchange.EXPECT().
				Exchange(ctx, mock.Anything).
				Return(common.ErrNotFound)
		})

		It("should handle missing exchange rate", func() {
			_, err := interestManager.ProcessInterest(ctx, interests)
			Expect(err).To(Equal(common.ErrNotFound))
		})
	})

	Context("Multiple Interests", func() {
		var (
			interests []tax.Interest
		)

		BeforeEach(func() {
			interests = []tax.Interest{
				{
					Symbol: "AAPL",
					Date:   "2024-01-15",
					Amount: 100.00,
					Tax:    10.00,
					Net:    90.00,
				},
				{
					Symbol: "MSFT",
					Date:   "2024-01-16",
					Amount: 200.00,
					Tax:    20.00,
					Net:    180.00,
				},
			}

			// Verify that exchangeables passed contain correct interest amounts
			mockExchange.EXPECT().
				Exchange(ctx, mock.AnythingOfType("[]tax.Exchangeable")).
				Run(func(_ context.Context, exchangeables []tax.Exchangeable) {
					Expect(exchangeables).To(HaveLen(2))
					Expect(exchangeables[0].GetUSDAmount()).To(Equal(100.00))
					Expect(exchangeables[1].GetUSDAmount()).To(Equal(200.00))
				}).
				Return(nil)
		})

		It("should process multiple interests correctly", func() {
			inrInterests, err := interestManager.ProcessInterest(ctx, interests)

			Expect(err).ToNot(HaveOccurred())
			Expect(inrInterests).To(HaveLen(2))

			// Verify first interest
			Expect(inrInterests[0].Symbol).To(Equal("AAPL"))
			Expect(inrInterests[0].Amount).To(Equal(100.00))

			// Verify second interest
			Expect(inrInterests[1].Symbol).To(Equal("MSFT"))
			Expect(inrInterests[1].Amount).To(Equal(200.00))
		})
	})
})
