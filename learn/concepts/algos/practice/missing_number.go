package practice

import (
	"sort"
)

/**
Find smallest postive mising number in array.
*/
func MissingNumbers(A []int) (missingNumber int) {
	var positiveNumbers []int
	for _, v := range A {
		if v > 0 {
			positiveNumbers = append(positiveNumbers, v)
		}
	}
	sort.Ints(positiveNumbers)

	if len(positiveNumbers) == 0 || positiveNumbers[0] != 1 {
		missingNumber = 1
		return
	}

	missingNumber = 1
	for _, v := range positiveNumbers {
		if missingNumber == v {
			continue
		} else if v == missingNumber+1 {
			missingNumber = missingNumber + 1
		} else {

			break
		}
	}
	missingNumber = missingNumber + 1
	return
}
