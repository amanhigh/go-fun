package practice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/practice"
)

var _ = Describe("Anagrams", func() {
	It("should identify anagrams", func() {
		words := []string{"dcbac", "bacdc"}
		for _, anagrams := range AnagramGroups(words) {
			Expect(len(anagrams)).To(Equal(2))
			Expect(anagrams).To(Equal(words))
		}
	})

	It("should ignore non anagrams", func() {
		words := []string{"bacdc", "dcbad"}
		for _, anagrams := range AnagramGroups(words) {
			Expect(len(anagrams)).To(Equal(1))
		}
	})
})
