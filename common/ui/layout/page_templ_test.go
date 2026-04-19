package layout_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	It("uses a wider desktop page shell", func() {
		Expect(doc.Find("section").First().AttrOr("class", "")).To(ContainSubstring("flex w-full flex-col gap-8"))
		Expect(html).To(ContainSubstring("rounded-[2rem]"))
		Expect(html).To(ContainSubstring("shadow-[0_24px_80px_-48px_rgba(15,23,42,0.85)]"))
	})

	It("keeps the hero card compact and left aligned", func() {
		Expect(html).To(ContainSubstring("justify-start"))
		Expect(html).To(ContainSubstring("xl:w-1/2"))
		Expect(doc.Find("h1").First().AttrOr("class", "")).To(ContainSubstring("max-w-4xl"))
		Expect(doc.Find("hgroup p").First().AttrOr("class", "")).To(ContainSubstring("max-w-3xl"))
	})

	It("uses breadcrumb reuse with current-page semantics", func() {
		Expect(html).To(ContainSubstring(`href="/"`))
		Expect(html).To(MatchRegexp(`aria-current="page"[^>]*>\s*Journal\s*<`))
	})

	It("omits optional hero content when props are empty", func() {
		emptyHTML, emptyDoc := renderPage(ctx, layout.PageProps{CurrentPage: "Home", Heading: "Home"})

		Expect(emptyHTML).ToNot(ContainSubstring("tracking-[0.32em]"))
		Expect(emptyDoc.Find("hgroup p").Length()).To(BeZero())
		Expect(emptyDoc.Find("div.flex.flex-wrap.gap-2").Length()).To(BeZero())
	})
})
