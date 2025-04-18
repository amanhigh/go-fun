package sort_test

import (
	sort2 "sort"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/algos/sort"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sort Tests", func() {
	const count, upperBound = 10, 10
	var (
		original []int
		input    []int
		expected []int
	)

	BeforeEach(func() {
		original = util.RandomInts(count, upperBound)

		/* Make Copies to avoid changing original */
		input = make([]int, count)
		expected = make([]int, count)
		copy(input, original)
		copy(expected, original)

		sort2.Ints(expected)
	})

	It("should bubble sort", func() {
		swaps := sort.BubbleSort(input, count)
		Expect(input).To(Equal(expected))
		Expect(swaps).To(BeNumerically(">", 0))
	})

	It("should quick sort", func() {
		sort.QuickSort(input, 0, count-1)
		Expect(input).To(Equal(expected))
	})

	It("should merge sort", func() {
		inversions := sort.MergeSort(input, 0, count-1)
		Expect(input).To(Equal(expected))
		Expect(inversions).To(BeNumerically(">", 0))
	})

})
