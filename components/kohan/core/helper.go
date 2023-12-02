package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
)

// Printf takes a template file, a parameter file, and a marker as input,
// and prints the formatted templates by replacing markers with parameters.
//
// templateFile: the path to the file containing the templates
// paramFile: the path to the file containing the parameters
// marker: the marker used to identify the placeholders in the templates
func Printf(templateFile string, paramFile string, marker string) {

	for _, templateSplit := range util.ReadAllLines(templateFile) {
		for _, paramLine := range util.ReadAllLines(paramFile) {

			template := templateSplit
			for i, param := range strings.Split(paramLine, " ") {
				template = strings.Replace(template, marker+strconv.Itoa(i+1), param, -1)
			}
			fmt.Println(template)
		}
		fmt.Println("")
	}
}
