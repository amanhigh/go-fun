package command

import (
	"github.com/amanhigh/go-fun/apps/models/config"
	"github.com/fatih/color"
	"os"

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
		color.Red(err.Error())
		os.Exit(1)
	}
}
