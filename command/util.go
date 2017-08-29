package commander

import "os"

func AppendFile(path string, content string) {
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600); err == nil {
		defer f.Close()
		if _, err = f.WriteString(content); err != nil {
			PrintRed(err.Error())
		}
	} else {
		PrintRed(err.Error())
	}
}
