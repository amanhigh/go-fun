package tax

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Valuation", func() {
	Describe("Position", func() {
		It("should calculate USDValue correctly", func() {
			p := Position{
				Date:     time.Now(),
				Quantity: 10,
				USDPrice: 100,
			}
			Expect(p.USDValue()).To(Equal(float64(1000)))
		})
	})

	Describe("Trade", func() {
		Context("IsValid validation", func() {
			It("should accept valid trades with uppercase trade types", func() {
				validTrade := Trade{
					Symbol: "AAPL",
					Date:   "2024-01-01",
					Type:   "BUY",
				}
				Expect(validTrade.IsValid()).To(BeTrue())

				validSell := Trade{
					Symbol: "AAPL",
					Date:   "2024-01-01",
					Type:   "SELL",
				}
				Expect(validSell.IsValid()).To(BeTrue())
			})

			// This test should FAIL initially - reproduces the production issue
			It("should accept mixed case trade types from CSV files", func() {
				// This reproduces the exact data that's failing in production
				mixedCaseBuy := Trade{
					Symbol: "VTWO",
					Date:   "2025-04-03",
					Type:   "Buy", // Mixed case - this currently fails validation
				}

				// This should pass but currently fails with Trade.IsValid()
				Expect(mixedCaseBuy.IsValid()).To(BeTrue())

				mixedCaseSell := Trade{
					Symbol: "AAPL",
					Date:   "2024-01-01",
					Type:   "Sell", // Mixed case
				}

				// This should pass but currently fails with Trade.IsValid()
				Expect(mixedCaseSell.IsValid()).To(BeTrue())
			})

			It("should reject invalid trades", func() {
				invalidTrades := []Trade{
					{Symbol: "", Date: "2024-01-01", Type: "BUY"},         // Empty symbol
					{Symbol: "AAPL", Date: "", Type: "BUY"},               // Empty date
					{Symbol: "AAPL", Date: "2024-01-01", Type: ""},        // Empty type
					{Symbol: "AAPL", Date: "2024-01-01", Type: "INVALID"}, // Invalid type
				}

				for _, trade := range invalidTrades {
					Expect(trade.IsValid()).To(BeFalse())
				}
			})
		})
	})
})
