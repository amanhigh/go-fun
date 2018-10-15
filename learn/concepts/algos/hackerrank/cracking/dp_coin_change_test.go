package cracking_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/concepts/algospts/algos/hackerrank/cracking"
)

var _ = Describe("DpCoinChange", func() {
	var (
		coins         = []int{1, 2, 3}
		money         = 4
		selectedCoins []int
	)
	It("should compute possibilities", func() {
		Expect(Split(money, coins, selectedCoins)).To(Equal(4))
	})
})
