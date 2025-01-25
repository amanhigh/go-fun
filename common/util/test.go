package util

import (
	"bufio"
)

const DIMENSION_COUNT = 2 // Number of values needed for matrix dimensions (n, m)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = ReadInt(scanner)
	ints = ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n, m int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = ReadInts(scanner, m)
	}
	return
}

func ReadMatrixWithDimensions(scanner *bufio.Scanner) (matrix [][]int, n, m int) {
	ints := ReadInts(scanner, DIMENSION_COUNT)
	n = ints[0]
	m = ints[1]
	matrix = ReadMatrix(scanner, n, m)
	return
}
