package cracking_test

import (
	"github.com/amanhigh/go-fun/learn/concepts/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		max := cracking.FindRegion(cells, size, n, m)
		Expect(max).To(Equal(5))
	})
})
