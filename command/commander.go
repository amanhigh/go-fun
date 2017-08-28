package commander

import (
	"os/exec"
	"fmt"
)

func RunCommand(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	return string(output), err
}

func PrintCommand(cmd string) {
	if output, err := RunCommand(cmd); err != nil {
		fmt.Printf("Error Executing Command: CMD:%v ERROR:%+v\n", cmd, err)
	} else {
		PrintWhite(output)
	}
}

func RunIf(cmd string, lambda func()) {
	if _, err := RunCommand(cmd); err == nil {
		lambda()
	}
}
