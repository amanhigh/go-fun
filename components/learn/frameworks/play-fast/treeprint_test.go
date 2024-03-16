package play_fast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xlab/treeprint"
)

var _ = Describe("Treeprint", func() {

	var (
		tree treeprint.Tree
	)
	BeforeEach(func() {
		tree = treeprint.New()
	})

	It("should build", func() {
		Expect(tree).To(Not(BeNil()))
	})

	Context("Branch", func() {
		BeforeEach(func() {
			// create a new branch in the root
			one := tree.AddBranch("1")

			// add some nodes
			one.AddNode("1.1").AddNode("1.2")

			// create a new sub-branch
			one.AddBranch("1.3").
				AddNode("1.3.1").AddNode("1.3.2").    // add some nodes
				AddBranch("1.3.3").                   // add a new sub-branch
				AddNode("1.3.3.1").AddNode("1.3.3.2") // add some nodes too

			// add one more node that should surround the inner branch
			one.AddNode("1.4")

			// add a new node to the root
			tree.AddNode("2")
		})

		It("should print", func() {
			output := tree.String()
			Expect(output).To(ContainSubstring("2"))
		})
	})

})
