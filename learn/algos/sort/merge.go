package sort

func MergeSort(input []int, start int, end int) {

	if end-start > 0 {
		/* Break Problem */
		mid := (start + end) / 2
		/* Solve Subproblems */
		MergeSort(input, start, mid)
		MergeSort(input, mid+1, end)
		Merge(input, start, mid, end)
	}
}

func Merge(input []int, start int, mid int, end int) {
	result := make([]int, end-start+1)
	i, j, k := start, mid+1, 0

	/* Copy Minimum of Left & Right */
	for ; i <= mid && j <= end; {
		if input[i] < input[j] {
			result[k] = input[i]
			i++
			k++
		} else {
			result[k] = input[j]
			j++
			k++
		}
	}

	/* Copy Remaining */
	for ; i <= mid; i++ {
		result[k] = input[i]
		k++
	}

	for ; j <= end; j++ {
		result[k] = input[j]
		k++
	}

	copy(input[start:end+1], result)
}