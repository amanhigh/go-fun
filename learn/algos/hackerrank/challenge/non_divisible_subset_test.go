package challenge_test

import (
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	challenge2 "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NonDivisibleSubset", func() {

	It("should work case 1", func() {
		ints, k := readInputNonDivisible(`4 3
1 7 2 4`)
		Expect(challenge2.NonDivisibleSubset(ints, k)).To(Equal(3))
	})

	It("should work case 2", func() {
		ints, k := readInputNonDivisible(`5 5
2 7 12 17 22`)
		Expect(challenge2.NonDivisibleSubset(ints, k)).To(Equal(5))
	})

	It("should work case 3", func() {
		ints, k := readInputNonDivisible(`10 4
1 2 3 4 5 6 7 8 9 10`)
		Expect(challenge2.NonDivisibleSubset(ints, k)).To(Equal(5))
	})
})

func readInputNonDivisible(input string) (ints []int, k int) {
	scanner := util2.NewStringScanner(input)
	n := util2.ReadInt(scanner)
	k = util2.ReadInt(scanner)
	ints = util2.ReadInts(scanner, n)
	return
}
