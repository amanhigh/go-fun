package tax_test

import (
	. "github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interest", func() {
	Context("INRInterest", func() {
		Describe("INRValue", func() {
			It("should calculate the correct INR value with 2 decimal places", func() {
				interest := INRInterest{
					Interest: Interest{Amount: 100.555}, // USD amount
					TTRate:   83.456,                    // Exchange rate
				}
				expectedINRValue := 8391.92
				Expect(interest.INRValue()).To(Equal(expectedINRValue))
			})

			It("should handle zero amount", func() {
				interest := INRInterest{
					Interest: Interest{Amount: 0},
					TTRate:   83.456,
				}
				expectedINRValue := 0.0
				Expect(interest.INRValue()).To(Equal(expectedINRValue))
			})

			It("should handle zero exchange rate", func() {
				interest := INRInterest{
					Interest: Interest{Amount: 100.555},
					TTRate:   0,
				}
				expectedINRValue := 0.0
				Expect(interest.INRValue()).To(Equal(expectedINRValue))
			})

			It("should handle negative amount", func() {
				interest := INRInterest{
					Interest: Interest{Amount: -50.25},
					TTRate:   80.0,
				}
				// Expected: -50.25 * 80.0 = -4020.0
				expectedINRValue := -4020.0
				Expect(interest.INRValue()).To(Equal(expectedINRValue))
			})
		})
	})
})
