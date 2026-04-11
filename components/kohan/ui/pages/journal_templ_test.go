//nolint:dupl
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
		Expect(html).To(ContainSubstring("Journal Browser"))
		Expect(html).To(ContainSubstring("Journal entries"))
		Expect(html).To(ContainSubstring("Browse journals with client-side loading powered by Alpine.js."))
		Expect(html).To(ContainSubstring("Loading journals..."))
	})

	It("binds dynamic status and type badge classes", func() {
		Expect(html).To(ContainSubstring("x-bind:class=\"statusBadgeClass(journal.status)\""))
		Expect(html).To(ContainSubstring("x-bind:class=\"typeBadgeClass(journal.type)\""))
		Expect(html).To(ContainSubstring("x-text=\"normalizeStatus(journal.status)\""))
	})

	It("links journal id to the detail page", func() {
		Expect(html).To(ContainSubstring("x-bind:href=\"'/journal/' + journal.id\""))
		Expect(html).To(ContainSubstring("x-text=\"journal.id\""))
		Expect(html).To(ContainSubstring("x-text=\"journal.ticker\""))
	})
})
