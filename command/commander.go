package commander

import (
	"os/exec"
	"fmt"
	"strings"
)

func RunCommandPrintError(cmd string) (string) {
	if output, err := runCommand(cmd); err == nil {
		return string(output)
	} else {
		PrintRed(err.Error())
		return ""
	}
}

func PrintCommand(cmd string) {
	if output, err := runCommand(cmd); err != nil {
		PrintWhite(output)
		PrintRed(fmt.Sprintf("Error Executing ssh: %v\n CMD:%v\n", err, cmd))
	} else {
		PrintWhite(output)
	}
}

func RunIf(cmd string, lambda func(output string)) bool {
	if output, err := runCommand(cmd); err == nil {
		lambda(output)
		return true
	}
	return false
}

func runCommand(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}
