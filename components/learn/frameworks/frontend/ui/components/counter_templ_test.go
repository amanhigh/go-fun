package components_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Counter Component Tests", func() {
	var (
		ctx     context.Context
		render  strings.Builder
		html    string
		doc     *goquery.Document
		buttons *goquery.Selection
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := components.Counter().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()

		// Parse HTML once for all tests
		doc, _ = goquery.NewDocumentFromReader(strings.NewReader(html))
		buttons = doc.Find("button")
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Component Structure", func() {
		It("should render counter with proper structure", func() {
			container := doc.Find("div[x-data]")
			Expect(container.Length()).To(Equal(1))
			Expect(container.AttrOr("x-data", "")).To(ContainSubstring("count: 0"))
			Expect(container.AttrOr("class", "")).To(Equal("flex items-center gap-4"))
			Expect(buttons.Length()).To(Equal(3))

			display := doc.Find("span[x-text]")
			Expect(display.Length()).To(Equal(1))
			Expect(display.AttrOr("x-text", "")).To(Equal("count"))
			Expect(display.AttrOr("class", "")).To(Equal("counter-display text-2xl font-semibold min-w-[4rem] text-center"))
			Expect(display.Text()).To(Equal("0"))
		})
	})

	Context("Button Functionality", func() {
		It("should render buttons with proper functionality and styling", func() {
			// Decrement button
			decrementBtn := buttons.Eq(0)
			Expect(decrementBtn.AttrOr("@click", "")).To(Equal("count--"))
			Expect(decrementBtn.AttrOr("class", "")).To(ContainSubstring("bg-red-500"))
			Expect(decrementBtn.AttrOr("class", "")).To(ContainSubstring("hover:bg-red-600"))
			Expect(decrementBtn.Text()).To(Equal("-"))

			// Increment button
			incrementBtn := buttons.Eq(1)
			Expect(incrementBtn.AttrOr("@click", "")).To(Equal("count++"))
			Expect(incrementBtn.AttrOr("class", "")).To(ContainSubstring("bg-green-500"))
			Expect(incrementBtn.AttrOr("class", "")).To(ContainSubstring("hover:bg-green-600"))
			Expect(incrementBtn.Text()).To(Equal("+"))

			// Reset button
			resetBtn := buttons.Eq(2)
			Expect(resetBtn.AttrOr("@click", "")).To(Equal("count = 0"))
			Expect(resetBtn.AttrOr("class", "")).To(ContainSubstring("bg-gray-500"))
			Expect(resetBtn.AttrOr("class", "")).To(ContainSubstring("hover:bg-gray-600"))
			Expect(resetBtn.Text()).To(Equal("Reset"))
		})
	})

	Context("Layout and Styling", func() {
		It("should use proper layout and styling", func() {
			// Container layout
			container := doc.Find("div[x-data]")
			Expect(container.AttrOr("class", "")).To(ContainSubstring("flex"))
			Expect(container.AttrOr("class", "")).To(ContainSubstring("items-center"))
			Expect(container.AttrOr("class", "")).To(ContainSubstring("gap-4"))

			// Display styling
			display := doc.Find("span[x-text]")
			Expect(display.AttrOr("class", "")).To(ContainSubstring("text-2xl"))
			Expect(display.AttrOr("class", "")).To(ContainSubstring("font-semibold"))
			Expect(display.AttrOr("class", "")).To(ContainSubstring("min-w-[4rem]"))
			Expect(display.AttrOr("class", "")).To(ContainSubstring("text-center"))

			// Consistent button styling
			buttons.Each(func(i int, s *goquery.Selection) {
				class := s.AttrOr("class", "")
				Expect(class).To(ContainSubstring("counter-button"))
				Expect(class).To(ContainSubstring("px-4 py-2"))
				Expect(class).To(ContainSubstring("text-white"))
				Expect(class).To(ContainSubstring("rounded-md"))
				Expect(class).To(ContainSubstring("transition-colors"))
				Expect(class).To(ContainSubstring("font-medium"))
			})
		})

		It("should have accessible elements with proper content", func() {
			// Buttons already tested for text content in functionality test
			Expect(buttons.Length()).To(Equal(3))

			// Counter display accessibility
			display := doc.Find("span[x-text]")
			Expect(display.Length()).To(Equal(1))
			Expect(display.Text()).ToNot(BeEmpty())
		})
	})
})
