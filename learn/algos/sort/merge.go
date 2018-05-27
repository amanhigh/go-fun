package sort

func MergeSort(input []int, start int, end int) (inversion int) {
	/* End If we have Single Element Left */
	if start < end {
		/* Break Problem */
		mid := (start + end) / 2
		/* Solve Subproblems */
		inversion += MergeSort(input, start, mid)
		inversion += MergeSort(input, mid+1, end)
		inversion += Merge(input, start, mid, end)
	}
	return
}

func Merge(input []int, start int, mid int, end int) (inversion int) {
	result := make([]int, end-start+1)
	i, j, k := start, mid+1, 0

	/* Copy Minimum of Left & Right */
	for i <= mid && j <= end {
		/* Made Less than equal to consider inversions only if not equal */
		if input[i] <= input[j] {
			result[k] = input[i]
			i++
			k++
		} else {
			result[k] = input[j]
			/* https://www.youtube.com/watch?v=k9RQh21KrH8 */
			inversion += mid - i + 1
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

	//fmt.Println("Range:", start, mid, end, input[start:end+1], input[start:mid+1], input[mid+1:end+1], result, inversion)
	copy(input[start:end+1], result)
	return
}
