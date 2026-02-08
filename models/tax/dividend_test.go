package tax_test

import (
	. "github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dividend", func() {
	Context("INRDividend", func() {
		Describe("INRValue", func() {
			It("should calculate INR value: 33.41 * 86.95 = 2905.00", func() {
				dividend := INRDividend{Dividend: Dividend{Amount: 33.41}, TTRate: 86.95}
				Expect(dividend.INRValue()).To(Equal(2905.0))
			})

			It("should round with precision errors: 156.83 * 84.15 = 13197.24", func() {
				dividend := INRDividend{Dividend: Dividend{Amount: 156.83}, TTRate: 84.15}
				Expect(dividend.INRValue()).To(Equal(13197.24))
			})
		})
	})
})
