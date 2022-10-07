package cracking_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
)

var _ = Describe("MedianFinder", func() {
	var (
		finder cracking.MedianFinder
		set1   = []int{60, 35, 58, 32}
		set2   = []int{42, 40, 50}
	)
	BeforeEach(func() {
		finder = cracking.NewMedianFinder()
		for _, i := range set1 {
			finder.Add(i)
		}
	})

	It("should build", func() {
		Expect(finder).To(Not(BeNil()))
	})

	It("should compute", func() {
		Expect(finder.GetMedian()).To(BeEquivalentTo(46.5))

		By("Adding Second Set")
		for _, i := range set2 {
			finder.Add(i)
		}
		Expect(finder.GetMedian()).To(BeEquivalentTo(42))
	})

})
