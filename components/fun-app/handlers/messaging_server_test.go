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
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(
				mock.Anything,
				mock.MatchedBy(func(evt fun.EnrollmentCancelledEvtV1) bool {
					return evt.EnrollmentID == cmd.EnrollmentID &&
						evt.PersonID == cmd.PersonID &&
						evt.Reason == fun.EnrollmentCancellationReasonSeatAllocationFailed
				}),
			).Run(func(_ context.Context, _ fun.EnrollmentCancelledEvtV1) {
				close(cancelled)
			}).Return(nil).Once()

			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, message.NewMessage("allocate-1", payload))).To(Succeed())
			Eventually(cancelled, 10*time.Second).Should(BeClosed())
		})

		It("retries allocation twice before compensating once", func() {
			seatMock.AssertNumberOfCalls(GinkgoT(), "AllocateSeat", 3)
			enrollmentMock.AssertNumberOfCalls(GinkgoT(), "CancelEnrollmentAndPublish", 1)
		})
	})

	Context("malformed allocation", func() {
		var poisonMessages <-chan *message.Message

		BeforeEach(func() {
			var err error
			poisonMessages, err = channel.Subscribe(routerCtx, util.PoisonTopic(fun.TopicAllocateSeatCmd))
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, message.NewMessage("allocate-malformed", []byte("not-json")))).To(Succeed())
			var poisonMessage *message.Message
			Eventually(poisonMessages).Should(Receive(&poisonMessage))
		})

		It("acknowledges permanently malformed allocation without retry or compensation", func() {
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("terminal enrollment confirmation DLQ", func() {
		var poisonMessage *message.Message

		BeforeEach(func() {
			poisonMessages, err := channel.Subscribe(routerCtx, util.PoisonTopic(fun.TopicEnrollmentConfirmedEvt))
			Expect(err).ToNot(HaveOccurred())
			payload := []byte("not-json")
			Expect(channel.Publish(fun.TopicEnrollmentConfirmedEvt, message.NewMessage("confirmation-malformed", payload))).To(Succeed())
			Eventually(poisonMessages).Should(Receive(&poisonMessage))
		})

		It("publishes one terminal poison record with source metadata", func() {
			Expect(string(poisonMessage.Payload)).To(Equal("not-json"))
			Expect(poisonMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicEnrollmentConfirmedEvt))
			enrollmentMock.AssertNotCalled(GinkgoT(), "OnEnrollmentConfirmedEvt", mock.Anything, mock.Anything)
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
		})
	})
})
