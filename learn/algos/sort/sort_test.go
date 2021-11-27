package sort_test

import (
	sort3 "github.com/amanhigh/go-fun/learn/algos/sort"
	"github.com/amanhigh/go-fun/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	sort2 "sort"
)

var _ = Describe("Sort Tests", func() {
	const count, max = 10, 10
	var (
		original []int
		input    []int
		expected []int
	)

	BeforeEach(func() {
		original = util.RandomInts(count, max)

		/* Make Copies to avoid changing original */
		input = make([]int, count)
		expected = make([]int, count)
		copy(input, original)
		copy(expected, original)

		sort2.Ints(expected)
	})

	It("should bubble sort", func() {
		sort3.BubbleSort(input, count)
		Expect(input).To(Equal(expected))
	})

	It("should quick sort", func() {
		sort3.QuickSort(input, 0, count-1)
		Expect(input).To(Equal(expected))
	})

	It("should merge sort", func() {
		sort3.MergeSort(input, 0, count-1)
		Expect(input).To(Equal(expected))
	})

})
