package util

import (
	"bufio"
	"net/http"

	"github.com/amanhigh/go-fun/models/config"
)

// Helper to create server with timeouts from DefaultHttpConfig
func NewTestServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		ReadTimeout:       config.DefaultHttpConfig.ReadTimeout,
		WriteTimeout:      config.DefaultHttpConfig.WriteTimeout,
		IdleTimeout:       config.DefaultHttpConfig.IdleTimeout,
		ReadHeaderTimeout: config.DefaultHttpConfig.ReadHeaderTimeout,
	}
}

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
	ints := ReadInts(scanner, 2)
	n = ints[0]
	m = ints[1]
	matrix = ReadMatrix(scanner, n, m)
	return
}
