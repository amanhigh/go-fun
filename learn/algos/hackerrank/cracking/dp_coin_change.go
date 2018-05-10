package cracking

import "fmt"

var dpMem = map[string]int{}

/**
Find number of ways to split money, with given denominations.
Keep track of Selected Coins and print when valid solution is found.

Returns number of ways
Eg. 	fmt.Println(Split(4, []int{1,2,3}, []int{})) is 4
*/
func Split(money int, denominations []int, selectedCoins []int) (p int) {
	/* All amount was used hence one possible solution */
	if money == 0 {
		p = 1
		//fmt.Println("Selection: ", selectedCoins)
		/* If Money Remains and also valid coins to select recurse */
	} else if money > 0 && len(denominations) > 0 {
		/* Choose First Coin */
		coin := denominations[0]
		/* Form key for mem store */
		key := fmt.Sprintf("%d-%d", money, coin)
		var ok bool
		/* Read Value from Mem store & compute if not found*/
		if p, ok = dpMem[key]; !ok {
			/* Compute Balance after choosing that */
			balance := money - coin
			/*
				Recurse into two branches:
				Left (Reduce Denominations, Don't reuse Coin): Remove selected coin (at 0th index) and find ways using same amount of money without using this coin.
				Eg. 3 = 2+1
				Right (Reduce Amount,Reuse Coin) : Use selected coin (at 0th index) and for remaining balance find ways to split using all coins.
				Eg. 3= 1 + 1 + 1
			*/
			p += Split(money, denominations[1:], selectedCoins) + Split(balance, denominations, append(selectedCoins, coin))
			dpMem[key] = p
		}
	}
	return
}
