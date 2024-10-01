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

			It("should use exponential backoff for retries", func() {
				initialDelay := time.Millisecond
				maxDelay := time.Millisecond * 8
				allowedExponentialBackoff := 3

				retryPolicy := retrypolicy.Builder[string]().
					WithBackoff(initialDelay, maxDelay).
					WithMaxRetries(allowedExponentialBackoff).
					Build()

				expectedDelays := []time.Duration{
					time.Millisecond,
					time.Millisecond * 2,
					time.Millisecond * 4,
				}

				attempts := 0
				startTimes := make([]time.Time, allowedExponentialBackoff+1)
				alwaysFailingFunction := func() (string, error) {
					startTimes[attempts] = time.Now()
					attempts++
					return "", errors.New("persistent error")
				}

				_, err := failsafe.Get(alwaysFailingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(attempts).To(Equal(allowedExponentialBackoff + 1)) // 1 initial attempt + 3 retries

				for i := 1; i <= allowedExponentialBackoff; i++ {
					delay := startTimes[i].Sub(startTimes[i-1])
					expected := expectedDelays[i-1]
					Expect(delay).To(BeNumerically("~", expected, time.Millisecond/2))
				}
			})

			It("should apply jitter to retry delays", FlakeAttempts(3), func() {
				initialDelay := time.Millisecond * 10
				jitter := time.Millisecond * 5
				allowedRetries := 10

				retryPolicy := retrypolicy.Builder[string]().
					WithDelay(initialDelay).
					WithJitter(jitter).
					WithMaxRetries(allowedRetries).
					Build()

				attempts := 0
				startTimes := make([]time.Time, allowedRetries+1)
				alwaysFailingFunction := func() (string, error) {
					startTimes[attempts] = time.Now()
					attempts++
					return "", errors.New("persistent error")
				}

				_, err := failsafe.Get(alwaysFailingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(attempts).To(Equal(allowedRetries + 1)) // 1 initial attempt + allowedRetries

				delaysWithinRange := 0
				for i := 1; i <= allowedRetries; i++ {
					delay := startTimes[i].Sub(startTimes[i-1])
					if delay >= initialDelay-jitter && delay <= initialDelay+jitter {
						delaysWithinRange++
					}
				}

				// Expect at least 80% of delays to be within the jitter range
				Expect(float64(delaysWithinRange) / float64(allowedRetries)).To(BeNumerically(">=", 0.8))
			})
		})
	})
})
