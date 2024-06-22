package command

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment Based Commands",
	Args:  cobra.ExactArgs(1),
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Enables/Disables Debug <true/false>",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var enable bool
		enable, err = util.ParseBool(args[0])
		util.DebugControl(enable)
		return
	},
}

func init() {
	envCmd.AddCommand(debugCmd)
	RootCmd.AddCommand(envCmd)
}
