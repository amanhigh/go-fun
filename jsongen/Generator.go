package main

import (
	"io"
	"text/template"
	"errors"
)

const JSON = "json"

type Metadata struct {
	PackageName   string
	Type          string
}

type Generator struct {
	Format string
}

func (g *Generator) Generate(writer io.Writer, metadata Metadata) error {
	tmpl, err := g.template()
	if err != nil {
		return nil
	}

	return tmpl.Execute(writer, metadata)
}

func (g *Generator) template() (*template.Template, error) {
	if g.Format != JSON {
		return nil, errors.New("Unsupported format")
	}

	tmpl, e := template.ParseFiles("/Users/amanpreet.singh/IdeaProjects/GoArena/src/github.com/amanhigh/go-fun/jsongen/tmpl/write_to_json.tmpl")
	return tmpl, e
}
