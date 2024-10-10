package play_fast_test

import (
	"errors"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/hedgepolicy"
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
	failFuncWithTimes := func(startTimes []time.Time) func() (string, error) {
		attempts := 0
		return func() (string, error) {
			startTimes[attempts] = time.Now()
			attempts++
			return "", errors.New("persistent error")
		}
	}

	successfulFunction := func() (string, error) {
		return "success", nil
	}

	failingFunction := func() (string, error) {
		return "", errors.New("persistent error")
	}

	Describe("Failsafe", func() {
		Context("RetryPolicy", func() {
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
				alwaysFailingFunction := failFuncWithTimes(startTimes)

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
				alwaysFailingFunction := failFuncWithTimes(startTimes)

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

		Context("CircuitBreaker", func() {
			var breakerBuilder circuitbreaker.CircuitBreakerBuilder[string]

			BeforeEach(func() {
				breakerBuilder = circuitbreaker.Builder[string]()
			})

			It("should build", func() {
				breaker := breakerBuilder.Build()
				Expect(breaker).NotTo(BeNil())
			})

			It("should keep the circuit breaker closed for successful executions", func() {
				breaker := breakerBuilder.Build()

				for i := 0; i < 10; i++ {
					result, err := failsafe.Get(successfulFunction, breaker)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal("success"))
				}

				Expect(breaker.State()).To(Equal(circuitbreaker.ClosedState))
			})

			It("should open the circuit before reaching the failure threshold", func() {
				breaker := breakerBuilder.
					WithFailureThreshold(maxRetries).
					WithDelay(10 * time.Millisecond). // Add a short delay before closing the circuit
					Build()

				// Execute the failing function up to one before failure threshold
				for i := 0; i < maxRetries-1; i++ {
					_, err := failsafe.Get(failingFunction, breaker)
					Expect(err).To(HaveOccurred())
					Expect(breaker.State()).To(Equal(circuitbreaker.ClosedState))
				}

				// The next execution should open the circuit
				_, err := failsafe.Get(failingFunction, breaker)
				By("At threshold")
				Expect(err).To(HaveOccurred())
				Expect(breaker.State()).To(Equal(circuitbreaker.OpenState))

				// Subsequent executions should immediately return ErrOpen
				_, err = failsafe.Get(failingFunction, breaker)
				By("After threshold")
				Expect(err).To(MatchError(circuitbreaker.ErrOpen))
				Expect(breaker.State()).To(Equal(circuitbreaker.OpenState))
			})

			It("should close the circuit after reaching the success threshold", func() {
				successThreshold := uint(3)
				delay := 50 * time.Millisecond
				breaker := breakerBuilder.
					WithFailureThreshold(1).
					WithSuccessThreshold(successThreshold).
					WithDelay(delay).
					Build()

				// Open the circuit
				_, err := failsafe.Get(failingFunction, breaker)
				Expect(err).To(HaveOccurred())
				Expect(breaker.State()).To(Equal(circuitbreaker.OpenState))

				// Wait for the remaining delay
				for breaker.RemainingDelay() > 0 {
					time.Sleep(10 * time.Millisecond)
				}

				// Execute successful functions up to success threshold
				for i := uint(0); i < successThreshold; i++ {
					result, err := failsafe.Get(successfulFunction, breaker)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal("success"))

					if i == 0 {
						Expect(breaker.State()).To(Equal(circuitbreaker.HalfOpenState), "Should transition to half-open on first success")
					} else if i < successThreshold-1 {
						Expect(breaker.State()).To(Equal(circuitbreaker.HalfOpenState), "Should remain half-open until threshold is met")
					} else {
						Expect(breaker.State()).To(Equal(circuitbreaker.ClosedState), "Should close after meeting success threshold")
					}
				}

				// Verify that subsequent calls succeed and keep the circuit closed
				for i := 0; i < 3; i++ {
					result, err := failsafe.Get(successfulFunction, breaker)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal("success"))
					Expect(breaker.State()).To(Equal(circuitbreaker.ClosedState), "Should remain closed for subsequent successful calls")
				}
			})
		})

		Context("HedgePolicy", func() {
			var (
				hedgeDelay  = 50 * time.Millisecond
				hedgePolicy hedgepolicy.HedgePolicy[string]
			)

			BeforeEach(func() {
				hedgePolicy = hedgepolicy.BuilderWithDelay[string](hedgeDelay).
					WithMaxHedges(2).
					Build()
			})

			It("should not trigger hedge if execution completes before delay", func() {
				attempts := 0
				failingFunction := func() (string, error) {
					attempts++
					// Simulate quick execution that completes before hedge delay
					time.Sleep(time.Millisecond * 10) // Less than hedgeDelay
					return "success", nil
				}

				result, err := failsafe.Get(failingFunction, hedgePolicy)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("success"))
				Expect(attempts).To(Equal(1)) // Only one attempt should have been made
			})

			It("should retry if latency exceeds the threshold", func() {
				attempts := 0
				executionTimes := make([]time.Duration, 0)

				// Define a function that simulates high latency
				highLatencyFunction := func() (string, error) {
					start := time.Now()
					attempts++
					if attempts == 1 {
						time.Sleep(time.Millisecond * 100) // Simulating high latency for first attempt
					} else {
						time.Sleep(time.Millisecond * 10) // Subsequent attempts are faster
					}
					executionTimes = append(executionTimes, time.Since(start))
					return "success", nil
				}

				// Run with hedge policy
				result, err := failsafe.Get(highLatencyFunction, hedgePolicy)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("success"))
				Expect(attempts).To(BeNumerically(">", 1))  // Should be more than one attempt
				Expect(attempts).To(BeNumerically("<=", 3)) // But not more than 3 (1 initial + 2 hedges)

				// Check if the fastest execution time is returned
				fastestTime := executionTimes[0]
				for _, t := range executionTimes[1:] {
					if t < fastestTime {
						fastestTime = t
					}
				}
				Expect(fastestTime).To(BeNumerically("<", time.Millisecond*100)) // The fastest should be less than the initial slow attempt
			})
		})

	})
})
