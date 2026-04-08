package layout_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/layout"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Base Template Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := layout.Base("Shadow Gate").Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	It("renders the base html shell", func() {
		Expect(html).To(ContainSubstring("<html lang=\"en\" class=\"h-full\">"))
		Expect(html).To(ContainSubstring("<body class=\"h-full bg-gray-50 text-gray-900\">"))
		Expect(html).To(ContainSubstring("<main class=\"flex-1 container mx-auto px-4 py-8\">"))
	})

	It("renders header, main, and footer", func() {
		Expect(html).To(ContainSubstring("Shadow Gate"))
		Expect(html).To(ContainSubstring("Home"))
		Expect(html).To(ContainSubstring("Built with Templ & Tailwind CSS"))
	})

	It("escapes the title and includes meta tags", func() {
		Expect(html).To(ContainSubstring("<title>Shadow Gate</title>"))
		Expect(html).To(ContainSubstring("<meta charset=\"UTF-8\">"))
		Expect(html).To(ContainSubstring("name=\"viewport\""))
		Expect(html).To(ContainSubstring("cdn.tailwindcss.com"))
		Expect(html).To(ContainSubstring("cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"))
		Expect(html).To(ContainSubstring("templui/js/"))
	})
})
