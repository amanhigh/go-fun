package components_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubmitButton", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.SubmitButtonProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.SubmitButtonProps{
			Label:        "Save",
			BusyExpr:     "submitter.isBusy()",
			DisabledExpr: "!form.isValid",
			OnClickExpr:  "submitter.submit()",
			Class:        "btn-primary",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Default props contract", func() {
		BeforeEach(func() {
			err := components.SubmitButton(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders label/busy text, click, disabled, opacity, custom class, and button type", func() {
			Expect(html).To(ContainSubstring(`submitter.isBusy()`))
			Expect(html).To(ContainSubstring(`Submitting...`))
			Expect(html).To(ContainSubstring(`Save`))
			Expect(html).To(ContainSubstring(`x-on:click="submitter.submit()"`))
			Expect(html).To(ContainSubstring(`x-bind:disabled="!form.isValid"`))
			Expect(html).To(ContainSubstring(`opacity-70`))
			Expect(html).To(ContainSubstring("btn-primary"))
			Expect(html).To(ContainSubstring(`type="button"`))
		})
	})
})
