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
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("Management"))
			Expect(html).To(ContainSubstring("Note"))
			Expect(html).To(ContainSubstring("Tags"))
			Expect(html).To(ContainSubstring(`x-on:click="deleteJournal()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.toggleReview()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.applyQuickReviewStatus()"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.hasManagementBar()"`))
			Expect(html).To(ContainSubstring(`x-for="preset in sidebar.managementTagPresets"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.submitManagementTag(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-bind:disabled="sidebar.managementTagSubmitting || sidebar.hasManagementTag(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-bind:class="sidebar.managementTagButtonClass(preset.value)"`))
			Expect(html).To(ContainSubstring(`x-model="sidebar.reasonTagInput"`))
			Expect(html).To(ContainSubstring(`x-model="sidebar.reasonTagOverride"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="sidebar.focusReasonTagOverride()"`))
			Expect(html).To(ContainSubstring(`x-ref="reasonTagOverride"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="sidebar.submitReasonTag()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.submitReasonTag()"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.deletableTags().length"`))
			Expect(html).To(ContainSubstring(`x-for="tag in sidebar.deletableTags()"`))
			Expect(html).To(ContainSubstring(`x-on:click="sidebar.deleteTag(tag.id)"`))
			Expect(html).To(ContainSubstring(`x-show="sidebar.sortedNotes().length"`))
			Expect(html).To(ContainSubstring(`x-for="note in sidebar.sortedNotes()"`))
			Expect(html).To(ContainSubstring(`tag.type].filter(Boolean).join(`))
			Expect(html).To(ContainSubstring(`x-bind:class="reviewQueueItemClass(item.type)"`))
			Expect(html).To(ContainSubstring(`x-text="formatReviewQueueDate(item.created_at)"`))
			Expect(html).To(ContainSubstring(`aria-label="Delete Note"`))
			Expect(html).To(ContainSubstring("h-4 w-4"))
		})
	})

	Context("Header Summary", func() {
		It("should render a compact summary card with new two-column layout", func() {
			// Identity unchanged
			Expect(html).To(ContainSubstring(`x-text="journal.ticker"`))
			Expect(html).To(ContainSubstring(`x-text="'ID: ' + journal.id"`))
			// Delete action
			Expect(html).To(ContainSubstring(`x-on:click="deleteJournal()"`))

			// Primary info row: type + status + timeframe
			Expect(html).To(ContainSubstring(`📅`))
			Expect(html).To(ContainSubstring(`📈`))
			Expect(html).To(ContainSubstring(`FAIL`))

			// Right metadata: created + pending/review
			Expect(html).To(ContainSubstring(`🕒`))
			Expect(html).To(ContainSubstring(`x-text="formatTimestamp(journal.created_at)"`))
			Expect(html).To(ContainSubstring(`x-show="!journal.reviewed_at"`))
			Expect(html).To(ContainSubstring(`x-show="journal.reviewed_at"`))

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
			Expect(html).To(ContainSubstring(`x-for="tag in sidebar.deletableTags()"`))
			Expect(html).To(ContainSubstring(`REASON&#39;,&#39;MANAGEMENT`))
			Expect(html).To(ContainSubstring(`DIRECTION&#39;,&#39;LEGACY`))
		})
	})

	Context("Image Preview Modal", func() {
		It("should render keyboard navigation bindings for preview mode", func() {
			Expect(html).To(ContainSubstring(`x-on:keydown.escape.window="closeImagePreview()"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.arrow-left.window="prevImage()"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.arrow-right.window="nextImage()"`))
		})

		It("should render mouse navigation bindings for preview mode", func() {
			Expect(html).To(ContainSubstring(`x-on:click="prevImage(true)"`))
			Expect(html).To(ContainSubstring(`x-on:click="nextImage(true)"`))
			Expect(html).To(ContainSubstring(`x-on:click.stop="nextImage(true)"`))
			Expect(html).To(ContainSubstring(`x-on:contextmenu.prevent.stop="prevImage(true)"`))
			Expect(html).To(ContainSubstring(`aria-label="Preview Image Navigation Overlay"`))
		})
	})
})
