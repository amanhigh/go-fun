package commander

import (
	"os"
	"io/ioutil"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
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

func ReadAllFiles(dirPath string) ([]string) {
	contents := []string{}
	contentMap := ReadFileMap(dirPath)
	for _, value := range contentMap {
		contents = append(contents, value...)
	}
	return contents
}

func ReadFileMap(dirPath string) (map[string][]string) {
	contents := map[string][]string{}
	if fileInfos, err := ioutil.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			contents[info.Name()] = ReadAllLines(filePath)
		}
	} else {
		log.WithFields(log.Fields{"Directory": dirPath, "Error": err}).Error("Error Reading Directory")
	}
	return contents
}

func ReadAllLines(filePath string) []string {
	if content, err := ioutil.ReadFile(filePath); err == nil {
		return strings.Split(string(content), "\n")
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Error Reading File")
		return []string{}
	}
}

func ClearDirectory(dirPath string) {
	if fileInfos, err := ioutil.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			os.Remove(filePath)
		}
	}
}
