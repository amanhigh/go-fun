package challenge

import (
	"math"
)

/*
*
cooridantes = top,left,bottom,right
*/
func MaximumSumSubRectangle(input [][]int, n, m int) (coordinates []int, maxSum int) {
	return MaximumSumSubRectangleSmart(input, n, m)
}

/*
*
n - Rows
m - Columns
O(nm^2)

https://www.youtube.com/watch?v=yCQN096CwWM
*/
func MaximumSumSubRectangleSmart(input [][]int, n, m int) (coordinates []int, maxSum int) {
	maxSum = -math.MaxInt32
	for jStart := range m {
		columnSum := make([]int, n)
		for jEnd := jStart; jEnd < m; jEnd++ {
			/* Sum Column to previous Sum O(n)*/
			for i := range n {
				columnSum[i] += input[i][jEnd]
			}
			/* O(n) */
			sum, _, iStart, iEnd := KadensAlgorithm(columnSum)
			if sum > maxSum {
				maxSum = sum
				coordinates = []int{iStart, jStart, iEnd, jEnd}
			}
			// fmt.Println(iStart, jStart, iEnd, jEnd, columnSum, sum)
		}
	}
	return
}

/*
*
Vary top-left and bottom-right for all possible combinations and sum.
O(n^4) or O(n^2m^2) incase fo rectangle
*/
func MaximumSumSubRectangleBruteForce(input [][]int, n, m int) (coordinates []int, maxSum int) {
	maxSum = -math.MaxInt32
	for iStart := range n {
		for jStart := range m {
			for iEnd := range n {
				for jEnd := range m {
					cords := []int{iStart, jStart, iEnd, jEnd}
					sum := SumSquare(input, cords)
					if sum > maxSum {
						coordinates = cords
						maxSum = sum
					}
				}
			}
		}
	}
	return
}

/*
*
cooridantes = top,left,bottom,right
O(n^2)
*/
func SumSquare(input [][]int, coordinates []int) (sum int) {
	// Top -> Bottom Row
	for i := coordinates[0]; i <= coordinates[2]; i++ {
		// Left -> Right Column
		for j := coordinates[1]; j <= coordinates[3]; j++ {
			sum += input[i][j]
		}
	}
	// fmt.Println(coordinates, sum)
	return
}
