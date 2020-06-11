package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rbac", func() {
	var (
		e, err = casbin.NewEnforcer("model.conf", "policy.csv")
	)

	It("should build", func() {
		fmt.Println(">>ERROR<<", err)

		Expect(e).To(Not(BeNil()))
		Expect(err).To(BeNil())
	})

	Context("Alice", func() {
		It("should be allowed to read image", func() {
			ok, err := e.Enforce("alice", "image", "read")
			Expect(ok).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("should not be allowed to write image", func() {
			ok, err := e.Enforce("alice", "image", "write")
			Expect(ok).To(BeFalse())
			Expect(err).To(BeNil())
		})

		It("should give admin for alice", func() {
			roles, err := e.GetRolesForUser("alice")
			Expect(roles).To(Equal([]string{"image-reader"}))
			Expect(err).To(BeNil())
		})
	})

	Context("Bob", func() {
		It("should be able to write image", func() {
			ok, err := e.Enforce("bob", "image", "write")
			Expect(ok).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("should give admin for alice", func() {
			roles, err := e.GetRolesForUser("bob")
			Expect(roles).To(Equal([]string{"image-writer"}))
			Expect(err).To(BeNil())
		})
	})
})
