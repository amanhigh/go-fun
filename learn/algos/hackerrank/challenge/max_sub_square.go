package challenge

import (
	"fmt"
	"math"
)

/**
cooridantes = top,left,bottom,right
*/
func MaximumSumSubSquare(input [][]int, n int) (coordinates []int, maxSum int) {
	return MaximumSumSubSquareSmart(input, n)
}

func MaximumSumSubSquareSmart(input [][]int, n int) (coordinates []int, maxSum int) {
	maxSum = -math.MaxInt32
	for jStart := 0; jStart < n; jStart++ {
		columnSum := make([]int, n)
		for jEnd := jStart; jEnd < n; jEnd++ {
			/* Sum Column to previous Sum */
			for i := 0; i < n; i++ {
				columnSum[i] += input[i][jEnd]
			}
			sum := KadensAlgorithm(columnSum)[0]
			if sum > maxSum {
				maxSum = sum
				coordinates = []int{0, jStart, 0, jEnd}
			}
			fmt.Println(jStart, jEnd, columnSum, sum)
		}
	}
	return
}

/**
Vary top-left and bottom-right for all possible combinations and sum.
O(n^4) or O(n^2m^2) incase fo rectangle
*/
func MaximumSumSubSquareBruteForce(input [][]int, n int) (coordinates []int, maxSum int) {
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
