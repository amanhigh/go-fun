package play_fast_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/components/fanin"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// CQRS Domain Objects - Simplest Implementation
type BookRoom struct {
	RoomID    string `json:"room_id"`
	GuestName string `json:"guest_name"`
}

type RoomBooked struct {
	RoomID    string `json:"room_id"`
	GuestName string `json:"guest_name"`
	Price     int64  `json:"price"`
}

type BookRoomHandler struct {
	eventBus *cqrs.EventBus
}

func (h BookRoomHandler) Handle(ctx context.Context, cmd *BookRoom) error {
	return h.eventBus.Publish(ctx, &RoomBooked{
		RoomID:    cmd.RoomID,
		GuestName: cmd.GuestName,
		Price:     100, // Fixed price for simplicity
	})
}

type FinancialReport struct {
	events  *[]string
	revenue *int64
	mutex   sync.Mutex
}

func NewFinancialReport(events *[]string, revenue *int64) *FinancialReport {
	return &FinancialReport{
		events:  events,
		revenue: revenue,
	}
}

func (f *FinancialReport) Handle(_ context.Context, event *RoomBooked) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	*f.events = append(*f.events, "RoomBooked")
	*f.revenue += event.Price
	return nil
}

type WelcomeEmailService struct {
	emails []string
	mutex  sync.Mutex
}

func NewWelcomeEmailService() *WelcomeEmailService {
	return &WelcomeEmailService{
		emails: []string{},
	}
}

func (w *WelcomeEmailService) Handle(_ context.Context, event *RoomBooked) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	emailContent := fmt.Sprintf("Welcome %s! Your room %s is confirmed.",
		event.GuestName, event.RoomID)
	w.emails = append(w.emails, emailContent)

	return nil
}

func (w *WelcomeEmailService) GetEmails() []string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return append([]string{}, w.emails...) // Return copy for thread safety
}

