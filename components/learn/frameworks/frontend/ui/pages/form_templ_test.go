package pages_test

import (
	"context"
	"strings"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Page Tests", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		err := pages.FormShowcasePage().Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Page Structure", func() {
		It("should render page with proper title and content", func() {
			Expect(html).To(ContainSubstring("<title>📝 Form Essentials</title>"))
			Expect(html).To(ContainSubstring("📝 Basic Components Showcase"))
			Expect(html).To(ContainSubstring("Master form inputs, validation, and user data collection patterns"))
		})

		It("should render navigation and breadcrumb", func() {
			Expect(html).To(ContainSubstring("Breadcrumb"))
			Expect(html).To(ContainSubstring("🏠 Home"))
			Expect(html).To(ContainSubstring("Form Essentials"))
			Expect(html).To(ContainSubstring("← Back to Home"))
		})

		It("should render multiple form sections", func() {
			Expect(html).To(ContainSubstring("🧾 Text Inputs"))
			Expect(html).To(ContainSubstring("⚙️ Selection Controls"))
			Expect(html).To(ContainSubstring("🚀 Form Actions"))
			Expect(html).To(ContainSubstring("🏷️ Status Indicators"))
			Expect(html).To(ContainSubstring("🪟 Modal Dialog"))
			Expect(strings.Count(html, "<article>")).To(Equal(5))
		})
	})

	Context("Form Components", func() {
		It("should include text input components", func() {
			Expect(html).To(ContainSubstring("username"))
			Expect(html).To(ContainSubstring("email"))
			Expect(html).To(ContainSubstring("notes"))
			Expect(html).To(ContainSubstring("Enter username"))
			Expect(html).To(ContainSubstring("name@example.com"))
			Expect(html).To(ContainSubstring("Share your goals..."))
		})

		It("should include selection controls", func() {
			Expect(html).To(ContainSubstring("selectbox"))
			Expect(html).To(ContainSubstring("country"))
			Expect(html).To(ContainSubstring("United States"))
			Expect(html).To(ContainSubstring("India"))
			Expect(html).To(ContainSubstring("United Kingdom"))
			Expect(html).To(ContainSubstring("Germany"))
			Expect(html).To(ContainSubstring("Subscription Plan"))
			Expect(html).To(ContainSubstring("Starter"))
			Expect(html).To(ContainSubstring("Pro"))
			Expect(html).To(ContainSubstring("terms and conditions"))
		})

		It("should include form action buttons", func() {
			Expect(html).To(ContainSubstring("Save Draft"))
			Expect(html).To(ContainSubstring("Submit Form"))
			Expect(html).To(ContainSubstring("Reset"))
			Expect(html).To(ContainSubstring("button"))
		})

		It("should include status badges", func() {
			Expect(html).To(ContainSubstring("Ready"))
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("Info"))
			Expect(html).To(ContainSubstring("Critical"))
		})

		It("should include modal dialog", func() {
			Expect(html).To(ContainSubstring("dialog"))
			Expect(html).To(ContainSubstring("showcase-dialog"))
			Expect(html).To(ContainSubstring("Preview Modal"))
			Expect(html).To(ContainSubstring("Submission Preview"))
			Expect(html).To(ContainSubstring("Cancel"))
			Expect(html).To(ContainSubstring("Confirm"))
		})
	})

	Context("Content and Styling", func() {
		It("should use proper form structure and styling", func() {
			Expect(html).To(ContainSubstring("<form>"))
			Expect(html).To(ContainSubstring("form"))
			Expect(html).To(ContainSubstring("label"))
			Expect(html).To(ContainSubstring("input"))
			Expect(html).To(ContainSubstring("textarea"))
			Expect(html).To(ContainSubstring("radio"))
			Expect(html).To(ContainSubstring("checkbox"))
		})

		It("should have semantic structure and accessibility", func() {
			Expect(html).To(ContainSubstring("<h1"))
			Expect(strings.Count(html, "<h2")).To(Equal(5))
			Expect(strings.Count(html, "<article")).To(Equal(5))
			Expect(html).To(ContainSubstring("<header>"))
			Expect(html).To(ContainSubstring("<footer>"))
			Expect(len(strings.TrimSpace(html))).To(BeNumerically(">", 2000))
		})

		It("should include validation and helper text", func() {
			Expect(html).To(ContainSubstring("3-20 characters: letters, numbers, underscore"))
			Expect(html).To(ContainSubstring("Use a valid email format"))
			Expect(html).To(ContainSubstring("Maximum 200 characters"))
			Expect(html).To(ContainSubstring("Select your primary region for services"))
		})
	})
})
