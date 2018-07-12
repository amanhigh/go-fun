package challenge

/**
https://www.hackerrank.com/challenges/non-divisible-subset/problem
https://www.geeksforgeeks.org/subset-no-pair-sum-divisible-k/
C(n,r) - http://www.mathwords.com/c/combination_formula.htm

Based on
If sum of two numbers is divisible by K, then if one of them gives remainder i, other will give remainder (K – i).
Eg. 3 + 7 = 10, k = 5 then 3%5=3, 7%5=2 thus 2(i)+3(k-i)==5(k)
*/
func NonDivisibleSubset(input []int, k int) (result int) {
	/* Build an array to count remainders */
	modulos := make([]int, k)

	/*
	* For each value increment count of its remainder.
	* Now we have count of input (numbers) giving remainder 0 to k-1.
	* No one will give remainder k as n%k < k always.
	 */
	for _, value := range input {
		modulos[value%k]++
	}

	if 1 < modulos[0] {
		result = 1
	} else {
		result = modulos[0]
	}
	//fmt.Println(result)
	//fmt.Println(k, modulos)

	/*
		Handling even values of k.
		For even values of K, the equal remainder is simular to the 0 case. For K = 6, pairs are 1+5, 2+4, 3+3.
		For values with remainder 3, at most one value can be added to the result set.
	*/
	if k%2 == 0 {
		if 1 < modulos[k/2] {
			modulos[k/2] = 1
		}
	}

	/* Go from 0 to k/2 as k/2 to k-1 are covered by computing k-i */
	for i := 1; i <= k/2; i++ {
		if modulos[i] < modulos[k-i] {
			result += modulos[k-i]
		} else {
			result += modulos[i]
		}
		//fmt.Println(modulos[i], modulos[k-i], result)
	}

	return
}
