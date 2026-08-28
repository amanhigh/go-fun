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

	Context("Main Flow", func() {
		It("binds the page root to the exact journal id and init", func() {
			root := doc.Find("section").First()
			Expect(root.AttrOr("x-data", "")).To(Equal(`journalDetailPage("jrn_1234abcd")`))
			Expect(root.AttrOr("x-init", "")).To(Equal("init()"))
		})

		It("places the detail header ahead of the loaded content grid", func() {
			// Exactly one detail header and one loaded content grid exist.
			Expect(doc.Find(".journal-detail-header").Length()).To(Equal(1))
			Expect(doc.Find("div.grid.gap-6").Length()).To(Equal(1))
			// The header is rendered outside the content grid (sibling subtree).
			Expect(doc.Find("div.grid.gap-6 .journal-detail-header").Length()).To(Equal(0))
			// Ordered union: the header precedes the content grid in document order.
			union := doc.Find(".journal-detail-header, div.grid.gap-6")
			Expect(union.Length()).To(Equal(2))
			headerIdx, gridIdx := -1, -1
			union.Each(func(i int, s *goquery.Selection) {
				switch {
				case s.HasClass("journal-detail-header"):
					headerIdx = i
				case s.HasClass("grid") && s.HasClass("gap-6"):
					gridIdx = i
				}
			})
			Expect(headerIdx).To(BeNumerically(">=", 0))
			Expect(gridIdx).To(BeNumerically(">=", 0))
			Expect(headerIdx).To(BeNumerically("<", gridIdx))
			// The outer content grid carries the responsive image-rail + capped-sidebar contract.
			contentGrid := doc.Find("div.grid.gap-6").First()
			Expect(contentGrid.AttrOr("class", "")).To(ContainSubstring("xl:grid-cols-[minmax(0,1fr)_minmax(20rem,22rem)]"))
			// It owns exactly two child column divs (main rail + sidebar).
			childCols := contentGrid.Children().Filter("div")
			Expect(childCols.Length()).To(Equal(2))
			// The last/sidebar child is the page-owned sticky rail.
			sidebarCol := childCols.Last()
			sidebarClass := sidebarCol.AttrOr("class", "")
			Expect(sidebarClass).To(ContainSubstring("xl:sticky"))
			Expect(sidebarClass).To(ContainSubstring("xl:top-6"))
		})
	})

	Context("Header Flow", func() {
		It("renders the journal identity and created-date/delete controls", func() {
			// Ticker + Journal ID badge.
			Expect(attrValueExists(doc, "x-text", "journal.detail.ticker")).To(BeTrue())
			Expect(doc.Find(":contains('Journal ID:')").Length()).To(BeNumerically(">", 0))
			Expect(attrValueExists(doc, "x-text", "journal.detail.id")).To(BeTrue())
			// Created-date and delete controls.
			Expect(attrValueExists(doc, "x-text", "present.date.format(journal.detail.created_at)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "header.deleteJournal()")).To(BeTrue())
			Expect(attrValueExists(doc, "aria-label", "Delete Journal")).To(BeTrue())
			// Review action loop and application bindings.
			Expect(attrValueExists(doc, "x-for", "action in sidebar.reviewBar.actions()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "action.apply()")).To(BeTrue())
			// Busy state text and disabled expression.
			Expect(html).To(ContainSubstring("Saving..."))
			Expect(attrValueExists(doc, "x-bind:disabled", "sidebar.reviewBar.submitter.isBusy()")).To(BeTrue())
		})

		It("renders the header metadata, tag groups and Punjabi quotation", func() {
			// Metadata row with type/status/timeframe badges.
			Expect(doc.Find(".journal-detail-header-meta").Length()).To(BeNumerically(">", 0))
			Expect(attrValueExists(doc, "x-text", "present.type.label(journal.detail.type)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "present.status.label(journal.detail.status)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "present.timeframe.label(journal.detail.top_timeframe)")).To(BeTrue())
			// Tag groups into reason/management/directional sections.
			Expect(attrValueExists(doc, "x-show", "sidebar.tags.all().length")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.reason()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.management()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.directional()")).To(BeTrue())
			// Punjabi quotation as a single-line muted italic truncating meta item.
			const quote = "ਜਿਹੜੇ ਲੋਕ ਇਤਿਹਾਸ ਨੂੰ ਯਾਦ ਨਹੀਂ ਰੱਖਦੇ, ਉਹ ਇਸਨੂੰ ਦੁਹਰਾਉਣ ਲਈ ਸਰਾਪੇ ਜਾਂਦੇ ਹਨ। – ਜਾਰਜ ਸਾਂਤਾਯਾਨਾ"
			// The quote lives inside the header meta row as a truncating item.
			quoteSel := doc.Find(".journal-detail-header-meta div").FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.AttrOr("title", "") == quote
			}).First()
			Expect(quoteSel.Length()).To(Equal(1))
			// Exact full quotation text is present and mirrored by the accessible title.
			Expect(strings.TrimSpace(quoteSel.Text())).To(Equal(quote))
			Expect(quoteSel.AttrOr("title", "")).To(Equal(quote))
			// Single-line / truncation treatment, reusing horizontal space.
			class := quoteSel.AttrOr("class", "")
			Expect(class).To(ContainSubstring("truncate"))
			// It sits inside the meta row, not as a dedicated vertical block.
			Expect(quoteSel.Parent().AttrOr("class", "")).To(ContainSubstring("journal-detail-header-meta"))
		})
	})

	Context("Images Flow", func() {
		It("renders the image gallery with sorted loop and full-image preview entry point", func() {
			// Section identity: the title uses the standard header h3.
			title := doc.Find("h3").FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.Text() == "Images"
			})
			Expect(title.Length()).To(Equal(1))
			// Metadata binding scoped to the header (count label + zoom guidance).
			headerRow := title.Parent()
			Expect(headerRow.Find("span[x-text=\"images.countLabel()\"]").Length()).To(Equal(1))
			Expect(headerRow.Text()).To(ContainSubstring("· click to zoom"))

			Expect(attrValueExists(doc, "x-show", "journal.detail.images.length")).To(BeTrue())
			Expect(attrValueExists(doc, "x-for", "(image, index) in images.sorted()")).To(BeTrue())
			// Journal-specific empty state when no images are present.
			Expect(html).To(ContainSubstring("No images available for this journal."))

			// Full-image tiles open the preview at their index and carry the file name.
			Expect(attrValueExists(doc, "x-on:click", "preview.open(index)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:title", "image.file_name")).To(BeTrue())
			// Image element bindings.
			Expect(attrValueExists(doc, "x-bind:src", "image.src")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:alt", "image.label")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:error", "$el.style.opacity='0.35'")).To(BeTrue())
		})
	})

	Context("Record Flow", func() {
		It("renders The Record header with note-count metadata and retains notes behavior", func() {
			// Section identity: the title uses the standard header h3.
			title := doc.Find("h3").FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.Text() == "The Record"
			})
			Expect(title.Length()).To(Equal(1))
			// Metadata binding scoped to the header (note count + pluralized label).
			headerRow := title.Parent()
			Expect(headerRow.Find("span[x-text=\"sidebar.notes.sorted().length\"]").Length()).To(Equal(1))
			Expect(headerRow.Text()).To(ContainSubstring("note(s)"))

			// Retained description/notes list behavior.
			Expect(attrValueExists(doc, "x-for", "note in sidebar.notes.sorted()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "sidebar.notes.hasItems()")).To(BeTrue())
			Expect(html).To(ContainSubstring("No notes available for this journal."))
			// Note deletion wiring per row.
			Expect(attrValueExists(doc, "x-on:click", "sidebar.notes.delete(note.id)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:disabled", "sidebar.notes.submitter.isBusy()")).To(BeTrue())
		})
	})

	Context("Sidebar Flow", func() {
		var (
			primaryDetails *goquery.Selection
			fpDetails      *goquery.Selection
			lineupDetails  *goquery.Selection
			noteDetails    *goquery.Selection
		)

		BeforeEach(func() {
			// The two primary workflow disclosures are the <details> elements
			// carrying the persisted Alpine open/toggle bindings (Fingerprints
			// + The Lineup). The Add review note disclosure is a separate,
			// non-persisted shared disclosure identified by its own binding.
			primaryDetails = doc.Find("details").FilterFunction(func(_ int, s *goquery.Selection) bool {
				open, ok := s.Attr("x-bind:open")
				return ok && (open == "sidebar.state.actionOpen" || open == "sidebar.state.reviewOpen")
			})
			fpDetails = primaryDetails.FilterFunction(func(_ int, s *goquery.Selection) bool {
				return strings.Contains(strings.TrimSpace(s.Find("summary").Text()), "Fingerprints")
			})
			lineupDetails = primaryDetails.FilterFunction(func(_ int, s *goquery.Selection) bool {
				return strings.Contains(strings.TrimSpace(s.Find("summary").Text()), "The Lineup")
			})
			noteDetails = doc.Find("details").FilterFunction(func(_ int, s *goquery.Selection) bool {
				open, ok := s.Attr("x-bind:open")
				return ok && open == "sidebar.state.noteOpen"
			})
		})

		It("orders the workflow Fingerprints → The Lineup → Add review note", func() {
			// Scoped, index-based ordering across the sidebar disclosures.
			var fpIdx, lineupIdx, noteIdx = -1, -1, -1
			doc.Find("details").Each(func(i int, s *goquery.Selection) {
				label := strings.TrimSpace(s.Find("summary").Text())
				switch {
				case strings.Contains(label, "Fingerprints"):
					fpIdx = i
				case strings.Contains(label, "The Lineup"):
					lineupIdx = i
				case strings.Contains(label, "Add review note"):
					noteIdx = i
				}
			})
			Expect(fpIdx).To(BeNumerically(">=", 0))
			Expect(lineupIdx).To(BeNumerically(">=", 0))
			Expect(noteIdx).To(BeNumerically(">=", 0))
			Expect(fpIdx).To(BeNumerically("<", lineupIdx))
			Expect(lineupIdx).To(BeNumerically("<", noteIdx))
		})

		It("wires the sidebar disclosure state for the two panels and the Add review note", func() {
			Expect(primaryDetails.Length()).To(Equal(2))
			Expect(fpDetails.Length()).To(Equal(1))
			Expect(lineupDetails.Length()).To(Equal(1))
			Expect(fpDetails.AttrOr("x-bind:open", "")).To(Equal("sidebar.state.actionOpen"))
			Expect(fpDetails.AttrOr("x-on:toggle", "")).To(Equal("sidebar.state.setActionOpen($el.open)"))
			Expect(lineupDetails.AttrOr("x-bind:open", "")).To(Equal("sidebar.state.reviewOpen"))
			Expect(lineupDetails.AttrOr("x-on:toggle", "")).To(Equal("sidebar.state.setReviewOpen($el.open)"))
			// Shared Add review note disclosure uses the non-persisted noteOpen binding.
			Expect(noteDetails.Length()).To(Equal(1))
			Expect(noteDetails.AttrOr("x-bind:open", "")).To(Equal("sidebar.state.noteOpen"))
			Expect(noteDetails.AttrOr("x-on:toggle", "")).To(Equal("sidebar.state.setNoteOpen($el.open)"))
		})

		It("wires Fingerprints tag entry and deletion", func() {
			// Reason + override inputs carry their respective models.
			Expect(attrValueExists(doc, "x-model", "sidebar.reasonTagForm.input")).To(BeTrue())
			Expect(attrValueExists(doc, "x-model", "sidebar.reasonTagForm.override")).To(BeTrue())
			// Both Enter-key handlers: focus override, then submit the form.
			Expect(attrValueExists(doc, "x-on:keydown.enter.prevent", "$refs.reasonTagOverride.focus()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:keydown.enter.prevent", "sidebar.reasonTagForm.submit()")).To(BeTrue())
			// Submit button click + busy/invalid disabled guard.
			Expect(attrValueExists(doc, "x-on:click", "sidebar.reasonTagForm.submit()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:disabled", "sidebar.reasonTagForm.submitter.isBusy() || !sidebar.reasonTagForm.canSubmit()")).To(BeTrue())
			// Tag list loop and per-tag deletion with busy guard.
			Expect(attrValueExists(doc, "x-for", "tag in sidebar.tags.all()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "sidebar.tags.delete(tag.id)")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:disabled", "sidebar.tags.submitter.isBusy()")).To(BeTrue())
		})

		It("wires The Lineup navigation and review-mode entry", func() {
			// Queue loop with ticker + review-mode entry on click.
			Expect(attrValueExists(doc, "x-for", "item in sidebar.reviewQueue.all()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:href", "'/journal/' + item.id")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "sidebar.state.enterReviewMode()")).To(BeTrue())
			// Ticker + formatted date bindings.
			Expect(attrValueExists(doc, "x-text", "item.ticker")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "present.date.formatReviewQueueDate(item.created_at)")).To(BeTrue())
			// Queue Loader lifecycle expressions.
			Expect(attrValueExists(doc, "x-show", "sidebar.reviewQueue.loader.isBusy()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "!sidebar.reviewQueue.loader.isBusy() && sidebar.reviewQueue.loader.hasError()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "!sidebar.reviewQueue.loader.isBusy() && !sidebar.reviewQueue.loader.hasError() && (sidebar.reviewQueue.hasItems())")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "!sidebar.reviewQueue.loader.isBusy() && !sidebar.reviewQueue.loader.hasError() && !(sidebar.reviewQueue.hasItems())")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "sidebar.reviewQueue.load()")).To(BeTrue())
		})

		It("wires Add review note submission and busy guard", func() {
			// Textarea model + submit click with busy/invalid disabled guard.
			Expect(attrValueExists(doc, "x-model", "sidebar.noteForm.content")).To(BeTrue())
			Expect(attrValueExists(doc, "x-on:click", "sidebar.noteForm.submit()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-bind:disabled", "sidebar.noteForm.submitter.isBusy() || !sidebar.noteForm.canSubmit()")).To(BeTrue())
			// Note submitter binding.
			Expect(attrValueExists(doc, "x-show", "sidebar.noteForm.submitter.isBusy()")).To(BeTrue())
		})

	})

	Context("Preview Flow", func() {
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

		It("keeps the preview modal outside loader-ready content", func() {
			// Exactly one preview modal element exists.
			previewSel := doc.Find("*[x-show=\"preview.hasPreview()\"]")
			Expect(previewSel.Length()).To(Equal(1))
			// Exactly one journal ready element carries the existing ready expression.
			readyExpr := "!journal.loader.isBusy() && !journal.loader.hasError() && (journal.detail)"
			readySel := doc.Find("*[x-show=\"" + readyExpr + "\"]")
			Expect(readySel.Length()).To(Equal(1))
			// The ready element must not contain the preview modal (sibling composition).
			Expect(readySel.Find("*[x-show=\"preview.hasPreview()\"]").Length()).To(Equal(0))
		})
	})

	Context("Loader Flow", func() {
		It("wires the journal loader lifecycle expressions and retry action", func() {
			Expect(attrValueExists(doc, "x-show", "journal.loader.isBusy()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && journal.loader.hasError()")).To(BeTrue())
			Expect(attrValueExists(doc, "x-text", "journal.loader.message")).To(BeTrue())
			// Ready gate uses the journal detail expression.
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && !journal.loader.hasError() && (journal.detail)")).To(BeTrue())
			// Empty visibility uses the negated journal detail expression.
			Expect(attrValueExists(doc, "x-show", "!journal.loader.isBusy() && !journal.loader.hasError() && !(journal.detail)")).To(BeTrue())
			// Empty message and retry action.
			Expect(html).To(ContainSubstring("No journal details available."))
			Expect(attrValueExists(doc, "x-on:click", "journal.loadJournal('jrn_1234abcd')")).To(BeTrue())
		})
	})
})
