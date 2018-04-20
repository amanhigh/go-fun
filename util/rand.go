package util

import (
	"math/rand"
	"time"
)

func RandomInts(n int,max int) (result []int) {
	seedRandom()
	for i := 0; i < n; i++ {
		result = append(result, rand.Intn(max))
	}
	return
}
func seedRandom() {
	rand.Seed(time.Now().UnixNano())
}
