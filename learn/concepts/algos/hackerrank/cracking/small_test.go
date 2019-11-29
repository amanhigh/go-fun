package cracking_test

import (
	"github.com/amanhigh/go-fun/learn/concepts/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Small", func() {
	It("should generate fibonacci", func() {
		Expect(cracking.Fibonacci(34)).To(Equal(5702887))
	})
})
