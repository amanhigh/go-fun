package play_fast_test

import (
	"errors"
	"time"

	"github.com/failsafe-go/failsafe-go"
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

			It("should retry three times before succeeding", func() {
				attempts := 0
				failingFunction := func() (string, error) {
					attempts++
					if attempts <= 3 {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}

				result, err := failsafe.Get(failingFunction, retryPolicy)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("success"))
				Expect(attempts).To(Equal(4)) // 1 initial attempt + 3 retries
			})

			It("should fail after exhausting all retry attempts", func() {
				attempts := 0
				alwaysFailingFunction := func() (string, error) {
					attempts++
					return "", errors.New("persistent error")
				}

				result, err := failsafe.Get(alwaysFailingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("retries exceeded"))
				Expect(err.Error()).To(ContainSubstring("persistent error"))
				Expect(result).To(BeEmpty())
				Expect(attempts).To(Equal(4)) // 1 initial attempt + 3 retries
			})
		})
	})
})
