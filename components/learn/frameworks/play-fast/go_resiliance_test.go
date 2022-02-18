package play_fast_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/timeout"
	"time"
)

var _ = Describe("GoResiliance", func() {

	Context("Timeout", func() {
		// Create our command.
		var (
			cmd     goresilience.Runner
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
})
