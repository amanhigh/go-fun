package challenge_test

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/challenge"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"bufio"
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
		n, m, matrix := readInput(scanner)
		coordinates, sum := challenge.MaximumSumSubRectangle(matrix, n, m)
		/* Top,Left,Bottom,Right = 1,1,2,2 */
		Expect(sum).To(Equal(16))
		Expect(coordinates).To(Equal([]int{1, 1, 2, 2}))
	})

	It("should compute case 2", func() {
		input := [][]int{
			{2, 1, -3, -4, 5},
			{0, 6, 3, 4, 1},
			{2, -2, -1, 4, -5},
			{-3, 3, 1, 0, 3},
		}
		coordinates, sum := challenge.MaximumSumSubRectangle(input, 4, 5)
		/* Top,Left,Bottom,Right = 1,1,2,2 */
		Expect(sum).To(Equal(18))
		Expect(coordinates).To(Equal([]int{1, 1, 3, 3}))
	})

	It("should compute case 3", func() {
		input := [][]int{
			{1, 2, -1, -4, -20},
			{-8, -3, 4, 2, 1},
			{3, 8, 10, 1, 3},
			{-4, -1, 1, 7, -6},
		}
		coordinates, sum := challenge.MaximumSumSubRectangle(input, 4, 5)
		/* Top,Left,Bottom,Right = 1,1,2,2 */
		Expect(sum).To(Equal(29))
		Expect(coordinates).To(Equal([]int{1, 1, 3, 3}))
	})

	It("should compute brute force", func() {
		input := [][]int{
			{1, 2, -1, -4, -20},
			{-8, -3, 4, 2, 1},
			{3, 8, 10, 1, 3},
			{-4, -1, 1, 7, -6},
		}
		coordinates, sum := challenge.MaximumSumSubRectangleBruteForce(input, 4, 5)
		/* Top,Left,Bottom,Right = 1,1,2,2 */
		Expect(sum).To(Equal(29))
		Expect(coordinates).To(Equal([]int{1, 1, 3, 3}))
	})

})

func readInput(scanner *bufio.Scanner) (n, m int, matrix [][]int) {
	n = util.ReadInt(scanner)
	m = util.ReadInt(scanner)
	matrix = util.ReadMatrix(scanner, n, m)
	return
}
