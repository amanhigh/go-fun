package layout_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/common/ui/layout"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test data for base template testing
var baseTestData = struct {
	Title string
}{
	Title: "Shadow Gate",
}

var _ = Describe("Base Template Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		doc    *goquery.Document
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := layout.Base(baseTestData.Title).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		// Parse HTML once for all tests
		doc, _ = goquery.NewDocumentFromReader(strings.NewReader(html))
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Template Structure Validation", func() {
		It("should render valid HTML structure", func() {
			Expect(doc.Find("html").Length()).To(Equal(1))
			Expect(doc.Find("html").AttrOr("lang", "")).To(Equal("en"))
			Expect(doc.Find("html").AttrOr("class", "")).To(Equal("h-full"))
			Expect(doc.Find("head").Length()).To(Equal(1))
			Expect(doc.Find("body").Length()).To(Equal(1))
			Expect(doc.Find("body").AttrOr("class", "")).To(Equal("h-full bg-background text-foreground"))
		})

		It("should include required meta tags", func() {
			Expect(doc.Find("meta[charset]").Length()).To(Equal(1))
			Expect(doc.Find("meta[charset]").AttrOr("charset", "")).To(Equal("UTF-8"))
			Expect(doc.Find("meta[name='viewport']").Length()).To(Equal(1))
			Expect(doc.Find("meta[name='viewport']").AttrOr("content", "")).To(Equal("width=device-width, initial-scale=1.0"))
		})
	})

	Context("Sections", func() {
		It("should render title correctly", func() {
			Expect(html).To(ContainSubstring("<title>" + baseTestData.Title + "</title>"))
		})

		It("should render header with brand title", func() {
			Expect(doc.Find("header").Length()).To(Equal(1))
			Expect(doc.Find("header").AttrOr("class", "")).To(ContainSubstring("border-b border-border bg-card"))
			Expect(doc.Find("nav").Length()).To(Equal(0))
		})

		It("should display page title in header", func() {
			h1 := doc.Find("h1")
			Expect(h1.Length()).To(Equal(1))
			Expect(h1.Text()).To(Equal(baseTestData.Title))
			Expect(h1.AttrOr("class", "")).To(Equal("text-xl font-semibold text-foreground"))
		})

		It("should not include global navigation links", func() {
			Expect(doc.Find("a").Length()).To(Equal(0))
		})

		It("should render main content area", func() {
			Expect(html).To(ContainSubstring("<main class=\"flex-1 container mx-auto px-4 py-8\"></main>"))
		})

		It("should render footer with attribution", func() {
			Expect(html).To(ContainSubstring("<footer class=\"border-t border-border bg-card mt-auto\">"))
			Expect(html).To(ContainSubstring("Built with TemplUI & Tailwind CSS, powered by AlpineJS"))
		})
	})

	Context("CSS and JavaScript Dependencies", func() {
		It("should include CSS and JS dependencies", func() {
			Expect(html).To(ContainSubstring("<link rel=\"stylesheet\" href=\"/assets/css/app.css\">"))
			Expect(html).To(ContainSubstring("selectbox.min.js"))
			Expect(html).To(ContainSubstring("<script defer src=\"/assets/js/app.js\"></script>"))
			Expect(html).To(ContainSubstring("<script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script>"))
		})
	})

	Context("Layout Structure", func() {
		It("should use proper layout structure", func() {
			Expect(html).To(ContainSubstring("<div class=\"min-h-screen flex flex-col\">"))
			Expect(html).To(ContainSubstring("<header"))
			Expect(html).To(ContainSubstring("<main"))
			Expect(html).To(ContainSubstring("<footer"))
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

	Context("Responsive Design Classes", func() {
		It("should include responsive classes", func() {
			Expect(html).To(ContainSubstring("container mx-auto"))
			Expect(html).To(ContainSubstring("px-4 py-4"))
			Expect(html).To(ContainSubstring("px-4 py-8"))
		})
	})

	Context("Theme and Styling", func() {
		It("should use proper theme color classes", func() {
			Expect(html).To(ContainSubstring("bg-background text-foreground"))
			Expect(html).To(ContainSubstring("text-muted-foreground"))
			Expect(html).To(ContainSubstring("bg-card"))
			Expect(html).To(ContainSubstring("border-border"))
		})
	})

})
