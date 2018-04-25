package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	typeName := flag.String("type", "", "Type that hosts io.WriterTo interface implementation")
	packageName := flag.String("package", "", "Package name")
	format := flag.String("format", "json", "Encoding format")
	flag.Parse()

	if *typeName == "" || *packageName == "" {
		flag.Usage()
		return
	}

	fmt.Println(os.Args[0])

	m := Metadata{PackageName: *packageName, Type: *typeName}
	g := &Generator{Format: *format}
	outputFile := getFileName(typeName)

	writer, _ := os.Create(filepath.Join(outputFile))
	defer writer.Close()

	if err := g.Generate(writer, m); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %s %s\n", *format, outputFile)
}

func getFileName(typeName *string) string {
	return strings.ToLower(*typeName) + "_json_writer.go"
}
