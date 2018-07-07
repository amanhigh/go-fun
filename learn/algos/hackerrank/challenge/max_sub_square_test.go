package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
)

var _ = Describe("MaxSubSquare", func() {
	var (
		input = `3 3
-1 -2 -4
-8 -2 5
-3 6 7`
	)
	It("should compute sum", func() {
		scanner := util.NewStringScanner(input)
		n := util.ReadInt(scanner)
		m := util.ReadInt(scanner)
		inputMatrix := helper.ReadMatrix(scanner, n, m)
		coordinates, sum := challenge.MaximumSumSubRectangle(inputMatrix, n, m)
		/* Top,Left,Bottom,Right = 1,1,2,2 */
		Expect(sum).To(Equal(16))
		Expect(coordinates).To(Equal([]int{1, 1, 2, 2}))
	})
})
