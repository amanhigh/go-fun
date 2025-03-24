package util

import (
	"crypto/rand"
	"math/big"
)

func RandomInts(n, upperBound int) (result []int) {
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(upperBound)))
		result = append(result, int(num.Int64()))
	}
	return
}

func RandomInt(lowerBound, upperBound int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(upperBound-lowerBound)))
	return int(num.Int64()) + lowerBound
}
