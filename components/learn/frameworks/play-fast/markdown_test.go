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
			headingText = "Sample Markdown File"
		)

		It("should start with root", func() {
			Expect(root).ShouldNot(BeNil())
			Expect(root).To(BeAssignableToTypeOf(&ast.Document{}))
			Expect(root.Type()).Should(Equal(ast.NodeType(3)))
			Expect(root.Text(data)).ShouldNot(BeNil())
			Expect(root.HasChildren()).Should(BeTrue())
			Expect(root.ChildCount()).Should(BeNumerically(">", 10))
		})

		It("should get first node", func() {
			node := root.FirstChild()
			Expect(node).ShouldNot(BeNil())
			Expect(node).To(BeAssignableToTypeOf(&ast.Heading{}))
			Expect(node.Type()).Should(Equal(ast.NodeType(1)))
			Expect(node.Text(data)).Should(Equal([]byte(headingText)))

			Expect(node.Parent()).Should(Equal(root))
			Expect(node.NextSibling()).ShouldNot(BeNil())
			Expect(node.HasChildren()).Should(BeTrue())
			Expect(node.ChildCount()).Should(Equal(1))

			_, exists := node.Attribute([]byte("id"))
			Expect(exists).Should(BeFalse())
		})

		It("should perform walk", func() {
			ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				// fmt.Println("\n Debug ----> ", string(node.Text(data)), entering, reflect.TypeOf(node))
				Expect(node).ShouldNot(BeNil())
				Expect(entering).Should(BeTrue())
				// Wait for First Node of Type Heading.
				switch n := node.(type) {
				case *ast.Heading:
					Expect(n).To(BeAssignableToTypeOf(&ast.Heading{}))
					Expect(node.Type()).Should(Equal(ast.NodeType(1)))
					Expect(node.Text(data)).Should(Equal([]byte(headingText)))
					return ast.WalkStop, nil
				}
				return ast.WalkContinue, nil
			})
		})
	})
})
