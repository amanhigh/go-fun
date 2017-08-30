package commander

import (
	"os"
	"io/ioutil"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func AppendFile(path string, content string) {
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600); err == nil {
		defer f.Close()
		if _, err = f.WriteString(content); err != nil {
			log.WithFields(log.Fields{"File": path, "Error": err}).Error("Error Appending Content to File")
		}
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Error Opening File for Append")
	}
}

func ReadAllFiles(dir string) ([]string) {
	contents := []string{}
	if fileInfos, err := ioutil.ReadDir(dir); err == nil {
		for _, info := range fileInfos {
			if content, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", dir, info.Name())); err == nil {
				contents = append(contents, string(content))
			} else {
				log.WithFields(log.Fields{"Error": err}).Error("Error Reading File")
			}
		}
	} else {
		log.WithFields(log.Fields{"Directory": dir, "Error": err}).Error("Error Reading Directory")
	}
	return contents
}
