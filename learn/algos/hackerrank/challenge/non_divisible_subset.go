package challenge

/**
https://www.hackerrank.com/challenges/non-divisible-subset/problem
https://www.geeksforgeeks.org/subset-no-pair-sum-divisible-k/
C(n,r) - http://www.mathwords.com/c/combination_formula.htm
*/
func NonDivisibleSubset(input []int, k int) (result int) {
	modulos := make([]int, k)
	for _, value := range input {
		modulos[value%k]++
	}

	if 1 < modulos[0] {
		result = 1
	} else {
		result = modulos[0]
	}

	for i := 1; i <= k/2; i++ {
		if modulos[i] < modulos[k-i] {
			result += modulos[k-i]
		} else {
			result += modulos[i]
		}
	}

	return
}
