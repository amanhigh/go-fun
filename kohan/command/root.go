package command

import (
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.PrintRed(err.Error())
		os.Exit(1)
	} else {
		util.PrintGreen("Command Successful")
	}
}