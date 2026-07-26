package tools_test

import (
	"github.com/amanhigh/go-fun/common/tools"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clipboard", Ordered, func() {
	var (
		text string
		err  error
	)

	BeforeEach(func() {
		text = "hello-clipboard"
	})

	// Round-trip: write then read back.
	Context("RoundTrip", func() {
		BeforeEach(func() {
			err = tools.ClipCopy(text)
		})

		It("should copy without error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should paste the copied text", func() {
			pasted, pasteErr := tools.ClipPaste()
			Expect(pasteErr).ToNot(HaveOccurred())
			Expect(pasted).To(Equal(text))
		})
	})

	// Overwrite previous value.
	Context("Overwrite", func() {
		BeforeEach(func() {
			Expect(tools.ClipCopy("first")).To(Succeed())
			err = tools.ClipCopy("second")
		})

		It("should overwrite without error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should paste the new value", func() {
			pasted, pasteErr := tools.ClipPaste()
			Expect(pasteErr).ToNot(HaveOccurred())
			Expect(pasted).To(Equal("second"))
		})
	})
})
