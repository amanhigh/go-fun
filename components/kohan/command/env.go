package command

import (
	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment Based Commands",
	Args:  cobra.ExactArgs(1),
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Enables/Disables Debug",
	Run: func(cmd *cobra.Command, args []string) {
		util2.DebugControl(enable)
	},
}

func init() {
	debugCmd.Flags().BoolVarP(&enable, "enable", "e", false, "Enables Debug Mode")
	debugCmd.MarkPersistentFlagRequired("enable")

	envCmd.AddCommand(debugCmd)
	RootCmd.AddCommand(envCmd)
}
