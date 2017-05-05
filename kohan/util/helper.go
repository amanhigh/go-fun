package util

import (
	"io/ioutil"
	"strings"
)

func ReadLines(path string) ([]string, error) {
	if content, e := ioutil.ReadFile(path); e == nil {
		lines := strings.Split(string(content), "\n")
		return lines, nil
	} else {
		return nil, e
	}
}
