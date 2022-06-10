package cracking_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
)

var _ = Describe("Primality", func() {
	It("should be true for prime numbers", func() {
		Expect(cracking.IsPrime(7)).To(BeTrue())
		Expect(cracking.IsPrime(13)).To(BeTrue())

		Expect(cracking.IsPrimeSmart(7)).To(BeTrue())
		Expect(cracking.IsPrimeSmart(13)).To(BeTrue())
	})

	It("should be false for non prime numbers", func() {
		Expect(cracking.IsPrime(8)).To(BeFalse())
		Expect(cracking.IsPrime(15)).To(BeFalse())

		Expect(cracking.IsPrimeSmart(8)).To(BeFalse())
		Expect(cracking.IsPrimeSmart(15)).To(BeFalse())
	})

})
