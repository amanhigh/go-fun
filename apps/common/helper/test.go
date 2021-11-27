package helper

import (
	"bufio"
	util2 "github.com/amanhigh/go-fun/apps/common/util"

	"fmt"
)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = util2.ReadInt(scanner)
	ints = util2.ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n, m int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = util2.ReadInts(scanner, m)
	}
	return
}

func ReadMatrixWithDimensions(scanner *bufio.Scanner) (matrix [][]int, n, m int) {
	ints := util2.ReadInts(scanner, 2)
	fmt.Println(ints)
	n = ints[0]
	m = ints[1]
	matrix = ReadMatrix(scanner, n, m)
	return
}
