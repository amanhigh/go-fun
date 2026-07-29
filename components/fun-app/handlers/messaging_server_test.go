package handlers_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	"github.com/amanhigh/go-fun/models/fun"
)

var _ = Describe("MessagingServer - Poison Consumer", func() {
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

	Context("when AllocateSeatCmd lands on the derived poison topic", func() {
		var cmd fun.AllocateSeatCmdV1

		BeforeEach(func() {
			cmd = fun.AllocateSeatCmdV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Grade:        3,
				RequestedAt:  time.Now().UTC(),
			}
		})

		It("invokes CancelEnrollmentAndPublish exactly once with matching enrollment and reason", func() {
			called := make(chan struct{})

			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(
				mock.Anything,
				mock.MatchedBy(func(evt fun.EnrollmentCancelledEvtV1) bool {
					return evt.EnrollmentID == cmd.EnrollmentID &&
						evt.PersonID == cmd.PersonID &&
						evt.Reason == fun.EnrollmentCancellationReasonSeatAllocationFailed
				}),
			).Run(func(_ context.Context, _ fun.EnrollmentCancelledEvtV1) {
				close(called)
			}).Return(nil).Once()

			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("poison-1", payload)

			Expect(channel.Publish(util.PoisonTopic(fun.TopicAllocateSeatCmd), msg)).To(Succeed())

			Eventually(called).Should(BeClosed())
		})
	})
})
