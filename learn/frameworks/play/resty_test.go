package play_test

import (
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resty", func() {

	var (
		client = resty.New()
	)

	It("should build", func() {
		Expect(client).To(Not(BeNil()))
	})
})
