package practice

import (
	"sort"
)

/*
*
Find smallest postive mising number in array.
*/
func MissingNumbers(a []int) (missingNumber int) {
	var positiveNumbers []int
	for _, v := range a {
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
		switch {
		case missingNumber == v:
			continue
		case v == missingNumber+1:
			missingNumber++
		default:
			break
		}
	}
	missingNumber++
	return
}

func TargetSum(input []int, target int) (i, j int) {
	numMap := map[int]int{}

	for i, v := range input {
		balance := target - v
		// Store Number with Index
		numMap[v] = i

		// Search Balance Required in Map
		j, ok := numMap[balance]
		if ok {
			// Return Index of Current and Balance as result
			return j, i
		}
	}

	return
}
