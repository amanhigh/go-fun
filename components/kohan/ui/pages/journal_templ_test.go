package pages_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Journal Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.JournalPage().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	It("renders empty journal page", func() {
		Expect(html).To(ContainSubstring("<title>Shadow Gate</title>"))
		Expect(html).To(ContainSubstring("Journal"))
		Expect(html).To(ContainSubstring("No entries yet"))
		Expect(html).To(ContainSubstring("This is an empty starting page for Journal."))
	})
})
