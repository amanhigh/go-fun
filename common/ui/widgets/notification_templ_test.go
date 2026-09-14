package widgets_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	widgets "github.com/amanhigh/go-fun/common/ui/widgets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification Template", func() {
	var (
		ctx context.Context

		render strings.Builder
		doc    *goquery.Document
	)

	BeforeEach(func() {
		ctx = context.Background()
		render.Reset()
		err := widgets.Notification().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(render.String()))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Prop-free rendering", func() {
		It("renders the viewport and template scaffold without props", func() {
			// Notification() takes no arguments and still renders a usable viewport.
			Expect(doc.Find("[data-notification-viewport]").Length()).To(Equal(1))
			Expect(doc.Find("template[data-notification-template]").Length()).To(Equal(1))
		})
	})

	Context("Semantic presentation hooks", func() {
		It("carries stable semantic hooks for the shared stylesheet", func() {
			// Presentation lives in the shared notification stylesheet; the
			// template must not freeze inline Tailwind utility strings.
			Expect(doc.Find("[data-notification-viewport]").AttrOr("class", "")).To(ContainSubstring("notification-viewport"))
			Expect(doc.Find("[data-slot=\"alert\"]").AttrOr("class", "")).To(ContainSubstring("notification-card"))
			Expect(doc.Find("[data-notification-title]").AttrOr("class", "")).To(ContainSubstring("notification-title"))
			Expect(doc.Find("[data-notification-message]").AttrOr("class", "")).To(ContainSubstring("notification-message"))
			Expect(doc.Find("[data-notification-action-scaffold]").AttrOr("class", "")).To(ContainSubstring("notification-action"))
			Expect(doc.Find("[data-notification-dismiss]").AttrOr("class", "")).To(ContainSubstring("notification-dismiss"))
		})
	})

	Context("Accessible live-region attributes", func() {
		It("exposes an aria-live region for screen readers", func() {
			viewport := doc.Find("[data-notification-viewport]")
			Expect(viewport.AttrOr("aria-live", "")).To(Equal("polite"))
			Expect(viewport.AttrOr("role", "")).To(Equal("region"))
			Expect(viewport.AttrOr("aria-label", "")).To(Equal("Notifications"))
		})
	})

	Context("Variant icon scaffolds", func() {
		It("renders one hidden badge scaffold per variant with the semantic icon class", func() {
			icons := doc.Find("[data-variant-icon]")
			Expect(icons.Length()).To(Equal(2))
			Expect(doc.Find(`[data-variant-icon="success"]`).Length()).To(Equal(1))
			Expect(doc.Find(`[data-variant-icon="error"]`).Length()).To(Equal(1))
			icons.Each(func(_ int, s *goquery.Selection) {
				_, hasHidden := s.Attr("hidden")
				Expect(hasHidden).To(BeTrue())
				Expect(s.AttrOr("class", "")).To(Equal("notification-icon"))
			})
		})
	})

	Context("Action and dismiss scaffolding", func() {
		It("renders the optional action button scaffold with a label slot", func() {
			scaffold := doc.Find("[data-notification-action-scaffold]")
			Expect(scaffold.Length()).To(Equal(1))
			Expect(scaffold.AttrOr("type", "")).To(Equal("button"))
			// The scaffold is hidden until cloned by the runtime for a real action.
			_, hasHidden := scaffold.Attr("hidden")
			Expect(hasHidden).To(BeTrue())
			Expect(doc.Find("[data-notification-action-label]").Length()).To(Equal(1))
		})

		It("renders the explicit dismiss control with an accessible label", func() {
			dismiss := doc.Find("[data-notification-dismiss]")
			Expect(dismiss.Length()).To(Equal(1))
			Expect(dismiss.AttrOr("type", "")).To(Equal("button"))
			Expect(dismiss.AttrOr("aria-label", "")).To(Equal("Dismiss"))
		})
	})
})
