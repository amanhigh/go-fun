package fun

type Inner struct {
	Name string
}

type Metadata struct {
	PackageName Inner
	Imports     []string
	Type        string
}

type Entry struct {
	Name string
	Done bool
}

type ToDo struct {
	User string
	List []Entry
}

// Constants
const TextTemplate = `
//This is Aman's Generated File
//Request you not to mess with it :)

package {{ .PackageName.Name }}

import ({{range .Imports}}
{{.}}{{end}}
)

{{if .Type}}
func (obj {{ .Type }}) WriteTo(writer io.Writer) (int64, error) {
	data, err := json.Marshal(&obj)
	if err != nil {
		return 0, err
	}
	length, err := writer.Write(data)
	return int64(length), err
}
{{ else }}
	You Missed Supplying Type Variable
{{end}}.
`
