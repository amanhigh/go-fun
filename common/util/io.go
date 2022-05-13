package util

import (
	"bufio"
	"fmt"
	"os"
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

func ReadStrings(scanner *bufio.Scanner, n int) []string {
	a := make([]string, n)
	for i := 0; i < n && scanner.Scan(); i++ {
		a[i] = scanner.Text()
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

func NewFileScanner(path string) (scanner *bufio.Scanner, err error) {
	var file *os.File
	if file, err = os.Open(path); err == nil {
		scanner = bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
	}
	return
}
