package play_fast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.design/x/clipboard"
)

var _ = FDescribe("Clipboard", func() {
	var (
		err      error
		testData = "CopyThis!!"
	)

	BeforeEach(func() {
		err = clipboard.Init()
	})

	It("should build", func() {
		Expect(err).To(BeNil())
	})

	Context("Paste", func() {
		var (
			ch <-chan struct{}
		)
		BeforeEach(func() {
			ch = clipboard.Write(clipboard.FmtText, []byte(testData))
		})

		It("should be pasted", func() {
			pastedData := clipboard.Read(clipboard.FmtText)
			Expect(string(pastedData)).To(Equal(testData))
		})

		It("should signal overwrite", func() {
			// Overwrite Clipboard
			clipboard.Write(clipboard.FmtText, []byte(""))
			Eventually(ch, 1).Should(Receive())
		})
	})
})
