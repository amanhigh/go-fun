package play

import (
	"context"
	"os"

	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.design/x/clipboard"
)

var _ = Describe("Clipboard", Label(models.GINKGO_SLOW), func() {
	var (
		err      error
		testData = "CopyThis!!"
		ctx      = context.Background()
	)

	BeforeEach(func() {
		err = clipboard.Init()
	})

	It("should build", func() {
		Expect(err).To(BeNil())
	})

	Context("Text Copy", func() {
		var ch <-chan struct{}
		BeforeEach(func() {
			ch = clipboard.Write(clipboard.FmtText, []byte(testData))
		})

		It("should be pasted", func() {
			pastedData := clipboard.Read(clipboard.FmtText)
			Expect(string(pastedData)).To(Equal(testData))
		})

		It("should signal overwrite", func() {
			// Overwrite Clipboard
			clipboard.Write(clipboard.FmtText, []byte("Overwrite"))
			Eventually(ch, 1).Should(Receive())
		})
	})

	Context("Image Copy", func() {
		var ch <-chan struct{}
		var imgData []byte

		BeforeEach(func() {
			imgData, err = os.ReadFile("../res/flower.jpg")
			Expect(err).To(BeNil())
			ch = clipboard.Write(clipboard.FmtImage, imgData)
		})

		It("should be pasted", func() {
			pastedImageData := clipboard.Read(clipboard.FmtImage)
			Expect(pastedImageData).To(Equal(imgData))
			// os.WriteFile("../res/test.jpg", imgData, 0644)
		})

		It("should signal overwrite", func() {
			// Overwrite Clipboard with a different image
			clipboard.Write(clipboard.FmtImage, []byte("newImage"))
			Eventually(ch, 1).Should(Receive())
		})
	})

	Context("Watch", func() {
		var (
			ch        <-chan []byte
			watchText = "I am Watching!!"
		)

		BeforeEach(func() {
			ch = clipboard.Watch(ctx, clipboard.FmtText)
		})

		It("should work", func() {
			clipboard.Write(clipboard.FmtText, []byte(watchText))
			Eventually(ch, "2s").Should(Receive(Equal([]byte(watchText))))
		})
	})
})
