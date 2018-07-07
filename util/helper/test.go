package helper

import (
	"bufio"

	"github.com/amanhigh/go-fun/util"
)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = util.ReadInt(scanner)
	ints = util.ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = util.ReadInts(scanner, n)
	}
	return
}
