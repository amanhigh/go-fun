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
				upperBound := 10
				result := util.RandomInts(n, upperBound)

				Expect(result).To(HaveLen(n))
				for _, r := range result {
					Expect(r).To(BeNumerically("<", upperBound))
				}
			})
		})
	})

	Describe("RandomInt", func() {
		Context("when min is less than max", func() {
			It("returns a random integer between min and max", func() {
				lowerBound := 5
				upperBound := 10
				result := util.RandomInt(lowerBound, upperBound)

				Expect(result).To(BeNumerically(">=", lowerBound))
				Expect(result).To(BeNumerically("<", upperBound))
			})
		})
	})
})
