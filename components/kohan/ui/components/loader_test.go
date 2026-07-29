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

var _ = Describe("Loader", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		props  components.LoaderProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.LoaderProps{
			Loader:       "state",
			ReadyExpr:    "state.hasItems()",
			EmptyMessage: "No items.",
			RetryExpr:    "state.reload()",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	Context("State bindings", func() {
		BeforeEach(func() {
			err := components.Loader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders isBusy, hasError, message, and ReadyExpr bindings", func() {
			Expect(html).To(ContainSubstring(`x-show="state.isBusy()"`))
			Expect(html).To(ContainSubstring("state.hasError()"))
			Expect(html).To(ContainSubstring(`x-text="state.message"`))
			Expect(html).To(ContainSubstring("state.hasItems()"))
		})
	})

	Context("Retry present", func() {
		BeforeEach(func() {
			err := components.Loader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders a retry button with the click binding", func() {
			Expect(html).To(ContainSubstring(`x-on:click="state.reload()"`))
			Expect(html).To(ContainSubstring("Retry"))
		})
	})

	Context("Retry omitted", func() {
		BeforeEach(func() {
			props.RetryExpr = ""
			err := components.Loader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("does not emit a retry button", func() {
			Expect(html).NotTo(ContainSubstring("Retry"))
		})
	})

	Context("Empty message present", func() {
		BeforeEach(func() {
			err := components.Loader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("renders the empty state with correct visibility binding", func() {
			Expect(html).To(ContainSubstring("No items."))
			Expect(html).To(ContainSubstring("statebox-empty"))
			Expect(html).To(ContainSubstring(`!(state.hasItems())`))
		})
	})

	Context("Empty message omitted", func() {
		BeforeEach(func() {
			props.EmptyMessage = ""
			err := components.Loader(props).Render(ctx, &render)
			Expect(err).ToNot(HaveOccurred())
			html = render.String()
		})

		It("does not emit the empty state div", func() {
			Expect(html).NotTo(ContainSubstring("statebox-empty"))
		})
	})

	Context("Child rendering", func() {
		var childHTML string

		BeforeEach(func() {
			childCtx := templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<div id="child-content">Hello</div>`)
				return err
			}))
			err := components.Loader(props).Render(childCtx, &render)
			Expect(err).ToNot(HaveOccurred())
			childHTML = render.String()
		})

		It("renders child content inside the ready state div", func() {
			Expect(childHTML).To(ContainSubstring("child-content"))
			Expect(childHTML).To(ContainSubstring("Hello"))
		})
	})
})
