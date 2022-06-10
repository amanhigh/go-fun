package cracking_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
)

var _ = Describe("MakingAnagram", func() {

	It("should fingerprint", func() {
		Expect(cracking.FingerPrint("santa")).To(Equal(cracking.FingerPrint("satan")))
		Expect(cracking.FingerPrint("restful")).To(Equal(cracking.FingerPrint("fluster")))

		Expect(cracking.FingerPrint("fullfil")).To(Not(Equal(cracking.FingerPrint("fluster"))))
		Expect(cracking.FingerPrint("dome")).To(Not(Equal(cracking.FingerPrint("rome"))))
	})

	It("should compute diff", func() {
		Expect(cracking.AnagramDiff([]int{1, 2, 1, 3}, []int{1, 2, 1, 3})).To(Equal(0))
		Expect(cracking.AnagramDiff([]int{1, 2, 1, 3}, []int{1, 2, 0, 3})).To(Equal(1))
		Expect(cracking.AnagramDiff([]int{1, 2, 1, 3}, []int{1, 2, 0, 0})).To(Equal(4))
		Expect(cracking.AnagramDiff([]int{1, 2, 1, 3}, []int{1, 2, 0, 1})).To(Equal(3))

	})

})
