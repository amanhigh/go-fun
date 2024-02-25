package play_test

import (
	"bytes"
	"text/template"

	"github.com/amanhigh/go-fun/models/fun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Generate", func() {
	var (
		metadata fun.Metadata
		buffer   *bytes.Buffer
		tmpl     = template.New("gen-test")
		template string
		expected string
	)

	BeforeEach(func() {
		metadata = fun.Metadata{
			PackageName: fun.Inner{Name: "com.test.gen"},
			Type:        "string",
			Imports:     []string{"encoding/json", "io"},
		}
		buffer = &bytes.Buffer{}
	})

	AfterEach(func() {
		tmpl, err := tmpl.Parse(template)
		Expect(err).To(BeNil())

		err = tmpl.Execute(buffer, metadata)
		Expect(err).To(BeNil())

		Expect(buffer.String()).To(Equal(expected))
	})

	It("should have Inner Template", func() {
		template = "package {{ .PackageName.Name }}"
		expected = "package com.test.gen"
	})

	It("should work for Range Template", func() {
		template = "import ({{range .Imports}}{{.}}, {{end}})"
		expected = "import (encoding/json, io, )"
	})

	Context("If template", func() {

		BeforeEach(func() {
			template = "->{{if .Type}} fmt.Println({{.Type}}) {{else}} You Missed Supplying Type Variable {{end}}<-"
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
