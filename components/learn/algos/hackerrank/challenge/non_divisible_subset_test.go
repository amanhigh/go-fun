package challenge_test

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/challenge"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NonDivisibleSubset", func() {

	It("should work case 1", func() {
		ints, k := readInputNonDivisible(`4 3
1 7 2 4`)
		Expect(challenge.NonDivisibleSubset(ints, k)).To(Equal(3))
	})

	It("should work case 2", func() {
		ints, k := readInputNonDivisible(`5 5
2 7 12 17 22`)
		Expect(challenge.NonDivisibleSubset(ints, k)).To(Equal(5))
	})

	It("should work case 3", func() {
		ints, k := readInputNonDivisible(`10 4
1 2 3 4 5 6 7 8 9 10`)
		Expect(challenge.NonDivisibleSubset(ints, k)).To(Equal(5))
	})
})

func readInputNonDivisible(input string) (ints []int, k int) {
	scanner := util.NewStringScanner(input)
	n := util.ReadInt(scanner)
	k = util.ReadInt(scanner)
	ints = util.ReadInts(scanner, n)
	return
}
