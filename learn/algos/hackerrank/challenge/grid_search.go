package challenge

import (
	"fmt"
	"strings"
)

/**
https://www.hackerrank.com/challenges/the-grid-search/problem
*/

func GridSearch(grid, search []string) (ok bool) {
	fmt.Println(grid, search)
	found := 0
	for _, pattern := range search {
		for _, line := range grid {
			if strings.Contains(line, pattern) {
				found++
				fmt.Println(line, pattern, found)
			}

			if found == len(search) {
				return true
			}
		}
	}
	return
}
