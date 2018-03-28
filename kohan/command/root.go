package command

import (
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		util.PrintRed(err.Error())
		os.Exit(1)
	} else {
		util.PrintGreen("Command Successful")
	}
}