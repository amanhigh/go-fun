package layout_test

import (
	"context"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/common/ui/layout"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func renderPage(ctx context.Context, props layout.PageProps) (string, *goquery.Document) {
	var render strings.Builder
	err := layout.Page(props).Render(ctx, &render)
	Expect(err).ToNot(HaveOccurred())

	html := render.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	Expect(err).ToNot(HaveOccurred())

	return html, doc
}

var _ = Describe("Page Template Tests", func() {
	var (
		ctx   context.Context
		html  string
		doc   *goquery.Document
		props layout.PageProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = layout.PageProps{
			CurrentPage: "Journal",
			Eyebrow:     "Kohan Portal",
			Heading:     "Journal Detail",
			Description: "View complete journal entry with all associated data.",
			Tags:        []string{"Public", "Read-only"},
		}

		html, doc = renderPage(ctx, props)
	})

	Context("Page shell", func() {
		It("uses a wider desktop page shell", func() {
			Expect(doc.Find("section").First().AttrOr("class", "")).To(ContainSubstring("flex w-full flex-col gap-8"))
			Expect(html).To(ContainSubstring("rounded-[2rem]"))
			Expect(html).To(ContainSubstring("shadow-[0_24px_80px_-48px_rgba(15,23,42,0.85)]"))
		})

		It("keeps the hero card compact and left aligned", func() {
			Expect(html).To(ContainSubstring("justify-start"))
			Expect(html).To(ContainSubstring("xl:w-1/2"))
		})
	})

	Context("PageBreadcrumb", func() {
		It("uses breadcrumb reuse with current-page semantics", func() {
			Expect(html).To(ContainSubstring(`href="/"`))
			Expect(html).To(MatchRegexp(`aria-current="page"[^>]*>\s*Journal\s*<`))
		})
	})

	Context("PageMeta", func() {
		It("renders heading and description when provided", func() {
			Expect(html).To(ContainSubstring("Kohan Portal"))
			Expect(doc.Find("h1").First().Text()).To(Equal("Journal Detail"))
			Expect(doc.Find("h1").First().AttrOr("class", "")).To(ContainSubstring("max-w-4xl"))
			Expect(doc.Find("hgroup p").First().Text()).To(Equal("View complete journal entry with all associated data."))
			Expect(doc.Find("hgroup p").First().AttrOr("class", "")).To(ContainSubstring("max-w-3xl"))
		})

		It("omits optional eyebrow and description when empty", func() {
			emptyHTML, emptyDoc := renderPage(ctx, layout.PageProps{CurrentPage: "Home", Heading: "Home"})

			Expect(emptyHTML).ToNot(ContainSubstring("tracking-[0.32em]"))
			Expect(emptyDoc.Find("hgroup p").Length()).To(BeZero())
		})
	})

	Context("PageTagList", func() {
		It("renders all tags when provided", func() {
			tags := doc.Find("div.flex.flex-wrap.gap-2")

			Expect(tags.Length()).To(Equal(1))
			Expect(tags.Text()).To(ContainSubstring("Public"))
			Expect(tags.Text()).To(ContainSubstring("Read-only"))
		})

		It("omits the tag list wrapper when tags are empty", func() {
			_, emptyDoc := renderPage(ctx, layout.PageProps{CurrentPage: "Home", Heading: "Home"})

			Expect(emptyDoc.Find("div.flex.flex-wrap.gap-2").Length()).To(BeZero())
		})
	})

	Context("Page children", func() {
		It("renders child content after the hero card", func() {
			var childRender strings.Builder
			childCtx := templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<section id="content-marker">Hello child</section>`)
				return err
			}))

			err := layout.Page(props).Render(childCtx, &childRender)
			Expect(err).ToNot(HaveOccurred())

			childHTML := childRender.String()
			childDoc, err := goquery.NewDocumentFromReader(strings.NewReader(childHTML))
			Expect(err).ToNot(HaveOccurred())
			Expect(childDoc.Find("#content-marker").Length()).To(Equal(1))
			Expect(strings.Index(childHTML, "Journal Detail")).To(BeNumerically("<", strings.Index(childHTML, "content-marker")))
		})
	})
})
