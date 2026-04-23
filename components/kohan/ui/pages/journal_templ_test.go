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
		Expect(html).To(ContainSubstring("x-data=\"journalPage()\""))
		Expect(html).To(ContainSubstring("x-init=\"init()\""))
		Expect(html).To(ContainSubstring("Loading journals..."))
	})

	It("binds dynamic status and type badge classes", func() {
		Expect(html).To(ContainSubstring("x-bind:class=\"statusBadgeClass(journal.status)\""))
		Expect(html).To(ContainSubstring("x-bind:class=\"typeBadgeClass(journal.type)\""))
		Expect(html).To(ContainSubstring("normalizeStatus(journal.status)"))
		Expect(html).To(ContainSubstring("journal.type"))
	})

	It("links journal id to the detail page", func() {
		Expect(html).To(ContainSubstring("x-bind:href=\"'/journal/' + journal.id\""))
		Expect(html).To(ContainSubstring("x-text=\"journal.id\""))
		Expect(html).To(ContainSubstring("x-text=\"journal.ticker\""))
	})

	It("renders review presets and removes the reviewed dropdown", func() {
		Expect(html).To(ContainSubstring("Review"))
		Expect(html).To(ContainSubstring("x-for=\"reviewPreset in reviewPresets\""))
		Expect(html).To(ContainSubstring("x-text=\"reviewPreset.label\""))
		Expect(html).To(ContainSubstring("x-on:click=\"applyReviewPreset(reviewPreset)\""))
		Expect(html).To(ContainSubstring("x-bind:class=\"reviewPresetClass(reviewPreset)\""))
		Expect(html).ToNot(ContainSubstring("id=\"journal-reviewed\""))
	})

	It("shows a month-only badge for the active review preset", func() {
		Expect(html).To(ContainSubstring("activeReviewPreset"))
		Expect(html).To(ContainSubstring("bg-amber-100"))
	})

	It("renders active filter chips for multiple filter types", func() {
		Expect(html).To(ContainSubstring("filter.createdBefore !== ''"))
		Expect(html).To(ContainSubstring("filter.reviewed !== '' && activeReviewPreset === ''"))
		Expect(html).To(ContainSubstring("filter.sortBy !== '' || filter.sortOrder !== ''"))
		Expect(html).To(ContainSubstring("'Created: ' + filter.createdAfter + ' → ' + filter.createdBefore"))
		Expect(html).To(ContainSubstring("'Review: ' + activeReviewPreset"))
		Expect(html).To(ContainSubstring("'Sort: ' + [filter.sortBy, filter.sortOrder].filter(Boolean).join(' · ')"))
	})

	It("renders shared error state with retry action", func() {
		Expect(html).To(ContainSubstring("x-show=\"hasError()\""))
		Expect(html).To(ContainSubstring("x-text=\"errorMessage\""))
		Expect(html).To(ContainSubstring("x-on:click=\"loadJournals()\""))
		Expect(html).To(ContainSubstring(">Retry<"))
	})

	It("renders themed quick bars with shared title sections", func() {
		Expect(html).To(ContainSubstring(">Date<"))
		Expect(html).To(ContainSubstring(">Review<"))
		Expect(html).To(ContainSubstring(">Quick<"))
		Expect(html).To(ContainSubstring("text-sky-700"))
		Expect(html).To(ContainSubstring("text-cyan-700"))
		Expect(html).To(ContainSubstring("text-slate-700"))
	})

	It("keeps shared quick bar action bindings", func() {
		Expect(html).To(ContainSubstring("applyCreatedPreset"))
		Expect(html).To(ContainSubstring("last7"))
		Expect(html).To(ContainSubstring("last30"))
		Expect(html).To(ContainSubstring("toggleType()"))
		Expect(html).To(ContainSubstring("x-bind:class=\"typeToggleClass()\""))
		Expect(html).To(ContainSubstring("x-text=\"typeToggleLabel()\""))
	})
})
