package main

import (
	"flag"
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

func main() {
	typeName := flag.String("type", "", "Type that hosts io.WriterTo interface implementation")
	packageName := flag.String("package", "", "Package name")
	format := flag.String("format", "json", "Encoding format")
	flag.Parse()

	outputFile := strings.ToLower(*typeName) + "_json_writer.go"
	writer, _ := os.Create(filepath.Join(outputFile))
	defer writer.Close()

	jsonGenerator := &Generator{Format: *format}

	m := Metadata{PackageName: *packageName, Type: *typeName}
	if err := jsonGenerator.Generate(writer, m); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %s %s\n", *format, outputFile)
}
