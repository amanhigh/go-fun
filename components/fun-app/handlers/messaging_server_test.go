package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	common "github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

var _ = Describe("MessagingServer learning scenarios", func() {
	var (
		channel        *gochannel.GoChannel
		ms             *handlers.MessagingServer
		enrollmentMock *managermocks.EnrollmentManagerInterface
		seatMock       *managermocks.SeatManagerInterface
		logger         watermill.LoggerAdapter
		routerCtx      context.Context
		routerCancel   context.CancelFunc
	)

	BeforeEach(func() {
		logger = watermill.NewStdLogger(false, false)
		channel = gochannel.NewGoChannel(gochannel.Config{}, logger)
		enrollmentMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		seatMock = managermocks.NewSeatManagerInterface(GinkgoT())

		enrollmentHandler := handlers.NewEnrollmentMessageHandler(enrollmentMock)
		seatHandler := handlers.NewSeatMessageHandler(seatMock, enrollmentMock)

		var err error
		ms, err = handlers.NewMessagingServer(logger, channel, channel, enrollmentHandler, seatHandler)
		Expect(err).ToNot(HaveOccurred())

		routerCtx, routerCancel = context.WithCancel(context.Background())
		go func() { _ = ms.Router().Run(routerCtx) }()
		<-ms.Router().Running()
	})

	AfterEach(func() {
		routerCancel()
		_ = ms.Router().Close()
		_ = channel.Close()
	})

	Context("allocation compensation", func() {
		var (
			cmd       fun.AllocateSeatCmdV1
			cancelled chan struct{}
		)

		BeforeEach(func() {
			cmd = fun.AllocateSeatCmdV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Grade:        3,
				RequestedAt:  time.Now().UTC(),
			}
			cancelled = make(chan struct{})
			seatMock.EXPECT().AllocateSeat(mock.Anything, cmd).
				Return(common.NewHttpError("seat service unavailable", http.StatusInternalServerError)).Times(3)
			seatMock.EXPECT().PublishSeatAllocationFailed(
				mock.Anything,
				fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID},
				mock.MatchedBy(func(reason string) bool {
					return reason != ""
				}),
			).Run(func(_ context.Context, enrollment fun.Enrollment, reason string) {
				evtPayload, err := json.Marshal(fun.SeatAllocationFailedEvtV1{
					EnrollmentID: enrollment.ID,
					PersonID:     enrollment.PersonID,
					Reason:       reason,
					FailedAt:     time.Now().UTC(),
				})
				if err != nil {
					return
				}
				Expect(channel.Publish(fun.TopicSeatAllocationFailedEvt, message.NewMessage("allocation-failed-1", evtPayload))).To(Succeed())
			}).Return(nil).Once()
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(mock.Anything, mock.MatchedBy(func(evt fun.EnrollmentCancelledEvtV1) bool {
				return evt.EnrollmentID == cmd.EnrollmentID && evt.PersonID == cmd.PersonID
			})).Run(func(context.Context, fun.EnrollmentCancelledEvtV1) {
				close(cancelled)
			}).Return(nil).Once()

			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, message.NewMessage("allocate-1", payload))).To(Succeed())
			Eventually(cancelled, 10*time.Second).Should(BeClosed())
		})

		It("retries allocation twice before publishing one failure", func() {
			seatMock.AssertNumberOfCalls(GinkgoT(), "AllocateSeat", 3)
			seatMock.AssertNumberOfCalls(GinkgoT(), "PublishSeatAllocationFailed", 1)
		})
	})

	Context("malformed allocation", func() {
		var (
			deadLetterMessages <-chan *message.Message
			quarantineMessages <-chan *message.Message
			deadLetterMessage  *message.Message
			quarantineMessage  *message.Message
		)

		BeforeEach(func() {
			var err error
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicAllocateSeatCmd))
			Expect(err).ToNot(HaveOccurred())
			quarantineMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(util.DeadLetterTopic(fun.TopicAllocateSeatCmd)))
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, message.NewMessage("allocate-malformed", []byte("not-json")))).To(Succeed())
			Eventually(deadLetterMessages).Should(Receive(&deadLetterMessage))
			Eventually(quarantineMessages).Should(Receive(&quarantineMessage))
		})

		It("quarantines malformed allocation without retry or compensation", func() {
			Expect(string(deadLetterMessage.Payload)).To(Equal("not-json"))
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicAllocateSeatCmd))
			Expect(string(quarantineMessage.Payload)).To(Equal("not-json"))
			Expect(quarantineMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(util.DeadLetterTopic(fun.TopicAllocateSeatCmd)))
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
			seatMock.AssertNotCalled(GinkgoT(), "PublishSeatAllocationFailed", mock.Anything, mock.Anything, mock.Anything)
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("allocation failure compensation", func() {
		var (
			deadLetterMessages <-chan *message.Message
			deadLetterMessage  *message.Message
		)

		BeforeEach(func() {
			failed := fun.SeatAllocationFailedEvtV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Reason:       "capacity unavailable",
				FailedAt:     time.Now().UTC(),
			}
			payload, err := json.Marshal(failed)
			Expect(err).ToNot(HaveOccurred())
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(mock.Anything, mock.Anything).
				Return(common.NewHttpError("database unavailable", http.StatusInternalServerError)).Times(3)
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicSeatAllocationFailedEvt))
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicSeatAllocationFailedEvt, message.NewMessage("failed-1", payload))).To(Succeed())
			Eventually(deadLetterMessages, 10*time.Second).Should(Receive(&deadLetterMessage))
		})

		It("bounds cancellation retries and preserves a terminal dead-letter record", func() {
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicSeatAllocationFailedEvt))
			enrollmentMock.AssertNumberOfCalls(GinkgoT(), "CancelEnrollmentAndPublish", 3)
		})
	})

	Context("terminal enrollment confirmation DLQ", func() {
		var (
			deadLetterMessages <-chan *message.Message
			deadLetterMessage  *message.Message
			err                error
		)

		BeforeEach(func() {
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicEnrollmentConfirmedEvt))
			Expect(err).ToNot(HaveOccurred())
			payload := []byte("not-json")
			Expect(channel.Publish(fun.TopicEnrollmentConfirmedEvt, message.NewMessage("confirmation-malformed", payload))).To(Succeed())
			Eventually(deadLetterMessages).Should(Receive(&deadLetterMessage))
		})

		It("publishes one terminal dead-letter record with source metadata", func() {
			Expect(string(deadLetterMessage.Payload)).To(Equal("not-json"))
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicEnrollmentConfirmedEvt))
			enrollmentMock.AssertNotCalled(GinkgoT(), "OnEnrollmentConfirmedEvt", mock.Anything, mock.Anything)
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
		})
	})
})
