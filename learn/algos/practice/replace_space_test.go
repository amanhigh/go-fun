package practice_test

import (
	practice2 "github.com/amanhigh/go-fun/learn/algos/practice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReplaceSpace", func() {
	It("should encode to %20", func() {
		result := practice2.ReplaceSpace("Aman Preet Singh")
		Expect(result).To(Not(BeNil()))
		Expect(result).To(Equal("Aman%20Preet%20Singh"))
	})
})
