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

	Context("Main Flow", func() {
		It("should bootstrap the page with Alpine init", func() {
			Expect(html).To(ContainSubstring("<title>Shadow Gate</title>"))
			Expect(html).To(ContainSubstring("Journal Browser"))
			Expect(html).To(ContainSubstring("Journal entries"))
			Expect(html).To(ContainSubstring("x-data=\"journalPage()\""))
			Expect(html).To(ContainSubstring("x-init=\"init()\""))
		})

		It("should wire the initial page load flow", func() {
			Expect(html).To(ContainSubstring("x-init=\"init()\""))
			Expect(html).To(ContainSubstring("Loading journals..."))
			Expect(html).To(ContainSubstring("journal in table.journals"))
		})
	})

	Context("Filter Flow", func() {
		It("should wire quick date and type actions", func() {
			Expect(html).To(ContainSubstring("applyCreatedPreset"))
			Expect(html).To(ContainSubstring("last7"))
			Expect(html).To(ContainSubstring("last30"))
			Expect(html).To(ContainSubstring("toggleType()"))
			Expect(html).To(ContainSubstring("typeToggle().label"))
			Expect(html).To(ContainSubstring("typeToggle().className"))
		})

		It("should wire review preset actions", func() {
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("reviewPreset in presets.reviewPresets"))
			Expect(html).To(ContainSubstring("reviewPreset.label"))
			Expect(html).To(ContainSubstring("presets.applyReviewPreset(reviewPreset)"))
			Expect(html).To(ContainSubstring("presets.reviewPresetClass(reviewPreset)"))
		})

		It("should expose important active filter integrations", func() {
			Expect(html).To(ContainSubstring("datePreset"))
			Expect(html).To(ContainSubstring("reviewed"))
			Expect(html).To(ContainSubstring("Sort: "))
			Expect(html).To(ContainSubstring(`presentation.display(presentation.type(filter.type))`))
			Expect(html).To(ContainSubstring(`presentation.display(presentation.status(filter.status))`))
			Expect(html).To(ContainSubstring(`presentation.display(presentation.sequence(filter.sequence))`))
		})
	})

	Context("Table Flow", func() {
		It("should wire sortable table columns", func() {
			Expect(html).To(ContainSubstring("Ticker"))
			Expect(html).To(ContainSubstring("Sequence"))
			Expect(html).To(ContainSubstring("Created"))
		})

		It("should render journal row integration points", func() {
			Expect(html).To(ContainSubstring("journal.id"))
			Expect(html).To(ContainSubstring("x-text=\"journal.ticker\""))
			Expect(html).To(ContainSubstring("journal.sequence"))
		})

		It("should wire row status and type expressions", func() {
			Expect(html).To(ContainSubstring("journal.status"))
			Expect(html).To(ContainSubstring("journal.type"))
		})

		It("should use created_at descending as the default sort", func() {
			Expect(html).To(ContainSubstring("created_at"))
			Expect(html).To(ContainSubstring("desc"))
		})

		It("should render mutually exclusive table section states", func() {
			Expect(html).To(ContainSubstring("Loading journals..."))
			Expect(html).To(ContainSubstring("No journals found."))
		})

		It("should render pagination bindings and page summary", func() {
			Expect(html).To(ContainSubstring("x-on:click=\"pagination.previousPage()\""))
			Expect(html).To(ContainSubstring("x-on:click=\"pagination.nextJournalPage()\""))
			Expect(html).To(ContainSubstring("x-bind:disabled=\"!pagination.hasPrev()\""))
			Expect(html).To(ContainSubstring("x-bind:disabled=\"!pagination.hasNext()\""))
			Expect(html).To(ContainSubstring("Page"))
			Expect(html).To(ContainSubstring("of"))
		})
	})
})
