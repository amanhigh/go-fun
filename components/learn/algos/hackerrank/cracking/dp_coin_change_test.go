package cracking_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
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
		Expect(cracking.Split(money, coins, selectedCoins)).To(Equal(4))
	})
})
