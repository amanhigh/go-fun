package components_test

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmptyState", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.EmptyStateProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.EmptyStateProps{
			Message:  "No results found.",
			ShowExpr: "items.length === 0",
			Class:    "mt-4",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Full contract with optional attributes", func() {
		BeforeEach(func() {
			props.Attributes = templ.Attributes{"data-testid": "empty-state"}
			err := components.EmptyState(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders message, x-show, classes, and custom attributes", func() {
			Expect(html).To(ContainSubstring("No results found."))
			Expect(html).To(ContainSubstring(`x-show="items.length === 0"`))
			Expect(html).To(ContainSubstring("border-dashed"))
			Expect(html).To(ContainSubstring("text-muted-foreground"))
			Expect(html).To(ContainSubstring("mt-4"))
			Expect(html).To(ContainSubstring(`data-testid="empty-state"`))
		})
	})
})
