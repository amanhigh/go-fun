package components_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterChip", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.FilterChipProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.FilterChipProps{
			TextExpr: "activeFilter.label",
			Tone:     components.ToneBlue,
			ShowExpr: "",
			Class:    "",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Tone, class, and text rendering", func() {
		BeforeEach(func() {
			props.Class = "rounded-full px-3"
			err := components.FilterChip(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders tone class, custom class, and x-text binding", func() {
			Expect(html).To(ContainSubstring("filter-chip-blue"))
			Expect(html).To(ContainSubstring("rounded-full"))
			Expect(html).To(ContainSubstring("px-3"))
			Expect(html).To(ContainSubstring(`x-text="activeFilter.label"`))
		})
	})

	Context("ShowExpr present", func() {
		BeforeEach(func() {
			props.ShowExpr = "activeFilter.visible"
			err := components.FilterChip(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("binds the x-show expression", func() {
			Expect(html).To(ContainSubstring(`x-show="activeFilter.visible"`))
		})
	})

	Context("ShowExpr absent", func() {
		BeforeEach(func() {
			err := components.FilterChip(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("does not emit x-show", func() {
			Expect(html).NotTo(ContainSubstring("x-show"))
		})
	})
})
