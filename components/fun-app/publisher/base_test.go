package publisher_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
)

// mockPublisher follows testify style used across the repo.
type mockPublisher struct{ mock.Mock }

func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	args := m.Called(topic, messages)
	return args.Error(0)
}
func (m *mockPublisher) Close() error { return nil }

func TestBasePublisher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BasePublisher Suite")
}

var _ = Describe("BasePublisher.PublishWithExtras", func() {
	var (
		pub     *mockPublisher
		base    publisher.BasePublisher
		ctx     context.Context
		result  common.HttpError
		extras  map[string]string
		topic   string
		payload any
	)

	BeforeEach(func() {
		pub = &mockPublisher{}
		base = publisher.NewBasePublisher(pub)
		topic = "enrollments.enroll_cmd"
		payload = struct{ ID string }{ID: "p-1"}
		extras = nil
		ctx = common.WithCorrelation(context.Background(), "corr-123")
	})

	Context("without causation and extras", func() {
		BeforeEach(func() {
			msgMatcher := mock.MatchedBy(func(msgs []*message.Message) bool {
				if len(msgs) == 0 || msgs[0] == nil {
					return false
				}
				md := msgs[0].Metadata
				return md.Get(common.MetadataCorrelationIDKey) == "corr-123" &&
					md.Get(common.MetadataCausationIDKey) == "" &&
					md.Get(common.MetadataMessageIDKey) != ""
			})
			pub.On("Publish", topic, msgMatcher).Return(nil)
			result = base.PublishWithExtras(ctx, topic, payload, extras)
		})

		It("publishes with correlation metadata and message id", func() {
			Expect(result).ToNot(HaveOccurred())
			pub.AssertExpectations(GinkgoT())
		})
	})

	Context("with causation and extras", func() {
		BeforeEach(func() {
			ctx = common.WithCausation(ctx, "cause-456")
			extras = map[string]string{"k1": "v1", "empty": ""}
			msgMatcher := mock.MatchedBy(func(msgs []*message.Message) bool {
				if len(msgs) == 0 || msgs[0] == nil {
					return false
				}
				md := msgs[0].Metadata
				return md.Get(common.MetadataCorrelationIDKey) == "corr-123" &&
					md.Get(common.MetadataCausationIDKey) == "cause-456" &&
					md.Get("k1") == "v1" &&
					md.Get("empty") == ""
			})
			pub.On("Publish", topic, msgMatcher).Return(nil)
			result = base.PublishWithExtras(ctx, topic, payload, extras)
		})

		It("includes correlation, causation and merges non-empty extras", func() {
			Expect(result).ToNot(HaveOccurred())
			pub.AssertExpectations(GinkgoT())
		})
	})

	Context("missing correlation in context", func() {
		BeforeEach(func() {
			// no correlation set
			ctx = context.Background()
			result = base.PublishWithExtras(ctx, topic, payload, extras)
		})

		It("returns 500 and does not publish", func() {
			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			pub.AssertNotCalled(GinkgoT(), "Publish", mock.Anything, mock.Anything)
		})
	})

	Context("payload marshal failure", func() {
		BeforeEach(func() {
			// json cannot marshal functions/channels
			payload = func() {}
			result = base.PublishWithExtras(ctx, topic, payload, extras)
		})

		It("returns 500 and does not publish", func() {
			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			pub.AssertNotCalled(GinkgoT(), "Publish", mock.Anything, mock.Anything)
		})
	})

	Context("publisher returns error", func() {
		BeforeEach(func() {
			pub.On("Publish", topic, mock.Anything).Return(errors.New("pub-fail"))
			result = base.PublishWithExtras(ctx, topic, payload, extras)
		})

		It("wraps and returns server error", func() {
			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			pub.AssertExpectations(GinkgoT())
		})
	})
})
