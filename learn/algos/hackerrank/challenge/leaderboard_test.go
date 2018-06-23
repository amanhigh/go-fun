package challenge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/hackerrank/challenge"
)

var _ = Describe("Leaderboard", func() {
	var (
		leaderBoard      = []int{100, 100, 50, 40, 40, 20, 10}
		games            = []int{5, 25, 50, 120}
		expectedPosition = []int{6, 4, 2, 1}
	)

	It("should compute leader board", func() {
		Expect(LeaderBoard(leaderBoard, games)).To(Equal(expectedPosition))
	})
})
