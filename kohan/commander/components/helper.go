package components

import (
	"fmt"
	"strconv"
	"strings"
)

func Printf(template string, params string, marker string) {
	paramSplit := strings.Split(params, "\n")

	for _, param := range paramSplit {
		command := template
		split := strings.Split(param, " ")

		for i, value := range split {
			command = strings.Replace(command, marker+strconv.Itoa(i+1), value, -1)
		}
		fmt.Printf("%+v\n", command)
	}
}
