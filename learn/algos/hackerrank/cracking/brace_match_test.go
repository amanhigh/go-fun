package cracking_test

import (
	cracking2 "github.com/amanhigh/go-fun/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BraceMatch", func() {
	It("should be success", func() {
		Expect(cracking2.MatchBrace("[({()})]")).To(BeTrue())
	})

	It("should fail", func() {
		Expect(cracking2.MatchBrace("[({}}]")).To(BeFalse())
	})
})
