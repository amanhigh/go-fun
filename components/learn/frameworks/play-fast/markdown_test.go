package play_fast

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var _ = FDescribe("Markdown", func() {

	var (
		filePath = "./res/play.md"
		data     []byte
		err      error
		root     ast.Node
	)

	BeforeEach(func() {
		// Read File
		data, err = os.ReadFile(filePath)
		Expect(err).ShouldNot(HaveOccurred())

		// Parse File
		root = goldmark.DefaultParser().Parse(text.NewReader(data))
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("file should be read", func() {
		var buf bytes.Buffer
		err = goldmark.Convert(data, &buf)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should be parsed", func() {
		Expect(root).ShouldNot(BeNil())
	})

	Context("Traverse", func() {
		var (
			rootText = "Sample Markdown File"
		)
		It("should read first node", func() {
			ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				// fmt.Println("\nName", string(node.Text(data)), entering)
				if entering {
					Expect(node).ShouldNot(BeNil())
					switch n := node.(type) {
					case *ast.Heading:
						if n.Level == 1 {
							Expect(node.Text(data)).Should(Equal([]byte(rootText)))
							Expect(node.Type()).Should(Equal(ast.NodeType(1)))
							return ast.WalkStop, nil
						}
					}
				}
				return ast.WalkContinue, nil
			})
		})
	})
})
