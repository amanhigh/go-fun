package ui_test

import (
	"context"
	"strings"

	ui "github.com/amanhigh/go-fun/common/ui"
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
		err := ui.Spinner().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		Expect(html).To(ContainSubstring("flex items-center justify-center"))
		Expect(html).To(ContainSubstring("h-6 w-6"))
		Expect(html).To(ContainSubstring("animate-spin text-primary"))
	})

	It("renders custom size and class", func() {
		var render strings.Builder
		err := ui.SpinnerWithProps(ui.SpinnerProps{Size: ui.SpinnerSizeSm, Class: "mt-2"}).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		Expect(html).To(ContainSubstring("mt-2"))
		Expect(html).To(ContainSubstring("h-4 w-4"))
	})
})
