package cracking

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IcecreamParlour", func() {
	var ()

	It("should run 1", func() {
		// The first time, they pool together  money=4 dollars.
		// There are five flavors available that day and flavors 1,4  and  have a total cost of 1+3=4.

		cost, indices := FindIcecreams(ToIcecreams([]int{1, 4, 5, 3, 2}), 4)
		Expect(cost).To(Equal([]int{1, 3}))
		Expect(indices).To(Equal([]int{1, 4}))
	})

	It("should run 2", func() {
		// The first time, they pool together  money=4 dollars.
		// There are five flavors available that day and flavors 1,2  and  have a total cost of 2+2=4.
		cost, indices := FindIcecreams(ToIcecreams([]int{2, 2, 4, 3}), 4)
		Expect(cost).To(Equal([]int{2, 2}))
		Expect(indices).To(Equal([]int{1, 2}))
	})

})
