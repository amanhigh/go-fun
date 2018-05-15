package cracking

var stairMem = map[int]int{}

/**
Eg. For n 1,3,7 answer is 1,4,44
https://www.hackerrank.com/challenges/ctci-recursive-staircase/problem
*/
func Possibilities(n int, steps, taken []int) (result int) {
	/* If Steps are remaining try to take using available steps */
	if n > 0 {
		var ok bool
		/* Try to Search Memory */
		if result, ok = stairMem[n]; !ok {
			/* Try each step, n may go negative if step is greater but it will return */
			for _, step := range steps {
				result += Possibilities(n-step, steps, append(taken, step))
			}
			stairMem[n] = result
		}
		//result += Possibilities(n, steps[1:], taken)
	} else if n == 0 {
		/* N becomes exactly zero is a valid combination count as result */
		result = 1
		/* Print Steps just to visualize */
		//fmt.Println(taken)
	}
	/* If n is negative it was invalid solution don't count it ,hence result zero */
	return
}
