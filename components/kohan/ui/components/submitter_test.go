package components_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Submitter", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Busy spinner binding", func() {
		BeforeEach(func() {
			err := components.Submitter(components.SubmitterProps{Submitter: "header.submitter"}).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders the supplied busy expression on the spinner root", func() {
			Expect(html).To(ContainSubstring(`x-show="header.submitter.isBusy()"`))
		})
	})
})
