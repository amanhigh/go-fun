package challenge

import (
	"math"
)

/**
cooridantes = top,left,bottom,right
*/
func MaximumSumSubSquare(input [][]int, n, m int) (coordinates []int, maxSum int) {
	return MaximumSumSubSquareSmart(input, n, m)
}

/**
O(nm^2)

https://www.youtube.com/watch?v=yCQN096CwWM
*/
func MaximumSumSubSquareSmart(input [][]int, n, m int) (coordinates []int, maxSum int) {
	maxSum = -math.MaxInt32
	for jStart := 0; jStart < n; jStart++ {
		columnSum := make([]int, n)
		for jEnd := jStart; jEnd < n; jEnd++ {
			/* Sum Column to previous Sum O(n)*/
			for i := 0; i < n; i++ {
				columnSum[i] += input[i][jEnd]
			}
			/* O(n) */
			sum, _, iStart, iEnd := KadensAlgorithm(columnSum)
			if sum > maxSum {
				maxSum = sum
				coordinates = []int{iStart, jStart, iEnd, jEnd}
			}
			//fmt.Println(iStart, jStart, iEnd, jEnd, columnSum, sum)
		}
	}
	return
}

/**
Vary top-left and bottom-right for all possible combinations and sum.
O(n^4) or O(n^2m^2) incase fo rectangle
*/
func MaximumSumSubSquareBruteForce(input [][]int, n, m int) (coordinates []int, maxSum int) {
	maxSum = -math.MaxInt32
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				for l := 0; l < n; l++ {
					cords := []int{i, j, k, l}
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

/**
cooridantes = top,left,bottom,right
O(n^2)
*/
func SumSquare(input [][]int, coordinates []int) (sum int) {
	//Top -> Bottom Row
	for i := coordinates[0]; i <= coordinates[2]; i++ {
		//Left -> Right Column
		for j := coordinates[1]; j <= coordinates[3]; j++ {
			sum += input[i][j]
		}
	}
	//fmt.Println(coordinates, sum)
	return
}
