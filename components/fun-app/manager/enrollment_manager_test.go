package manager_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	daomocks "github.com/amanhigh/go-fun/components/fun-app/dao/mocks"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	pubmocks "github.com/amanhigh/go-fun/components/fun-app/publisher/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

func TestEnrollmentManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EnrollmentManager Suite")
}

var _ = Describe("EnrollmentManager", func() {
	var (
		personMgr *managermocks.PersonManagerInterface
		dao       *daomocks.EnrollmentDaoInterface
		publisher *pubmocks.EnrollmentPublisher
		seatMgr   *managermocks.SeatManagerInterface
		em        *manager.EnrollmentManager
	)

	BeforeEach(func() {
		personMgr = managermocks.NewPersonManagerInterface(GinkgoT())
		dao = daomocks.NewEnrollmentDaoInterface(GinkgoT())
		publisher = pubmocks.NewEnrollmentPublisher(GinkgoT())
		seatMgr = managermocks.NewSeatManagerInterface(GinkgoT())
		em = manager.NewEnrollmentManager(personMgr, dao, publisher, seatMgr)
	})

	Context("EnrollCmd", func() {
		var (
			ctx context.Context
			cmd fun.EnrollCmdV1
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctx = common.WithMetadata(ctx, common.NewChildMetadata("child-456", common.NewRootMetadata("cause-456", "corr-123")))
			cmd = fun.EnrollCmdV1{EnrollmentID: "enr-101", PersonID: "person-1", Grade: 5, RequestedAt: time.Now().UTC()}
		})

		It("delegates to SeatManager with ctx unchanged", func() {
			ctxMatcher := mock.MatchedBy(func(c context.Context) bool {
				metadata := common.MetadataFromContext(c)
				return metadata.CorrelationID == "corr-123" && metadata.CausationID == "cause-456"
			})
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool {
				return e.ID == cmd.EnrollmentID && e.PersonID == cmd.PersonID && e.Grade == cmd.Grade
			})
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(nil)

			err := em.EnrollCmd(ctx, cmd)
			Expect(err).ToNot(HaveOccurred())
		})

		It("works with nil context (uses non-nil ctx, no stamps)", func() {
			ctxMatcher := mock.MatchedBy(func(c context.Context) bool {
				return c != nil && common.MetadataFromContext(c) == (common.Metadata{})
			})
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool {
				return e.ID == cmd.EnrollmentID && e.PersonID == cmd.PersonID && e.Grade == cmd.Grade
			})
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(nil)

			var nilCtx context.Context
			err := em.EnrollCmd(nilCtx, cmd)
			Expect(err).ToNot(HaveOccurred())
		})

		It("propagates SeatManager error", func() {
			ctx = common.WithMetadata(context.Background(), common.NewChildMetadata("child-err", common.NewRootMetadata("cause-err", "corr-err")))
			expected := common.NewHttpError("seat-fail", 500)

			ctxMatcher := mock.MatchedBy(func(_ context.Context) bool { return true })
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool { return e.ID == cmd.EnrollmentID })
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(expected)

			err := em.EnrollCmd(ctx, cmd)
			Expect(err).To(Equal(expected))
		})
	})

	Context("CancelEnrollmentAndPublish", func() {
		It("persists and publishes only on the first cancellation", func() {
			evt := fun.EnrollmentCancelledEvtV1{
				EnrollmentID: "enr-101",
				PersonID:     "person-1",
				Reason:       "seat allocation failed",
			}
			findCount := 0

			dao.EXPECT().UseOrCreateTx(mock.Anything, mock.Anything).RunAndReturn(
				func(ctx context.Context, run util.DbRun, _ ...bool) common.HttpError {
					return run(ctx)
				},
			).Twice()
			dao.EXPECT().FindById(mock.Anything, evt.EnrollmentID, mock.Anything).Run(
				func(_ context.Context, _ any, entity any) {
					findCount++
					status := fun.EnrollmentStatusSeatAllocationInitiated
					if findCount == 2 {
						status = fun.EnrollmentStatusCancelled
					}
					enrollment, ok := entity.(*fun.Enrollment)
					Expect(ok).To(BeTrue())
					*enrollment = fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Status: status}
				},
			).Return(nil).Twice()
			dao.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()
			publisher.EXPECT().EnrollmentCancelledEvt(mock.Anything, mock.MatchedBy(func(enrollment fun.Enrollment) bool {
				return enrollment.ID == evt.EnrollmentID && enrollment.Status == fun.EnrollmentStatusCancelled
			}), evt.Reason).Return(nil).Once()

			Expect(em.CancelEnrollmentAndPublish(context.Background(), evt)).ToNot(HaveOccurred())
			Expect(em.CancelEnrollmentAndPublish(context.Background(), evt)).ToNot(HaveOccurred())
			Expect(findCount).To(Equal(2))
		})
	})
})
