package libraries

import (
	"fmt"
	htemplate "html/template"
	"os"
	"text/template"
)

// HACK: #C GenerateFun as Ginkgo Test.
func GenerateFun() {
	printHtml()
	codeInjection()
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

/* Security */
func codeInjection() {
	fmt.Println("\nText Template, Script Injection")
	if t, err := template.New("text").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`); err == nil {
		err = t.ExecuteTemplate(os.Stdout, "T", "<script>alert('you have been pwned')</script>")
	}

	fmt.Println("\n\nHTML Template, Script Injection")
	if t, err := htemplate.New("html").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`); err == nil {
		err = t.ExecuteTemplate(os.Stdout, "T", "<script>alert('you have been pwned')</script>")
	}
}
