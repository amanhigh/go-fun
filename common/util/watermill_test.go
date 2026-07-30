package util_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	common "github.com/amanhigh/go-fun/models/common"
)

var _ = Describe("Watermill Utilities", func() {
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
			Expect(shouldRetry(middleware.RetryParams{Err: fmt.Errorf("wrapped: %w", util.ErrInvalidMessageMetadata)})).To(BeFalse())
		})
	})

	Context("ContextFromMessage", func() {
		It("uses Watermill correlation and the consumed message ID for causation", func() {
			msg := message.NewMessage("msg-1", []byte("payload"))
			middleware.SetCorrelationID("corr-1", msg)
			msg.Metadata.Set(common.MetadataMessageIDKey, "msg-1")
			msg.Metadata.Set(common.MetadataCausationIDKey, "upstream-msg")

			result, err := util.ContextFromMessage(msg)

			Expect(err).NotTo(HaveOccurred())
			Expect(common.CorrelationFrom(result)).To(Equal("corr-1"))
			Expect(common.CausationFrom(result)).To(Equal("msg-1"))
		})

		It("rejects missing or mismatched identity and missing correlation", func() {
			cases := []struct {
				name string
				msg  *message.Message
			}{
				{name: "missing correlation", msg: message.NewMessage("msg-1", nil)},
				{name: "missing identity", msg: message.NewMessage("", nil)},
				{name: "mismatched identity", msg: message.NewMessage("msg-1", nil)},
			}
			middleware.SetCorrelationID("corr-1", cases[1].msg)
			middleware.SetCorrelationID("corr-1", cases[2].msg)
			cases[2].msg.Metadata.Set(common.MetadataMessageIDKey, "msg-2")

			for _, testCase := range cases {
				_, err := util.ContextFromMessage(testCase.msg)
				Expect(err).To(HaveOccurred(), testCase.name)
				Expect(errors.Is(err, util.ErrInvalidMessageMetadata)).To(BeTrue(), testCase.name)
			}
		})
	})

})
