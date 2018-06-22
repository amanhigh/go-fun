package main

import (
	"bufio"
	"fmt"
	"os"

	"strconv"
)

func main() {
	//var n int
	//file:= os.Stdin
	file, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	/* Read Query Count */
	//scanner.Scan()
	//fmt.Sscanf(scanner.Text(), "%d", &n)
	ints := ReadInts(scanner, 4)
	meet := KangarooMeet(ints)
	if meet {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

func ReadInts(scanner *bufio.Scanner, n int) []int {
	a := make([]int, n)
	for i := 0; i < n && scanner.Scan(); i++ {
		if value, err := strconv.Atoi(scanner.Text()); err == nil {
			a[i] = value
		}
	}
	return a
}
