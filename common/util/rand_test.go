package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
)

var _ = Describe("Rand", func() {
	Describe("RandomInts", func() {
		Context("when n is greater than 0", func() {
			It("returns a slice of n random integers", func() {
				n := 5
				max := 10
				result := util.RandomInts(n, max)

				Expect(len(result)).To(Equal(n))
				for _, r := range result {
					Expect(r).To(BeNumerically("<", max))
				}
			})
		})
	})

	Describe("RandomInt", func() {
		Context("when min is less than max", func() {
			It("returns a random integer between min and max", func() {
				min := 5
				max := 10
				result := util.RandomInt(min, max)

				Expect(result).To(BeNumerically(">=", min))
				Expect(result).To(BeNumerically("<", max))
			})
		})
	})
})
