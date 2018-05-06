package cracking

import "sort"

type IceCream struct {
	index int
	price int
}

/* Encapsulate original index so we can sort without losing it */
func ToIcecreams(prices []int) (icecreams []IceCream) {
	for i, price := range prices {
		icecreams = append(icecreams, IceCream{i, price})
	}
	return
}

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
		key := money - current.price

		/* Binary Search for key */
		j := sort.Search(l, func(k int) bool {
			/* Since array is sorted consider items only greater than i
			It also prevents icecream i to be selected twice */
			return k > i && icecreams[k].price >= key
		})

		/* If a key is found supply back values and original indices */
		if j < l && icecreams[j].price == key {
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
