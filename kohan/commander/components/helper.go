package components

import (
	"fmt"
	"strconv"
	"strings"
)

func Printf(template string, paramPara string, marker string) {

	for _, templateSplit := range strings.Split(template, "\n") {
		for _, paramLine := range strings.Split(paramPara, "\n") {

			template := templateSplit
			for i, param := range strings.Split(paramLine, " ") {
				template = strings.Replace(template, marker+strconv.Itoa(i+1), param, -1)
			}
			fmt.Println(template)
		}
		fmt.Println("")
	}
}
