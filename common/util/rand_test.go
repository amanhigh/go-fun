package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
)

var _ = Describe("Rand", func() {
	Describe("RandomInts", func() {
		It("should generate array of random integers (SDK wrapper)", func() {
			result := util.RandomInts(5, 10)
			Expect(result).To(HaveLen(5))
			for _, r := range result {
				Expect(r).To(BeNumerically(">=", 0))
				Expect(r).To(BeNumerically("<", 10))
			}

			// Edge cases
			Expect(util.RandomInts(0, 10)).To(BeEmpty())

			ones := util.RandomInts(3, 1)
			for _, r := range ones {
				Expect(r).To(Equal(0))
			}
		})
	})

	Describe("RandomInt", func() {
		It("should generate random integer in range (SDK wrapper)", func() {
			result := util.RandomInt(5, 10)
			Expect(result).To(BeNumerically(">=", 5))
			Expect(result).To(BeNumerically("<", 10))

			// Range of 1
			Expect(util.RandomInt(5, 6)).To(Equal(5))

			// Negative ranges
			result = util.RandomInt(-10, -5)
			Expect(result).To(BeNumerically(">=", -10))
			Expect(result).To(BeNumerically("<", -5))
		})

		It("should panic when bounds are equal (implementation limitation)", func() {
			Expect(func() {
				util.RandomInt(42, 42)
			}).To(Panic())
		})
	})
})
