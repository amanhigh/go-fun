package challenge

import "github.com/thoas/go-funk"

/**
https://www.hackerrank.com/challenges/non-divisible-subset/problem

C(n,r) - http://www.mathwords.com/c/combination_formula.htm
*/
func NonDivisibleSubset(input []int, k int) (result int) {
	return NonDivisibleSubsetRecursive(input, []int{}, k, 2)
}

func NonDivisibleSubsetRecursive(input, permute []int, k, r int) (result int) {
	if r == len(permute) {
		//fmt.Println(permute)
		sum := funk.SumInt(permute)
		if sum%k == 0 {
			result++
		}
	}

	for i, val := range input {
		result += NonDivisibleSubsetRecursive(input[i+1:], append(permute, val), k, r)
	}
	return
}
