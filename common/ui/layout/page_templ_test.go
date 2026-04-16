package layout_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/common/ui/layout"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

		var render strings.Builder
		err := layout.Page(props).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		Expect(err).ToNot(HaveOccurred())
	})

	It("uses a wider desktop page shell", func() {
		Expect(doc.Find("section").First().AttrOr("class", "")).To(ContainSubstring("flex w-full flex-col gap-8"))
		Expect(html).ToNot(ContainSubstring("max-w-6xl"))
	})

	It("keeps the hero card compact and left aligned", func() {
		Expect(html).To(ContainSubstring("justify-start"))
		Expect(html).To(ContainSubstring("xl:w-1/2"))
		Expect(doc.Find("h1").First().AttrOr("class", "")).To(ContainSubstring("max-w-4xl"))
		Expect(doc.Find("hgroup p").First().AttrOr("class", "")).To(ContainSubstring("max-w-3xl"))
	})
})
