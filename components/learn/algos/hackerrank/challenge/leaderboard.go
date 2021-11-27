package challenge

import "math"

/**

 *   The player with the highest score is ranked number on the leaderboard.

*    Players who have equal scores receive the same ranking number, and the next player(s) receive the immediately following ranking number.


https://www.hackerrank.com/challenges/climbing-the-leaderboard/problem
*/
func LeaderBoard(leaderBoard, games []int) (gameRankings []int) {
	/* Extract Ranks and Scores, O(n) */
	ranks := getRanks(leaderBoard)

	/* Remember Last Rank Index that was not beaten */
	nextRank := len(ranks) - 1
	/* Evaluate each Game, O(m) */
	for _, game := range games {
		/* Assume you are top player */
		gameRank := 1

		/* Start from last rank which was not beaten until end */
		for i := nextRank; i >= 0; i-- {
			/* If your score is less than of some player */
			if game < ranks[i] {
				/* Take Rank below his, his rank is i+1, so yours is i+2 */
				gameRank = i + 2
				/* Remember this index so next game score can be compared to him */
				nextRank = i
				break
			}
		}

		/* Collect Computed Rank */
		gameRankings = append(gameRankings, gameRank)
	}

	return
}

func getRanks(leaderBoard []int) []int {
	var ranks []int
	topScore := math.MaxInt32
	for _, score := range leaderBoard {
		if score < topScore {
			ranks = append(ranks, score)
			topScore = score
		}
	}
	return ranks
}
