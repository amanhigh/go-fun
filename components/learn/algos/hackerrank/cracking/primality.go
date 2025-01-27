package cracking

func IsPrime(n int) bool {
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsPrimeSmart(n int) bool {
	/* Base Cases */
	if n <= 1 { //-ve,0,1 can't be prime
		return false
	}

	if n < 3 { // 2,3 are both prime
		return true
	}

	// Handles cases from 4 to 25 (except primes like 5,7,11)
	// all cases hit this condition
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	/*
		All primes are of the form 6k+-1
		https://www.youtube.com/watch?v=AaNUzEHiDpI
	*/
	for i := 5; i*i <= n; i += 6 {
		sixKMinus1 := i
		sixKPlus1 := sixKMinus1 + 2
		// fmt.Println("N=", n, "K=", i/6+1, sixKMinus1, sixKPlus1)

		/* If its perfectly divisible it means its not a prime */
		if n%sixKMinus1 == 0 || n%sixKPlus1 == 0 {
			return false
		}
	}

	return true
}
