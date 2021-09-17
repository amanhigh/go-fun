package learn_test

import (
	"github.com/amanhigh/go-fun/apps/models/learn"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SafeRead", func() {
	var (
		safe  learn.SafeReadWrite
		start = 1
		len   = 2
	)
	BeforeEach(func() {
		safe = learn.SafeReadWrite{start, make(chan int, len)}
	})

	It("should build", func() {
		Expect(safe).To(Not(BeNil()))
	})
})
