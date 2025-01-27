package cracking

var stairMem = map[int]int{}

/*
*
Davis has a number of staircases in his house and he likes to climb each staircase 1, 2, or 3 steps at a time.
Being a very precocious child, he wonders how many ways there are to reach the top of the staircase.

Eg.
For n 1,3,7 answer is 1,4,44
[]steps: []int{1,2,3} as described in problem above.
[]taken: Purely for debugging

https://www.hackerrank.com/challenges/ctci-recursive-staircase/problem
*/
func Staircase(n int, steps, taken []int) (result int) {
	/* If Steps are remaining try to take using available steps */
	if n > 0 {
		var ok bool
		/* Try to Search Memory */
		if result, ok = stairMem[n]; !ok {
			/* Try each step, n may go negative if step is greater but it will return */
			for _, step := range steps {
				result += Staircase(n-step, steps, append(taken, step))
			}
			stairMem[n] = result
		}
		// result += Staircase(n, steps[1:], taken)
	} else if n == 0 {
		/* N becomes exactly zero is a valid combination count as result */
		result = 1
		/* Print Steps just to visualize */
		// fmt.Println(taken)
	}
	/* If n is negative it was invalid solution don't count it ,hence result zero */
	return
}

func StaircaseDp(n int) int {
	/* Base Case */
	base := []int{0, 1, 2, 4}
	var store []int

	/* Init Dp Array if not base case */
	if n > 3 {
		store = make([]int, n+1)
		copy(store, base)
	} else {
		store = base
	}

	/* Compute for the rest */
	for i := 4; i <= n; i++ {
		store[i] = store[i-1] + store[i-2] + store[i-3]
	}

	return store[n]
}
