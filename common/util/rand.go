package util

import (
	"crypto/rand"
	"math/big"
)

func RandomInts(n int, max int) (result []int) {
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		result = append(result, int(num.Int64()))
	}
	return
}

func RandomInt(min int, max int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(num.Int64()) + min
}
