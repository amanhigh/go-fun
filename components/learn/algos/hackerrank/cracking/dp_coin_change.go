package cracking

import "fmt"

var dpMem = map[string]int{}

/*
*
Find number of ways to split money, with given denominations.
Keep track of Selected Coins and print when valid solution is found.

Returns number of ways
Eg. 	fmt.Println(Split(4, []int{1,2,3}, []int{})) is 4
*/
func Split(money int, denominations, selectedCoins []int) (p int) {
	/* All amount was used hence one possible solution */
	if money == 0 {
		p = 1
		// fmt.Println("Selection: ", selectedCoins)
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

var coinTable []int

/*
*
Dynamic Programming Based
Find number of ways to split money, with given denominations.

Returns number of ways
Eg. 	fmt.Println(Split(4, []int{1,2,3}, []int{})) is 4
*/
func SplitDp(money int, denominations []int) (p int) {
	/* Make Dp Table for possible money*/
	m := len(denominations)
	coinTable = make([]int, money+1)

	/* Base Case, if all money is used up there is one solution */
	coinTable[0] = 1

	/* Compute for all Coins */
	for i := 0; i < m; i++ {
		/* Take Selected Coin */
		coin := denominations[i]
		/*
			Start Updating possibilties from coin money onwards
			because below that using that coin will result in negative balance.
			Eg. If Coin is 3 start from index 3 as for money 1,2 coin 3 can't be used.
		*/
		for j := coin; j <= money; j++ {
			/*
				Hit Base case when coin is equal to money\
				Eg. j=coin -coin =0 hence for money 3 and coin 3 increment its value by 1 using coin 3.
				Before that coin 2 & 1 would have its value 1+1=2 hence total possiblties 3.

				Eg2. For money = 4
				Coin 1 - Increases it by 1 as only 1 of it can be used in conjuction with 3.
				Coin 2 - Increases it by 2 (total 3) = 1 from [Coin 1] i.e 1,1,1,1 + 2 from [Coin 2] i.e 2,2 & 1,1,2.
				Coin 3 - Increase it by 1 (total 4) = 1,3
			*/
			coinTable[j] += coinTable[j-coin]
			/* In summary this pass each coin passes over all monetary values */
			// fmt.Println("Progrses:", coin, coinTable, j-coin, coinTable[j-coin])
		}
	}
	return coinTable[money]
}
