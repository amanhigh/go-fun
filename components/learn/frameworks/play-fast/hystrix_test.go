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
	const (
		initialDelayMs = 10
		jitterMs       = 5
		maxRetries     = 3
		allowedRetries = 20
	)

	// Helper function to create an always failing function
	createAlwaysFailingFunction := func(startTimes []time.Time) func() (string, error) {
		attempts := 0
		return func() (string, error) {
			startTimes[attempts] = time.Now()
			attempts++
			return "", errors.New("persistent error")
		}
	}

	Describe("Failsafe", func() {
		Describe("RetryPolicy", func() {
			var retryBuilder retrypolicy.RetryPolicyBuilder[string]

			BeforeEach(func() {
				retryBuilder = retrypolicy.Builder[string]()
			})

			It("should build", func() {
				retryPolicy := retryBuilder.Build()
				Expect(retryPolicy).NotTo(BeNil())
			})

			It("should retry three times before succeeding", func() {
				attempts := 0
				failingFunction := func() (string, error) {
					attempts++
					if attempts <= maxRetries {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}

				retryPolicy := retryBuilder.
					WithDelay(time.Duration(initialDelayMs) * time.Millisecond).
					WithMaxRetries(maxRetries).
					Build()

				result, err := failsafe.Get(failingFunction, retryPolicy)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("success"))
				Expect(attempts).To(Equal(maxRetries + 1)) // 1 initial attempt + maxRetries
			})

			It("should fail after exhausting all retry attempts", func() {
				startTimes := make([]time.Time, maxRetries+1)
				alwaysFailingFunction := createAlwaysFailingFunction(startTimes)

				retryPolicy := retryBuilder.
					WithDelay(time.Duration(initialDelayMs) * time.Millisecond).
					WithMaxRetries(maxRetries).
					Build()

				result, err := failsafe.Get(alwaysFailingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("retries exceeded"))
				Expect(err.Error()).To(ContainSubstring("persistent error"))
				Expect(result).To(BeEmpty())
				Expect(len(startTimes)).To(Equal(maxRetries + 1)) // 1 initial attempt + maxRetries
			})

			It("should use exponential backoff", func() {
				retryPolicy := retryBuilder.
					WithBackoff(time.Millisecond, time.Millisecond*8).
					WithMaxRetries(maxRetries).
					Build()

				expectedDelays := []time.Duration{
					time.Millisecond,
					time.Millisecond * 2,
					time.Millisecond * 4,
				}

				startTimes := make([]time.Time, maxRetries+1)
				alwaysFailingFunction := createAlwaysFailingFunction(startTimes)

				_, err := failsafe.Get(alwaysFailingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(len(startTimes)).To(Equal(maxRetries + 1)) // 1 initial attempt + maxRetries

				for i := 1; i <= maxRetries; i++ {
					delay := startTimes[i].Sub(startTimes[i-1])
					expected := expectedDelays[i-1]
					Expect(delay).To(BeNumerically("~", expected, time.Millisecond/2))
				}
			})

			It("should apply jitter to retry delays", FlakeAttempts(3), func() {
				initialDelay := time.Millisecond * 10
				jitter := time.Millisecond * 5
				allowedRetries := 10

				retryPolicy := retryBuilder.
					WithDelay(time.Duration(initialDelayMs) * time.Millisecond).
					WithJitter(time.Duration(jitterMs) * time.Millisecond).
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

			It("should trigger OnRetry event listener", func() {
				attempts := 0
				retryCount := 0

				failingFunction := func() (string, error) {
					attempts++
					if attempts <= maxRetries {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}

				retryPolicy := retryBuilder.
					WithMaxRetries(maxRetries).
					OnRetry(func(e failsafe.ExecutionEvent[string]) {
						retryCount++
						Expect(e.Retries()).To(Equal(retryCount))
						Expect(e.LastResult()).To(BeEmpty())
						Expect(e.LastError()).To(MatchError("temporary error"))
					}).
					Build()

				result, err := failsafe.Get(failingFunction, retryPolicy)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("success"))
				Expect(attempts).To(Equal(maxRetries + 1)) // 1 initial attempt + maxRetries
				Expect(retryCount).To(Equal(maxRetries))
			})

			It("should handle specific errors", func() {
				attempts := 0

				// Define custom error types
				type RetryableError struct{ error }
				type NonRetryableError struct{ error }

				failingFunction := func() (string, error) {
					attempts++
					switch attempts {
					case 1, 2:
						return "", RetryableError{errors.New("retryable error")}
					case 3:
						return "", NonRetryableError{errors.New("non-retryable error")}
					default:
						return "success", nil
					}
				}

				retryPolicy := retryBuilder.
					WithMaxRetries(maxRetries).
					HandleErrorTypes(RetryableError{}).
					Build()

				result, err := failsafe.Get(failingFunction, retryPolicy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(NonRetryableError{}))
				Expect(err.Error()).To(ContainSubstring("non-retryable error"))
				Expect(result).To(BeEmpty())
				Expect(attempts).To(Equal(3)) // 1 initial attempt + 2 retries
			})
		})

	})
})
