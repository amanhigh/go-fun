package challenge_test

import (
	helper "github.com/amanhigh/go-fun/apps/common/helper"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	challenge2 "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"

	"fmt"
)

var _ = Describe("MaxSubArray", func() {
	var (
		input = `2
4
1 2 3 4
6
2 -1 2 3 4 -5`
		output = `
10 10 0 3
10 11 0 4
`
	)

	var (
		inScan  *bufio.Scanner
		outScan *bufio.Scanner
	)

	AfterEach(func() {
		n := util2.ReadInt(inScan)
		for i := 0; i < n; i++ {
			_, ints := helper.ReadCountInts(inScan)
			expected := util2.ReadInts(outScan, 4)
			arraySum, segmentSum, start, end := challenge2.MaxSubArray(ints)
			Expect([]int{arraySum, segmentSum, start, end}).To(Equal(expected), fmt.Sprintf("Input: %v Expected: %v", ints, expected))
		}
	})

	It("should be correct", func() {
		inScan = util2.NewStringScanner(input)
		outScan = util2.NewStringScanner(output)
	})

	It("should be correct 1", func() {
		inScan = util2.NewStringScanner(`6
1
1
6
-1 -2 -3 -4 -5 -6
2
1 -2
3
1 2 3
1
-10
6
1 -1 -1 -1 -1 5`)
		outScan = util2.NewStringScanner(`1 1 0 0
-1 -1 0 0
1 1 0 0
6 6 0 2
-10 -10 0 0
5 6 5 5`)
	})

	It("should be correct nonContigous Sum", func() {
		inScan = util2.NewStringScanner(`
1
5
-1 2 -3 4 5`)
		outScan = util2.NewStringScanner(`
9 11 3 4`)
	})
})
