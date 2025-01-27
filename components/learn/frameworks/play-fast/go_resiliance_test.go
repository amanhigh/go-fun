package play_fast_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/chaos"
	"github.com/slok/goresilience/circuitbreaker"
	errors2 "github.com/slok/goresilience/errors"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
)

var _ = Describe("GoResiliance", func() {

	var (
		cmd goresilience.Runner
	)

	Context("Timeout", func() {
		var (
			TIMEOUT = 2 * time.Millisecond
		)
		BeforeEach(func() {
			cmd = timeout.New(timeout.Config{
				Timeout: TIMEOUT,
			})
		})

		It("should build", func() {
			Expect(cmd).To(Not(BeNil()))
		})

		It("should not timeout", func() {
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("should timeout", func() {
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				time.Sleep(TIMEOUT * 2)
				return nil
			})
			Expect(err).To(Not(BeNil()))
		})

	})

	Context("Retry", func() {
		var (
			RETRY = 2
		)
		BeforeEach(func() {
			cmd = retry.New(retry.Config{
				Times: RETRY,
			})
		})

		It("should work in first go", func() {
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("should run twice", func() {
			count := 0
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				if count < RETRY {
					count++
					return errors.New("First Call Failed")
				}
				return nil
			})
			Expect(err).To(BeNil())
			Expect(count).To(Equal(RETRY))
		})

		It("should throw error after repeated retries", func() {
			count := 0
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				count++
				return errors.New("Call Failed")
			})
			Expect(err).To(Not(BeNil()))
			Expect(count).To(Equal(RETRY + 1))
		})

	})

	Context("Circuit", func() {
		var (
			CIRCUIT_OPEN = 2
		)
		BeforeEach(func() {
			cmd =
				goresilience.RunnerChain(
					retry.NewMiddleware(retry.Config{
						Times: CIRCUIT_OPEN,
					}),
					circuitbreaker.NewMiddleware(circuitbreaker.Config{
						//ErrorPercentThresholdToOpen:        50,
						MinimumRequestToOpen:         CIRCUIT_OPEN,
						SuccessfulRequiredOnHalfOpen: CIRCUIT_OPEN / 2,
						//WaitDurationInOpenState:            5 * time.Second,
						//MetricsSlidingWindowBucketQuantity: 10,
						//MetricsBucketDuration:              1 * time.Second,
					}),
				)
		})

		It("should open after failures", func() {
			count := 0
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				if count < CIRCUIT_OPEN {
					count++
					return errors.New("Call Failed")
				}
				return nil
			})
			Expect(err).To(Equal(errors2.ErrCircuitOpen))
			Expect(count).To(Equal(CIRCUIT_OPEN))
		})

	})

	Context("Chaos", func() {
		var (
			LATENCY          = 2 * time.Millisecond
			ERROR_PERCENTILE = 10
		)
		BeforeEach(func() {
			injector := chaos.Injector{}
			injector.SetLatency(LATENCY)
			injector.SetErrorPercent(ERROR_PERCENTILE)
			cmd = chaos.New(chaos.Config{
				Injector: &injector,
			})

		})

		It("should fail due to chaos", func() {
			err := cmd.Run(context.TODO(), func(_ context.Context) error {
				return nil
			})
			Expect(err).To(Equal(errors2.ErrFailureInjected))
		})

	})
})
