package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

func RunCommandPrintError(cmd string) string {
	if output, err := runCommand(cmd); err == nil {
		return output
	} else {
		log.Error().Str("CMD", cmd).Err(err).Msg("Error Running Command")
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

// RunBackgroundProcess runs a background command and returns an error and a cancel function.
//
// The command parameter specifies the command to be executed.
// The function returns an error if the command fails to start.
// The cancel function can be used to kill the command and all of its child processes.
func RunBackgroundProcess(command string) (cancel util.CancelFunc, err error) {
	cmd := exec.Command("sh", "-c", command)
	// Ensure any Child Process are Killed As Well.
	// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()

	cancel = func() (err error) {
		//Kill Command with Subprocess
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		log.Debug().Str("CMD", command).Int("PID", cmd.Process.Pid).Msg("Killing Background Process")
		return
	}
	return
}

func RunCommandIgnoreError(cmd string) string {
	output, _ := runCommand(cmd)
	return output
}

func PrintCommand(cmd string) {
	if output, err := runCommand(cmd); err != nil {
		color.White(output)
		log.Error().Str("CMD", cmd).Err(err).Msg("Error Running Command")
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

func RunProcess(cmd string) (output string, err error) {
	output, err = script.Exec(fmt.Sprintf("sh -c '%v'", cmd)).String()
	return
}

func runCommand(cmd string) (string, error) {
	if util.IsDebugMode() {
		log.Debug().Str("CMD", cmd).Msg("Running Command")
	}
	output, err := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimSpace(string(output)), err
}

func LiveCommand(cmd string) {
	command := exec.Command("sh", "-c", cmd)
	if util.IsDebugMode() {
		log.Debug().Str("CMD", cmd).Msg("Running Command")
	}
	/* Connect Command Outputs */
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	/* Start Command Wait for Finish */
	if err := command.Start(); err != nil {
		log.Warn().
			Str("CMD", cmd).
			Err(err).
			Msg("Failed to start command")
		return
	}

	if err := command.Wait(); err != nil {
		log.Warn().
			Str("CMD", cmd).
			Err(err).
			Msg("Command execution failed")
	}
}
