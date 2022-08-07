package cracking_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
)

var _ = Describe("Tries", func() {
	var (
		node = cracking.NewNode()
	)

	It("should build", func() {
		Expect(node).To(Not(BeNil()))
	})

	Context("Add Nodes", func() {
		BeforeEach(func() {
			cracking.Add(node, "hack")
			cracking.Add(node, "hackerrank")
		})

		It("should find hac", func() {
			Expect(cracking.Find(node, "hac")).To(Equal(2))
		})

		It("should not find hak", func() {
			Expect(cracking.Find(node, "hak")).To(Equal(0))
		})
	})

})
