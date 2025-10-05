package manager_test

import (
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tax Utils", func() {
	Describe("MatchDividendWithTax", func() {
		var (
			taxMap   map[string]map[string]float64
			dividend *tax.Dividend
		)

		BeforeEach(func() {
			taxMap = map[string]map[string]float64{
				"MSFT": {
					"2024-01-20": 7.5,
					"2024-02-20": 8.0,
				},
				"AAPL": {
					"2024-01-15": 5.0,
				},
			}
		})

		Context("when tax exists for dividend", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 50.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
			})

			It("should match tax correctly", func() {
				Expect(dividend.Tax).To(Equal(7.5))
			})

			It("should calculate net amount", func() {
				Expect(dividend.Net).To(Equal(42.5))
			})

			It("should remove tax from pool", func() {
				_, exists := taxMap["MSFT"]["2024-01-20"]
				Expect(exists).To(BeFalse())
			})

			It("should keep other taxes in pool", func() {
				Expect(taxMap["MSFT"]["2024-02-20"]).To(Equal(8.0))
				Expect(taxMap["AAPL"]["2024-01-15"]).To(Equal(5.0))
			})
		})

		Context("when no tax exists for dividend", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "GOOGL",
					Date:   "2024-01-20",
					Amount: 100.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
			})

			It("should set tax to 0", func() {
				Expect(dividend.Tax).To(Equal(0.0))
			})

			It("should set net equal to amount", func() {
				Expect(dividend.Net).To(Equal(100.0))
			})

			It("should not modify tax pool", func() {
				Expect(taxMap).To(HaveLen(2))
			})
		})

		Context("when symbol exists but date doesn't match", func() {
			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-03-20",
					Amount: 50.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
			})

			It("should set tax to 0", func() {
				Expect(dividend.Tax).To(Equal(0.0))
			})

			It("should set net equal to amount", func() {
				Expect(dividend.Net).To(Equal(50.0))
			})

			It("should keep existing taxes for that symbol", func() {
				Expect(taxMap["MSFT"]).To(HaveLen(2))
			})
		})

		Context("when multiple dividends share same tax pool", func() {
			var dividend2 *tax.Dividend

			BeforeEach(func() {
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 30.0,
				}
				dividend2 = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 20.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
				manager.MatchDividendWithTax(dividend2, taxMap)
			})

			It("should match first dividend with tax", func() {
				Expect(dividend.Tax).To(Equal(7.5))
				Expect(dividend.Net).To(Equal(22.5))
			})

			It("should not match second dividend (tax already consumed)", func() {
				Expect(dividend2.Tax).To(Equal(0.0))
				Expect(dividend2.Net).To(Equal(20.0))
			})

			It("should have removed tax from pool after first match", func() {
				_, exists := taxMap["MSFT"]["2024-01-20"]
				Expect(exists).To(BeFalse())
			})
		})

		Context("with empty tax map", func() {
			BeforeEach(func() {
				taxMap = map[string]map[string]float64{}
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 50.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
			})

			It("should set tax to 0", func() {
				Expect(dividend.Tax).To(Equal(0.0))
			})

			It("should set net equal to amount", func() {
				Expect(dividend.Net).To(Equal(50.0))
			})
		})

		Context("with nil dividend date map", func() {
			BeforeEach(func() {
				taxMap = map[string]map[string]float64{
					"MSFT": nil,
				}
				dividend = &tax.Dividend{
					Symbol: "MSFT",
					Date:   "2024-01-20",
					Amount: 50.0,
				}
				manager.MatchDividendWithTax(dividend, taxMap)
			})

			It("should handle gracefully and set tax to 0", func() {
				Expect(dividend.Tax).To(Equal(0.0))
			})

			It("should set net equal to amount", func() {
				Expect(dividend.Net).To(Equal(50.0))
			})
		})
	})
})
