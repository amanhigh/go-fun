package tools

import (
	"fmt"
	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func RunCommandPrintError(cmd string) string {
	if output, err := runCommand(cmd); err == nil {
		return output
	} else {
		log.WithFields(log.Fields{"CMD": cmd, "Error": err}).Error("Error Running Command")
		return ""
	}
}

func RunAsyncCommand(heading string, cmd string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		output, err := runCommand(cmd)
		color.Cyan(heading)
		if err == nil {
			color.White(output)
		} else {
			color.White(err.Error())
		}
		wg.Done()
	}()
}

func RunCommandIgnoreError(cmd string) string {
	output, _ := runCommand(cmd)
	return output
}

func PrintCommand(cmd string) {
	if output, err := runCommand(cmd); err != nil {
		color.White(output)
		color.Red(fmt.Sprintf("Error Executing: %v\n CMD:%v\n", err, cmd))
	} else {
		color.White(output)
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
	if util2.IsDebugMode() {
		color.Magenta(cmd)
	}
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}

func LiveCommand(cmd string) {
	command := exec.Command("sh", "-c", cmd)
	if util2.IsDebugMode() {
		color.Magenta(cmd)
	}
	/* Connect Command Outputs */
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	/* Start Command Wait for Finish */
	command.Start()
	command.Wait()
}
