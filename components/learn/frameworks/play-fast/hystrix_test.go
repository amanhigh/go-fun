package play_fast_test

import (
	"time"

	"github.com/failsafe-go/failsafe-go/retrypolicy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Hystrix", func() {
	Describe("Failsafe", func() {

		Describe("RetryPolicy", func() {
			var retryPolicy retrypolicy.RetryPolicy[string]
			BeforeEach(func() {
				retryPolicy = retrypolicy.Builder[string]().
					WithDelay(time.Millisecond * 10).
					WithMaxRetries(3).
					Build()
			})

			It("should build", func() {
				Expect(retryPolicy).NotTo(BeNil())
			})
		})
	})
})
