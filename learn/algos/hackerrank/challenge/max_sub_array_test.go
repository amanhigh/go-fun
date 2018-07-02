package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"

	"fmt"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
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
10 10
10 11
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
			expected := util.ReadInts(outScan, 2)
			Expect(MaxSubArray(ints)).To(Equal(expected), fmt.Sprintf("Input: %v Expected: %v", ints, expected))
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
		outScan = util.NewStringScanner(`1 1
-1 -1
1 1
6 6
-10 -10
5 6`)
	})

	It("should be correct nonContigous Sum", func() {
		inScan = util.NewStringScanner(`
1
5
-1 2 -3 4 5`)
		outScan = util.NewStringScanner(`
9 11`)
	})
})
