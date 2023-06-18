package practice_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/practice"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReplaceSpace", func() {
	It("should encode to %20", func() {
		result := practice.ReplaceSpace("Aman Preet Singh")
		Expect(result).To(Not(BeNil()))
		Expect(result).To(Equal("Aman%20Preet%20Singh"))
	})
})
