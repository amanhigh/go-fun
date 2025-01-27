package learn_test

import (
	"time"

	"github.com/amanhigh/go-fun/models/learn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SafeRead", func() {
	var (
		safe     learn.SafeReadWrite
		start    = 1
		len      = 2
		timeUnit = 10 * time.Millisecond
	)

	BeforeEach(func() {
		safe = learn.SafeReadWrite{start, make(chan int, len)}
	})

	It("should build", func() {
		Expect(safe).To(Not(BeNil()))
	})

	It("should return start value", func() {
		Expect(safe.Read()).To(Equal(start))
	})

	Context("On Write", func() {
		var (
			writeValue = 5
		)
		BeforeEach(func() {
			go func() {
				// Wait Sometime and send Channel Write
				time.Sleep(4 * timeUnit)
				safe.Write(writeValue)
			}()
		})

		It("should update safely", func() {
			Expect(safe.Read()).To(Equal(start))
			Eventually(safe.Read()).Should(Equal(writeValue))
			Eventually(safe.Intc).Should(Not(BeClosed()))
		})

		Context("post close", func() {
			BeforeEach(func() {
				safe.Write(10)
				safe.Close()
			})

			It("write should panic", func() {
				Expect(func() { safe.Write(10) }).Should(Panic())
			})

			It("should not update", func() {
				Eventually(safe.Intc).Should(BeClosed())
				// Should Return Old Value as Channel Close will reject Updates
				Expect(safe.Read()).To(Equal(start))
			})
		})
	})

})
