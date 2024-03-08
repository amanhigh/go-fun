package command

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automation Related Commands",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cluster = args[0]
	},
}

var runOrFocusCmd = &cobra.Command{
	Use:   "run-or-focus [Title]",
	Short: "Focus Window or Start Program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = tools.RunOrFocus(args[0])
		fmt.Println("-->>", err)
		return
	},
}

func init() {
	autoCmd.AddCommand(runOrFocusCmd)
	RootCmd.AddCommand(autoCmd)
}
