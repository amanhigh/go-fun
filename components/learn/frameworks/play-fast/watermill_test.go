package play_fast_test

import (
	"context"
	"fmt"
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
		pubSub     *gochannel.GoChannel
	)

	const (
		testTopic   = "test-topic"
		testPayload = "test-payload"
		inputTopic  = "input-topic"
		outputTopic = "output-topic"
	)

	BeforeEach(func() {
		logger = watermill.NewStdLogger(false, false)
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		pubSub = gochannel.NewGoChannel(
			gochannel.Config{},
			logger,
		)

		publisher = pubSub
		subscriber = pubSub
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
		if pubSub != nil {
			pubSub.Close()
		}
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

		Context("Multiple Messages", func() {
			Context("Batch Publishing", func() {
				var (
					msg1, msg2, msg3 *message.Message
					receivedMsgs     []*message.Message
					expectedPayloads []string
				)

				BeforeEach(func() {
					expectedPayloads = []string{"payload-1", "payload-2", "payload-3"}
					msg1 = message.NewMessage(watermill.NewUUID(), []byte(expectedPayloads[0]))
					msg2 = message.NewMessage(watermill.NewUUID(), []byte(expectedPayloads[1]))
					msg3 = message.NewMessage(watermill.NewUUID(), []byte(expectedPayloads[2]))

					By("Publishing multiple messages at once")
					err = publisher.Publish(testTopic, msg1, msg2, msg3)
					Expect(err).NotTo(HaveOccurred())

					By("Receiving all published messages")
					receivedMsgs = make([]*message.Message, 0, 3)
					for i := 0; i < 3; i++ {
						select {
						case receivedMsg := <-messages:
							Expect(receivedMsg).NotTo(BeNil())
							receivedMsgs = append(receivedMsgs, receivedMsg)
						case <-ctx.Done():
							Fail("Timeout waiting for message")
						}
					}
				})

				// BATCH: Publishing multiple messages in single operation (order not guaranteed)
				// Important: Different pub/sub implementations have different ordering guarantees:
				// - GoChannel (in-memory): No FIFO guarantee for batch operations
				// - Kafka: Preserves order within partitions
				// - RabbitMQ: Order depends on configuration
				// Batch publishing is about throughput and efficiency, not ordering guarantees
				It("should publish and receive all messages", func() {
					Expect(receivedMsgs).To(HaveLen(3))

					receivedPayloads := make([]string, len(receivedMsgs))
					for i, msg := range receivedMsgs {
						receivedPayloads[i] = string(msg.Payload)
					}
					Expect(receivedPayloads).To(ContainElements(expectedPayloads))
				})
			})

			Context("Sequential Processing", func() {
				var (
					seqMsgs      []*message.Message
					receivedMsgs []*message.Message
				)

				BeforeEach(func() {
					seqMsgs = make([]*message.Message, 3)
					for i := 0; i < 3; i++ {
						payload := fmt.Sprintf("sequence-%d", i+1)
						seqMsgs[i] = message.NewMessage(watermill.NewUUID(), []byte(payload))
					}

					By("Publishing messages with sequence numbers")
					err = publisher.Publish(testTopic, seqMsgs...)
					Expect(err).NotTo(HaveOccurred())

					By("Receiving messages for processing")
					receivedMsgs = make([]*message.Message, 0, 3)
					for i := 0; i < 3; i++ {
						select {
						case receivedMsg := <-messages:
							Expect(receivedMsg).NotTo(BeNil())
							receivedMsgs = append(receivedMsgs, receivedMsg)
						case <-ctx.Done():
							Fail("Timeout waiting for message")
						}
					}
				})

				// SEQUENCE: All sequence messages received
				It("should process all sequence messages", func() {
					receivedPayloads := make([]string, len(receivedMsgs))
					for i, msg := range receivedMsgs {
						receivedPayloads[i] = string(msg.Payload)
					}
					expectedSequences := []string{"sequence-1", "sequence-2", "sequence-3"}
					Expect(receivedPayloads).To(ContainElements(expectedSequences))
				})

				// MIXED RESPONSES: Real-world ack/nack scenarios
				It("should handle mixed ack/nack scenarios", func() {
					By("Acking first received message")
					receivedMsgs[0].Ack()

					By("Nacking second received message")
					receivedMsgs[1].Nack()

					By("Acking third received message")
					receivedMsgs[2].Ack()

					By("Verifying all operations completed without error")
					Expect(receivedMsgs).To(HaveLen(3))
					for _, msg := range receivedMsgs {
						payload := string(msg.Payload)
						Expect(payload).To(MatchRegexp("sequence-[1-3]"))
					}
				})
			})
		})
	})

	Context("Router", func() {
		var (
			router         *message.Router
			outputMessages <-chan *message.Message
		)

		BeforeEach(func() {
			router, err = message.NewRouter(message.RouterConfig{}, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(router).NotTo(BeNil())

			// Subscribe to output topic to capture produced messages
			outputMessages, err = pubSub.Subscribe(ctx, outputTopic)
			Expect(err).NotTo(HaveOccurred())

			// TRANSFORM: Handler that transforms input to output message
			transformHandler := func(msg *message.Message) ([]*message.Message, error) {
				outputPayload := fmt.Sprintf("processed-%s", string(msg.Payload))
				outputMsg := message.NewMessage(watermill.NewUUID(), []byte(outputPayload))
				return []*message.Message{outputMsg}, nil
			}

			router.AddHandler(
				"transform-handler",
				inputTopic,
				pubSub,
				outputTopic,
				pubSub,
				transformHandler,
			)
		})

		AfterEach(func() {
			router.Close()
		})

		// TRANSFORM: Router handlers can transform and route messages to other topics
		// This test covers: router creation, handler registration, message processing, and topic routing
		It("should transform and route messages", func() {
			By("Starting the router")
			go router.Run(ctx)
			defer router.Close()

			// Wait for router to start
			<-router.Running()

			By("Publishing message to input topic")
			inputMsg := message.NewMessage(watermill.NewUUID(), []byte("transform-me"))
			err = pubSub.Publish(inputTopic, inputMsg)
			Expect(err).NotTo(HaveOccurred())

			By("Receiving transformed message from output topic")
			select {
			case outputMsg := <-outputMessages:
				Expect(outputMsg).NotTo(BeNil())
				Expect(string(outputMsg.Payload)).To(Equal("processed-transform-me"))
				outputMsg.Ack()
			case <-ctx.Done():
				Fail("Timeout waiting for output message")
			}
		})
	})
})
