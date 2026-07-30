package handlers_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

func TestEnrollmentMessageHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EnrollmentMessageHandler Suite")
}

var _ = Describe("EnrollmentMessageHandler", func() {
	var (
		managerMock *managermocks.EnrollmentManagerInterface
		handler     *handlers.EnrollmentMessageHandlerImpl
	)

	BeforeEach(func() {
		managerMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		handler = handlers.NewEnrollmentMessageHandler(managerMock)
	})

	Context("HandleEnrollCmd", func() {
		var (
			cmd         fun.EnrollCmdV1
			msg         *message.Message
			resultErr   error
			expectedErr common.HttpError
			capturedCtx context.Context
			capturedCmd fun.EnrollCmdV1
		)

		BeforeEach(func() {
			cmd = fun.EnrollCmdV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Grade:        3,
				RequestedAt:  time.Now().UTC(),
			}
			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			msg = message.NewMessage("msg-uuid", payload)
			msg.SetContext(common.WithMetadata(context.Background(), common.NewRootMetadata(msg.UUID, "meta-corr")))
		})

		Context("with required metadata", func() {
			BeforeEach(func() {
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.HandleEnrollCmd(msg)
			})

			It("sets correlation and causation from the consumed message", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(common.MetadataFromContext(capturedCtx)).To(Equal(common.NewRootMetadata(msg.UUID, "meta-corr")))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
				Expect(capturedCmd.PersonID).To(Equal(cmd.PersonID))
				Expect(capturedCmd.Grade).To(Equal(cmd.Grade))
			})
		})

		Context("with metadata, overrides correlation and causation", func() {
			BeforeEach(func() {
				msg.SetContext(common.WithMetadata(context.Background(), common.Metadata{
					MessageID: msg.UUID, CorrelationID: "meta-corr", CausationID: "meta-cause",
				}))
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.HandleEnrollCmd(msg)
			})

			It("passes typed metadata from the message context", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(common.MetadataFromContext(capturedCtx)).To(Equal(common.Metadata{MessageID: msg.UUID, CorrelationID: "meta-corr", CausationID: "meta-cause"}))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
			})
		})

		Context("message context with root metadata", func() {
			BeforeEach(func() {
				msg.SetContext(common.WithMetadata(context.TODO(), common.NewRootMetadata(msg.UUID, "meta-corr")))
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.HandleEnrollCmd(msg)
			})

			It("passes the message context metadata to the manager", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(capturedCtx).ToNot(BeNil())
				Expect(common.MetadataFromContext(capturedCtx)).To(Equal(common.NewRootMetadata(msg.UUID, "meta-corr")))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
			})
		})

		Context("invalid JSON payload", func() {
			BeforeEach(func() {
				msg.Payload = []byte("not-json")
				resultErr = handler.HandleEnrollCmd(msg)
			})

			It("returns unmarshal error and does not call manager.EnrollCmd", func() {
				Expect(resultErr).To(HaveOccurred())
				Expect(resultErr.Error()).To(ContainSubstring("unmarshal enroll cmd"))
				managerMock.AssertNotCalled(GinkgoT(), "EnrollCmd", mock.Anything, mock.Anything)
			})
		})

		Context("manager error propagation", func() {
			BeforeEach(func() {
				expectedErr = common.NewHttpError("seat-fail", 500)
				ctxMatcher := mock.MatchedBy(func(_ context.Context) bool { return true })
				cmdMatcher := mock.MatchedBy(func(in fun.EnrollCmdV1) bool { return in.EnrollmentID == cmd.EnrollmentID })
				managerMock.EXPECT().EnrollCmd(ctxMatcher, cmdMatcher).Return(expectedErr)
				resultErr = handler.HandleEnrollCmd(msg)
			})

			It("returns the same HttpError returned by manager.EnrollCmd", func() {
				Expect(resultErr).To(HaveOccurred())
				Expect(resultErr).To(Equal(expectedErr))
			})
		})
	})
})
