package util

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

/**
Scanner must be split on words
*/
func ReadInts(scanner *bufio.Scanner, n int) []int {
	a := make([]int, n)
	for i := 0; i < n && scanner.Scan(); i++ {
		if value, err := strconv.Atoi(scanner.Text()); err == nil {
			a[i] = value
		}
	}
	return a
}

func ReadInt(scanner *bufio.Scanner) (n int) {
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &n)
	return
}

func NewStringScanner(s string) (scanner *bufio.Scanner) {
	scanner = bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	return
}
