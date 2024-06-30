package command

import (
	"github.com/amanhigh/go-fun/components/kohan/tui"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/spf13/cobra"
)

var dariusCmd = &cobra.Command{
	Use:   "darius",
	Short: "Kohan Commander TUI",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		config := config.DariusConfig{
			MakeFileDir: "/home/aman/Projects/go-fun/Kubernetes/services",
		}
		darius, berr := tui.NewDariusInjector(config).BuildApp()
		if berr != nil {
			err = berr
		} else {
			err = darius.Run()
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(dariusCmd)
}
