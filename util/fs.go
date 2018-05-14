package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"strconv"

	log "github.com/Sirupsen/logrus"
)

const DEFAULT_PERM = os.FileMode(0644)     //Owner RW,Group R,Other R
const DIR_DEFAULT_PERM = os.FileMode(0755) //Owner RWX,Group RX,Other RX
/*
	Helpfull File Related Cheatsheet
	https://www.devdungeon.com/content/working-files-go#read_quick
*/

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

func ReadAllFiles(dirPath string) []string {
	contents := []string{}
	contentMap := ReadFileMap(dirPath, false)
	for _, value := range contentMap {
		contents = append(contents, value...)
	}
	return contents
}

func ReadFileMap(dirPath string, readEmpty bool) map[string][]string {
	contents := map[string][]string{}
	for _, filePath := range ListFiles(dirPath) {
		if lines := ReadAllLines(filePath); len(lines) > 0 || readEmpty {
			contents[filePath] = lines
		}
	}
	return contents
}

func PrintFile(title string, filepath string) {
	PrintSkyBlue(title)
	fmt.Println(strings.Join(ReadAllLines(filepath), "\n"))
}

func ListFiles(dirPath string) []string {
	filePaths := []string{}
	if fileInfos, err := ioutil.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			filePaths = append(filePaths, filePath)
		}
	} else {
		log.WithFields(log.Fields{"Directory": dirPath, "Error": err}).Error("Error Reading Directory")
	}
	return filePaths
}

func ReplaceContent(path string, findRegex string, replace string) {
	if bytes, err := ioutil.ReadFile(path); err == nil {
		if reg, err := regexp.Compile(findRegex); err == nil {
			newContent := reg.ReplaceAll(bytes, []byte(replace))
			ioutil.WriteFile(path, newContent, DEFAULT_PERM)
		} else {
			log.WithFields(log.Fields{"Error": err}).Error("Invalid Regex")
		}
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Missing File")
	}
}

/**
Reads all Lines from a File.
*/
func ReadAllLines(filePath string) (lines []string) {
	if file, err := os.Open(filePath); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Error Reading File")
	}
	return
}

/**
Scanner must be split on words
*/
func ReadInts(scanner *bufio.Scanner, n int) []int {
	a := make([]int, n)
	for i := 0; i < n && scanner.Scan(); i++ {
		if value, err := strconv.Atoi(scanner.Text()); err == nil {
			a[i] = value
		}
	}
	return a
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RecreateDir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, DIR_DEFAULT_PERM)
}

func ClearDirectory(dirPath string) {
	if fileInfos, err := ioutil.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			os.Remove(filePath)
		}
	}
}
