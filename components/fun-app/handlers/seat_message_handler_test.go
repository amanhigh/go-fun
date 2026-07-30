package handlers_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

var _ = Describe("SeatMessageHandler", func() {
	var (
		enrollmentMock  *managermocks.EnrollmentManagerInterface
		seatManagerMock *managermocks.SeatManagerInterface
		seatHandler     *handlers.SeatMessageHandlerImpl
		resultErr       error
	)

	BeforeEach(func() {
		enrollmentMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		seatManagerMock = managermocks.NewSeatManagerInterface(GinkgoT())
		seatHandler = handlers.NewSeatMessageHandler(
			seatManagerMock,
			enrollmentMock,
		)
	})

	Context("HandleDeadLetteredAllocateSeatCmd with malformed payload", func() {
		BeforeEach(func() {
			msg := message.NewMessage("poison-msg", []byte("not-json"))
			msg.SetContext(common.WithMetadata(context.Background(), common.NewRootMetadata(msg.UUID, "corr-1")))
			msg.Metadata.Set(middleware.PoisonedTopicKey, fun.TopicAllocateSeatCmd)
			msg.Metadata.Set(middleware.PoisonedHandlerKey, fun.TopicAllocateSeatCmd)
			msg.Metadata.Set(middleware.ReasonForPoisonedKey, "invalid payload")
			resultErr = seatHandler.HandleDeadLetteredAllocateSeatCmd(msg)
		})

		It("returns an unmarshal error without compensating", func() {
			Expect(resultErr).To(HaveOccurred())
			Expect(resultErr.Error()).To(ContainSubstring("unmarshal dead-lettered allocate seat cmd v1"))
			seatManagerMock.AssertNotCalled(GinkgoT(), "PublishSeatAllocationFailed", mock.Anything, mock.Anything, mock.Anything)
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("HandleDeadLetteredAllocateSeatCmd with valid payload", func() {
		BeforeEach(func() {
			cmd := fun.AllocateSeatCmdV1{EnrollmentID: "enr-1", PersonID: "person-1"}
			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			msg := message.NewMessage("dead-letter-msg", payload)
			msg.SetContext(common.WithMetadata(context.Background(), common.NewRootMetadata(msg.UUID, "corr-1")))
			msg.Metadata.Set(middleware.ReasonForPoisonedKey, "capacity service unavailable")
			seatManagerMock.EXPECT().PublishSeatAllocationFailed(mock.Anything, fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID}, "capacity service unavailable").Return(nil)
			resultErr = seatHandler.HandleDeadLetteredAllocateSeatCmd(msg)
		})

		It("publishes an allocation failure through the seat manager", func() {
			Expect(resultErr).ToNot(HaveOccurred())
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("HandleSeatAllocationFailedEvt with malformed payload", func() {
		BeforeEach(func() {
			msg := message.NewMessage("failed-msg", []byte("not-json"))
			resultErr = seatHandler.HandleSeatAllocationFailedEvt(msg)
		})

		It("returns an unmarshal error without compensating", func() {
			Expect(resultErr).To(HaveOccurred())
			Expect(resultErr.Error()).To(ContainSubstring("unmarshal seat allocation failed evt"))
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("HandleSeatAllocationFailedEvt with valid payload", func() {
		var (
			evt       fun.SeatAllocationFailedEvtV1
			startedAt time.Time
		)

		BeforeEach(func() {
			evt = fun.SeatAllocationFailedEvtV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Reason:       "capacity unavailable",
				FailedAt:     time.Now().UTC(),
			}
			payload, err := json.Marshal(evt)
			Expect(err).ToNot(HaveOccurred())
			msg := message.NewMessage("failed-msg", payload)
			msg.SetContext(common.WithMetadata(context.Background(), common.Metadata{
				MessageID:     msg.UUID,
				CorrelationID: "corr-1",
				CausationID:   "cause-1",
			}))
			startedAt = time.Now().UTC()
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(
				mock.MatchedBy(func(ctx context.Context) bool {
					metadata := common.MetadataFromContext(ctx)
					return metadata.CorrelationID == "corr-1" && metadata.CausationID == "cause-1"
				}),
				mock.MatchedBy(func(cancelled fun.EnrollmentCancelledEvtV1) bool {
					return cancelled.EnrollmentID == evt.EnrollmentID &&
						cancelled.PersonID == evt.PersonID &&
						cancelled.Reason == evt.Reason &&
						!cancelled.CancelledAt.Before(startedAt) &&
						!cancelled.CancelledAt.After(time.Now().UTC())
				}),
			).Return(nil)
			resultErr = seatHandler.HandleSeatAllocationFailedEvt(msg)
		})

		It("cancels the enrollment through the enrollment manager", func() {
			Expect(resultErr).ToNot(HaveOccurred())
		})
	})
})
