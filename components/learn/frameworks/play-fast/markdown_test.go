package play_fast

import (
	"bytes"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var _ = Describe("Markdown", func() {

	var (
		filePath = "../res/play.md"
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
			Expect(root.Kind()).Should(Equal(ast.KindDocument))
			Expect(root.Text(data)).ShouldNot(BeNil())
			Expect(root.HasChildren()).Should(BeTrue())
			Expect(root.ChildCount()).Should(BeNumerically(">", 10))
		})

		It("should perform walk", func() {
			ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				Expect(node).ShouldNot(BeNil())
				Expect(entering).Should(BeTrue())
				// Wait for First Node of Type Heading.
				switch n := node.(type) {
				case *ast.Heading:
					Expect(n).To(BeAssignableToTypeOf(&ast.Heading{}))
					Expect(node.Kind()).Should(Equal(ast.KindHeading))
					Expect(node.Text(data)).Should(Equal([]byte(headingText)))
					return ast.WalkStop, nil
				}
				return ast.WalkContinue, nil
			})
		})

		Context("First Node", func() {
			var (
				node ast.Node
			)

			BeforeEach(func() {
				node = root.FirstChild()
			})

			It("should be read", func() {
				Expect(node).ShouldNot(BeNil())
				Expect(node).To(BeAssignableToTypeOf(&ast.Heading{}))
				Expect(node.Kind()).Should(Equal(ast.KindHeading))
				Expect(node.Text(data)).Should(Equal([]byte(headingText)))

				Expect(node.Parent()).Should(Equal(root))
				Expect(node.NextSibling()).ShouldNot(BeNil())
				Expect(node.HasChildren()).Should(BeTrue())
				Expect(node.ChildCount()).Should(Equal(1))

				_, exists := node.Attribute([]byte("id"))
				Expect(exists).Should(BeFalse())
			})

			It("should have text", func() {
				text := node.FirstChild()

				Expect(text).ShouldNot(BeNil())
				Expect(text).To(BeAssignableToTypeOf(&ast.Text{}))
				Expect(text.Kind()).Should(Equal(ast.KindText))
				Expect(text.Text(data)).Should(Equal([]byte(headingText)))

				Expect(text.Parent()).Should(Equal(node))
				Expect(text.NextSibling()).Should(BeNil())
				Expect(text.HasChildren()).Should(BeFalse())
				Expect(text.ChildCount()).Should(Equal(0))
			})
		})

		Context("List", func() {
			var (
				list *ast.List
			)

			Context("Ordered", func() {
				BeforeEach(func() {
					ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
						switch n := node.(type) {
						case *ast.List:
							if n.IsOrdered() {
								list = n
								return ast.WalkStop, nil
							}
						}
						return ast.WalkContinue, nil
					})
				})

				It("should exist", func() {
					Expect(list).ShouldNot(BeNil())
					Expect(list.IsOrdered()).Should(BeTrue())
					Expect(list.ChildCount()).Should(Equal(3))

					// Sub List
					Expect(list.FirstChild().Kind()).Should(Equal(ast.KindListItem))
					Expect(list.FirstChild()).Should(BeAssignableToTypeOf(&ast.ListItem{}))
					Expect(list.FirstChild().Text(data)).Should(Equal([]byte("Level 1 Item 1")))
					Expect(list.LastChild().Text(data)).Should(Equal([]byte("Level 1 Item 3")))
				})

				Context("Sub Lists", func() {
					var (
						subList *ast.ListItem
					)

					BeforeEach(func() {
						ast.Walk(list, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
							switch n := node.(type) {
							case *ast.ListItem:
								if entering { // Only process the node when entering, not when exiting.
									txt := string(n.Text(data))
									if strings.HasPrefix(txt, "Level 1 Item 2") {
										subList = n
									}
								}
							}
							return ast.WalkContinue, nil
						})
					})

					It("should exist", func() {
						Expect(subList).ShouldNot(BeNil())
						Expect(subList).Should(BeAssignableToTypeOf(&ast.ListItem{}))
						Expect(string(subList.FirstChild().Text(data))).Should(Equal("Level 1 Item 2"))
						Expect(subList.ChildCount()).Should(Equal(2))
					})

					It("should have child Level 2 Item 2a", func() {
						subChild := subList.FirstChild()
						//FIXME: Add Sub Item Text Check
						Expect(subChild.ChildCount()).Should(Equal(1))
						Expect(subChild.FirstChild().Kind()).Should(Equal(ast.KindText))
					})
				})
			})

		})
	})
})
