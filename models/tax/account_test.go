package tax_test

import (
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Context("FromValuations", func() {
		It("should convert valuations to accounts", func() {
			valuations := []tax.Valuation{
				{
					Ticker: "GOOG",
					YearEndPosition: tax.Position{
						Quantity: 10,
						USDPrice: 150,
					},
				},
				{
					Ticker: "AAPL",
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

			Expect(accounts[1].Symbol).To(Equal("AAPL"))
			Expect(accounts[1].Quantity).To(Equal(float64(20)))
			Expect(accounts[1].Cost).To(Equal(float64(3500)))
			Expect(accounts[1].MarketValue).To(Equal(float64(3500)))
		})
	})
})
