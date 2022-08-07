package cracking_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Small", func() {
	It("should generate fibonacci", func() {
		Expect(cracking.Fibonacci(34)).To(Equal(5702887))
	})

	It("should generate fibonacci recursively", func() {
		Expect(cracking.FibonacciRecursive(34)).To(Equal(5702887))
	})

	It("should find lonely", func() {
		Expect(cracking.FindLonely([]int{1, 1, 2})).To(Equal(2))
	})

	It("should KangaroMeet", func() {
		Expect(cracking.KangarooMeet([]int{0, 3, 4, 2})).To(BeTrue())
		Expect(cracking.KangarooMeet([]int{0, 2, 5, 3})).To(BeFalse())

		// Same Speed
		Expect(cracking.KangarooMeet([]int{0, 2, 2, 2})).To(BeFalse())
	})

	It("should Left Rotate", func() {
		Expect(cracking.LeftRotate([]int{1, 2, 3}, 2)).To(Equal([]int{3, 1, 2}))
	})
})
