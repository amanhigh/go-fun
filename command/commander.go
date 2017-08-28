package commander

import (
	"os/exec"
	"fmt"
)

func RunCommand(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Error Executing Command: CMD:%v ERROR:%+v\n", cmd, err)
	}
	return string(output), err
}

func PrintCommand(cmd string) {
	output, _ := RunCommand(cmd)
	PrintWhite(output)
}
