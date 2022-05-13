package cracking_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Small", func() {
	It("should generate fibonacci", func() {
		Expect(cracking.Fibonacci(34)).To(Equal(5702887))
	})
})
