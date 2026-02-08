package tax_test

import (
	. "github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gains", func() {
	Context("INRGains", func() {
		Describe("INRValue", func() {
			It("should convert negative PNL to INR with rounding", func() {
				gains := INRGains{
					Gains:  Gains{PNL: -3134.91},
					TTRate: 82.50,
				}
				Expect(gains.INRValue()).To(Equal(-258630.08))
			})

			It("should convert positive PNL to INR with rounding", func() {
				gains := INRGains{
					Gains:  Gains{PNL: 1234.56},
					TTRate: 85.75,
				}
				Expect(gains.INRValue()).To(Equal(105863.52))
			})
		})
	})
})
