package libraries

import (
	"fmt"
	htemplate "html/template"
	"os"
)

func codeInjection() {
	fmt.Println("\n\nHTML Template, Script Injection")
	if t, err := htemplate.New("html").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`); err == nil {
		err = t.ExecuteTemplate(os.Stdout, "T", "<script>alert('you have been pwned')</script>")
	}
}
