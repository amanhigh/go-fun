package cracking_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
)

var _ = Describe("Staircase", func() {
	var (
		steps = []int{1, 2, 3}
		taken = []int{}
	)
	It("should run brute force", func() {
		Expect(cracking.Staircase(1, steps, taken)).To(Equal(1))
		Expect(cracking.Staircase(3, steps, taken)).To(Equal(4))
		Expect(cracking.Staircase(7, steps, taken)).To(Equal(44))
	})

	It("should run dynamic programming", func() {
		Expect(cracking.StaircaseDp(1)).To(Equal(1))
		Expect(cracking.StaircaseDp(3)).To(Equal(4))
		Expect(cracking.StaircaseDp(7)).To(Equal(44))
	})

})
