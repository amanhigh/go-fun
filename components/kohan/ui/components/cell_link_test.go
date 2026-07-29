package components_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CellLink", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.CellLinkProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.CellLinkProps{
			HrefExpr:    "",
			OnClickExpr: "handleClick()",
			TextExpr:    "item.name",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("HrefExpr branch", func() {
		BeforeEach(func() {
			props.HrefExpr = "item.url"
			err := components.CellLink(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders an anchor with href and text bindings", func() {
			Expect(html).To(ContainSubstring("<a "))
			Expect(html).To(ContainSubstring("</a>"))
			Expect(html).To(ContainSubstring(`x-bind:href="item.url"`))
			Expect(html).To(ContainSubstring(`x-text="item.name"`))
		})
	})

	Context("Button branch", func() {
		BeforeEach(func() {
			err := components.CellLink(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders a button with click and text bindings", func() {
			Expect(html).NotTo(ContainSubstring("<a "))
			Expect(html).To(ContainSubstring("<button"))
			Expect(html).To(ContainSubstring(`x-on:click="handleClick()"`))
			Expect(html).To(ContainSubstring(`x-text="item.name"`))
		})
	})
})
