package pages_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// attrValueExists reports whether any element in doc carries the attribute
// name with exactly value. Used for scoped Alpine binding checks.
func attrValueExists(doc *goquery.Document, name, value string) bool {
	return doc.Find("*").FilterFunction(func(_ int, s *goquery.Selection) bool {
		v, ok := s.Attr(name)
		return ok && v == value
	}).Length() > 0
}

var _ = Describe("Journal Detail Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		doc    *goquery.Document
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.JournalDetailPage("jrn_1234abcd").Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Page Composition", func() {
		It("binds the page root to the exact journal id and init", func() {
			root := doc.Find("section").First()
			Expect(root.AttrOr("x-data", "")).To(Equal(`journalDetailPage("jrn_1234abcd")`))
			Expect(root.AttrOr("x-init", "")).To(Equal("init()"))
		})

		It("places the detail header inside the hero, ahead of the 70/30 grid", func() {
			// Detail header lives within the hero card.
			Expect(doc.Find(".w-full.max-w-none .journal-detail-header").Length()).To(Equal(1))
			// The 70/30 content grid is rendered outside the hero (sibling subtree).
			Expect(doc.Find(".grid.gap-6").Length()).To(BeNumerically(">", 0))
			Expect(doc.Find(".w-full.max-w-none .grid.gap-6").Length()).To(Equal(0))
			// Loaded detail is gated by the journal.detail x-if template.
			Expect(doc.Find("template[x-if=\"journal.detail\"]").Length()).To(BeNumerically(">", 0))
			// The outer content grid carries the responsive 70/30 contract.
			contentGrid := doc.Find(".grid.gap-6").First()
			Expect(contentGrid.AttrOr("class", "")).To(ContainSubstring("xl:grid-cols-[minmax(0,7fr)_minmax(320px,3fr)]"))
		})
	})

	Context("Header Visual Structure", func() {
		It("shows the ticker and Journal ID badge", func() {
			Expect(attrValueExists(doc, "x-text", "journal.detail.ticker")).To(BeTrue())
			Expect(doc.Find(":contains('Journal ID:')").Length()).To(BeNumerically(">", 0))
			Expect(attrValueExists(doc, "x-text", "journal.detail.id")).To(BeTrue())
		})

		It("uses a centered balanced identity grid", func() {
			// Balanced outer tracks around a centered identity column.
			Expect(html).To(ContainSubstring("lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)]"))
			// Centered identity column wraps the ticker.
			identity := doc.Find(".flex.flex-col.items-center")
			Expect(identity.Length()).To(BeNumerically(">", 0))
			Expect(identity.Find("h2[x-text=\"journal.detail.ticker\"]").Length()).To(BeNumerically(">", 0))
		})

		It("exposes created-date and delete controls", func() {
			Expect(attrValueExists(doc, "x-text", "present.date.format(journal.detail.created_at)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "header.deleteJournal()")).To(BeTrue())
			Expect(attrValueExists(doc, "aria-label", "Delete Journal")).To(BeTrue())
		})

		It("renders the metadata row with type/status/timeframe badges", func() {
			Expect(doc.Find(".journal-detail-header-meta").Length()).To(BeNumerically(">", 0))
			Expect(attrValueExists(doc, "x-text", "present.type.label(journal.detail.type)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "present.status.label(journal.detail.status)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "present.timeframe.label(journal.detail.top_timeframe)")).To(BeTrue())
		})

		It("groups tags into reason/management/directional sections", func() {
			Expect(attrValueExists(doc, "x-show", "sidebar.tags.all().length")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.reason()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.management()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.directional()")).To(BeTrue())
		})
	})

	Context("Sidebar Visual Structure", func() {
		It("renders the major visible sections", func() {
			Expect(doc.Find(":contains('Command Center')").Length()).To(BeNumerically(">", 0))
			Expect(doc.Find(":contains('Review')").Length()).To(BeNumerically(">", 0))
			Expect(doc.Find(":contains('The Lineup')").Length()).To(BeNumerically(">", 0))
			Expect(doc.Find(":contains('Speak Now')").Length()).To(BeNumerically(">", 0))
			Expect(doc.Find(":contains('The Record')").Length()).To(BeNumerically(">", 0))
		})

		It("wires the two panel disclosure bindings", func() {
			Expect(attrValueExists(doc, "x-bind:open", "sidebar.state.actionOpen")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:toggle", "sidebar.state.setActionOpen($el.open)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:open", "sidebar.state.reviewOpen")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:toggle", "sidebar.state.setReviewOpen($el.open)")).To(BeTrue())
		})

		It("should use a horizontal badge-row and remaining-height image sizing", func() {
			// Dedicated horizontal badge-row for the timeframe/image-type chips
			Expect(html).To(ContainSubstring(`class="flex items-center gap-2"`))
			// Image should size to the remaining modal height, not a fixed viewport height
			Expect(html).To(ContainSubstring(`max-h-full max-w-full`))
			Expect(html).ToNot(ContainSubstring(`max-h-[82vh]`))
		})
	})

	Context("Gallery and Preview", func() {
		It("renders the Images section with a sorted collection loop", func() {
			Expect(doc.Find(":contains('Images')").Length()).To(BeNumerically(">", 0))
			Expect(attrValueExists(doc, "x-text", "images.countLabel()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "journal.detail.images.length")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "(image, index) in images.sorted()")).To(BeTrue())
			// Journal-specific empty state when no images are present.
			Expect(html).To(ContainSubstring("No images available for this journal."))
		})

		It("renders full-image tiles with preview, source, alt and error bindings", func() {
			// Tile opens the preview at its index and carries the file name.
			Expect(attrValueExists(doc, "x-on:click", "preview.open(index)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:title", "image.file_name")).To(BeTrue())
			// Image element bindings.
			Expect(attrValueExists(doc, "x-bind:src", "image.src")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:alt", "image.label")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:error", "$el.style.opacity='0.35'")).To(BeTrue())
			// Full-bleed tile layout contract (scoped to the preview-open button).
			tile := doc.Find("button").FilterFunction(func(_ int, s *goquery.Selection) bool {
				v, ok := s.Attr("x-on:click")
				return ok && v == "preview.open(index)"
			})
			Expect(tile.Length()).To(Equal(1))
			tileClass := tile.AttrOr("class", "")
			Expect(tileClass).To(ContainSubstring("w-full"))
			Expect(tileClass).To(ContainSubstring("h-auto"))
			Expect(tileClass).To(ContainSubstring("overflow-hidden"))
			Expect(tileClass).To(ContainSubstring("rounded-2xl"))
			Expect(tileClass).To(ContainSubstring("bg-muted"))
			Expect(tileClass).To(ContainSubstring("text-left"))
		})

		It("renders the preview modal media, close and navigation structure", func() {
			Expect(attrValueExists(doc, "x-show", "preview.hasPreview()")).To(BeTrue())
			// Cloak hides the modal until Alpine initializes.
			Expect(attrValueExists(doc, "x-cloak", "")).To(BeTrue())
			// Keyboard navigation.
			Expect(attrValueExists(doc, "x-on:keydown.escape.window", "preview.close()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:keydown.arrow-left.window", "preview.prev()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:keydown.arrow-right.window", "preview.next()")).To(BeTrue())
			// Close + overlay navigation.
			Expect(attrValueExists(doc, "x-on:click", "preview.close()")).To(BeTrue())
			Expect(attrValueExists(doc, "aria-label", "Preview Image Navigation Overlay")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click.stop", "preview.wrapNext()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:contextmenu.prevent.stop", "preview.wrapPrev()")).To(BeTrue())
			// Visible media + badges.
			Expect(attrValueExists(doc, "x-bind:src", "preview.src()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:alt", "preview.label()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "preview.timeframe()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "preview.imageType()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "preview.counter()")).To(BeTrue())
			// Filename binding surfaces the active preview file name.
			Expect(attrValueExists(doc, "x-text", "preview.fileName()")).To(BeTrue())
		})
	})

	Context("Loader Integration", func() {
		It("wires journal-specific busy/error/ready expressions", func() {
			Expect(attrValueExists(doc, "x-show", "journal.loader.isBusy()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && journal.loader.hasError()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "journal.loader.message")).To(BeTrue())
			// Ready gate uses the journal detail expression.
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && !journal.loader.hasError() && (journal.detail)")).To(BeTrue())
			// Empty visibility uses the negated journal detail expression.
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && !journal.loader.hasError() && !(journal.detail)")).To(BeTrue())
		})

		It("renders the empty message and retry action", func() {
			Expect(html).To(ContainSubstring("No journal details available."))
			Expect(attrValueExists(doc, "x-on:click", "journal.loadJournal('jrn_1234abcd')")).To(BeTrue())
		})
	})
})
