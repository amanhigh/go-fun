package play_test

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/amanhigh/go-fun/models/fun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Generate", func() {
	var (
		metadata fun.Metadata
		buffer   *bytes.Buffer
	)

	BeforeEach(func() {
		metadata = fun.Metadata{
			PackageName: fun.Inner{Name: "com.test.gen"},
			Type:        "string",
			Imports:     []string{"encoding/json", "io"},
		}
		buffer = &bytes.Buffer{}
	})

	It("should have Inner Template", func() {
		tmpl := template.New("jsonTemplate")
		tmpl, err := tmpl.Parse(fun.InnerTemplate)
		Expect(err).To(BeNil())

		err = tmpl.Execute(buffer, metadata)
		Expect(err).To(BeNil())

		expectedOutput := "package com.test.gen"
		Expect(buffer.String()).To(Equal(expectedOutput))
	})

	It("should work for Range Template", func() {
		tmpl := template.New("jsonTemplate")
		tmpl, err := tmpl.Parse(fun.RangeTemplate)
		Expect(err).To(BeNil())

		err = tmpl.Execute(buffer, metadata)
		Expect(err).To(BeNil())

		expectedOutput := "import ( encoding/json,  io, )"
		fmt.Println(buffer.String())
		Expect(buffer.String()).To(Equal(expectedOutput))
	})
})
