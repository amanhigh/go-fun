package helper

import (
	"bufio"
	"github.com/amanhigh/go-fun/common/util"

	"fmt"
)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = util.ReadInt(scanner)
	ints = util.ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n, m int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = util.ReadInts(scanner, m)
	}
	return
}

func ReadMatrixWithDimensions(scanner *bufio.Scanner) (matrix [][]int, n, m int) {
	ints := util.ReadInts(scanner, 2)
	fmt.Println(ints)
	n = ints[0]
	m = ints[1]
	matrix = ReadMatrix(scanner, n, m)
	return
}
