package challenge

import (
	"math"
)

/*
*

		We define subsequence as any subset of an array. We define a subarray as a contiguous subsequence in an array.

		Given an array, find the maximum possible sum among:
	    * all nonempty subarrays.
	    * all nonempty subsequences.

		https://www.hackerrank.com/challenges/maxsubarray/problem
*/
// nolint: revive
func MaxSubArray(input []int) (arraySum, segmentSum, start, end int) {
	return KadensAlgorithm(input)
}

/*
*
https://www.youtube.com/watch?v=86CQq3pKSUw
*/
func calculateContiguousSum(input []int) (sum, start, end int) {
	sum = input[0]
	currentSum := 0

	for i, value := range input {
		/*
			If current sum plus current value is less than current value,
			we should start a new subarray from this point
		*/
		if currentSum+value < value {
			/* To Handle decrementing array -1,-2,-3,-4 */
			if currentSum < value {
				/* Mark start only if current value is greater than previous sum */
				start = i
			}
			currentSum = value
		} else {
			/*
				This element is part of max subarray hence previous max
				subarray plus this element
			*/
			currentSum += value
		}

		/*
			Global Sum can be more than sum in between
			as its not max sum at this index, its max sum
			till now over any index.
		*/
		if sum < currentSum {
			sum = currentSum
			end = i
		}
	}
	return
}

func calculateNonContiguousSum(input []int) int {
	sum := 0
	for _, value := range input {
		if value > 0 {
			sum += value
		}
	}
	return sum
}

// nolint: revive
func KadensAlgorithm(input []int) (contigousSum, nonContigousSum, start, end int) {
	contigousSum, start, end = calculateContiguousSum(input)
	nonContigousSum = calculateNonContiguousSum(input)

	/*
		#Mistake 2 Fix
		If all values are negative then we would have not included anything
		in contigous sum and it must be equal to global sum.

		If there is even one positive value then contigousSum can't be less than 0
	*/
	if contigousSum < 0 {
		nonContigousSum = contigousSum
	}
	return
}

/*
*
Brute Force O(n^2)
*/
// nolint: revive
func MaxSubArrayBruteForce(input []int) (contigousSum, nonContigousSum, start, end int) {
	/*
		Mistake #1 as array can have negative numbers
		Sums should start negative.
	*/
	contigousSum, nonContigousSum = -math.MaxInt32, 0
	n := len(input)
	for i := 0; i < n; i++ {
		sum := 0
		/* Consider Segment from i to j over all possiblities */
		for j := i; j < n; j++ {
			sum += input[j]
			if sum > contigousSum {
				/* Subarry elements must be placed next to each other */
				start = i
				end = j
				contigousSum = sum
			}
			// fmt.Println(i, j, input[i:j+1], sum, contigousSum)
		}

		/*
			#Mistake 2 Check this Commit
			Didn't correctly find nonContigous Sum
			Include only positives values and handle all negative values case
			at end.
		*/
		if input[i] > 0 {
			nonContigousSum += input[i]
		}
	}

	if contigousSum < 0 {
		nonContigousSum = contigousSum
	}

	return
}
