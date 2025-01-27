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
})
