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
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
)

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

var _ = Describe("BasePublisher", func() {
	var (
		pub       *mockPublisher
		base      publisher.BasePublisher
		published *message.Message
		result    common.HttpError
	)

	BeforeEach(func() {
		pub = &mockPublisher{}
		base = publisher.NewBasePublisher(pub)
	})

	Context("root publication", func() {
		BeforeEach(func() {
			pub.On("Publish", "topic", mock.Anything).Return(nil)
			result = base.PublishRoot(context.Background(), "topic", struct{ ID string }{ID: "p-1"})
			args := pub.Calls[0].Arguments
			publishedMessages, ok := args.Get(1).([]*message.Message)
			Expect(ok).To(BeTrue())
			published = publishedMessages[0]
		})

		It("publishes with fresh message and correlation IDs", func() {
			Expect(result).ToNot(HaveOccurred())
			Expect(published.UUID).NotTo(BeEmpty())
			Expect(middleware.MessageCorrelationID(published)).NotTo(BeEmpty())
			Expect(published.Metadata.Get("causation_id")).To(BeEmpty())
			pub.AssertExpectations(GinkgoT())
		})
	})

	Context("child publication", func() {
		var parent common.Metadata

		BeforeEach(func() {
			parent = common.NewRootMetadata("parent-message", "correlation-1")
			ctx := common.WithMetadata(context.Background(), parent)
			pub.On("Publish", "topic", mock.Anything).Return(nil)
			result = base.PublishChild(ctx, "topic", struct{ ID string }{ID: "p-1"})
			args := pub.Calls[0].Arguments
			publishedMessages, ok := args.Get(1).([]*message.Message)
			Expect(ok).To(BeTrue())
			published = publishedMessages[0]
		})

		It("uses the parent correlation and message IDs", func() {
			Expect(result).ToNot(HaveOccurred())
			Expect(published.UUID).NotTo(Equal(parent.MessageID))
			Expect(middleware.MessageCorrelationID(published)).To(Equal(parent.CorrelationID))
			Expect(published.Metadata.Get("causation_id")).To(Equal(parent.MessageID))
		})
	})

	Context("degraded child publication", func() {
		BeforeEach(func() {
			pub.On("Publish", "topic", mock.Anything).Return(nil)
			result = base.PublishChild(context.Background(), "topic", struct{}{})
			args := pub.Calls[0].Arguments
			publishedMessages, ok := args.Get(1).([]*message.Message)
			Expect(ok).To(BeTrue())
			published = publishedMessages[0]
		})

		It("uses fresh identity metadata", func() {
			Expect(result).ToNot(HaveOccurred())
			Expect(published.UUID).NotTo(BeEmpty())
			Expect(middleware.MessageCorrelationID(published)).NotTo(BeEmpty())
			Expect(published.Metadata.Get("causation_id")).To(BeEmpty())
			pub.AssertExpectations(GinkgoT())
		})
	})

	Context("degraded child with a parent message ID", func() {
		var parent common.Metadata

		BeforeEach(func() {
			parent = common.NewRootMetadata("parent-message", "")
			ctx := common.WithMetadata(context.Background(), parent)
			pub.On("Publish", "topic", mock.Anything).Return(nil)
			result = base.PublishChild(ctx, "topic", struct{}{})
			args := pub.Calls[0].Arguments
			publishedMessages, ok := args.Get(1).([]*message.Message)
			Expect(ok).To(BeTrue())
			published = publishedMessages[0]
		})

		It("keeps the parent message ID as child causation", func() {
			Expect(result).ToNot(HaveOccurred())
			Expect(middleware.MessageCorrelationID(published)).NotTo(BeEmpty())
			Expect(published.Metadata.Get("causation_id")).To(Equal(parent.MessageID))
			pub.AssertExpectations(GinkgoT())
		})
	})

	Context("publisher failure", func() {
		BeforeEach(func() {
			pub.On("Publish", "topic", mock.Anything).Return(errors.New("pub-fail"))
			result = base.PublishRoot(context.Background(), "topic", struct{}{})
		})

		It("returns a server error", func() {
			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			pub.AssertExpectations(GinkgoT())
		})
	})
})
