package play_fast_test

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Watermill", func() {
	var (
		publisher  message.Publisher
		subscriber message.Subscriber
		logger     watermill.LoggerAdapter
		err        error
		ctx        context.Context
		cancel     context.CancelFunc
	)

	const (
		testTopic   = "test-topic"
		testPayload = "test-payload"
	)

	BeforeEach(func() {
		logger = watermill.NewStdLogger(false, false)
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		pubSub := gochannel.NewGoChannel(
			gochannel.Config{},
			logger,
		)

		publisher = pubSub
		subscriber = pubSub
	})

	AfterEach(func() {
		cancel()
		publisher.Close()
		subscriber.Close()
	})

	It("should build publisher and subscriber", func() {
		Expect(publisher).NotTo(BeNil())
		Expect(subscriber).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Basic Operations", func() {
		var (
			messages <-chan *message.Message
		)

		BeforeEach(func() {
			messages, err = subscriber.Subscribe(ctx, testTopic)
			Expect(err).NotTo(HaveOccurred())
			Expect(messages).NotTo(BeNil())
		})

		Context("Publish", func() {
			var (
				msg         *message.Message
				receivedMsg *message.Message
			)

			BeforeEach(func() {
				msg = message.NewMessage(watermill.NewUUID(), []byte(testPayload))

				By("Publishing a message")
				err = publisher.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Receiving the message")
				select {
				case receivedMsg = <-messages:
					Expect(receivedMsg).NotTo(BeNil())
					Expect(string(receivedMsg.Payload)).To(Equal(testPayload))
					Expect(receivedMsg.UUID).To(Equal(msg.UUID))
				case <-ctx.Done():
					Fail("Timeout waiting for message")
				}
			})

			// RECEIVE: Gets message from topic/queue for processing
			It("should publish and receive a message", func() {
				Expect(receivedMsg).NotTo(BeNil())
				Expect(string(receivedMsg.Payload)).To(Equal(testPayload))
				Expect(receivedMsg.UUID).To(Equal(msg.UUID))
			})

			// ACK: "I successfully processed this message" - message deleted/marked done
			It("should handle message acknowledgment", func() {
				By("Acknowledging the message")
				receivedMsg.Ack()

				By("Verifying message context is done after ack")
				Eventually(receivedMsg.Context().Done()).Should(BeClosed())
			})

			// NACK: "I failed to process this message" - message redelivered or dead letter
			It("should handle message nack", func() {
				By("Nacking the message")
				receivedMsg.Nack()

				By("Verifying message can be nacked without error")
				Expect(receivedMsg.UUID).To(Equal(msg.UUID))
				Expect(string(receivedMsg.Payload)).To(Equal(testPayload))
			})
		})
	})
})
