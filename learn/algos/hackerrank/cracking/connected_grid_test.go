package cracking_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/cracking"
)

var _ = Describe("ConnectedGrid", func() {
	var (
		n, m  = 4, 4
		cells = [][]int{
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{1, 0, 0, 0},
		}
		size = [][]int{
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
		}
	)
	It("should find correct regions", func() {

		max := FindRegion(cells, size, n, m)
		Expect(max).To(Equal(5))
	})
})
