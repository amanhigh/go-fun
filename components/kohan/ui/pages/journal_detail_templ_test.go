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
			Expect(html).To(ContainSubstring("Note"))
			Expect(html).To(ContainSubstring("Reason Tags"))
			Expect(html).To(ContainSubstring(`x-on:click="toggleReview()"`))
			Expect(html).To(ContainSubstring(`x-on:click="applyQuickReviewStatus()"`))
			Expect(html).To(ContainSubstring(`x-show="hasQuickReviewAction()"`))
			Expect(html).To(ContainSubstring(`x-text="reviewSubmitting ? &#39;Saving...&#39; : quickReviewLabel()"`))
			Expect(html).To(ContainSubstring(`x-model="reasonTagInput"`))
			Expect(html).To(ContainSubstring(`x-model="reasonTagOverride"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="focusReasonTagOverride()"`))
			Expect(html).To(ContainSubstring(`x-ref="reasonTagOverride"`))
			Expect(html).To(ContainSubstring(`x-on:keydown.enter.prevent="submitReasonTag()"`))
			Expect(html).To(ContainSubstring(`x-on:click="submitReasonTag()"`))
			Expect(html).To(ContainSubstring(`x-show="reasonTags().length"`))
			Expect(html).To(ContainSubstring(`x-on:click="deleteReasonTag(tag.id)"`))
			Expect(html).To(ContainSubstring(`x-bind:class="reviewQueueItemClass(item.type)"`))
			Expect(html).To(ContainSubstring(`x-text="formatReviewQueueDate(item.created_at)"`))
			Expect(html).To(ContainSubstring(`aria-label="Delete Note"`))
			Expect(html).To(ContainSubstring("h-4 w-4"))
		})
	})
})
