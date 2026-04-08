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

var _ = Describe("Base Template Tests", func() {
	var (
		ctx   context.Context
		title string

		render strings.Builder
		html   string
		doc    *goquery.Document
	)

	BeforeEach(func() {
		ctx = context.Background()
		title = "Shadow Gate"
		err := layout.Base(title).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		// Parse HTML once for all tests
		doc, _ = goquery.NewDocumentFromReader(strings.NewReader(html))
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Core structure", func() {
		It("renders html shell and metadata", func() {
			Expect(strings.ToLower(html)).To(HavePrefix("<!doctype html>"))
			Expect(doc.Find("html").AttrOr("lang", "")).To(Equal("en"))
			Expect(doc.Find("meta[charset]").AttrOr("charset", "")).To(Equal("UTF-8"))
			Expect(doc.Find("meta[name='viewport']").AttrOr("content", "")).To(Equal("width=device-width, initial-scale=1.0"))
		})

		It("renders title, header, main and footer", func() {
			Expect(doc.Find("title").Text()).To(Equal(title))
			Expect(doc.Find("h1").Text()).To(Equal(title))
			Expect(doc.Find("main").Length()).To(Equal(1))
			Expect(doc.Find("footer").Length()).To(Equal(1))
			Expect(html).To(ContainSubstring("Built with TemplUI & Tailwind CSS, powered by AlpineJS"))
		})

		It("does not render global header navigation", func() {
			Expect(doc.Find("header nav").Length()).To(Equal(0))
			Expect(doc.Find("header a").Length()).To(Equal(0))
		})
	})

	Context("Asset dependencies", func() {
		It("includes stylesheet and required scripts in expected order", func() {
			Expect(doc.Find("link[href='/assets/css/app.css']").Length()).To(Equal(1))
			Expect(doc.Find("script[src='/assets/js/app.js']").Length()).To(Equal(1))
			Expect(doc.Find("script[src='https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js']").Length()).To(Equal(1))
			Expect(html).To(ContainSubstring("selectbox.min.js"))

			appScriptPos := strings.Index(html, "/assets/js/app.js")
			alpineScriptPos := strings.Index(html, "cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js")
			Expect(appScriptPos).To(BeNumerically(">=", 0))
			Expect(alpineScriptPos).To(BeNumerically(">", appScriptPos))
		})
	})

	Context("Edge Cases", func() {
		It("should handle empty title gracefully", func() {
			err := layout.Base("").Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			Expect(render.String()).To(ContainSubstring("<title></title>"))
		})

		It("should handle special characters in title", func() {
			specialTitle := "Test & Demo <Script>alert('xss')</Script>"
			err := layout.Base(specialTitle).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			Expect(render.String()).To(ContainSubstring("<title>Test &amp; Demo &lt;Script&gt;alert(&#39;xss&#39;)&lt;/Script&gt;</title>"))
		})
	})

	Context("Children rendering", func() {
		It("renders child content inside main", func() {
			var childRender strings.Builder
			childCtx := templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<section id="content-marker">Hello child</section>`)
				return err
			}))

			err := layout.Base(title).Render(childCtx, &childRender)
			Expect(err).ToNot(HaveOccurred())

			childDoc, err := goquery.NewDocumentFromReader(strings.NewReader(childRender.String()))
			Expect(err).ToNot(HaveOccurred())
			Expect(childDoc.Find("main #content-marker").Length()).To(Equal(1))
		})
	})

})
