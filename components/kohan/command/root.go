package command

import (
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"

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
