package cracking_test

import (
	"github.com/amanhigh/go-fun/components/learn/algos/hackerrank/cracking"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BraceMatch", func() {
	It("should be success", func() {
		Expect(cracking.MatchBrace("[({()})]")).To(BeTrue())
	})

	It("should fail", func() {
		Expect(cracking.MatchBrace("[({}}]")).To(BeFalse())
	})
})
