package play_fast

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yuin/goldmark"
)

var _ = FDescribe("Markdown", func() {

	var (
		filePath = "./res/play.md"
		data     []byte
		err      error
		buf      bytes.Buffer
	)

	BeforeEach(func() {
		data, err = os.ReadFile(filePath)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("file should be read", func() {
		err = goldmark.Convert(data, &buf)
		Expect(err).ShouldNot(HaveOccurred())
	})
})
