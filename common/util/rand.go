package util

import (
	"math/rand"
	"time"
)

func RandomInts(n int, max int) (result []int) {
	for i := 0; i < n; i++ {
		result = append(result, rand.Intn(max))
	}
	return
}

func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func SeedRandom() {
	rand.Seed(time.Now().UnixNano())
}
