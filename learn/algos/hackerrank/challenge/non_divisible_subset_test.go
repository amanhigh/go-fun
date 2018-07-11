package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	"github.com/amanhigh/go-fun/util"
)

var _ = Describe("NonDivisibleSubset", func() {

	It("should work case 1", func() {
		ints, k := readInputNonDivisible(`4 3
1 7 2 4`)
		Expect(NonDivisibleSubset(ints, k)).To(Equal(3))
	})

	It("should work case 2", func() {
		ints, k := readInputNonDivisible(`5 5
2 7 12 17 22`)
		Expect(NonDivisibleSubset(ints, k)).To(Equal(5))
	})
})

func readInputNonDivisible(input string) (ints []int, k int) {
	scanner := util.NewStringScanner(input)
	n := util.ReadInt(scanner)
	k = util.ReadInt(scanner)
	ints = util.ReadInts(scanner, n)
	return
}
