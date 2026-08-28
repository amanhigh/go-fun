package components_test

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SectionHeader", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.SectionHeaderProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.SectionHeaderProps{
			Title:       "My Section",
			Description: "",
			Class:       "",
			Attributes:  nil,
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Description present", func() {
		BeforeEach(func() {
			props.Description = "This is a description"
			props.Attributes = templ.Attributes{"data-testid": "section-header", "id": "sec-1"}
			err := components.SectionHeader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders description, title, and forwards custom attributes", func() {
			Expect(html).To(ContainSubstring("<p "))
			Expect(html).To(ContainSubstring("This is a description"))
			Expect(html).To(ContainSubstring("text-muted-foreground"))
			Expect(html).To(ContainSubstring("My Section"))
			Expect(html).To(ContainSubstring(`data-testid="section-header"`))
			Expect(html).To(ContainSubstring(`id="sec-1"`))
			Expect(html).To(ContainSubstring(`class="space-y-0.5 `))
			Expect(html).To(ContainSubstring(`text-base sm:text-lg`))
		})
	})

	Context("Description omitted", func() {
		BeforeEach(func() {
			err := components.SectionHeader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("omits the paragraph and still renders the title", func() {
			Expect(html).NotTo(ContainSubstring("<p "))
			Expect(html).To(ContainSubstring("My Section"))
		})
	})

	Context("Meta present", func() {
		BeforeEach(func() {
			props.Meta = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<span class="meta-badge">Updated 2h ago</span>`)
				return err
			})
			err := components.SectionHeader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders metadata alongside the title in a space-between row", func() {
			Expect(html).To(ContainSubstring("My Section"))
			Expect(html).To(ContainSubstring(`class="flex flex-wrap items-center justify-between gap-x-2 gap-y-0.5"`))
			Expect(html).To(ContainSubstring(`class="meta-badge"`))
			Expect(html).To(ContainSubstring("Updated 2h ago"))
		})
	})
})
