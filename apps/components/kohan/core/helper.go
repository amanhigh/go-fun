package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/util"
)

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
