package practice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/algos/practice"
)

var _ = Describe("ReplaceSpace", func() {
	It("should encode to %20", func() {
		result := ReplaceSpace("Aman Preet Singh")
		Expect(result).To(Not(BeNil()))
		Expect(result).To(Equal("Aman%20Preet%20Singh"))
	})
})
