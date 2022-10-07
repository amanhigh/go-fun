package cracking

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShortestReach", func() {
	var (
		nodeMap   map[int]*GraphNode
		startNode *GraphNode
	)

	JustBeforeEach(func() {
		startNode.Distance = 0

		/* Find Paths to Remaining Nodes */
		Travel(startNode)
	})

	Context("Case 1", func() {

		BeforeEach(func() {
			// Create Graph
			nodeMap = generateNodes(4)
			//Link Nodes
			linkNodes(nodeMap, 1, 2)
			linkNodes(nodeMap, 1, 3)

			//Set Start Node with Distance Set to Zero
			startNode = nodeMap[1]
		})

		It("should have no distance to start node", func() {
			Expect(startNode.Distance).To(Equal(0))
		})

		It("should have reachable Nodes", func() {
			Expect(nodeMap[2].Distance).To(Equal(EDGE_WEIGHT))
			Expect(nodeMap[3].Distance).To(Equal(EDGE_WEIGHT))
		})

		It("should unreachable nodes", func() {
			Expect(nodeMap[4].Distance).To(Equal(-1))
		})
	})

	Context("Case 2", func() {

		BeforeEach(func() {
			// Create Graph
			nodeMap = generateNodes(3)
			//Link Nodes
			linkNodes(nodeMap, 2, 3)

			//Set Start Node with Distance Set to Zero
			startNode = nodeMap[2]
		})

		It("should have no distance to start node", func() {
			Expect(startNode.Distance).To(Equal(0))
		})

		It("should have reachable Nodes", func() {
			Expect(nodeMap[3].Distance).To(Equal(EDGE_WEIGHT))
		})

		It("should unreachable nodes", func() {
			Expect(nodeMap[1].Distance).To(Equal(-1))
		})
	})

})
