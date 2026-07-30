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
		props  components.SubmitterProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.SubmitterProps{Submitter: "header.submitter"}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Default prefix rendering", func() {
		BeforeEach(func() {
			err := components.Submitter(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders busy, message, error/success class bindings, and message output", func() {
			Expect(html).To(ContainSubstring(`x-show="header.submitter.isBusy()"`))
			Expect(html).To(ContainSubstring(`header.submitter.hasMessage()`))
			Expect(html).To(ContainSubstring(`header.submitter.hasError()`))
			Expect(html).To(ContainSubstring(`statebox-error`))
			Expect(html).To(ContainSubstring(`statebox-success`))
			Expect(html).To(ContainSubstring(`x-text="header.submitter.message"`))
		})
	})

	Context("Prefix propagation", func() {
		BeforeEach(func() {
			props = components.SubmitterProps{Submitter: "form.submitter"}
			err := components.Submitter(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("uses the provided prefix for all bindings", func() {
			Expect(html).To(ContainSubstring(`x-show="form.submitter.isBusy()"`))
			Expect(html).To(ContainSubstring(`x-text="form.submitter.message"`))
		})
	})
})
