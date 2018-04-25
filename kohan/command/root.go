package command

import (
	"os"

	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		util.PrintRed(err.Error())
		os.Exit(1)
	}
}
