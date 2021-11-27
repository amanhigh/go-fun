package sort

/**
O(n^2) Runtime, O(1) Space Complexity
Each iteration Bubble Largest Element to end of array

Time: n^2, Space: O(1)
*/
func BubbleSort(ints []int, n int) ([]int, int) {
	swapCount := 0
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if ints[i] > ints[j] {
				ints[i], ints[j] = ints[j], ints[i]
				swapCount++
			}
		}
	}
	return ints, swapCount
}
