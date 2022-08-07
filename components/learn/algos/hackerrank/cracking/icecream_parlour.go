package cracking

import "sort"

type IceCream struct {
	index int //holds orignal index of icecream, so its not lost in sort.1
	price int
}

/* Encapsulate original index so we can sort without losing it */
func ToIcecreams(prices []int) (icecreams []IceCream) {
	for i, price := range prices {
		icecreams = append(icecreams, IceCream{i, price})
	}
	return
}

/**
Given the value of and the of each flavor for trips to the Ice Cream Parlor,
help Sunny and Johnny choose two distinct flavors such that they spend their entire pool of money during each visit

https://www.hackerrank.com/challenges/ctci-ice-cream-parlor/problem
*/
func FindIcecreams(icecreams []IceCream, money int) (values []int, indices []int) {
	/* Sort Icecreams on price */
	sort.Slice(icecreams, func(i, j int) bool {
		return icecreams[i].price < icecreams[j].price
	})

	l := len(icecreams)
	for i := 0; i < l; i++ {
		/* Consider each icecream */
		current := icecreams[i]

		/* Find amount for pair to be equal to money */
		balance := money - current.price

		/* Binary Search for balance */
		j := sort.Search(l, func(k int) bool {
			/* Since array is sorted consider items only greater than i
			It also prevents icecream i to be selected twice */
			return k > i && icecreams[k].price >= balance
		})

		/* If a icecream is found supply back values and original indices */
		if j < l && icecreams[j].price == balance {
			ith := icecreams[i]
			jth := icecreams[j]

			/* Since original index of i can be less than j ensure.
			   ith item always less than j in index*/
			if jth.index < ith.index {
				ith, jth = jth, ith
			}
			values = []int{ith.price, jth.price}
			indices = []int{ith.index + 1, jth.index + 1}
		}
	}
	return
}
