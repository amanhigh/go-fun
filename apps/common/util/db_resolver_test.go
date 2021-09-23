package util_test

import (
	"errors"
	"github.com/amanhigh/go-fun/apps/common/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("DbResolver", func() {
	var (
		interval  = time.Second
		pingTable = "test"
		err       = errors.New("connect failed")
		policy    *util.FallBackPolicy
	)
	BeforeEach(func() {
		policy = util.NewFallBackPolicy(nil, interval, pingTable)
	})

	It("should build", func() {
		Expect(policy).To(Not(BeNil()))
	})

	Context("Default Pool", func() {
		It("should be PRIMARY", func() {
			Expect(policy.GetPool()).To(Equal(util.POOL_PRIMARY))
		})

		Context("On Error", func() {
			BeforeEach(func() {
				policy.ReportError(err)
			})

			It("should be FALLBACK", func() {
				Expect(policy.GetPool()).To(Equal(util.POOL_FALLBACK))
			})
		})
	})
})