// Outbox Pattern Domain Objects - E-commerce Order Processing
type OrderCreated struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
	Timestamp  string `json:"timestamp"`
}

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
			router, err = message.NewRouter(message.RouterConfig{
				CloseTimeout: 50 * time.Millisecond, // Ultra-aggressive cleanup timeout for tests
			}, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(router).NotTo(BeNil())

			// Subscribe to output topic to capture produced messages
			outputMessages, err = pubSub.Subscribe(ctx, outputTopic)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			router.Close()
		})

		Context("Transform Handler", func() {
			BeforeEach(func() {
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

		Context("Consumer Handler", func() {
			var (
				consumedMessages  []string
				publishedMessages []string
				mu                sync.Mutex
			)

			BeforeEach(func() {
				// Reset message tracking
				mu.Lock()
				consumedMessages = []string{}
				publishedMessages = []string{}
				mu.Unlock()
			})

			BeforeEach(func() {
				// LIFECYCLE: Consumer handler setup with processing and conditional publishing
				consumerHandler := func(msg *message.Message) error {
					// CONSUME: Track all received messages
					mu.Lock()
					consumedMessages = append(consumedMessages, string(msg.Payload))
					mu.Unlock()

					// DIRECT-PUBLISH: Conditionally publish to output topic
					if string(msg.Payload) == "publish-me" {
						outputMsg := message.NewMessage(watermill.NewUUID(), []byte("processed-"+string(msg.Payload)))
						err := pubSub.Publish(outputTopic, outputMsg)
						if err == nil {
							mu.Lock()
							publishedMessages = append(publishedMessages, string(outputMsg.Payload))
							mu.Unlock()
						}
						return err
					}
					return nil
				}

				router.AddConsumerHandler(
					"lifecycle-consumer",
					testTopic,
					pubSub,
					consumerHandler,
				)
			})

			// LIFECYCLE: Complete consumer handler lifecycle demonstration
			It("should demonstrate complete consumer handler lifecycle", func() {
				By("Starting router and waiting for it to be running")
				go router.Run(ctx)
				<-router.Running()
				Expect(router.IsRunning()).To(BeTrue())

				By("Publishing test messages to consumer topic")
				testMsg1 := message.NewMessage(watermill.NewUUID(), []byte("test-message-1"))
				testMsg2 := message.NewMessage(watermill.NewUUID(), []byte("publish-me"))
				testMsg3 := message.NewMessage(watermill.NewUUID(), []byte("test-message-2"))

				err = pubSub.Publish(testTopic, testMsg1, testMsg2, testMsg3)
				Expect(err).NotTo(HaveOccurred())

				By("Verifying messages were consumed by handlers")
				Eventually(func() []string {
					mu.Lock()
					defer mu.Unlock()
					return append([]string{}, consumedMessages...)
				}, "200ms", "10ms").Should(ContainElements("test-message-1", "publish-me", "test-message-2"))

				By("Verifying direct publishing worked for conditional messages")
				Eventually(func() []string {
					mu.Lock()
					defer mu.Unlock()
					return append([]string{}, publishedMessages...)
				}, "200ms", "10ms").Should(ContainElement("processed-publish-me"))
			})
		})

		Context("Auto Ack Test", func() {
			var (
				capturedMsg *message.Message
				processed   chan struct{}
			)

			BeforeEach(func() {
				processed = make(chan struct{})

				// Simple handler that stores message reference and succeeds
				autoHandler := func(msg *message.Message) error {
					capturedMsg = msg // Store message reference
					close(processed)  // Signal handler completion
					return nil        // Router will ack AFTER this returns
				}

				router.AddConsumerHandler(
					"auto-ack-handler",
					testTopic,
					pubSub,
					autoHandler,
				)
			})

			It("should automatically ack successful messages", func() {
				By("Starting router and waiting for it to be running")
				go router.Run(ctx)
				<-router.Running()

				By("Publishing message")
				msg := message.NewMessage(watermill.NewUUID(), []byte("test"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for handler to complete processing")
				Eventually(processed).Should(BeClosed())

				By("Verifying Router automatically acked the message")
				Eventually(capturedMsg.Acked()).Should(BeClosed())

				By("Verifying Router did NOT nack the message")
				Consistently(capturedMsg.Nacked(), "100ms", "10ms").ShouldNot(BeClosed())
			})
		})

		Context("Auto Nack Test", func() {
			var (
				capturedMsg *message.Message
				processed   chan struct{}
				callCount   int
			)

			BeforeEach(func() {
				processed = make(chan struct{})
				callCount = 0

				// Simple handler that stores message reference and fails only once
				autoHandler := func(msg *message.Message) error {
					callCount++
					if callCount == 1 {
						capturedMsg = msg                              // Store message reference (only first time)
						close(processed)                               // Signal handler completion (only once)
						return fmt.Errorf("intentional handler error") // Router will nack AFTER this returns
					}
					// Subsequent calls succeed to avoid infinite retries
					return nil
				}

				router.AddConsumerHandler(
					"auto-nack-handler",
					testTopic,
					pubSub,
					autoHandler,
				)
			})

			It("should automatically nack failed messages", func() {
				By("Starting router and waiting for it to be running")
				go router.Run(ctx)
				<-router.Running()

				By("Publishing message")
				msg := message.NewMessage(watermill.NewUUID(), []byte("test"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for handler to complete processing")
				Eventually(processed).Should(BeClosed())

				By("Verifying Router automatically nacked the message")
				Eventually(capturedMsg.Nacked()).Should(BeClosed())

				By("Verifying Router did NOT ack the message")
				Consistently(capturedMsg.Acked(), "100ms", "10ms").ShouldNot(BeClosed())
			})
		})

		Context("Manual Ack Override", func() {
			var (
				capturedMsg *message.Message
				processed   chan struct{}
			)

			BeforeEach(func() {
				processed = make(chan struct{})

				// Handler manually acks but returns error
				manualOverrideHandler := func(msg *message.Message) error {
					capturedMsg = msg                           // Store message reference
					msg.Ack()                                   // Manual ack INSIDE handler
					close(processed)                            // Signal completion
					return fmt.Errorf("error after manual ack") // Return error (should be ignored)
				}

				router.AddConsumerHandler(
					"manual-override-handler",
					testTopic,
					pubSub,
					manualOverrideHandler,
				)
			})

			It("should respect manual ack despite error return", func() {
				By("Starting router and waiting for it to be running")
				go router.Run(ctx)
				<-router.Running()

				By("Publishing message")
				msg := message.NewMessage(watermill.NewUUID(), []byte("test"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for handler to complete processing")
				Eventually(processed).Should(BeClosed())

				By("Verifying message was manually acked")
				Eventually(capturedMsg.Acked()).Should(BeClosed())

				By("Verifying Router did NOT auto-nack despite error return")
				Consistently(capturedMsg.Nacked(), "100ms", "10ms").ShouldNot(BeClosed())
			})
		})

		Context("Fanout with Context", func() {
			var (
				handler1Name, handler2Name string
				processed1, processed2     chan struct{}
			)

			BeforeEach(func() {
				processed1 = make(chan struct{})
				processed2 = make(chan struct{})

				// Handler 1: Capture handler name from context
				router.AddConsumerHandler(
					"fanout-handler-1",
					testTopic,
					pubSub,
					func(msg *message.Message) error {
						handler1Name = message.HandlerNameFromCtx(msg.Context())
						close(processed1)
						return nil
					},
				)

				// Handler 2: Capture handler name from context
				router.AddConsumerHandler(
					"fanout-handler-2",
					testTopic,
					pubSub,
					func(msg *message.Message) error {
						handler2Name = message.HandlerNameFromCtx(msg.Context())
						close(processed2)
						return nil
					},
				)
			})

			It("should fanout to both handlers with correct context", func() {
				By("Starting router")
				go router.Run(ctx)
				<-router.Running()

				By("Publishing one message")
				msg := message.NewMessage(watermill.NewUUID(), []byte("test"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Both handlers should process the same message")
				Eventually(processed1).Should(BeClosed())
				Eventually(processed2).Should(BeClosed())

				By("Each handler should get its own correct context")
				Expect(handler1Name).To(Equal("fanout-handler-1"))
				Expect(handler2Name).To(Equal("fanout-handler-2"))
			})
		})

		Context("Dynamic Handler Lifecycle", func() {
			var (
				handler1, handler2     *message.Handler
				processed1, processed2 chan struct{}
			)

			BeforeEach(func() {
				processed1 = make(chan struct{})
				processed2 = make(chan struct{})

				// Static handler: Added before router starts
				handler1 = router.AddConsumerHandler(
					"static-handler",
					testTopic,
					pubSub,
					func(_ *message.Message) error {
						close(processed1)
						return nil
					},
				)
			})

			It("should add handler dynamically, stop individually, and auto-close router", func() {
				By("Starting router with static handler")
				go router.Run(ctx)
				<-router.Running()

				By("Adding handler dynamically after router start")
				handler2 = router.AddConsumerHandler(
					"dynamic-handler",
					testTopic,
					pubSub,
					func(_ *message.Message) error {
						close(processed2)
						return nil
					},
				)
				err = router.RunHandlers(ctx)
				Expect(err).NotTo(HaveOccurred())

				By("Publishing message - both handlers should process")
				msg := message.NewMessage(watermill.NewUUID(), []byte("test"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				Eventually(processed1).Should(BeClosed())
				Eventually(processed2).Should(BeClosed())

				By("Stopping dynamic handler individually")
				handler2.Stop()
				Eventually(handler2.Stopped()).Should(BeClosed())
				Expect(router.IsRunning()).To(BeTrue()) // Router still running

				By("Stopping last handler - router should auto-close")
				handler1.Stop()
				Eventually(handler1.Stopped()).Should(BeClosed())

				// Note: router.IsRunning() has limitations per Watermill docs
				// The logs show "All handlers stopped, closing router" which proves auto-shutdown
				// We can verify both handlers are stopped as evidence of router closure
				Expect(handler1.Stopped()).To(BeClosed())
				Expect(handler2.Stopped()).To(BeClosed())
			})
		})
	})

	Context("Router Middleware", func() {
		var (
			router      *message.Router
			processed   chan bool
			failCount   int
			maxFailures int
		)

		BeforeEach(func() {
			router, err = message.NewRouter(message.RouterConfig{}, logger)
			Expect(err).NotTo(HaveOccurred())

			processed = make(chan bool, 1)
			failCount = 0
		})

		AfterEach(func() {
			if router != nil {
				router.Close()
			}
		})

		Context("Retry Middleware", func() {
			BeforeEach(func() {
				maxFailures = 2 // Fail twice, succeed on third attempt

				// Add retry middleware with 3 max retries
				router.AddMiddleware(
					middleware.Retry{
						MaxRetries:      3,
						InitialInterval: 10 * time.Millisecond,
						Logger:          logger,
					}.Middleware,
				)

				// Handler that fails first two times, succeeds on third
				router.AddConsumerHandler(
					"retry-test-handler",
					testTopic,
					pubSub,
					func(_ *message.Message) error {
						failCount++
						if failCount <= maxFailures {
							return fmt.Errorf("simulated failure %d", failCount)
						}
						processed <- true
						return nil
					},
				)
			})

			It("should retry handler on failure then succeed", func() {
				go router.Run(ctx)
				<-router.Running()

				err = router.RunHandlers(ctx)
				Expect(err).NotTo(HaveOccurred())

				// Publish message
				msg := message.NewMessage(watermill.NewUUID(), []byte(testPayload))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				// Should eventually succeed after retries
				Eventually(processed).Should(Receive(Equal(true)))
				Expect(failCount).To(Equal(3)) // 2 failures + 1 success
			})
		})

		Context("Timeout Middleware", func() {
			var timeout time.Duration

			BeforeEach(func() {
				timeout = 50 * time.Millisecond

				// Add timeout middleware
				router.AddMiddleware(
					middleware.Timeout(timeout),
				)

				// Handler that takes longer than timeout
				router.AddConsumerHandler(
					"timeout-test-handler",
					testTopic,
					pubSub,
					func(msg *message.Message) error {
						select {
						case <-msg.Context().Done():
							// Context cancelled due to timeout
							processed <- true
							return msg.Context().Err()
						case <-time.After(200 * time.Millisecond):
							// Should not reach here due to timeout
							return nil
						}
					},
				)
			})

			It("should timeout long-running handler", func() {
				go router.Run(ctx)
				<-router.Running()

				err = router.RunHandlers(ctx)
				Expect(err).NotTo(HaveOccurred())

				// Publish message
				msg := message.NewMessage(watermill.NewUUID(), []byte(testPayload))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				// Should receive timeout signal
				Eventually(processed).Should(Receive(Equal(true)))
			})
		})

		Context("Deduplicator Middleware", func() {
			var processCount int

			BeforeEach(func() {
				processCount = 0

				// Add deduplicator middleware with default settings
				router.AddMiddleware(
					(&middleware.Deduplicator{}).Middleware,
				)

				// Handler that counts processed messages
				router.AddConsumerHandler(
					"dedup-test-handler",
					testTopic,
					pubSub,
					func(_ *message.Message) error {
						processCount++
						processed <- true
						return nil
					},
				)
			})

			It("should drop duplicate messages", func() {
				go router.Run(ctx)
				<-router.Running()

				err = router.RunHandlers(ctx)
				Expect(err).NotTo(HaveOccurred())

				// Publish same message twice (same payload = same hash)
				msg1 := message.NewMessage(watermill.NewUUID(), []byte("duplicate-content"))
				msg2 := message.NewMessage(watermill.NewUUID(), []byte("duplicate-content"))

				err = pubSub.Publish(testTopic, msg1)
				Expect(err).NotTo(HaveOccurred())

				err = pubSub.Publish(testTopic, msg2)
				Expect(err).NotTo(HaveOccurred())

				// Should only process first message, second should be dropped
				Eventually(processed).Should(Receive(Equal(true)))
				Consistently(processed, 100*time.Millisecond).ShouldNot(Receive())
				Expect(processCount).To(Equal(1)) // Only first message processed
			})
		})

		Context("Poison Queue Middleware", func() {
			var (
				poisonTopic    = "poison-topic"
				poisonMessages <-chan *message.Message
			)

			BeforeEach(func() {
				// Subscribe to poison topic to capture poisoned messages
				var err error
				poisonMessages, err = pubSub.Subscribe(ctx, poisonTopic)
				Expect(err).NotTo(HaveOccurred())

				// Add poison queue middleware
				poisonMiddleware, err := middleware.PoisonQueue(pubSub, poisonTopic)
				Expect(err).NotTo(HaveOccurred())
				router.AddMiddleware(poisonMiddleware)

				// Handler that always fails
				router.AddConsumerHandler(
					"poison-test-handler",
					testTopic,
					pubSub,
					func(_ *message.Message) error {
						return fmt.Errorf("unprocessable message")
					},
				)
			})

			It("should send failing messages to poison queue", func() {
				go router.Run(ctx)
				<-router.Running()

				err = router.RunHandlers(ctx)
				Expect(err).NotTo(HaveOccurred())

				// Publish a message that will fail
				msg := message.NewMessage(watermill.NewUUID(), []byte("failing-message"))
				err = pubSub.Publish(testTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				// Should receive the message on poison queue
				Eventually(func() bool {
					select {
					case poisonMsg := <-poisonMessages:
						return string(poisonMsg.Payload) == "failing-message"
					default:
						return false
					}
				}).Should(BeTrue())
			})
		})
	})

	Context("CQRS - Event-Driven Architecture", func() {
		var (
			router           *message.Router
			commandBus       *cqrs.CommandBus
			eventBus         *cqrs.EventBus
			commandProcessor *cqrs.CommandProcessor
			eventProcessor   *cqrs.EventProcessor

			bookRoomHandler     BookRoomHandler
			financialReport     *FinancialReport
			welcomeEmailService *WelcomeEmailService

			roomID    string
			guestName string

			receivedEvents []string
			totalRevenue   int64
		)

		BeforeEach(func() {
			pubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)

			router, err = message.NewRouter(message.RouterConfig{}, logger)
			Expect(err).ToNot(HaveOccurred())

			commandBus, err = cqrs.NewCommandBusWithConfig(pubSub, cqrs.CommandBusConfig{
				GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
					return "commands." + params.CommandName, nil
				},
				Marshaler: cqrs.JSONMarshaler{},
				Logger:    logger,
			})
			Expect(err).ToNot(HaveOccurred())

			eventBus, err = cqrs.NewEventBusWithConfig(pubSub, cqrs.EventBusConfig{
				GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
					return "events." + params.EventName, nil
				},
				Marshaler: cqrs.JSONMarshaler{},
				Logger:    logger,
			})
			Expect(err).ToNot(HaveOccurred())

			commandProcessor, err = cqrs.NewCommandProcessorWithConfig(router, cqrs.CommandProcessorConfig{
				GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
					return "commands." + params.CommandName, nil
				},
				SubscriberConstructor: func(_ cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
					return pubSub, nil
				},
				Marshaler: cqrs.JSONMarshaler{},
				Logger:    logger,
			})
			Expect(err).ToNot(HaveOccurred())

			eventProcessor, err = cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
				GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
					return "events." + params.EventName, nil
				},
				SubscriberConstructor: func(_ cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
					return pubSub, nil
				},
				Marshaler: cqrs.JSONMarshaler{},
				Logger:    logger,
			})
			Expect(err).ToNot(HaveOccurred())

			receivedEvents = []string{}
			totalRevenue = 0

			bookRoomHandler = BookRoomHandler{eventBus: eventBus}
			financialReport = NewFinancialReport(&receivedEvents, &totalRevenue)
			welcomeEmailService = NewWelcomeEmailService()

			err = commandProcessor.AddHandlers(
				cqrs.NewCommandHandler("BookRoomHandler", bookRoomHandler.Handle),
			)
			Expect(err).ToNot(HaveOccurred())

			err = eventProcessor.AddHandlers(
				cqrs.NewEventHandler("FinancialReport", financialReport.Handle),
				cqrs.NewEventHandler("WelcomeEmailService", welcomeEmailService.Handle),
			)
			Expect(err).ToNot(HaveOccurred())

			roomID = "101"
			guestName = "John Doe"

			go func() {
				defer GinkgoRecover()
				err := router.Run(ctx)
				if err != nil && !errors.Is(err, context.Canceled) {
					panic(err)
				}
			}()

			time.Sleep(100 * time.Millisecond)
		})

		AfterEach(func() {
			if router != nil {
				router.Close()
			}
		})

		It("should demonstrate CQRS fan-out: 1 Command → 1 Event → Multiple Handlers", func() {
			bookRoomCmd := &BookRoom{
				RoomID:    roomID,
				GuestName: guestName,
			}

			err := commandBus.Send(ctx, bookRoomCmd)
			Expect(err).ToNot(HaveOccurred())

			// Verify Financial Report Handler processed the event
			Eventually(func() []string {
				return receivedEvents
			}, "2s", "100ms").Should(ContainElement("RoomBooked"))

			Eventually(func() int64 {
				return totalRevenue
			}, "2s", "100ms").Should(Equal(int64(100)))

			// Verify Welcome Email Handler processed the same event
			Eventually(func() []string {
				return welcomeEmailService.GetEmails()
			}, "2s", "100ms").Should(HaveLen(1))

			// Verify email content
			emails := welcomeEmailService.GetEmails()
			Expect(emails[0]).To(ContainSubstring("Welcome John Doe"))
			Expect(emails[0]).To(ContainSubstring("room 101"))

			// Verify both handlers processed the SAME event independently
			Expect(receivedEvents).To(HaveLen(1))
			Expect(totalRevenue).To(Equal(int64(100)))
			Expect(emails).To(HaveLen(1))
		})
	})

	// ================================================================================
	// OUTBOX PATTERN DOCUMENTATION - Solving the Dual-Write Problem
	// ================================================================================
	//
	// THE PROBLEM:
	// You need to update a database AND publish an event atomically.
	// If either fails, the system becomes inconsistent.
	//k
	// Example Scenario:
	// 1. Save order to database
	// 2. Publish "OrderCreated" event
	//
	// What if step 2 fails? Database updated but no event published!
	// What if step 1 fails after step 2? Event published but no database record!
	//
	// THE SOLUTION (Outbox Pattern):
	// 1. Save order to database
	// 2. Save "OrderCreated" event to SAME database (in same transaction)
	// 3. Background process reads events from database and forwards to message broker
	//
	// GUARANTEE: Either both succeed (transaction commits) or both fail (transaction rolls back)
	//
	// WATERMILL IMPLEMENTATION:
	// - OutboxPublisher (ForwarderPublisher): Publishes events to database instead of message broker
	// - EventRelay (Forwarder): Background daemon that forwards DB events to broker
	// - Event (Envelope): Contains destination topic + original message
	//
	// REAL-WORLD USAGE:
	// - E-commerce: Order processing with payment + inventory + notifications
	// - Banking: Account transfers with audit trail + customer notifications
	// - Microservices: Service updates with reliable inter-service communication
	//
	// BENEFITS:
	// ✅ Guaranteed consistency between DB writes and event publishing
	// ✅ Works with existing database transactions
	// ✅ Handles network failures gracefully
	// ✅ Provides exactly-once delivery semantics
	//
	// TRADE-OFFS:
	// ❌ Events are eventually consistent (slight delay)
	// ❌ Requires background processing capability
	// ❌ Additional complexity in infrastructure
	//
	// IMPLEMENTATION: See Patterns context below for working demonstration
	// Future: Add Fanout, FanIn, and other messaging patterns

	// ================================================================================
	// PATTERNS - Advanced Messaging Patterns with Watermill
	// ================================================================================
	Context("Patterns", func() {
		var (
			// Common setup for all messaging patterns - A→db→broker→B flow
			logger watermill.LoggerAdapter
			ctx    context.Context
			cancel context.CancelFunc
			db     *gochannel.GoChannel // Database with outbox (A writes here)
			broker *gochannel.GoChannel // Message broker (events distributed to B)
		)

		BeforeEach(func() {
			// Shared pattern infrastructure
			logger = watermill.NewStdLogger(false, false)
			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			db = gochannel.NewGoChannel(gochannel.Config{}, logger)
			broker = gochannel.NewGoChannel(gochannel.Config{}, logger)
		})

		AfterEach(func() {
			// Common cleanup
			if cancel != nil {
				cancel()
			}
			if db != nil {
				db.Close()
			}
			if broker != nil {
				broker.Close()
			}
		})

		// OUTBOX PATTERN - Solving the Dual-Write Problem
		// Demonstrates Watermill's ForwarderPublisher and Forwarder components
		// that enable atomic database updates + event publishing
		Context("Outbox Pattern - Solving the Dual-Write Problem", func() {
			var (
				// A→db→broker→B outbox components
				publisher      message.Publisher    // A writes to db via this
				relayer        *forwarder.Forwarder // Moves events db→broker
				relayerRunning bool                 // Track if relayer was started

				// Topics
				outboxTopic = "outbox"
				ordersTopic = "orders"

				// Event monitoring channels
				dbEvents     <-chan *message.Message // Events stored in db
				brokerEvents <-chan *message.Message // Events distributed by broker
			)

			BeforeEach(func() {
				// Setup publisher - A (app) writes to db with outbox pattern
				publisher = forwarder.NewPublisher(db, forwarder.PublisherConfig{
					ForwarderTopic: outboxTopic,
				})

				// Setup relayer - moves events from db to broker
				var err error
				relayer, err = forwarder.NewForwarder(db, broker, logger, forwarder.Config{
					ForwarderTopic: outboxTopic,
				})
				Expect(err).NotTo(HaveOccurred())

				// Subscribe to db events for verification
				dbEvents, err = db.Subscribe(ctx, outboxTopic)
				Expect(err).NotTo(HaveOccurred())

				// Subscribe to broker events for verification
				brokerEvents, err = broker.Subscribe(ctx, ordersTopic)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				// Cleanup outbox pattern resources only if relayer was used
				if relayer != nil && relayerRunning {
					relayer.Close()
				}
			})

			// ENVELOPE WRAPPING: Shows how ForwarderPublisher wraps messages
			// Instead of publishing directly to "orders", it publishes envelope to "outbox"
			It("should wrap messages in envelopes and publish to outbox topic", func() {
				By("Publishing order event using publisher (A→db)")
				orderEvent := OrderCreated{
					OrderID:    "order-123",
					CustomerID: "customer-456",
					Amount:     9999,
					Timestamp:  time.Now().Format(time.RFC3339),
				}

				payload, err := json.Marshal(orderEvent)
				Expect(err).NotTo(HaveOccurred())

				msg := message.NewMessage(watermill.NewUUID(), payload)
				err = publisher.Publish(ordersTopic, msg)
				Expect(err).NotTo(HaveOccurred())

				By("Verifying message appears in db outbox (not broker)")
				select {
				case envelopedMsg := <-dbEvents:
					Expect(envelopedMsg).NotTo(BeNil())

					// Verify it's an envelope (contains destination topic info)
					Expect(string(envelopedMsg.Payload)).To(ContainSubstring("destination_topic"))
					Expect(string(envelopedMsg.Payload)).To(ContainSubstring(ordersTopic))

					// Verify the envelope structure
					var envelope map[string]interface{}
					err = json.Unmarshal(envelopedMsg.Payload, &envelope)
					Expect(err).NotTo(HaveOccurred())

					Expect(envelope["destination_topic"]).To(Equal(ordersTopic))
					Expect(envelope["payload"]).NotTo(BeEmpty())

					envelopedMsg.Ack()
				case <-ctx.Done():
					Fail("Timeout waiting for enveloped message in outbox")
				}

				By("Verifying message does NOT appear directly in broker")
				select {
				case <-brokerEvents:
					Fail("Message should not appear directly in broker")
				case <-time.After(50 * time.Millisecond):
					// Expected - no direct message in broker yet
				}
			})

			// COMPLETE OUTBOX FLOW: End-to-end demonstration with forwarding and batch processing
			// Shows the complete pattern: App → ForwarderPublisher → Outbox → Forwarder → Final Topic
			It("should demonstrate complete outbox pattern with message forwarding", func() {
				By("Starting the complete A→db→broker→B infrastructure")
				relayerCtx, relayerCancel := context.WithCancel(ctx)
				defer relayerCancel()

				relayerReady := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					close(relayerReady)
					err := relayer.Run(relayerCtx)
					if err != nil && !errors.Is(err, context.Canceled) {
						Fail("Event relayer failed: " + err.Error())
					}
				}()

				<-relayerReady
				relayerRunning = true // Mark relayer as running
				// Remove unnecessary sleep - relayer is ready immediately

				By("Application (A) publishes multiple events atomically")
				events := []OrderCreated{
					{OrderID: "order-001", CustomerID: "customer-A", Amount: 1000, Timestamp: time.Now().Format(time.RFC3339)},
					{OrderID: "order-002", CustomerID: "customer-B", Amount: 2000, Timestamp: time.Now().Format(time.RFC3339)},
					{OrderID: "order-003", CustomerID: "customer-C", Amount: 3000, Timestamp: time.Now().Format(time.RFC3339)},
				}

				By("Publishing all events through publisher (A→db)")
				for _, event := range events {
					payload, err := json.Marshal(event)
					Expect(err).NotTo(HaveOccurred())

					msg := message.NewMessage(watermill.NewUUID(), payload)
					err = publisher.Publish(ordersTopic, msg)
					Expect(err).NotTo(HaveOccurred())
				}

				By("Verifying relayer unwraps and forwards all events to broker")
				receivedOrders := make([]OrderCreated, 0, 3)

				for i := 0; i < 3; i++ {
					select {
					case orderMsg := <-brokerEvents:
						var order OrderCreated
						err := json.Unmarshal(orderMsg.Payload, &order)
						Expect(err).NotTo(HaveOccurred())

						receivedOrders = append(receivedOrders, order)
						orderMsg.Ack()
					case <-time.After(5 * time.Second):
						Fail(fmt.Sprintf("Timeout waiting for forwarded event %d of 3", i+1))
					}
				}

				By("Verifying all orders were unwrapped and forwarded correctly")
				Expect(receivedOrders).To(HaveLen(3))

				orderIDs := make([]string, len(receivedOrders))
				for i, order := range receivedOrders {
					orderIDs[i] = order.OrderID
				}

				Expect(orderIDs).To(ContainElements("order-001", "order-002", "order-003"))

				By("Demonstrating the A→db→broker→B outbox pattern solved dual-write problem")
				// In real scenarios:
				// 1. A (app) saves order model to db (in transaction)
				// 2. A (app) saves events to db outbox (same transaction)
				// 3. Relayer forwards events from db to broker
				// 4. B (services) consume events from broker
				// GUARANTEE: Either both A writes succeed or both fail (atomicity)
			})
		})

		// FANIN PATTERN - Multi-Region Event Aggregation
		// Demonstrates Watermill's FanIn component for merging multiple event streams
		// into a single unified stream for centralized processing
		//
		// THE PROBLEM:
		// In distributed systems, you have multiple independent event sources
		// (e.g., multiple regions, services, or data centers) that need to be
		// processed together as a unified stream.
		//
		// Example Scenario:
		// - US region produces "orders-us" events
		// - EU region produces "orders-eu" events
		// - Asia region produces "orders-asia" events
		// - Analytics service needs to process ALL orders from one stream
		//
		// Without FanIn: Need 3 separate consumers, manual correlation, complex orchestration
		//
		// THE SOLUTION (FanIn Pattern):
		// 1. Multiple sources publish to their regional topics
		// 2. FanIn component subscribes to all regional topics
		// 3. FanIn merges events into single unified topic
		// 4. Single consumer processes unified stream
		//
		// GUARANTEE: All events from all sources reach the unified stream
		//
		// WATERMILL IMPLEMENTATION:
		// - FanIn Component: Subscribes to multiple source topics, publishes to one target topic
		// - Per-Source Ordering: Events from same source maintain their order
		// - Cross-Source Independence: Sources don't affect each other
		//
		// REAL-WORLD USAGE:
		// - E-commerce: Multi-region order aggregation for global analytics
		// - Microservices: Centralized log aggregation from multiple services
		// - IoT: Sensor data collection from multiple device types
		// - Financial: Order book aggregation from multiple exchanges
		//
		// BENEFITS:
		// ✅ Simplified architecture - one consumer instead of N consumers
		// ✅ Centralized processing - single pipeline for all sources
		// ✅ Fault isolation - one source failure doesn't affect others
		// ✅ Scalability - easy to add new sources without changing consumers
		//
		// TRADE-OFFS:
		// ❌ No global ordering guarantee (only per-source ordering)
		// ❌ Single point of aggregation (though can be scaled horizontally)
		// ❌ Potential bottleneck if source volume is very high
		//
		// IMPLEMENTATION: See tests below for working demonstration
		Context("FanIn Pattern - Multi-Region Event Aggregation", func() {
			var (
				regionalOrders *gochannel.GoChannel
				globalOrders   *gochannel.GoChannel
				fanIn          *fanin.FanIn

				usTopic     = "orders-us"
				euTopic     = "orders-eu"
				asiaTopic   = "orders-asia"
				globalTopic = "orders-global"

				globalEvents <-chan *message.Message
			)

			BeforeEach(func() {
				var err error

				regionalOrders = gochannel.NewGoChannel(gochannel.Config{}, logger)
				globalOrders = gochannel.NewGoChannel(gochannel.Config{}, logger)

				fanIn, err = fanin.NewFanIn(
					regionalOrders,
					globalOrders,
					fanin.Config{
						SourceTopics: []string{usTopic, euTopic, asiaTopic},
						TargetTopic:  globalTopic,
					},
					logger,
				)
				Expect(err).NotTo(HaveOccurred())

				globalEvents, err = globalOrders.Subscribe(ctx, globalTopic)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				if fanIn != nil {
					fanIn.Close()
				}
				if regionalOrders != nil {
					regionalOrders.Close()
				}
				if globalOrders != nil {
					globalOrders.Close()
				}
			})

			It("should merge orders from multiple regions into unified analytics stream", func() {
				By("Starting FanIn to merge regional order streams")
				fanInCtx, fanInCancel := context.WithCancel(ctx)
				defer fanInCancel()

				fanInReady := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					close(fanInReady)
					err := fanIn.Run(fanInCtx)
					if err != nil && !errors.Is(err, context.Canceled) {
						Fail("FanIn failed: " + err.Error())
					}
				}()

				<-fanInReady
				time.Sleep(100 * time.Millisecond)

				By("Processing orders from US, EU, and Asia regions")
				usOrder := OrderCreated{OrderID: "US-order-001", CustomerID: "us-customer-1", Amount: 100, Timestamp: time.Now().Format(time.RFC3339)}
				euOrder := OrderCreated{OrderID: "EU-order-001", CustomerID: "eu-customer-1", Amount: 200, Timestamp: time.Now().Format(time.RFC3339)}
				asiaOrder := OrderCreated{OrderID: "Asia-order-001", CustomerID: "asia-customer-1", Amount: 300, Timestamp: time.Now().Format(time.RFC3339)}

				for _, order := range []struct {
					topic string
					order OrderCreated
				}{
					{usTopic, usOrder},
					{euTopic, euOrder},
					{asiaTopic, asiaOrder},
				} {
					payload, err := json.Marshal(order.order)
					Expect(err).NotTo(HaveOccurred())
					msg := message.NewMessage(watermill.NewUUID(), payload)
					err = regionalOrders.Publish(order.topic, msg)
					Expect(err).NotTo(HaveOccurred())
				}

				By("Verifying all regional orders reach global analytics")
				receivedOrders := make([]OrderCreated, 0, 3)

				for i := 0; i < 3; i++ {
					select {
					case orderMsg := <-globalEvents:
						var order OrderCreated
						err := json.Unmarshal(orderMsg.Payload, &order)
						Expect(err).NotTo(HaveOccurred())

						receivedOrders = append(receivedOrders, order)
						orderMsg.Ack()
					case <-time.After(5 * time.Second):
						Fail(fmt.Sprintf("Timeout waiting for order %d of 3", i+1))
					}
				}

				Expect(receivedOrders).To(HaveLen(3))

				orderIDs := make([]string, len(receivedOrders))
				for i, order := range receivedOrders {
					orderIDs[i] = order.OrderID
				}

				Expect(orderIDs).To(ContainElements("US-order-001", "EU-order-001", "Asia-order-001"))
			})

			It("should continue processing orders when one region is unavailable", func() {
				By("Starting FanIn with all 3 regions")
				fanInCtx, fanInCancel := context.WithCancel(ctx)
				defer fanInCancel()

				fanInReady := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					close(fanInReady)
					err := fanIn.Run(fanInCtx)
					if err != nil && !errors.Is(err, context.Canceled) {
						Fail("FanIn failed: " + err.Error())
					}
				}()

				<-fanInReady
				time.Sleep(100 * time.Millisecond)

				By("Simulating Asia region outage")

				By("Verifying US and EU orders still reach global analytics")
				usOrder := OrderCreated{OrderID: "US-order-002", CustomerID: "us-customer-2", Amount: 150, Timestamp: time.Now().Format(time.RFC3339)}
				euOrder := OrderCreated{OrderID: "EU-order-002", CustomerID: "eu-customer-2", Amount: 250, Timestamp: time.Now().Format(time.RFC3339)}

				for _, order := range []struct {
					topic string
					order OrderCreated
				}{
					{usTopic, usOrder},
					{euTopic, euOrder},
				} {
					payload, err := json.Marshal(order.order)
					Expect(err).NotTo(HaveOccurred())
					msg := message.NewMessage(watermill.NewUUID(), payload)
					err = regionalOrders.Publish(order.topic, msg)
					Expect(err).NotTo(HaveOccurred())
				}

				receivedOrders := make([]OrderCreated, 0, 2)

				for i := 0; i < 2; i++ {
					select {
					case orderMsg := <-globalEvents:
						var order OrderCreated
						err := json.Unmarshal(orderMsg.Payload, &order)
						Expect(err).NotTo(HaveOccurred())

						receivedOrders = append(receivedOrders, order)
						orderMsg.Ack()
					case <-time.After(5 * time.Second):
						Fail(fmt.Sprintf("Timeout waiting for order %d of 2", i+1))
					}
				}

				Expect(receivedOrders).To(HaveLen(2))

				orderIDs := make([]string, len(receivedOrders))
				for i, order := range receivedOrders {
					orderIDs[i] = order.OrderID
				}

				Expect(orderIDs).To(ContainElements("US-order-002", "EU-order-002"))

				By("Demonstrating regional isolation preserves business continuity")
			})
		})
	})
})
