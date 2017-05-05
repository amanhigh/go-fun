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

func WriteLines(path string, lines []string) {
	ips := strings.Join(lines, "\n")
	ioutil.WriteFile(path, []byte(ips), 0644)
}
