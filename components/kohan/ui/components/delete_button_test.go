package components_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteButton", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.DeleteButtonProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.DeleteButtonProps{
			Label:        "Delete entry",
			DisabledExpr: "!canDelete",
			OnClickExpr:  "deleteItem()",
			Class:        "",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Delete button contract", func() {
		BeforeEach(func() {
			err := components.DeleteButton(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders accessibility labels, click, disabled, and delete-button behavior", func() {
			Expect(html).To(ContainSubstring(`aria-label="Delete entry"`))
			Expect(html).To(ContainSubstring(`title="Delete entry"`))
			Expect(html).To(ContainSubstring(`x-on:click="deleteItem()"`))
			Expect(html).To(ContainSubstring(`x-bind:disabled="!canDelete"`))
			Expect(html).To(ContainSubstring(`opacity-70`))
			Expect(html).To(ContainSubstring(`data-lucide="icon"`))
		})
	})

	Context("Custom class", func() {
		BeforeEach(func() {
			props.Class = "ml-2"
			err := components.DeleteButton(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("appends the caller-provided class", func() {
			Expect(html).To(ContainSubstring("ml-2"))
		})
	})
})
