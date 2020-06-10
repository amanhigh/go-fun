package casbin

import (
	"github.com/casbin/casbin/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rbac", func() {
	var (
		e, err = casbin.NewEnforcer("model.conf", "policy.csv")
	)

	It("should build", func() {
		Expect(e).To(Not(BeNil()))
		Expect(err).To(BeNil())
	})

	Context("IsAuthorized", func() {
		It("should be allowed", func() {
			ok, err := e.Enforce("alice", "data1", "read")
			Expect(ok).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})

	Context("Roles", func() {
		Context("Get", func() {
			It("should give admin for alice", func() {
				roles, err := e.GetRolesForUser("alice")
				Expect(roles).To(Equal([]string{"admin"}))
				Expect(err).To(BeNil())
			})
		})

	})
})
