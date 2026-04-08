package pages_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Index Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.IndexPage().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	It("renders the Shadow Gate landing page", func() {
		Expect(html).To(ContainSubstring("<title>Shadow Gate</title>"))
		Expect(html).To(ContainSubstring("Shadow Gate"))
		Expect(html).To(ContainSubstring("Welcome to Shadow Gate"))
	})

	It("includes the base layout", func() {
		Expect(html).To(ContainSubstring("<html lang=\"en\" class=\"h-full\">"))
		Expect(html).To(ContainSubstring("Shadow Gate"))
		Expect(html).To(ContainSubstring("Built with TemplUI & Tailwind CSS, powered by AlpineJS"))
	})
})
