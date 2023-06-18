package practice_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/practice"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Number Test", func() {
	Context("Missing Number", func() {

		It("should work for no Input", func() {
			Expect(1).To(Equal(practice.MissingNumbers([]int{})))
		})
		It("should work for duplicate numbers", func() {
			Expect(1).To(Equal(practice.MissingNumbers([]int{0, 0})))
			Expect(5).To(Equal(practice.MissingNumbers([]int{1, 3, 6, 4, 1, 2})))
		})

		It("should work for other cases", func() {
			Expect(1).To(Equal(practice.MissingNumbers([]int{-1, 3, 2, 0})))
			Expect(1).To(Equal(practice.MissingNumbers([]int{-1, -2, -3})))
			Expect(1).To(Equal(practice.MissingNumbers([]int{-1, -3})))
			Expect(1).To(Equal(practice.MissingNumbers([]int{3, 2, 0})))

			Expect(4).To(Equal(practice.MissingNumbers([]int{1, 3, 6, 1, 2})))
			Expect(4).To(Equal(practice.MissingNumbers([]int{1, 2, 3})))
		})
	})

	It("should work for target sum", func() {
		i, j := practice.TargetSum([]int{2, 7, 11, 15}, 9)
		Expect(i).To(Equal(0))
		Expect(j).To(Equal(1))
	})
})
