package practice

import (
	"fmt"
	"strings"
	"unicode"
)

/**
	We consider two strings to be anagrams of each other if the first string's letters can be rearranged to form the second string.
    In other words, both strings must contain the same exact letters in the same exact frequency For example, bacdc and dcbac are anagrams,
	but bacdc and dcbad are not.

	https://www.hackerrank.com/challenges/ctci-making-anagrams/problem
*/
func AnagramGroups(words []string) map[string][]string {
	anagramMap := map[string][]string{}
	for _, word := range words {
		fingerPrint := strings.Replace(fmt.Sprint(fingerPrint(word)), " ", "", -1)
		anagramMap[fingerPrint] = append(anagramMap[fingerPrint], word)
	}
	return anagramMap
}

/**
Finger print will work if it contains only lowercase a-z
*/
func fingerPrint(word string) (fingerPrint []int) {
	fingerPrint = make([]int, 26)
	for _, c := range word {
		cc := unicode.ToLower(c)
		fingerPrint[cc-'a']++
	}
	return
}
