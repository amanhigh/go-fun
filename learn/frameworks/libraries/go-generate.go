package libraries

import (
	"html/template"
	"os"
)

type Metadata struct {
	PackageName string
	Type        string
}

func GenerateFun() {
	tmpl := template.New("jsonTemplate")
	if tmpl, err := tmpl.Parse(templateString()); err == nil {
		_ = tmpl.Execute(os.Stdout, Metadata{PackageName: "com.test.gen", Type: "string"})
	}
}

func templateString() string {
	return `//This is Aman's Generated File
//Request you not to mess with it :)

package {{ .PackageName }}

import (
	"encoding/json"
	"io"
)

func (obj {{ .Type }}) WriteTo(writer io.Writer) (int64, error) {
	data, err := json.Marshal(&obj)
	if err != nil {
		return 0, err
	}
	length, err := writer.Write(data)
	return int64(length), err
}`
}
