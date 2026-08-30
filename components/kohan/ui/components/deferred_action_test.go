package components_test

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/components/kohan/ui/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeferredAction", func() {
	var (
		ctx    context.Context
		render strings.Builder
		html   string
		doc    *goquery.Document
		props  components.DeferredActionProps
	)

	BeforeEach(func() {
		ctx = context.Background()
		props = components.DeferredActionProps{
			StateExpr:   "da",
			CancelLabel: "Cancel",
		}
	})

	AfterEach(func() {
		render.Reset()
	})

	// renderDeferred renders the component and parses the resulting HTML.
	renderDeferred := func() {
		err := components.DeferredAction(props).Render(ctx, &render)
		Expect(err).ToNot(HaveOccurred())
		html = render.String()
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
		Expect(err).ToNot(HaveOccurred())
	}

	// rootHasCloak reports whether the deferred-action root element (the one
	// carrying the x-show binding) also carries x-cloak.
	rootHasCloak := func() bool {
		sel := doc.Find("*[x-show=\"" + props.StateExpr + ".active\"]")
		if sel.Length() == 0 {
			return false
		}
		_, ok := sel.Attr("x-cloak")
		return ok
	}

	Context("Active banner contract", func() {
		BeforeEach(func() {
			renderDeferred()
		})

		It("renders only while active, exposes a polite live countdown, and binds cancel", func() {
			Expect(html).To(ContainSubstring(`x-show="da.active"`))
			Expect(html).To(ContainSubstring(`aria-live="polite"`))
			Expect(html).To(ContainSubstring(`x-text="da.message + &#39; (&#39; + da.remainingSeconds + &#39;s)&#39;"`))
			Expect(html).To(ContainSubstring(`x-on:click="da.cancel()"`))
			Expect(html).To(ContainSubstring("Cancel"))
			Expect(rootHasCloak()).To(BeTrue())
		})
	})

	Context("Configurable state expression", func() {
		BeforeEach(func() {
			props.StateExpr = "journal.da"
			renderDeferred()
		})

		It("uses the provided state expression for all bindings", func() {
			Expect(html).To(ContainSubstring(`x-show="journal.da.active"`))
			Expect(html).To(ContainSubstring(`x-text="journal.da.message + &#39; (&#39; + journal.da.remainingSeconds + &#39;s)&#39;"`))
			Expect(html).To(ContainSubstring(`x-on:click="journal.da.cancel()"`))
			Expect(rootHasCloak()).To(BeTrue())
		})
	})

	Context("Configurable cancel label", func() {
		BeforeEach(func() {
			props.CancelLabel = "Undo"
			renderDeferred()
		})

		It("renders the caller-provided cancel label", func() {
			Expect(html).To(ContainSubstring("Undo"))
		})
	})
})
