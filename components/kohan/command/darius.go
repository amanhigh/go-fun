package command

import (
	"github.com/amanhigh/go-fun/components/kohan/tui"
	"github.com/spf13/cobra"
)

var dariusCmd = &cobra.Command{
	Use:   "darius",
	Short: "Kohan Commander TUI",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		darius := tui.NewDarius()
		err = darius.Run()
		return
	},
}

func init() {
	RootCmd.AddCommand(dariusCmd)
}
