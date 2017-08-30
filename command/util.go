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