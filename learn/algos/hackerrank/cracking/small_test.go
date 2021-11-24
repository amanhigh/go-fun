package cracking_test

import (
	cracking2 "github.com/amanhigh/go-fun/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Small", func() {
	It("should generate fibonacci", func() {
		Expect(cracking2.Fibonacci(34)).To(Equal(5702887))
	})
})
