package handlers_test

import (
	"encoding/json"

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
		enrollmentMock *managermocks.EnrollmentManagerInterface
		seatHandler    *handlers.SeatMessageHandlerImpl
		resultErr      error
	)

	BeforeEach(func() {
		enrollmentMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		seatHandler = handlers.NewSeatMessageHandler(
			managermocks.NewSeatManagerInterface(GinkgoT()),
			enrollmentMock,
		)
	})

	Context("HandlePoisonedAllocateSeatCmd with malformed payload", func() {
		BeforeEach(func() {
			msg := message.NewMessage("poison-msg", []byte("not-json"))
			msg.Metadata.Set(middleware.PoisonedTopicKey, fun.TopicAllocateSeatCmd)
			msg.Metadata.Set(middleware.PoisonedHandlerKey, fun.TopicAllocateSeatCmd)
			msg.Metadata.Set(middleware.ReasonForPoisonedKey, "invalid payload")
			resultErr = seatHandler.HandlePoisonedAllocateSeatCmd(msg)
		})

		It("returns an unmarshal error without compensating", func() {
			Expect(resultErr).To(HaveOccurred())
			Expect(resultErr.Error()).To(ContainSubstring("unmarshal poisoned allocate seat cmd v1"))
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("HandlePoisonedAllocateSeatCmd with valid payload", func() {
		BeforeEach(func() {
			cmd := fun.AllocateSeatCmdV1{EnrollmentID: "enr-1", PersonID: "person-1"}
			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			msg := message.NewMessage("poison-msg", payload)
			expectedErr := common.NewHttpError("cancel failed", 500)
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(mock.Anything, mock.MatchedBy(func(evt fun.EnrollmentCancelledEvtV1) bool {
				return evt.EnrollmentID == cmd.EnrollmentID && evt.PersonID == cmd.PersonID
			})).Return(expectedErr)
			resultErr = seatHandler.HandlePoisonedAllocateSeatCmd(msg)
		})

		It("propagates the compensation manager error", func() {
			Expect(resultErr).To(HaveOccurred())
			Expect(resultErr.Error()).To(Equal("cancel failed"))
		})
	})
})
