package components_test

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuickBar", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.QuickBarProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.QuickBarProps{
			Label:    "Actions",
			IconName: "",
			Theme:    components.ToneBlue,
			Class:    "",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Icon present", func() {
		BeforeEach(func() {
			props.IconName = "zap"
			err := components.QuickBar(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders icon, label, and theme class", func() {
			Expect(html).To(ContainSubstring("data-lucide"))
			Expect(html).To(ContainSubstring("Actions"))
			Expect(html).To(ContainSubstring("quick-bar-blue"))
		})
	})

	Context("Icon absent", func() {
		BeforeEach(func() {
			err := components.QuickBar(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("omits the icon and still renders the label", func() {
			Expect(html).NotTo(ContainSubstring("data-lucide"))
			Expect(html).To(ContainSubstring("Actions"))
		})
	})

	Context("Child slot", func() {
		var childHTML string

		BeforeEach(func() {
			childCtx := templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<button id="slot-child">Click me</button>`)
				return err
			}))
			err := components.QuickBar(props).Render(childCtx, &render)
			Expect(err).ToNot(HaveOccurred())
			childHTML = render.String()
		})

		It("renders children inside the actions div", func() {
			Expect(childHTML).To(ContainSubstring("slot-child"))
			Expect(childHTML).To(ContainSubstring("Click me"))
		})
	})
})

var _ = Describe("QuickAction", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.QuickActionProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.QuickActionProps{
			TextExpr:     "action.label",
			OnClickExpr:  "doAction()",
			DisabledExpr: "action.disabled",
			ClassExpr:    "action.active ? 'bg-green-100' : ''",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("Expression forwarding", func() {
		BeforeEach(func() {
			err := components.QuickAction(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("binds text, click, disabled, and class expressions", func() {
			Expect(html).To(ContainSubstring(`x-text="action.label"`))
			Expect(html).To(ContainSubstring(`x-on:click="doAction()"`))
			Expect(html).To(ContainSubstring(`x-bind:disabled="action.disabled"`))
			Expect(html).To(ContainSubstring(`x-bind:class="action.active ? &#39;bg-green-100&#39; : &#39;&#39;"`))
		})
	})
})
