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
		PrintRed(fmt.Sprintf("Error Executing Command: CMD:%v ERROR:%+v\n", cmd, err))
	} else {
		PrintWhite(output)
	}
}

func RunIf(cmd string, lambda func()) {
	if _, err := runCommand(cmd); err == nil {
		lambda()
	}
}

func runCommand(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}
