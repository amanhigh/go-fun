package util_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	common "github.com/amanhigh/go-fun/models/common"
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

	// =========================================================================
	// 4. AddConsumerHandler — registers handler by topic name and recovers panics
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
			var panicReason string

			BeforeEach(func() {
				panicReason = "unexpected nil pointer"
			})

			It("should recover a panicking handler via the helper's Recoverer", func() {
				called := make(chan struct{}, 1)

				util.AddConsumerHandler(router, "panic-topic", pubSub, func(_ *message.Message) error {
					called <- struct{}{}
					panic(panicReason)
				})

				go func() {
					defer GinkgoRecover()
					_ = router.Run(ctx)
				}()
				<-router.Running()
				_ = router.RunHandlers(ctx)

				msg := message.NewMessage(watermill.NewUUID(), []byte("trigger"))
				Expect(pubSub.Publish("panic-topic", msg)).To(Succeed())

				// Handler was invoked — Recoverer caught the panic so the goroutine survived.
				Eventually(called).Should(Receive())
				// Router remains functional after recovery.
				Expect(ctx.Err()).NotTo(HaveOccurred())
			})
		})

		// =====================================================================
		// 5. Retry and DLQ configuration — unified ConsumerConfig tests
		// =====================================================================
		Context("retry and DLQ configuration", func() {
			var (
				startRouter func()
				retryConfig middleware.Retry
			)

			BeforeEach(func() {
				retryConfig = util.DefaultRetryConfig()
				retryConfig.InitialInterval = 10 * time.Millisecond
				startRouter = func() {
					go func() {
						defer GinkgoRecover()
						_ = router.Run(ctx)
					}()
					<-router.Running()
					_ = router.RunHandlers(ctx)
				}
			})

			Context("default retry config", func() {
				It("should retry twice (initial + 2 retries) before DLQ", func() {
					var attempts atomic.Int32
					poisonMessages, subErr := pubSub.Subscribe(ctx, util.PoisonTopic("default-retry"))
					Expect(subErr).NotTo(HaveOccurred())

					util.AddConsumerHandler(router, "default-retry", pubSub, func(_ *message.Message) error {
						attempts.Add(1)
						return errors.New("transient failure")
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("doomed"))
					Expect(pubSub.Publish("default-retry", msg)).To(Succeed())

					By("Waiting for exactly one message on the poison topic")
					var poisonMsg *message.Message
					Eventually(poisonMessages).Should(Receive(&poisonMsg))
					Expect(string(poisonMsg.Payload)).To(Equal("doomed"))

					By("Verifying handler was attempted 3 times (initial + 2 retries)")
					Expect(attempts.Load()).To(Equal(int32(3)))

					By("Verifying no additional messages appear on the poison topic")
					Consistently(poisonMessages, "200ms", "50ms").ShouldNot(Receive())
				})
			})

			Context("MaxRetries override", func() {
				It("should respect a custom MaxRetries value", func() {
					var attempts atomic.Int32
					poisonMessages, subErr := pubSub.Subscribe(ctx, util.PoisonTopic("custom-retry"))
					Expect(subErr).NotTo(HaveOccurred())

					retryConfig.MaxRetries = 4

					util.AddConsumerHandler(router, "custom-retry", pubSub, func(_ *message.Message) error {
						attempts.Add(1)
						return errors.New("transient failure")
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("doomed"))
					Expect(pubSub.Publish("custom-retry", msg)).To(Succeed())

					var poisonMsg *message.Message
					Eventually(poisonMessages).Should(Receive(&poisonMsg))

					Expect(attempts.Load()).To(Equal(int32(5))) // initial + 4 retries
				})
			})

			Context("permanent malformed-JSON pipeline", func() {
				It("should bypass retries for malformed JSON", func() {
					var attempts atomic.Int32
					poisonMessages, subErr := pubSub.Subscribe(ctx, util.PoisonTopic("perm-json"))
					Expect(subErr).NotTo(HaveOccurred())

					util.AddConsumerHandler(router, "perm-json", pubSub, func(_ *message.Message) error {
						attempts.Add(1)
						return &json.SyntaxError{Offset: 5}
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("bad-json"))
					Expect(pubSub.Publish("perm-json", msg)).To(Succeed())

					var poisonMsg *message.Message
					Eventually(poisonMessages).Should(Receive(&poisonMsg))

					Expect(attempts.Load()).To(Equal(int32(1)))
				})
			})

			Context("DLQ publication", func() {
				It("should publish exhausted messages to poison topic with metadata", func() {
					poisonMessages, subErr := pubSub.Subscribe(ctx, util.PoisonTopic("dlq-meta"))
					Expect(subErr).NotTo(HaveOccurred())

					retryConfig.MaxRetries = 1

					util.AddConsumerHandler(router, "dlq-meta", pubSub, func(_ *message.Message) error {
						return errors.New("fatal error")
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("payload"))
					Expect(pubSub.Publish("dlq-meta", msg)).To(Succeed())

					var poisonMsg *message.Message
					Eventually(poisonMessages).Should(Receive(&poisonMsg))

					By("Verifying payload is preserved")
					Expect(string(poisonMsg.Payload)).To(Equal("payload"))

					By("Verifying poison metadata keys are set")
					Expect(poisonMsg.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal("dlq-meta"))
					Expect(poisonMsg.Metadata.Get(middleware.PoisonedHandlerKey)).To(Equal("dlq-meta"))
					Expect(poisonMsg.Metadata.Get(middleware.ReasonForPoisonedKey)).To(ContainSubstring("fatal error"))
				})
			})

			Context("poison handler delivery", func() {
				It("should deliver poison messages to the registered poison handler", func() {
					poisonSeen := make(chan *message.Message, 1)

					retryConfig.MaxRetries = 1

					util.AddConsumerHandler(router, "poison-hdl", pubSub, func(_ *message.Message) error {
						return errors.New("fatal error")
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
						PoisonHandler: func(msg *message.Message) error {
							poisonSeen <- msg
							return nil
						},
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("poisoned"))
					Expect(pubSub.Publish("poison-hdl", msg)).To(Succeed())

					var received *message.Message
					Eventually(poisonSeen).Should(Receive(&received))
					Expect(string(received.Payload)).To(Equal("poisoned"))
				})
			})

			Context("retry without DLQ", func() {
				It("should retry and let errors propagate when no publisher is set", func() {
					var attempts atomic.Int32

					retryConfig.MaxRetries = 2

					util.AddConsumerHandler(router, "no-dlq", pubSub, func(_ *message.Message) error {
						attempts.Add(1)
						return errors.New("transient failure")
					}, util.ConsumerConfig{
						Retry: retryConfig,
						// No DeadLetterPublisher — errors propagate to subscriber for broker redelivery.
					})

					startRouter()

					msg := message.NewMessage(watermill.NewUUID(), []byte("retry-no-dlq"))
					Expect(pubSub.Publish("no-dlq", msg)).To(Succeed())

					// Initial cycle: 3 attempts (initial + 2 retries), then error propagates
					// and broker redelivers causing more cycles. Verify at least one full
					// retry cycle completed (>= 3 attempts).
					Eventually(attempts.Load).Should(BeNumerically(">=", 3))
				})
			})

			Context("recovered panic with retry", func() {
				It("should retry a panicking handler and deliver to DLQ", func() {
					var attempts atomic.Int32
					poisonMessages, subErr := pubSub.Subscribe(ctx, util.PoisonTopic("panic-retry"))
					Expect(subErr).NotTo(HaveOccurred())

					retryConfig.MaxRetries = 1

					util.AddConsumerHandler(router, "panic-retry", pubSub, func(_ *message.Message) error {
						attempts.Add(1)
						panic("deliberate panic")
					}, util.ConsumerConfig{
						Retry:               retryConfig,
						DeadLetterPublisher: pubSub,
					})

					startRouter()

					orig := message.NewMessage(watermill.NewUUID(), []byte("doomed"))
					Expect(pubSub.Publish("panic-retry", orig)).To(Succeed())

					var poisonMsg *message.Message
					Eventually(poisonMessages).Should(Receive(&poisonMsg))
					Expect(string(poisonMsg.Payload)).To(Equal("doomed"))
					Expect(poisonMsg.Metadata.Get(middleware.ReasonForPoisonedKey)).To(ContainSubstring("panic"))

					Expect(attempts.Load()).To(Equal(int32(2))) // initial + 1 retry
				})
			})
		})
	})
})
