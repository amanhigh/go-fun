package cracking_test

import (
	cracking2 "github.com/amanhigh/go-fun/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DpCoinChange", func() {
	var (
		coins         = []int{1, 2, 3}
		money         = 4
		selectedCoins []int
	)
	It("should compute possibilities", func() {
		Expect(cracking2.Split(money, coins, selectedCoins)).To(Equal(4))
	})
})
