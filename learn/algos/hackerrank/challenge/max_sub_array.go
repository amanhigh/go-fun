package challenge

/**
	We define subsequence as any subset of an array. We define a subarray as a contiguous subsequence in an array.

	Given an array, find the maximum possible sum among:
    * all nonempty subarrays.
    * all nonempty subsequences.

	https://www.hackerrank.com/challenges/maxsubarray/problem
*/
func MaxSubArray(input []int) (result []int) {
	return MaxSubArrayBruteForce(input)
}

/**
Brute Force O(n^2)
*/
func MaxSubArrayBruteForce(input []int) (result []int) {
	contigousSum := 0
	nonContigousSum := 0
	n := len(input)
	for i := 0; i < n; i++ {
		sum := 0
		/* Consider Segment from i to j over all possiblities */
		for j := i; j < n; j++ {
			sum += input[j]
			if sum > contigousSum {
				/* Subarry elements must be placed next to each other */
				contigousSum = sum
			}
			//fmt.Println(i, j, input[i:j+1], sum, contigousSum)
		}

		/* Sub Segment Sum may have gaps */
		if nonContigousSum < nonContigousSum+input[i] {
			nonContigousSum += input[i]
		}
	}
	return []int{contigousSum, nonContigousSum}
}
