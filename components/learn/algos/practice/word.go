package practice

import (
	"fmt"
	"strings"
	"unicode"
)

/*
*

		We consider two strings to be anagrams of each other if the first string's letters can be rearranged to form the second string.
	    In other words, both strings must contain the same exact letters in the same exact frequency For example, bacdc and dcbac are anagrams,
		but bacdc and dcbad are not.

		https://www.hackerrank.com/challenges/ctci-making-anagrams/problem
*/
func AnagramGroups(words []string) map[string][]string {
	anagramMap := map[string][]string{}
	for _, word := range words {
		fingerPrint := strings.ReplaceAll(fmt.Sprint(fingerPrint(word)), " ", "")
		anagramMap[fingerPrint] = append(anagramMap[fingerPrint], word)
	}
	return anagramMap
}

func CommonPrefix(words []string) (prefix string) {
	if len(words) == 0 {
		return ""
	}

	// Select First word as Prefix
	prefix = words[0]

	// Iterate through the rest of the words
	for _, word := range words[1:] {
		//Continue till Prefix doesn't disappear or word list ends
		for len(prefix) > 0 {
			//Slice Word is larger than prefix and Try a Match
			if len(word) >= len(prefix) && word[:len(prefix)] == prefix {
				//Move to Next Word on Match
				break
			}

			// Reduce the length of the prefix by 1 and try Rematch with Word.
			prefix = prefix[:len(prefix)-1]
		}
	}

	return
}

/*
*
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
