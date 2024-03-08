package command

import (
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

var winIsFoucsedCmd = &cobra.Command{
	Use:   "is-focused [Title]",
	Short: "Checks if Window is Focused",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = tools.IsWindowFocused(args[0])
		return
	},
}

func init() {
	autoCmd.AddCommand(winIsFoucsedCmd)
	RootCmd.AddCommand(autoCmd)
}
