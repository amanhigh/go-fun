package manager_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
			ctx = common.WithCorrelation(ctx, "corr-123")
			ctx = common.WithCausation(ctx, "cause-456")
			cmd = fun.EnrollCmdV1{EnrollmentID: "enr-101", PersonID: "person-1", Grade: 5, RequestedAt: time.Now().UTC()}
		})

		It("delegates to SeatManager with ctx unchanged", func() {
			ctxMatcher := mock.MatchedBy(func(c context.Context) bool {
				return common.CorrelationFrom(c) == "corr-123" && common.CausationFrom(c) == "cause-456"
			})
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool {
				return e.ID == cmd.EnrollmentID && e.PersonID == cmd.PersonID && e.Grade == cmd.Grade
			})
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(nil)

			err := em.EnrollCmd(ctx, cmd)
			Expect(err).To(BeNil())
		})

		It("works with nil context (uses non-nil ctx, no stamps)", func() {
			ctxMatcher := mock.MatchedBy(func(c context.Context) bool {
				return c != nil && common.CorrelationFrom(c) == "" && common.CausationFrom(c) == ""
			})
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool {
				return e.ID == cmd.EnrollmentID && e.PersonID == cmd.PersonID && e.Grade == cmd.Grade
			})
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(nil)

			var nilCtx context.Context
			err := em.EnrollCmd(nilCtx, cmd)
			Expect(err).To(BeNil())
		})

		It("propagates SeatManager error", func() {
			ctx = common.WithCausation(common.WithCorrelation(context.Background(), "corr-err"), "cause-err")
			expected := common.NewHttpError("seat-fail", 500)

			ctxMatcher := mock.MatchedBy(func(c context.Context) bool { return true })
			enrMatcher := mock.MatchedBy(func(e fun.Enrollment) bool { return e.ID == cmd.EnrollmentID })
			seatMgr.EXPECT().PublishAllocateSeat(ctxMatcher, enrMatcher).Return(expected)

			err := em.EnrollCmd(ctx, cmd)
			Expect(err).To(Equal(expected))
		})
	})
})
