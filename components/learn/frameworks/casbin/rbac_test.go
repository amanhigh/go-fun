package casbin

import (
	"github.com/casbin/casbin/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rbac", func() {
	var (
		e, err = casbin.NewEnforcer("model.conf", "policy.csv")
	)

	It("should build", func() {
		Expect(e).To(Not(BeNil()))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Alice", func() {
		It("should be allowed to read image", func() {
			ok, err := e.Enforce("alice", "image", "read")
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not be allowed to write image", func() {
			ok, err := e.Enforce("alice", "image", "write")
			Expect(ok).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should give roles for alice", func() {
			roles, err := e.GetRolesForUser("alice")
			Expect(roles).To(Equal([]string{"image-reader"}))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Bob", func() {
		It("should be able to read image", func() {
			ok, err := e.Enforce("bob", "image", "read")
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be able to write image", func() {
			ok, err := e.Enforce("bob", "image", "write")
			Expect(ok).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
