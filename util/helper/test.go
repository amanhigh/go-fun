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
