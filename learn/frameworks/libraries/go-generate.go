package libraries

import (
	"html/template"
	"os"
)

type Metadata struct {
	PackageName string
	Imports     []string
	Type        string
}

type entry struct {
	Name string
	Done bool
}

type ToDo struct {
	User string
	List []entry
}

func GenerateFun() {
	printText()
	printHtml()
}

/* Text Template */
func printText() {
	tmpl := template.New("jsonTemplate")
	if tmpl, err := tmpl.Parse(textTemplate()); err == nil {
		_ = tmpl.Execute(os.Stdout, Metadata{
			PackageName: "com.test.gen",
			Type:        "string",
			Imports: []string{"encoding/json",
				"io"}})
	}
}

func textTemplate() string {
	return `//This is Aman's Generated File
//Request you not to mess with it :)

package {{ .PackageName }}

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
}

/* HTML Template */
func printHtml() {
	tmpl := template.New("htmlTemplate")
	if tmpl, err := tmpl.Parse(htmlTemplate()); err == nil {
		_ = tmpl.Execute(os.Stdout, ToDo{
			User: "Amanpreet Singh",
			List: []entry{{
				Name: "Prepare Slides",
				Done: true,
			}, {
				Name: "Give Demo",
				Done: false,
			}},
		})
	}
}

func htmlTemplate() string {
	return `<!DOCTYPE html>
<html>
  <head>
    <title>Go To-Do list</title>
  </head>
  <body>
    <p>
      To-Do list for user: {{ .User }} 
    </p>
    <table>
      	<tr>
          <td>Task</td>
          <td>Done</td>
    	</tr>
      	{{ with .List }}
			{{ range . }}
      			<tr>
              		<td>{{ .Name }}</td>
              		<td>{{ if .Done }}Yes{{ else }}No{{ end }}</td>
      			</tr>
			{{ end }} 
      	{{ end }}
    </table>
  </body>
</html>`
}
