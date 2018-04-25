package cracking

type Node struct {
	NodeMap map[rune]*Node
	Words   int
}

func NewNode() *Node {
	return &Node{
		NodeMap: map[rune]*Node{},
	}
}

func Add(node *Node, query string) {
	node.Words++
	if len(query) > 0 {
		/* Until query finishes follow or create nodes */
		c := rune(query[0])
		cNode, ok := node.NodeMap[c]
		if !ok {
			cNode = NewNode()
			node.NodeMap[c] = cNode
		}

		Add(cNode, query[1:])
	}

	return
}

func Find(node *Node, query string) (matches int) {
	/* Follow query until complete */
	if len(query) > 0 {
		c := rune(query[0])
		if cNode, ok := node.NodeMap[c]; ok {
			matches = Find(cNode, query[1:])
		}
		/* Query didn't match hence Zero */
	} else {
		/* Reached query end this node holds count */
		matches = node.Words
	}
	return
}
