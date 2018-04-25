package cracking

import "math"

func FingerPrint(input string) (result []int) {
	result = make([]int, 26)
	for _, c := range input {
		result[c-'a']++
	}
	return
}

func AnagramDiff(f1 []int, f2 []int) (diff int) {
	for i, p1 := range f1 {
		p2:= f2[i]
		diff+=int(math.Abs(float64(p2-p1)))
	}
	return
}
