package commander

import "fmt"

func PrintRed(text string) {
	PrintColor(31, text)
}

func PrintGreen(text string) {
	PrintColor(32, text)
}

func PrintYellow(text string) {
	PrintColor(33, text)
}

func PrintBlue(text string) {
	PrintColor(34, text)
}

func PrintWhite(text string) {
	PrintColor(28, text)
}

func PrintColor(code int, text string) {
	fmt.Printf("\033[1;%vm %v \033[0m \n", code, text)
}
