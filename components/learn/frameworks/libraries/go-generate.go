package libraries

import (
	"fmt"
	htemplate "html/template"
	"os"
	"text/template"
)

// HACK: #C GenerateFun as Ginkgo Test.
func GenerateFun() {
	codeInjection()
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
