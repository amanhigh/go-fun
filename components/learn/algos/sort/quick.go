package sort

/*
*
The steps are:

	    Pick an element, called a pivot, from the array.
	    1. Partitioning: reorder the array so that all elements with values less than the pivot come before the pivot,
		while all elements with values greater than the pivot come after it (equal values can go either way).
		After this partitioning, the pivot is in its final position. This is called the partition operation.

		2.Recursively apply the above steps to the sub-array of elements with smaller values and
		separately to the sub-array of elements with greater values.

Time: nlog(n), n^2 (worst), Space: log(n)
http://bigocheatsheet.com/
*/
func QuickSort(input []int, start, end int) {
	// fmt.Println("Quick", start, end, input[start:end+1])
	/* Sort only if more than one element in Segment */
	if start < end {
		/* Split Problem */
		pIndex := Partition(input, start, end)
		/* Solve Subproblems */
		QuickSort(input, start, pIndex-1)
		QuickSort(input, pIndex+1, end)
	}
}

/*
*
Ensure everything less than pivot is moved left of Partition Index (pIndex)
Post this everything on left of pivot is less than pivot and right is greater than pivot
*/
func Partition(input []int, start, end int) (pIndex int) {
	pivot := input[end]
	pIndex = start

	/* Scan from start to end (-1 as last is pivot itself), finding anything less than pivot and moving it left of pIndex */
	for i := start; i < end; i++ {
		if input[i] <= pivot {
			input[i], input[pIndex] = input[pIndex], input[i]
			pIndex++
		}
	}

	// fmt.Println("Partition", start, pIndex, end, input[start:end+1])
	/* Place Pivot at end of partition
	(No increment of pIndex unlike in Loop as its last placement)
	*/
	input[pIndex], input[end] = input[end], input[pIndex]
	return
}
