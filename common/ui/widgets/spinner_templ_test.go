package widgets_test

import (
	"context"
	"strings"

	widgets "github.com/amanhigh/go-fun/common/ui/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spinner Template", func() {
	var (
		ctx  context.Context
		html string
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("renders the default spinner", func() {
		var render strings.Builder
		err := widgets.Spinner().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		Expect(html).To(ContainSubstring("flex"))
		Expect(html).To(ContainSubstring("items-center"))
		Expect(html).To(ContainSubstring("justify-center"))
		Expect(html).To(ContainSubstring("h-6 w-6"))
		Expect(html).To(ContainSubstring("animate-spin text-primary"))
	})

	It("renders custom size and class", func() {
		var render strings.Builder
		err := widgets.SpinnerWithProps(widgets.SpinnerProps{Size: widgets.SpinnerSizeSm, Class: "mt-2"}).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		Expect(html).To(ContainSubstring("mt-2"))
		Expect(html).To(ContainSubstring("h-4 w-4"))
	})
})
