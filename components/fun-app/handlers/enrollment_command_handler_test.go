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

func TestEnrollmentCommandHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EnrollmentCommandHandler Suite")
}

var _ = Describe("EnrollmentCommandHandler", func() {
	var (
		managerMock *managermocks.EnrollmentManagerInterface
		handler     *handlers.EnrollmentCommandHandlerImpl
	)

	BeforeEach(func() {
		managerMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		handler = handlers.NewEnrollmentCommandHandler(managerMock)
	})

	Context("EnrollCmd", func() {
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
		})

		Context("without metadata, uses default stamps", func() {
			BeforeEach(func() {
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.EnrollCmd(msg)
			})

			It("sets correlation to EnrollmentID and causation to message UUID", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(common.CorrelationFrom(capturedCtx)).To(Equal(cmd.EnrollmentID))
				Expect(common.CausationFrom(capturedCtx)).To(Equal(msg.UUID))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
				Expect(capturedCmd.PersonID).To(Equal(cmd.PersonID))
				Expect(capturedCmd.Grade).To(Equal(cmd.Grade))
			})
		})

		Context("with metadata, overrides correlation and causation", func() {
			BeforeEach(func() {
				msg.Metadata = message.Metadata{
					common.MetadataCorrelationIDKey: "meta-corr",
					common.MetadataCausationIDKey:   "meta-cause",
				}
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.EnrollCmd(msg)
			})

			It("uses Metadata[CorrelationID] and Metadata[CausationID] instead of defaults", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(common.CorrelationFrom(capturedCtx)).To(Equal("meta-corr"))
				Expect(common.CausationFrom(capturedCtx)).To(Equal("meta-cause"))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
			})
		})

		Context("nil message context falls back to background ctx", func() {
			BeforeEach(func() {
				msg.SetContext(context.TODO())
				managerMock.EXPECT().EnrollCmd(mock.Anything, mock.Anything).
					Run(func(c context.Context, in fun.EnrollCmdV1) {
						capturedCtx = c
						capturedCmd = in
					}).
					Return(nil)

				resultErr = handler.EnrollCmd(msg)
			})

			It("uses background ctx and default stamps when msg.Context() is nil", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(capturedCtx).ToNot(BeNil())
				Expect(common.CorrelationFrom(capturedCtx)).To(Equal(cmd.EnrollmentID))
				Expect(common.CausationFrom(capturedCtx)).To(Equal(msg.UUID))
				Expect(capturedCmd.EnrollmentID).To(Equal(cmd.EnrollmentID))
			})
		})

		Context("invalid JSON payload", func() {
			BeforeEach(func() {
				msg.Payload = []byte("not-json")
				resultErr = handler.EnrollCmd(msg)
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
				resultErr = handler.EnrollCmd(msg)
			})

			It("returns the same HttpError returned by manager.EnrollCmd", func() {
				Expect(resultErr).To(HaveOccurred())
				Expect(resultErr).To(Equal(expectedErr))
			})
		})
	})
})
