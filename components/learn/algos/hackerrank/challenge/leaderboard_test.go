package challenge_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/challenge"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Leaderboard", func() {
	var (
		leaderBoard      = []int{100, 100, 50, 40, 40, 20, 10}
		games            = []int{5, 25, 50, 120}
		expectedPosition = []int{6, 4, 2, 1}
	)

	It("should compute leader board", func() {
		Expect(challenge.LeaderBoard(leaderBoard, games)).To(Equal(expectedPosition))
	})
})
