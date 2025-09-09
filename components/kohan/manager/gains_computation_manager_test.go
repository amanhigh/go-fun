package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GainsComputationManager", func() {
	var (
		gainsManager manager.GainsComputationManager
		ctx          context.Context
	)

	BeforeEach(func() {
		gainsManager = manager.NewGainsComputationManager()
		ctx = context.Background()
	})

	Context("with simple BUY/SELL pairs", func() {
		var trades []tax.Trade

		BeforeEach(func() {
			trades = []tax.Trade{
				{Symbol: "AAPL", Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
				{Symbol: "AAPL", Date: "2024-01-17", Type: "SELL", Quantity: 100, USDPrice: 150.00, Commission: 10.00},
				{Symbol: "MSFT", Date: "2022-01-10", Type: "BUY", Quantity: 50, USDPrice: 200.00, Commission: 5.00},
				{Symbol: "MSFT", Date: "2024-02-15", Type: "SELL", Quantity: 50, USDPrice: 210.00, Commission: 5.00},
			}
		})

		It("should compute correct gains for simple pairs", func() {
			gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
			Expect(err).ToNot(HaveOccurred())
			Expect(gains).To(HaveLen(2))

			// Find gains by symbol using helper function
			aaplGain := findGainBySymbol(gains, "AAPL")
			msftGain := findGainBySymbol(gains, "MSFT")

			// Validate AAPL gain (STCG - 2 days holding)
			Expect(aaplGain).ToNot(BeNil())
			Expect(aaplGain.Symbol).To(Equal("AAPL"))
			Expect(aaplGain.BuyDate).To(Equal("2024-01-15"))
			Expect(aaplGain.SellDate).To(Equal("2024-01-17"))
			Expect(aaplGain.Quantity).To(Equal(100.0))
			Expect(aaplGain.PNL).To(Equal(980.00)) // (150-140)*100 - 20 commission
			Expect(aaplGain.Commission).To(Equal(20.00))
			Expect(aaplGain.Type).To(Equal(tax.GAIN_TYPE_STCG))

			// Validate MSFT gain (LTCG - 2+ years holding)
			Expect(msftGain).ToNot(BeNil())
			Expect(msftGain.Symbol).To(Equal("MSFT"))
			Expect(msftGain.BuyDate).To(Equal("2022-01-10"))
			Expect(msftGain.SellDate).To(Equal("2024-02-15"))
			Expect(msftGain.Quantity).To(Equal(50.0))
			Expect(msftGain.PNL).To(Equal(490.00)) // (210-200)*50 - 10 commission
			Expect(msftGain.Commission).To(Equal(10.00))
			Expect(msftGain.Type).To(Equal(tax.GAIN_TYPE_LTCG))
		})
	})

	Context("with FIFO multiple lots", func() {
		var trades []tax.Trade

		BeforeEach(func() {
			trades = []tax.Trade{
				// Buy multiple lots of AAPL
				{Symbol: "AAPL", Date: "2023-01-10", Type: "BUY", Quantity: 20, USDPrice: 150.00, Commission: 2.00},
				{Symbol: "AAPL", Date: "2023-07-10", Type: "BUY", Quantity: 30, USDPrice: 165.00, Commission: 3.00},
				{Symbol: "AAPL", Date: "2023-12-10", Type: "BUY", Quantity: 50, USDPrice: 180.00, Commission: 5.00},
				// Sell partial quantity - should match FIFO (oldest first)
				{Symbol: "AAPL", Date: "2023-10-20", Type: "SELL", Quantity: 15, USDPrice: 170.00, Commission: 1.50},
			}
		})

		It("should use FIFO for multiple buy lots", func() {
			gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
			Expect(err).ToNot(HaveOccurred())
			Expect(gains).To(HaveLen(1))

			gain := gains[0]
			Expect(gain.Symbol).To(Equal("AAPL"))
			Expect(gain.BuyDate).To(Equal("2023-01-10")) // Should match oldest buy
			Expect(gain.SellDate).To(Equal("2023-10-20"))
			Expect(gain.Quantity).To(Equal(15.0))

			// PNL calculation: (170-150)*15 - commission allocation
			// Buy commission portion: 2.00 * (15/20) = 1.50
			// Sell commission: 1.50
			// Total commission: 3.00
			// PNL: (170-150)*15 - 3.00 = 300 - 3 = 297
			Expect(gain.PNL).To(Equal(297.00))
			Expect(gain.Type).To(Equal(tax.GAIN_TYPE_STCG)) // Less than 2 years
		})
	})

	Context("with partial sales across multiple lots", func() {
		var trades []tax.Trade

		BeforeEach(func() {
			trades = []tax.Trade{
				{Symbol: "AAPL", Date: "2023-01-10", Type: "BUY", Quantity: 20, USDPrice: 150.00, Commission: 2.00},
				{Symbol: "AAPL", Date: "2023-07-10", Type: "BUY", Quantity: 30, USDPrice: 165.00, Commission: 3.00},
				// Sell more than first lot - should span multiple lots
				{Symbol: "AAPL", Date: "2024-01-15", Type: "SELL", Quantity: 35, USDPrice: 180.00, Commission: 3.50},
			}
		})

		It("should span multiple lots when sell quantity exceeds first lot", func() {
			gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
			Expect(err).ToNot(HaveOccurred())
			Expect(gains).To(HaveLen(2)) // Should create 2 gains - one per lot

			// Sort gains by buy date to ensure consistent testing
			if gains[0].BuyDate > gains[1].BuyDate {
				gains[0], gains[1] = gains[1], gains[0]
			}

			// First gain: Complete first lot (20 shares)
			gain1 := gains[0]
			Expect(gain1.Symbol).To(Equal("AAPL"))
			Expect(gain1.BuyDate).To(Equal("2023-01-10"))
			Expect(gain1.Quantity).To(Equal(20.0))
			// PNL: (180-150)*20 - commissions
			// Buy commission: 2.00 (full lot)
			// Sell commission portion: 3.50 * (20/35) = 2.00
			// Total commission: 4.00
			expectedPnl1 := (180.00-150.00)*20 - 4.00
			Expect(gain1.PNL).To(Equal(expectedPnl1))

			// Second gain: Partial second lot (15 out of 30 shares)
			gain2 := gains[1]
			Expect(gain2.Symbol).To(Equal("AAPL"))
			Expect(gain2.BuyDate).To(Equal("2023-07-10"))
			Expect(gain2.Quantity).To(Equal(15.0))
			// PNL: (180-165)*15 - commissions
			// Buy commission portion: 3.00 * (15/30) = 1.50
			// Sell commission portion: 3.50 * (15/35) = 1.50
			// Total commission: 3.00
			expectedPnl2 := (180.00-165.00)*15 - 3.00
			Expect(gain2.PNL).To(Equal(expectedPnl2))
		})
	})

	Context("STCG/LTCG classification for foreign stocks", func() {
		var trades []tax.Trade

		BeforeEach(func() {
			trades = []tax.Trade{
				// Less than 2 years - STCG
				{Symbol: "STCG1", Date: "2023-01-15", Type: "BUY", Quantity: 100, USDPrice: 100.00, Commission: 1.00},
				{Symbol: "STCG1", Date: "2024-01-14", Type: "SELL", Quantity: 100, USDPrice: 110.00, Commission: 1.00}, // 364 days

				// Exactly 2 years - LTCG
				{Symbol: "LTCG1", Date: "2022-01-15", Type: "BUY", Quantity: 100, USDPrice: 100.00, Commission: 1.00},
				{Symbol: "LTCG1", Date: "2024-01-15", Type: "SELL", Quantity: 100, USDPrice: 110.00, Commission: 1.00}, // Exactly 730 days

				// More than 2 years - LTCG
				{Symbol: "LTCG2", Date: "2021-01-15", Type: "BUY", Quantity: 100, USDPrice: 100.00, Commission: 1.00},
				{Symbol: "LTCG2", Date: "2024-01-16", Type: "SELL", Quantity: 100, USDPrice: 110.00, Commission: 1.00}, // 3+ years
			}
		})

		It("should classify gains correctly based on 2-year rule", func() {
			gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
			Expect(err).ToNot(HaveOccurred())
			Expect(gains).To(HaveLen(3))

			// Create map for easy lookup
			gainsBySymbol := make(map[string]tax.Gains)
			for _, gain := range gains {
				gainsBySymbol[gain.Symbol] = gain
			}

			// STCG - less than 2 years
			stcgGain := gainsBySymbol["STCG1"]
			Expect(stcgGain.Type).To(Equal(tax.GAIN_TYPE_STCG))

			// LTCG - exactly 2 years (730 days)
			ltcgGain1 := gainsBySymbol["LTCG1"]
			Expect(ltcgGain1.Type).To(Equal(tax.GAIN_TYPE_LTCG))

			// LTCG - more than 2 years
			ltcgGain2 := gainsBySymbol["LTCG2"]
			Expect(ltcgGain2.Type).To(Equal(tax.GAIN_TYPE_LTCG))
		})
	})

	Context("edge cases", func() {
		Context("with no sell transactions", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: "AAPL", Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
					{Symbol: "MSFT", Date: "2024-01-16", Type: "BUY", Quantity: 50, USDPrice: 200.00, Commission: 5.00},
				}
			})

			It("should return empty gains list", func() {
				gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).ToNot(HaveOccurred())
				Expect(gains).To(BeEmpty())
			})
		})

		Context("with sell without matching buy", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: "AAPL", Date: "2024-01-15", Type: "SELL", Quantity: 100, USDPrice: 150.00, Commission: 10.00},
				}
			})

			It("should return an error", func() {
				_, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no buy positions found"))
			})
		})

		Context("with sell quantity exceeding buy quantity", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: "AAPL", Date: "2024-01-15", Type: "BUY", Quantity: 50, USDPrice: 140.00, Commission: 5.00},
					{Symbol: "AAPL", Date: "2024-01-17", Type: "SELL", Quantity: 100, USDPrice: 150.00, Commission: 10.00},
				}
			})

			It("should return an error", func() {
				_, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("insufficient buy quantity"))
			})
		})

		Context("with same-day buy and sell", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: "AAPL", Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
					{Symbol: "AAPL", Date: "2024-01-15", Type: "SELL", Quantity: 100, USDPrice: 150.00, Commission: 10.00},
				}
			})

			It("should compute gain correctly for same-day trade", func() {
				gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).ToNot(HaveOccurred())
				Expect(gains).To(HaveLen(1))

				gain := gains[0]
				Expect(gain.Symbol).To(Equal("AAPL"))
				Expect(gain.BuyDate).To(Equal(gain.SellDate))
				Expect(gain.Type).To(Equal(tax.GAIN_TYPE_STCG)) // Same day = 0 days < 2 years
				Expect(gain.PNL).To(Equal(980.00))              // (150-140)*100 - 20 commission
			})
		})

		Context("with invalid date format", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					// Note: Date validation also occurs in Trade Repository during CSV parsing
					{Symbol: "AAPL", Date: "invalid-date", Type: "BUY", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
				}
			})

			It("should return an error for invalid date format", func() {
				_, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid trade date"))
			})
		})

		Context("with mixed case trade types", func() {
			var trades []tax.Trade

			BeforeEach(func() {
				trades = []tax.Trade{
					{Symbol: "AAPL", Date: "2024-01-15", Type: "Buy", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
					{Symbol: "AAPL", Date: "2024-01-17", Type: "Sell", Quantity: 100, USDPrice: 150.00, Commission: 10.00},
				}
			})

			It("should handle mixed case trade types correctly", func() {
				gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
				Expect(err).ToNot(HaveOccurred())
				Expect(gains).To(HaveLen(1))

				gain := gains[0]
				Expect(gain.Symbol).To(Equal("AAPL"))
				Expect(gain.PNL).To(Equal(980.00))
			})
		})
	})

	Context("with complex multi-symbol, multi-lot scenario", func() {
		var trades []tax.Trade

		BeforeEach(func() {
			// Reproduce the exact test data pattern
			trades = []tax.Trade{
				// AAPL transactions
				{Symbol: "AAPL", Date: "2024-01-15", Type: "BUY", Quantity: 100, USDPrice: 140.00, Commission: 10.00},
				{Symbol: "AAPL", Date: "2023-03-15", Type: "BUY", Quantity: 20, USDPrice: 150.00, Commission: 1.00},
				{Symbol: "AAPL", Date: "2023-07-10", Type: "BUY", Quantity: 30, USDPrice: 165.00, Commission: 1.00},
				{Symbol: "AAPL", Date: "2023-10-20", Type: "SELL", Quantity: 15, USDPrice: 170.00, Commission: 1.00},
				{Symbol: "AAPL", Date: "2024-01-17", Type: "SELL", Quantity: 100, USDPrice: 150.00, Commission: 10.00},

				// MSFT transactions
				{Symbol: "MSFT", Date: "2022-01-10", Type: "BUY", Quantity: 50, USDPrice: 200.00, Commission: 5.00},
				{Symbol: "MSFT", Date: "2023-05-01", Type: "BUY", Quantity: 20, USDPrice: 205.00, Commission: 2.00},
				{Symbol: "MSFT", Date: "2023-09-01", Type: "BUY", Quantity: 30, USDPrice: 215.00, Commission: 3.00},
				{Symbol: "MSFT", Date: "2024-02-15", Type: "SELL", Quantity: 50, USDPrice: 210.00, Commission: 5.00},
			}
		})

		It("should handle complex multi-symbol scenario correctly", func() {
			gains, err := gainsManager.ComputeGainsFromTrades(ctx, trades)
			Expect(err).ToNot(HaveOccurred())
			Expect(gains).To(HaveLen(5)) // FIFO algorithm generates detailed gain records

			// Group gains by symbol
			aaplGains := make([]tax.Gains, 0)
			msftGains := make([]tax.Gains, 0)

			for _, gain := range gains {
				switch gain.Symbol {
				case "AAPL":
					aaplGains = append(aaplGains, gain)
				case "MSFT":
					msftGains = append(msftGains, gain)
				}
			}

			// Validate AAPL gains (4 gain records due to FIFO matching)
			Expect(aaplGains).To(HaveLen(4))

			// All AAPL gains should be STCG (less than 2 years)
			for _, gain := range aaplGains {
				Expect(gain.Type).To(Equal(tax.GAIN_TYPE_STCG))
				Expect(gain.Symbol).To(Equal("AAPL"))
			}

			// Validate MSFT gains (should have 1 gain)
			Expect(msftGains).To(HaveLen(1))

			msftGain := msftGains[0]
			Expect(msftGain.Symbol).To(Equal("MSFT"))
			Expect(msftGain.BuyDate).To(Equal("2022-01-10")) // Should match oldest MSFT buy (FIFO)
			Expect(msftGain.Quantity).To(Equal(50.0))
			Expect(msftGain.Type).To(Equal(tax.GAIN_TYPE_LTCG)) // 2+ years

			// Validate total PNL makes sense
			var totalPnl float64
			for _, gain := range gains {
				totalPnl += gain.PNL
			}
			// Should have positive total PNL for this scenario
			Expect(totalPnl).To(BeNumerically(">", 0))
		})
	})
})

// Helper function for finding gains by symbol (could be moved to test helpers if needed)
func findGainBySymbol(gains []tax.Gains, symbol string) *tax.Gains {
	for i := range gains {
		if gains[i].Symbol == symbol {
			return &gains[i]
		}
	}
	return nil
}
