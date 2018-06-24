package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
)

var _ = Describe("MaxSubArray", func() {
	var (
		input = `2
4
1 2 3 4
4
1 2 3 4
6
2 -1 2 3 4 -5
5
2 -1 2 3 4
10 11`
	)

	It("should be correct", func() {
		scanner := util.NewStringScanner(input)
		n := util.ReadInt(scanner)
		for i := 0; i < n; i++ {
			_, ints := helper.ReadCountInts(scanner)
			_, expected := helper.ReadCountInts(scanner)
			Expect(MaxSubArray(ints)).To(Equal(expected))
		}
	})
})
