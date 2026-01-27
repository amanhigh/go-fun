package tax_test

import (
	"time"

	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Context("FromValuations", func() {
		It("should convert valuations to accounts with origin metadata", func() {
			yearEndDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
			valuations := []tax.Valuation{
				{
					Ticker: "AAPL",
					FirstPosition: tax.Position{
						Date:     time.Date(2021, 3, 5, 0, 0, 0, 0, time.UTC),
						Quantity: 50,
						USDPrice: 130.00,
					},
					YearEndPosition: tax.Position{
						Date:     yearEndDate,
						Quantity: 100,
						USDPrice: 180.00,
					},
				},
			}

			accounts := tax.FromValuations(valuations)
			Expect(accounts).To(HaveLen(1))

			account := accounts[0]
			Expect(account.Symbol).To(Equal("AAPL"))
			Expect(account.Quantity).To(Equal(100.0))
			Expect(account.OriginDate).To(Equal("2021-03-05"))
			Expect(account.OriginQty).To(Equal(50.0))
			Expect(account.OriginPrice).To(Equal(130.00))
		})

		It("should handle multiple valuations", func() {
			valuations := []tax.Valuation{
				{
					Ticker: "GOOG",
					FirstPosition: tax.Position{
						Date:     time.Date(2021, 3, 5, 0, 0, 0, 0, time.UTC),
						Quantity: 10,
						USDPrice: 150,
					},
					YearEndPosition: tax.Position{
						Quantity: 10,
						USDPrice: 150,
					},
				},
				{
					Ticker: "AAPL",
					FirstPosition: tax.Position{
						Date:     time.Date(2020, 8, 15, 0, 0, 0, 0, time.UTC),
						Quantity: 20,
						USDPrice: 175,
					},
					YearEndPosition: tax.Position{
						Quantity: 20,
						USDPrice: 175,
					},
				},
			}

			accounts := tax.FromValuations(valuations)
			Expect(accounts).To(HaveLen(2))

			Expect(accounts[0].Symbol).To(Equal("GOOG"))
			Expect(accounts[0].Quantity).To(Equal(float64(10)))
			Expect(accounts[0].Cost).To(Equal(float64(1500)))
			Expect(accounts[0].MarketValue).To(Equal(float64(1500)))
			Expect(accounts[0].OriginDate).To(Equal("2021-03-05"))
			Expect(accounts[0].OriginQty).To(Equal(10.0))
			Expect(accounts[0].OriginPrice).To(Equal(150.0))

			Expect(accounts[1].Symbol).To(Equal("AAPL"))
			Expect(accounts[1].Quantity).To(Equal(float64(20)))
			Expect(accounts[1].Cost).To(Equal(float64(3500)))
			Expect(accounts[1].MarketValue).To(Equal(float64(3500)))
			Expect(accounts[1].OriginDate).To(Equal("2020-08-15"))
			Expect(accounts[1].OriginQty).To(Equal(20.0))
			Expect(accounts[1].OriginPrice).To(Equal(175.0))
		})
	})

	Context("IsValid", func() {
		It("should validate account with required fields", func() {
			account := tax.Account{
				Symbol:      "AAPL",
				Quantity:    100,
				Cost:        15050.00,
				MarketValue: 16000.00,
			}
			Expect(account.IsValid()).To(BeTrue())
		})

		It("should invalidate account missing symbol", func() {
			account := tax.Account{
				Symbol:      "",
				Quantity:    100,
				Cost:        15050.00,
				MarketValue: 16000.00,
			}
			Expect(account.IsValid()).To(BeFalse())
		})

		It("should allow origin fields to be zero", func() {
			account := tax.Account{
				Symbol:      "AAPL",
				Quantity:    100,
				Cost:        15050.00,
				MarketValue: 16000.00,
				// Origin fields are zero/empty, which is allowed
			}
			Expect(account.IsValid()).To(BeTrue())
		})
	})
})
