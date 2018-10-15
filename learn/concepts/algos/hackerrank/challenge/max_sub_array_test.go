package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"

	"fmt"

	. "github.com/amanhigh/go-fun/learn/concepts/algospts/algos/hackerrank/challenge"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
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
		n := util.ReadInt(inScan)
		for i := 0; i < n; i++ {
			_, ints := helper.ReadCountInts(inScan)
			expected := util.ReadInts(outScan, 4)
			arraySum, segmentSum, start, end := MaxSubArray(ints)
			Expect([]int{arraySum, segmentSum, start, end}).To(Equal(expected), fmt.Sprintf("Input: %v Expected: %v", ints, expected))
		}
	})

	It("should be correct", func() {
		inScan = util.NewStringScanner(input)
		outScan = util.NewStringScanner(output)
	})

	It("should be correct 1", func() {
		inScan = util.NewStringScanner(`6
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
		outScan = util.NewStringScanner(`1 1 0 0
-1 -1 0 0
1 1 0 0
6 6 0 2
-10 -10 0 0
5 6 5 5`)
	})

	It("should be correct nonContigous Sum", func() {
		inScan = util.NewStringScanner(`
1
5
-1 2 -3 4 5`)
		outScan = util.NewStringScanner(`
9 11 3 4`)
	})
})
