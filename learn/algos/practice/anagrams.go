package practice

import (
	"fmt"
	"strings"
	"unicode"
)

func AnagramGroups(words []string) map[string][]string {
	anagramMap := map[string][]string{}
	for _, word := range words {
		fingerPrint := strings.Replace(fmt.Sprint(fingerPrint(word)), " ", "", -1)
		anagramMap[fingerPrint] = append(anagramMap[fingerPrint], word)
	}
	return anagramMap
}

func fingerPrint(word string) (fingerPrint []int) {
	fingerPrint = make([]int, 26)
	for _, c := range word {
		cc := unicode.ToLower(c)
		fingerPrint[cc-'a']++
	}
	return
}
