package sort

func QuickSort(input []int, start int, end int) {
	//fmt.Println("Quick", start, end, input[start:end+1])
	/* Sort only if more than one element in Segment */
	if start < end {
		/* Split Problem */
		pIndex := Partition(input, start, end)
		/* Solve Subproblems */
		QuickSort(input, start, pIndex-1)
		QuickSort(input, pIndex+1, end)
	}
}

/**
Ensure everything less than pivot is moved left of Partition Index (pIndex)
Post this everything on left of pivot is less than pivot and right is greater than pivot
*/
func Partition(input []int, start int, end int) (pIndex int) {
	pivot := input[end]
	pIndex = start

	/* Scan from start to end (-1 as last is pivot itself), finding anything less than pivot and moving it left of pIndex */
	for i := start; i < end; i++ {
		if input[i] <= pivot {
			input[i], input[pIndex] = input[pIndex], input[i]
			pIndex++
		}
	}

	//fmt.Println("Partition", start, pIndex, end, input[start:end+1])
	/* Place Pivot at end of partition
	(No increment of pIndex unlike in Loop as its last placement)
	*/
	input[pIndex], input[end] = input[end], input[pIndex]
	return
}
