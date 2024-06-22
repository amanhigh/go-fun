package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

const DEFAULT_PERM = os.FileMode(0644)     //Owner RW,Group R,Other R
const DIR_DEFAULT_PERM = os.FileMode(0755) //Owner RWX,Group RX,Other RX
/*
	Helpfull File Related Cheatsheet
	https://www.devdungeon.com/content/working-files-go#read_quick
*/

func OpenOrCreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, DEFAULT_PERM)
}

func AppendFile(path string, content string) {
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600); err == nil {
		defer f.Close()
		if _, err = f.WriteString(content); err != nil {
			log.Error().Str("File", path).Err(err).Msg("Error Appending Content to File")
		}
	} else {
		log.Error().Str("File", path).Err(err).Msg("Error Opening File for Append")
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

func FindReplaceFile(filePath string, find string, replace string) (err error) {
	var compile *regexp.Regexp
	var fileBytes []byte
	if fileBytes, err = os.ReadFile(filePath); err == nil {
		if compile, err = regexp.Compile(find); err == nil {
			replacedBytes := compile.ReplaceAll(fileBytes, []byte(replace))
			os.WriteFile(filePath, replacedBytes, DEFAULT_PERM)
		}
	}
	return
}

func PrintFile(title string, filepath string) {
	log.Info().Str("File", filepath).Msg("File Contents")
	fmt.Println(strings.Join(ReadAllLines(filepath), "\n"))
}

func ListFiles(dirPath string) []string {
	var filePaths []string
	if fileInfos, err := os.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			filePaths = append(filePaths, filePath)
		}
	} else {
		log.Error().Str("Directory", dirPath).Err(err).Msg("Error Reading Directory")
	}
	return filePaths
}

func ReplaceContent(path string, findRegex string, replace string) {
	if bytes, err := os.ReadFile(path); err == nil {
		if reg, err := regexp.Compile(findRegex); err == nil {
			newContent := reg.ReplaceAll(bytes, []byte(replace))
			os.WriteFile(path, newContent, DEFAULT_PERM)
		} else {
			log.Error().Str("Regex", findRegex).Err(err).Msg("Invalid Regex")
		}
	} else {
		log.Error().Str("File", path).Err(err).Msg("Error Reading File")
	}
}

/*
*
Reads all Lines from a File.
*/
func ReadAllLines(filePath string) (lines []string) {
	if file, err := os.Open(filePath); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		log.Error().Str("File", filePath).Err(err).Msg("Error Reading File")
	}
	return
}

func WriteLines(filePath string, lines []string) error {
	content := strings.Join(lines, "\n")
	return os.WriteFile(filePath, []byte(content), DEFAULT_PERM)
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
	if fileInfos, err := os.ReadDir(dirPath); err == nil {
		for _, info := range fileInfos {
			filePath := fmt.Sprintf("%v/%v", dirPath, info.Name())
			os.Remove(filePath)
		}
	}
}
