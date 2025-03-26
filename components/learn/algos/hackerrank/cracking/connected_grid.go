package cracking

/*
*
Visited all -1 before start
n-rows
m-columns

https://www.hackerrank.com/challenges/ctci-connected-cell-in-a-grid/problem
*/
func FindRegion(cells, visited [][]int, n, m int) (maxSize int) {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			current := FindRegionRecursive(cells, visited, n, m, i, j)
			if current > maxSize {
				maxSize = current
			}
		}
	}

	return
}

func FindRegionRecursive(cells, visited [][]int, n, m, j, i int) (maxSize int) {
	/* If we are within Cell Bound and current cell is 1 & never visited Proceed */
	if 0 <= i && i < n && 0 <= j && j < m && cells[i][j] == 1 && visited[i][j] == -1 {
		// fmt.Println(i, j, cells[i][j])
		maxSize = 1
		visited[i][j] = 1 // i,j = row,column
		/* Go Left, Right */
		maxSize += FindRegionRecursive(cells, visited, n, m, j-1, i)
		maxSize += FindRegionRecursive(cells, visited, n, m, j+1, i)
		/* Go Top & Bottom */
		maxSize += FindRegionRecursive(cells, visited, n, m, j, i+1)
		maxSize += FindRegionRecursive(cells, visited, n, m, j, i-1)
		/* Go Upper Left, Upper Right */
		maxSize += FindRegionRecursive(cells, visited, n, m, j-1, i-1)
		maxSize += FindRegionRecursive(cells, visited, n, m, j+1, i-1)
		/* Go Bottom Left, Bottom Right */
		maxSize += FindRegionRecursive(cells, visited, n, m, j-1, i+1)
		maxSize += FindRegionRecursive(cells, visited, n, m, j+1, i+1)
	}
	return
}
