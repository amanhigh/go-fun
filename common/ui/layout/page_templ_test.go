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

// heroCardClass returns the class attribute of the fixed full-width hero card
// element (the card.Card rendered with the w-full max-w-none class).
func heroCardClass(doc *goquery.Document) string {
	return doc.Find(".w-full.max-w-none").First().AttrOr("class", "")
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
		It("renders the standard full-width hero shell", func() {
			Expect(doc.Find("section").First().AttrOr("class", "")).To(ContainSubstring("flex w-full flex-col gap-8"))
			Expect(heroCardClass(doc)).To(ContainSubstring("rounded-[2rem]"))
			Expect(heroCardClass(doc)).To(ContainSubstring("shadow-[0_24px_80px_-48px_rgba(15,23,42,0.85)]"))
			// Fixed standard full-width hero card class remains for the regular render.
			Expect(heroCardClass(doc)).To(ContainSubstring("w-full"))
			Expect(heroCardClass(doc)).To(ContainSubstring("max-w-none"))
		})

		It("renders custom HeroContent", func() {
			customProps := props
			customProps.HeroContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<div id="custom-hero-content">Custom hero</div>`)
				return err
			})
			customHTML, customDoc := renderPage(ctx, customProps)

			// Standard root remains.
			Expect(customDoc.Find("section").First().AttrOr("class", "")).To(ContainSubstring("flex w-full flex-col gap-8"))
			// Custom hero marker is present.
			Expect(customHTML).To(ContainSubstring(`<div id="custom-hero-content">Custom hero</div>`))
			// Breadcrumb and eyebrow remain.
			Expect(customHTML).To(ContainSubstring("🏠 Home"))
			Expect(customHTML).To(ContainSubstring("Kohan Portal"))
			// Default heading is absent from the custom render.
			Expect(customHTML).ToNot(ContainSubstring("Journal Detail"))
		})
	})

	Context("Page root attributes", func() {
		It("renders custom attributes on the root section", func() {
			customAttrs := templ.Attributes{
				"x-data":       "journalDetailPage()",
				"data-test-id": "jrn_1234",
			}
			props.Attributes = customAttrs
			_, attrDoc := renderPage(ctx, props)

			section := attrDoc.Find("section").First()
			val, exists := section.Attr("x-data")
			Expect(exists).To(BeTrue())
			Expect(val).To(Equal("journalDetailPage()"))

			val, exists = section.Attr("data-test-id")
			Expect(exists).To(BeTrue())
			Expect(val).To(Equal("jrn_1234"))
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
			Expect(doc.Find("h1").First().Text()).To(Equal("Journal Detail"))
			Expect(doc.Find("hgroup p").First().Text()).To(Equal("View complete journal entry with all associated data."))
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

			// The hero card must precede the child content in document order.
			ordered := childDoc.Find(".w-full.max-w-none, #content-marker")
			Expect(ordered.Length()).To(Equal(2))
			firstClass, _ := ordered.First().Attr("class")
			Expect(firstClass).To(ContainSubstring("w-full"))
			lastID, _ := ordered.Last().Attr("id")
			Expect(lastID).To(Equal("content-marker"))
		})
	})
})
