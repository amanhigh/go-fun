package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	"github.com/amanhigh/go-fun/util"
)

var _ = Describe("NonDivisibleSubset", func() {
	var (
		input = `4 3
1 7 2 4`
	)

	It("should work case 1", func() {
		scanner := util.NewStringScanner(input)
		n := util.ReadInt(scanner)
		k := util.ReadInt(scanner)
		ints := util.ReadInts(scanner, n)
		Expect(NonDivisibleSubset(ints, n, k)).To(Equal(3))
	})
})
