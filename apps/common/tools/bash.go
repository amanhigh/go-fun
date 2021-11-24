package tools

import (
	"fmt"
	"github.com/amanhigh/go-fun/apps/components/kohan/core"
	"os"
	"os/exec"
	"strings"
	"sync"

	. "github.com/amanhigh/go-fun/util"
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
		PrintSkyBlue(heading)
		if err == nil {
			PrintWhite(output)
		} else {
			PrintWhite(err.Error())
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
		PrintWhite(output)
		PrintRed(fmt.Sprintf("Error Executing: %v\n CMD:%v\n", err, cmd))
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
	if core.IsDebugMode() {
		PrintPink(cmd)
	}
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}

func LiveCommand(cmd string) {
	command := exec.Command("sh", "-c", cmd)
	if core.IsDebugMode() {
		PrintPink(cmd)
	}
	/* Connect Command Outputs */
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	/* Start Command Wait for Finish */
	command.Start()
	command.Wait()
}
