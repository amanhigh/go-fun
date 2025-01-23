package cracking

type GraphNode struct {
	Data     int
	Nodes    []*GraphNode
	Distance int
}

const EDGE_WEIGHT = 6

/*
*
Perform Travel computing distances
*/
func Travel(start *GraphNode) {
	nextHopDistance := start.Distance + EDGE_WEIGHT
	for _, node := range start.Nodes {
		/* Update Non Visited Nodes If this is shortest Path */
		if node.Distance == -1 || node.Distance > nextHopDistance {
			node.Distance = nextHopDistance

			/* Traverse Further if we have discovered a new shortest path */
			Travel(node)
		}
	}
}

/* Helpers */
func linkNodes(nodeMap map[int]*GraphNode, u int, v int) {
	nodeMap[u].Nodes = append(nodeMap[u].Nodes, nodeMap[v])
	nodeMap[v].Nodes = append(nodeMap[v].Nodes, nodeMap[u])
}

/*
*
Generate N Nodes starting from 1 to N.
*/
func generateNodes(n int) (nodeMap map[int]*GraphNode) {
	nodeMap = map[int]*GraphNode{}
	/* Starting index 1 as per input data */
	for i := 1; i <= n; i++ {
		nodeMap[i] = &GraphNode{Data: i, Nodes: []*GraphNode{}, Distance: -1}
	}
	return
}
