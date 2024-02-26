package play_fast

import (
	"bytes"
	htemplate "html/template"
	"strings"
	"text/template"

	"github.com/amanhigh/go-fun/models/fun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {
	var (
		metadata   fun.Metadata
		buffer     *bytes.Buffer
		goTemplate string
		expected   string
	)

	BeforeEach(func() {
		metadata = fun.Metadata{
			PackageName: fun.Inner{Name: "com.test.gen"},
			Type:        "string",
			Imports:     []string{"encoding/json", "io"},
		}
		buffer = &bytes.Buffer{}
	})

	Context("Text Template", func() {
		var (
			tmpl = template.New("gen-text").Funcs(template.FuncMap{"toLower": toLower})
		)

		Context("Parse", func() {
			AfterEach(func() {
				tmpl, err := tmpl.Parse(goTemplate)
				Expect(err).To(BeNil())

				err = tmpl.Execute(buffer, metadata)
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(Equal(expected))
			})

			It("should have Inner Template", func() {
				goTemplate = "package {{ .PackageName.Name }}"
				expected = "package com.test.gen"
			})

			It("should work for Range Template", func() {
				goTemplate = "import ({{range .Imports}}{{.}}, {{end}})"
				expected = "import (encoding/json, io, )"
			})

			It("should work for With Template", func() {
				goTemplate = "{{with .Type}}Type: {{.}}{{end}}"
				expected = "Type: string"
			})

			It("should work with Functions", func() {
				metadata.PackageName.Name = "COM.TEST.GEN"
				goTemplate = "Package name is {{.PackageName.Name | toLower}}"
				expected = "Package name is com.test.gen"
			})

			Context("If template", func() {

				BeforeEach(func() {
					goTemplate = "->{{if .Type}} fmt.Println({{.Type}}) {{else}} You Missed Supplying Type Variable {{end}}<-"
				})

				It("should with Value", func() {
					expected = "-> fmt.Println(string) <-"
				})

				It("should with No Value", func() {
					metadata.Type = ""
					expected = "-> You Missed Supplying Type Variable <-"
				})
			})
		})

		It("should support injection", func() {
			goTemplate = "{{define \"T\"}}Hello, {{.}}{{end}}"
			expected = "Hello, <script>alert('you have been pwned')</script>"

			tmpl, err := tmpl.Parse(goTemplate)
			Expect(err).To(BeNil())

			err = tmpl.ExecuteTemplate(buffer, "T", "<script>alert('you have been pwned')</script>")
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(Equal(expected))
		})
	})

	Context("Html Template", func() {
		var (
			htmpl = htemplate.New("gen-htmpl")
		)

		It("should not support injection", func() {
			goTemplate = "{{define \"T\"}}Hello, {{.}}{{end}}"
			//The < and > characters are replaced with &lt; and &gt; respectively,
			// and the single quote ' is replaced with &#39;. This prevents the string from being interpreted as a script,
			// thus avoiding code injection.
			expected = "Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;"

			tmpl, err := htmpl.Parse(goTemplate)
			Expect(err).To(BeNil())

			err = tmpl.ExecuteTemplate(buffer, "T", "<script>alert('you have been pwned')</script>")
			Expect(err).To(BeNil())
			Expect(buffer.String()).To(Equal(expected))
		})
	})

})

func toLower(s string) string {
	return strings.ToLower(s)
}
