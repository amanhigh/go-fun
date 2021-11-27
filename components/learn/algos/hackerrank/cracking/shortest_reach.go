package cracking

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type GraphNode struct {
	Data     int
	Nodes    []*GraphNode
	Distance int
}

const EDGE_WAIT = 6

func main() {
	var q, n, m, s, u, v int
	//file:= os.Stdin
	file, _ := os.Open("inputl.txt")
	scanner := bufio.NewScanner(file)
	/* Read Query Count */
	scanner.Scan()
	q, _ = strconv.Atoi(scanner.Text())
	/* Process Queries */
	for i := 0; i < q; i++ {
		/* Read Graph Dimensions */
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)
		//fmt.Println("Count:", n, m)
		/* Generate Nodes */
		nodeMap := GenerateNodes(n)
		/* Generate Edges */
		for j := 0; j < m; j++ {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "%d %d", &u, &v)
			//fmt.Printf("%v --> %v\n", u, v)
			/* Link Both u & v as Graph is undirected */
			nodeMap[u].Nodes = append(nodeMap[u].Nodes, nodeMap[v])
			nodeMap[v].Nodes = append(nodeMap[v].Nodes, nodeMap[u])
		}
		/* Read Starting Node & Set Distance Zero */
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d", &s)
		startNode := nodeMap[s]
		startNode.Distance = 0
		/* Find Paths to Remaining Nodes */
		Travel(startNode)
		/* Print Distances in Increasing order */
		PrintDistances(startNode, nodeMap, n)
	}
}

/**
Print Distances except starting Node
*/
func PrintDistances(startNode *GraphNode, nodeMap map[int]*GraphNode, n int) {
	/* Print all nodes except start */
	for i := 1; i <= n; i++ {
		if i != startNode.Data {
			fmt.Print(nodeMap[i].Distance, " ")
		}
	}
	fmt.Println("")
}

/**
Perform Travel computing distances
*/
func Travel(start *GraphNode) {
	nextHopDistance := start.Distance + EDGE_WAIT
	for _, node := range start.Nodes {
		/* Update Non Visited Nodes If this is shortest Path */
		if node.Distance == -1 || node.Distance > nextHopDistance {
			node.Distance = nextHopDistance

			/* Traverse Further if we have discovered a new shortest path */
			Travel(node)
		}
	}
}

/**
Generate N Nodes starting from 1 to N.
*/
func GenerateNodes(n int) (nodeMap map[int]*GraphNode) {
	nodeMap = map[int]*GraphNode{}
	/* Starting index 1 as per input data */
	for i := 1; i <= n; i++ {
		nodeMap[i] = &GraphNode{Data: i, Nodes: []*GraphNode{}, Distance: -1}
	}
	return
}
