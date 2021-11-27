package util

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func PromptInput(promptText string) string {
	color.Cyan(promptText)
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

func DisplayMenu(msg string, options []string) (int, string) {
	color.Yellow(msg)
	for i, option := range options {
		color.White(fmt.Sprintf("%v. %v", i+1, option))
	}
	input := PromptInput("Please Select an Option.")
	if selection, err := strconv.Atoi(input); err == nil {
		return selection, options[selection-1]
	} else {
		return -1, "INVALID"
	}
}
