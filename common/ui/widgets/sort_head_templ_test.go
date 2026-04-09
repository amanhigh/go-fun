package widgets_test

import (
	"context"
	"strings"

	widgets "github.com/amanhigh/go-fun/common/ui/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortHead Template", func() {
	var (
		ctx  context.Context
		html string
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	JustBeforeEach(func() {
		var render strings.Builder
		err := widgets.SortHead(widgets.SortHeadProps{
			Label:         "Ticker",
			Field:         "ticker",
			SortByExpr:    "filterTracker.sortBy",
			SortOrderExpr: "filterTracker.sortOrder",
			OnClick:       "toggleSort('ticker')",
			Class:         "text-left",
		}).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	It("renders table head with aria sort expression", func() {
		Expect(html).To(ContainSubstring("x-bind:aria-sort=\"filterTracker.sortBy === &#39;ticker&#39; ? (filterTracker.sortOrder === &#39;asc&#39; ? &#39;ascending&#39; : &#39;descending&#39;) : &#39;none&#39;\""))
	})

	It("renders label, click action, and merged classes", func() {
		Expect(html).To(ContainSubstring("Ticker"))
		Expect(html).To(ContainSubstring("x-on:click=\"toggleSort(&#39;ticker&#39;)\""))
		Expect(html).To(ContainSubstring("inline-flex"))
		Expect(html).To(ContainSubstring("items-center"))
		Expect(html).To(ContainSubstring("gap-1"))
		Expect(html).To(ContainSubstring("text-left"))
	})

	It("renders unsorted expression using provided sortBy expression", func() {
		Expect(html).To(ContainSubstring("x-show=\"filterTracker.sortBy !== &#39;ticker&#39;\""))
		Expect(html).To(ContainSubstring("d=\"m21 16-4 4-4-4\""))
		Expect(html).To(ContainSubstring("text-sky-500"))
	})

	It("renders ascending and descending expressions using provided sort fields", func() {
		Expect(html).To(ContainSubstring("x-show=\"filterTracker.sortBy === &#39;ticker&#39; &amp;&amp; filterTracker.sortOrder === &#39;asc&#39;\""))
		Expect(html).To(ContainSubstring("x-show=\"filterTracker.sortBy === &#39;ticker&#39; &amp;&amp; filterTracker.sortOrder === &#39;desc&#39;\""))
		Expect(html).To(ContainSubstring("d=\"m5 12 7-7 7 7\""))
		Expect(html).To(ContainSubstring("d=\"M12 5v14\""))
		Expect(html).To(ContainSubstring("text-emerald-500"))
		Expect(html).To(ContainSubstring("text-rose-500"))
	})
})
