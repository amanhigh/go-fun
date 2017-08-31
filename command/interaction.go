package commander

import (
	"bufio"
	"os"
	"strings"
)

func PromptInput(promptText string) string {
	PrintWhite(promptText)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Confirm(msg string, runLamda func()) {
	input := PromptInput(msg + " Y/y to Continue")
	if strings.EqualFold(input, "Y") {
		runLamda()
	}
}

func NoConfirm(msg string, runLamda func()) {
	input := PromptInput(msg + " N/n to Abort")
	if !strings.EqualFold(input, "N") {
		runLamda()
	}
}
