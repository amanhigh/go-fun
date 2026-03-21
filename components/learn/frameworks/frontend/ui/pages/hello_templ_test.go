package pages_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hello Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.HelloPage().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Page Structure", func() {
		It("should render page with proper title and content", func() {
			Expect(html).To(ContainSubstring("<title>Hello World Showcase</title>"))
			Expect(html).To(ContainSubstring("Hello World Showcase"))
			Expect(html).To(ContainSubstring("Demonstrating both TemplUI components"))
		})

		It("should render two main sections", func() {
			Expect(html).To(ContainSubstring("TemplUI Component"))
			Expect(html).To(ContainSubstring("Native Counter Component"))
			Expect(strings.Count(html, "section class=\"mb-8 p-6 border rounded-lg bg-card\"")).To(Equal(2))
		})
	})

	Context("Component Integration", func() {
		It("should include TemplUI selectbox component", func() {
			Expect(html).To(ContainSubstring("selectbox"))
			Expect(html).To(ContainSubstring("country"))
			Expect(html).To(ContainSubstring("United States"))
			Expect(html).To(ContainSubstring("India"))
			Expect(html).To(ContainSubstring("United Kingdom"))
		})

		It("should include native counter component", func() {
			Expect(html).To(ContainSubstring("x-data=\"{ count: 0 }\""))
			Expect(html).To(ContainSubstring("@click=\"count--\""))
			Expect(html).To(ContainSubstring("@click=\"count++\""))
			Expect(html).To(ContainSubstring("x-text=\"count\""))
		})
	})

	Context("Content and Styling", func() {
		It("should use consistent styling and have proper content", func() {
			Expect(html).To(ContainSubstring("text-muted-foreground"))
			Expect(html).To(ContainSubstring("font-semibold"))
			Expect(html).To(ContainSubstring("A professional selectbox component"))
			Expect(html).To(ContainSubstring("A simple counter built with native HTML"))
			Expect(html).To(ContainSubstring("Custom Features:"))
		})

		It("should have semantic structure and accessibility", func() {
			Expect(html).To(ContainSubstring("<h1"))
			Expect(strings.Count(html, "<h2")).To(Equal(2))
			Expect(strings.Count(html, "<section")).To(Equal(2))
			Expect(len(strings.TrimSpace(html))).To(BeNumerically(">", 1000))
		})
	})
})
