package tools

import (
	"os/exec"
	"fmt"
	"strings"
	. "github.com/amanhigh/go-fun/util"
	log "github.com/Sirupsen/logrus"
	"os"
)

func RunCommandPrintError(cmd string) (string) {
	if output, err := runCommand(cmd); err == nil {
		return output
	} else {
		log.WithFields(log.Fields{"CMD": cmd, "Error": err}).Fatal("Error Running Command")
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

func RunNotIf(cmd string, lambda func(output string)) bool {
	if output, err := runCommand(cmd); err != nil {
		lambda(output)
		return true
	}
	return false
}

func runCommand(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}

func LiveCommand(cmd string) {
	command := exec.Command("sh", "-c", cmd)
	/* Connect Command Outputs */
	command.Stdout=os.Stdout
	command.Stderr=os.Stderr

	/* Start Command Wait for Finish */
	command.Start()
	command.Wait()
}
