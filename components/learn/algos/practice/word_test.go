package practice_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/practice"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Word Test", func() {
	Context("Anagrams", func() {

		It("should identify anagrams", func() {
			words := []string{"dcbac", "bacdc"}
			for _, anagrams := range practice.AnagramGroups(words) {
				Expect(anagrams).To(HaveLen(2))
				Expect(anagrams).To(Equal(words))
			}
		})

		It("should ignore non anagrams", func() {
			words := []string{"bacdc", "dcbad"}
			for _, anagrams := range practice.AnagramGroups(words) {
				Expect(anagrams).To(HaveLen(1))
			}
		})
	})

	It("should find common prefix", func() {
		Expect(practice.CommonPrefix([]string{"flower", "flow", "flight"})).To(Equal("fl"))
	})
})
