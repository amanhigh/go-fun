package pages_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Journal Detail Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.JournalDetailPage("jrn_1234abcd").Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Sidebar Actions", func() {
		It("should render review and note actions", func() {
			Expect(html).To(ContainSubstring(`x-bind:open="sidebar.state.actionOpen"`))
			Expect(html).To(ContainSubstring(`x-bind:open="sidebar.state.reviewOpen"`))
			Expect(html).To(ContainSubstring(`x-on:toggle="sidebar.state.setActionOpen($el.open)"`))
			Expect(html).To(ContainSubstring(`x-on:toggle="sidebar.state.setReviewOpen($el.open)"`))
			Expect(html).To(ContainSubstring("Actions"))
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("Quick actions"))
			Expect(html).To(ContainSubstring(`flex flex-wrap items-center gap-2 lg:gap-3`))
			Expect(html).To(ContainSubstring(`sidebar.reviewActions.toggleLabel()`))
			Expect(html).To(ContainSubstring(`sidebar.reviewActions.quickAction().label`))
			Expect(html).To(ContainSubstring(`Saving...`))
			Expect(html).To(ContainSubstring("Management"))
			Expect(html).To(ContainSubstring("Note"))
			Expect(html).To(ContainSubstring("Tags"))
			Expect(html).ToNot(ContainSubstring(`>Action<`))
			Expect(html).To(ContainSubstring(`x-on:click="header.deleteJournal()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.reviewActions.toggle()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.reviewActions.applyQuickStatus()"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.reviewActions.quickAction().status"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.managementTags.hasBar()"`))
			Expect(html).To(ContainSubstring(`x-for="preset in sidebar.managementTags.presets"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.managementTags.submit(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-bind:disabled="sidebar.managementTags.submitting || sidebar.managementTags.hasTag(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-bind:class="sidebar.managementTags.buttonClass(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-model="sidebar.reasonTagForm.input"`))
			Expect(html).To(ContainSubstring(`x-model="sidebar.reasonTagForm.override"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="sidebar.reasonTagForm.focusOverride()"`))
			Expect(html).To(ContainSubstring(`x-ref="reasonTagOverride"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="sidebar.reasonTagForm.submit()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.reasonTagForm.submit()"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.tags.hasItems()"`))
			Expect(html).To(ContainSubstring(`x-for="tag in sidebar.tags.all()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.tags.delete(tag.id)"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.notes.hasItems()"`))
			Expect(html).To(ContainSubstring(`x-for="note in sidebar.notes.sorted()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.state.enterReviewMode()"`))
			Expect(html).To(ContainSubstring(`tag.type].filter(Boolean).join(`))
			Expect(html).To(ContainSubstring(`x-bind:class="presentation.type(item.type).badgeClass"`))
			Expect(html).To(ContainSubstring(`x-text="presentation.formatReviewQueueDate(item.created_at)"`))
			Expect(html).To(ContainSubstring(`sidebar.reviewQueue.isLoading()`))
			Expect(html).To(ContainSubstring(`sidebar.reviewQueue.isError()`))
			Expect(html).To(ContainSubstring(`aria-label="Delete Note"`))
			Expect(html).To(ContainSubstring("h-4 w-4"))
		})
	})

	Context("Header Summary", func() {
		It("should render a compact summary card with new two-column layout", func() {
			// Identity unchanged
			Expect(html).To(ContainSubstring(`x-text="current.journal.ticker"`))
			Expect(html).To(ContainSubstring(`x-text="'ID: ' + current.journal.id"`))
			// Delete action
			Expect(html).To(ContainSubstring(`x-on:click="header.deleteJournal()"`))

			// Primary info row: type + status + sequence
			Expect(html).To(ContainSubstring(`presentation.sequence(current.journal.sequence).text`))
			Expect(html).To(ContainSubstring(`presentation.type(current.journal.type).text`))
			Expect(html).To(ContainSubstring(`presentation.status(current.journal.status).text`))

			// Right metadata: created + pending/review
			Expect(html).To(ContainSubstring(`x-text="presentation.formatTimestamp(current.journal.created_at)"`))
			Expect(html).To(ContainSubstring(`x-show="!current.journal.reviewed_at"`))
			Expect(html).To(ContainSubstring(`x-show="current.journal.reviewed_at"`))
			Expect(html).To(ContainSubstring(`presentation.reviewedAt(current.journal.reviewed_at).text`))
			Expect(html).To(ContainSubstring(`presentation.pendingReview().text`))

			// Tags rendered directly without section label
			Expect(html).ToNot(ContainSubstring(`Summary Tags`))
			Expect(html).ToNot(ContainSubstring(`Signal Tags`))

			// No old highlight card labels
			Expect(html).ToNot(ContainSubstring(`>STATUS</p>`))
			Expect(html).ToNot(ContainSubstring(`>TYPE</p>`))
			Expect(html).ToNot(ContainSubstring(`>SEQUENCE</p>`))
			Expect(html).ToNot(ContainSubstring(`>CREATED</p>`))
		})
	})
	Context("Header Tags", func() {
		It("should render separate primary and secondary tag sections", func() {
			Expect(html).To(ContainSubstring(`x-for="tag in sidebar.tags.reason()"`))
			Expect(html).To(ContainSubstring(`x-for="tag in sidebar.tags.directional()"`))
			Expect(html).To(ContainSubstring(`x-text="presentation.reasonTag(tag).text"`))
			Expect(html).To(ContainSubstring(`x-text="presentation.directionalTag(tag).text"`))
		})
	})

	Context("Image Preview Modal", func() {
		It("should render a visible timeframe chip and counter", func() {
			Expect(html).To(ContainSubstring(`x-bind:class="presentation.timeframe(preview.timeframe()).badgeClass"`))
			Expect(html).To(ContainSubstring(`x-text="preview.timeframe()"`))
			Expect(html).To(ContainSubstring(`x-text="preview.counter()"`))
		})

		It("should render keyboard navigation bindings for preview mode", func() {
			Expect(html).To(ContainSubstring(`x-on:keydown.escape.window="preview.close()"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.arrow-left.window="preview.prev()"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.arrow-right.window="preview.next()"`))
		})

		It("should render mouse navigation bindings for preview mode", func() {
			Expect(html).To(ContainSubstring(`x-on:click.stop="preview.next(true)"`))
			Expect(html).To(ContainSubstring(`x-on:contextmenu.prevent.stop="preview.prev(true)"`))
			Expect(html).To(ContainSubstring(`aria-label="Preview Image Navigation Overlay"`))
		})
	})

	Context("Image Tiles", func() {
		It("should render full-image friendly classes for journal screenshots", func() {
			Expect(html).To(ContainSubstring(`group`))
			Expect(html).To(ContainSubstring(`h-auto`))
			Expect(html).To(ContainSubstring(`p-0`))
			Expect(html).To(ContainSubstring(`items-start`))
			Expect(html).To(ContainSubstring(`justify-start`))
			Expect(html).To(ContainSubstring(`overflow-hidden`))
			Expect(html).To(ContainSubstring(`rounded-2xl`))
			Expect(html).To(ContainSubstring(`border-border`))
			Expect(html).To(ContainSubstring(`bg-muted`))
			Expect(html).To(ContainSubstring(`text-left`))
			Expect(html).To(ContainSubstring(`x-bind:class="presentation.timeframe(image.timeframe).badgeClass"`))
			Expect(html).To(ContainSubstring(`x-text="image.timeframe"`))
			Expect(html).To(ContainSubstring(`x-on:click="preview.open(index)"`))
			Expect(html).To(ContainSubstring(`x-bind:title="images.tileTitle(image)"`))
			Expect(html).To(ContainSubstring(`class="block h-auto w-full transition-transform duration-300 group-hover:scale-[1.01]"`))
			Expect(html).ToNot(ContainSubstring(`aspect-[15/10]`))
			Expect(html).ToNot(ContainSubstring(`object-cover`))
		})
	})
})
