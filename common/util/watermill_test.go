package util_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	common "github.com/amanhigh/go-fun/models/common"
)

var _ = Describe("Watermill Utilities", func() {
	Context("SagaMetadataMiddleware", func() {
		It("recovers missing UUID and correlation metadata", func() {
			msg := message.NewMessage("", []byte("payload"))
			var metadata common.Metadata
			handler := util.SagaMetadataMiddleware()(func(msg *message.Message) ([]*message.Message, error) {
				metadata = common.MetadataFromContext(msg.Context())
				return nil, nil
			})

			_, err := handler(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(metadata.MessageID).NotTo(BeEmpty())
			Expect(metadata.CorrelationID).NotTo(BeEmpty())
			Expect(msg.UUID).To(Equal(metadata.MessageID))
			Expect(middleware.MessageCorrelationID(msg)).To(Equal(metadata.CorrelationID))
		})

		It("keeps recovered metadata stable across repeated invocation", func() {
			msg := message.NewMessage("", []byte("payload"))
			var seen []common.Metadata
			handler := util.SagaMetadataMiddleware()(func(msg *message.Message) ([]*message.Message, error) {
				seen = append(seen, common.MetadataFromContext(msg.Context()))
				return nil, nil
			})

			_, firstErr := handler(msg)
			_, secondErr := handler(msg)
			Expect(firstErr).NotTo(HaveOccurred())
			Expect(secondErr).NotTo(HaveOccurred())
			Expect(seen).To(HaveLen(2))
			Expect(seen[1]).To(Equal(seen[0]))
		})
	})

	// =========================================================================
	// 1. DeadLetterTopic — pure string utility
	// =========================================================================
	Context("DeadLetterTopic", func() {
		It("should append .dead-letter suffix to the source topic", func() {
			Expect(util.DeadLetterTopic("orders")).To(Equal("orders.dead-letter"))
		})

		It("should handle an empty source topic", func() {
			Expect(util.DeadLetterTopic("")).To(Equal(".dead-letter"))
		})

		It("should handle topics that already contain dots", func() {
			Expect(util.DeadLetterTopic("a.b.c")).To(Equal("a.b.c.dead-letter"))
		})
	})

	// =========================================================================
	// 2. DefaultRetryConfig — centralized defaults
	// =========================================================================
	Context("DefaultRetryConfig", func() {
		It("should return two retries with two-second initial interval", func() {
			cfg := util.DefaultRetryConfig()
			Expect(cfg.MaxRetries).To(Equal(2))
			Expect(cfg.InitialInterval).To(Equal(2 * time.Second))
		})
	})

	// =========================================================================
	// 3. retry classification
	// =========================================================================
	Context("retry classification", func() {
		var shouldRetry func(middleware.RetryParams) bool

		BeforeEach(func() {
			shouldRetry = util.DefaultRetryConfig().ShouldRetry
		})

		It("should classify retryable and permanent errors", func() {
			httpCases := []struct {
				status int
				retry  bool
			}{
				{http.StatusBadRequest, false},
				{http.StatusRequestTimeout, true},
				{http.StatusTooManyRequests, true},
				{http.StatusInternalServerError, true},
			}
			for _, testCase := range httpCases {
				err := common.NewHttpError("http error", testCase.status)
				Expect(shouldRetry(middleware.RetryParams{Err: err})).To(Equal(testCase.retry))
			}

			payloadErrors := []error{
				&json.SyntaxError{Offset: 10},
				&json.UnmarshalTypeError{Value: "string", Type: nil},
			}
			for _, err := range payloadErrors {
				Expect(shouldRetry(middleware.RetryParams{Err: err})).To(BeFalse())
			}

			Expect(shouldRetry(middleware.RetryParams{Err: errors.New("transient failure")})).To(BeTrue())
		})
	})

})
