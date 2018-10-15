package cracking_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/amanhigh/go-fun/learn/concepts/algospts/algos/hackerrank/cracking"
)

var _ = Describe("BraceMatch", func() {
	It("should be success", func() {
		Expect(MatchBrace("[({()})]")).To(BeTrue())
	})

	It("should fail", func() {
		Expect(MatchBrace("[({}}]")).To(BeFalse())
	})
})
