package sort

/*
*

		Conceptually, a merge sort works as follows:

	    Divide the unsorted list into n sublists, each containing 1 element (a list of 1 element is considered sorted).
	    Repeatedly merge sublists to produce new sorted sublists until there is only 1 sublist remaining. This will be the sorted list.

		Produces Stable Sort, Space O(n), Time nlog(n)
*/
func MergeSort(input []int, start, end int) (inversion int) {
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

// nolint:funlen
func Merge(input []int, start, mid, end int) (inversion int) {
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
			/*
				Eg.
				1,2,8,4,5 has two inversions 8,4 and 8,5
				Inversion is a[i] > a[j] where i<j
				Inversion tells how far we are from sorted array.

				Since left and right segments are sorted.
				When input[i] > input[j] then number at location i forms inversion with
				all numbers right of it i.e. mid-i+1

				https://www.youtube.com/watch?v=k9RQh21KrH8
			*/
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

	// fmt.Println("Range:", start, mid, end, input[start:end+1], input[start:mid+1], input[mid+1:end+1], result, inversion)
	copy(input[start:end+1], result)
	return
}
