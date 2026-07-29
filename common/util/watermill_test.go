package util_test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
)

var _ = Describe("Watermill Utilities", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		logger watermill.LoggerAdapter
		pubSub *gochannel.GoChannel
	)

	BeforeEach(func() {
		logger = watermill.NewStdLogger(false, false)
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		pubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
		if pubSub != nil {
			pubSub.Close()
		}
	})

	// =========================================================================
	// 1. PoisonTopic — pure string utility
	// =========================================================================
	Context("PoisonTopic", func() {
		It("should append .poison suffix to the source topic", func() {
			Expect(util.PoisonTopic("orders")).To(Equal("orders.poison"))
		})

		It("should handle an empty source topic", func() {
			Expect(util.PoisonTopic("")).To(Equal(".poison"))
		})

		It("should handle topics that already contain dots", func() {
			Expect(util.PoisonTopic("a.b.c")).To(Equal("a.b.c.poison"))
		})
	})

	// =========================================================================
	// 2. AddConsumerHandler — registers handler by topic name and recovers panics
	// =========================================================================
	Context("AddConsumerHandler", func() {
		var (
			router *message.Router
		)

		BeforeEach(func() {
			var err error
			router, err = message.NewRouter(message.RouterConfig{}, logger)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if router != nil {
				router.Close()
			}
		})

		Context("handler registration", func() {
			var topic string

			BeforeEach(func() {
				topic = "events"
			})

			It("should register the handler using the topic as the Watermill handler ID", func() {
				handlerID := make(chan string, 1)

				handler := util.AddConsumerHandler(router, topic, pubSub, func(msg *message.Message) error {
					handlerID <- message.HandlerNameFromCtx(msg.Context())
					return nil
				})
				Expect(handler).ToNot(BeNil())

				go func() {
					defer GinkgoRecover()
					_ = router.Run(ctx)
				}()
				<-router.Running()
				_ = router.RunHandlers(ctx)

				msg := message.NewMessage(watermill.NewUUID(), []byte("trigger"))
				Expect(pubSub.Publish(topic, msg)).To(Succeed())

				var received string
				Eventually(handlerID).Should(Receive(&received))
				Expect(received).To(Equal(topic))
			})
		})

		Context("panic recovery", func() {
			var (
				recovered   chan error
				panicReason string
			)

			BeforeEach(func() {
				recovered = make(chan error, 1)
				panicReason = "unexpected nil pointer"
			})

			It("should recover a panicking handler via the helper's Recoverer, observed through the capture middleware", func() {
				// Capture middleware passed through the variadic argument.
				// It wraps the handler (outermost) and observes RecoveredPanicError
				// returned by the helper's internal Recoverer.
				captureMiddleware := func(h message.HandlerFunc) message.HandlerFunc {
					return func(msg *message.Message) ([]*message.Message, error) {
						events, err := h(msg)
						if err != nil {
							var rpe middleware.RecoveredPanicError
							if errors.As(err, &rpe) {
								recovered <- rpe
							}
						}
						return events, nil // swallow error so router doesn't NACK-loop
					}
				}

				util.AddConsumerHandler(router, "panic-topic", pubSub, func(_ *message.Message) error {
					panic(panicReason)
				}, captureMiddleware)

				go func() {
					defer GinkgoRecover()
					_ = router.Run(ctx)
				}()
				<-router.Running()
				_ = router.RunHandlers(ctx)

				msg := message.NewMessage(watermill.NewUUID(), []byte("trigger"))
				Expect(pubSub.Publish("panic-topic", msg)).To(Succeed())

				var rpe middleware.RecoveredPanicError
				Eventually(recovered).Should(Receive(&rpe))
				Expect(rpe.Error()).To(ContainSubstring(panicReason))
			})
		})
	})

	// =========================================================================
	// 3. AddRetryPoisonConsumerHandler — retry, Recoverer, and poison queue
	// =========================================================================
	Context("AddRetryPoisonConsumerHandler", func() {
		var (
			router *message.Router
		)

		BeforeEach(func() {
			var err error
			router, err = message.NewRouter(message.RouterConfig{}, logger)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if router != nil {
				router.Close()
			}
		})

		Context("retry + recoverer + poison queue composition", func() {
			var (
				attempts       int
				mu             sync.Mutex
				poisonSeen     chan *message.Message
				poisonMessages <-chan *message.Message
			)

			BeforeEach(func() {
				attempts = 0
				poisonSeen = make(chan *message.Message, 1)

				// Subscribe to the poison topic before the router starts
				// so the dedicated poison handler does not consume the
				// message before the test subscription exists.
				var subErr error
				poisonMessages, subErr = pubSub.Subscribe(ctx, util.PoisonTopic("retry-source"))
				Expect(subErr).NotTo(HaveOccurred())
			})

			It("should retry twice, catch a panic via Recoverer, and deliver one message with poison metadata", func() {
				poisonHandler := func(msg *message.Message) error {
					poisonSeen <- msg
					return nil
				}

				retry := middleware.Retry{
					MaxRetries:      2,
					InitialInterval: 10 * time.Millisecond,
				}

				err := util.AddRetryPoisonConsumerHandler(
					router,
					pubSub, // publisher (for poison queue)
					pubSub, // subscriber
					util.RetryPoisonConsumerConfig{
						Topic: "retry-source",
						Retry: retry,
						Handler: func(_ *message.Message) error {
							mu.Lock()
							attempts++
							mu.Unlock()
							panic("deliberate handler panic")
						},
						PoisonHandler: poisonHandler,
					},
				)
				Expect(err).NotTo(HaveOccurred())

				go func() {
					defer GinkgoRecover()
					_ = router.Run(ctx)
				}()
				<-router.Running()
				_ = router.RunHandlers(ctx)

				By("Publishing one original message that will always panic")
				orig := message.NewMessage(watermill.NewUUID(), []byte("doomed"))
				Expect(pubSub.Publish("retry-source", orig)).To(Succeed())

				By("Waiting for exactly one message on the poison topic")
				var poisonMsg *message.Message
				Eventually(poisonMessages).Should(Receive(&poisonMsg))

				By("Verifying the original payload is preserved")
				Expect(string(poisonMsg.Payload)).To(Equal("doomed"))

				By("Verifying Watermill poison metadata keys are set")
				Expect(poisonMsg.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal("retry-source"))
				Expect(poisonMsg.Metadata.Get(middleware.PoisonedHandlerKey)).To(Equal("retry-source"))
				Expect(poisonMsg.Metadata.Get(middleware.ReasonForPoisonedKey)).To(ContainSubstring("panic"))

				By("Verifying the dedicated poison handler receives the message")
				var received *message.Message
				Eventually(poisonSeen).Should(Receive(&received))
				Expect(string(received.Payload)).To(Equal("doomed"))

				By("Verifying handler was attempted 3 times (initial + 2 retries)")
				mu.Lock()
				finalAttempts := attempts
				mu.Unlock()
				Expect(finalAttempts).To(Equal(3))

				By("Verifying no additional messages appear on the poison topic")
				Consistently(poisonMessages, "200ms", "50ms").ShouldNot(Receive())
			})
		})
	})
})
