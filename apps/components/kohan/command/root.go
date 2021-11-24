package command

import (
	"github.com/amanhigh/go-fun/apps/models/config"
	"os"

	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{}
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&config.KOHAN_DEBUG, "debug", "d", config.KOHAN_DEBUG, "Enable Debug")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		util.PrintRed(err.Error())
		os.Exit(1)
	}
}
