package command

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/fatih/color"
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
		return
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor [IdleCmd] [ClipPath]",
	Short: "System Monitoring",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		color.Green("Monitoring System: Wait -> %v, Idle -> %v, Now -> %v", wait, idle, time.Now())
		go core.MonitorIdle(args[0], wait, idle)
		go core.MonitorClipboard(ctx, args[1])
		core.MonitorInternetConnection(wait)
		return
	},
}

func init() {
	//Flags
	monitorCmd.Flags().DurationVarP(&wait, "wait", "w", wait, "Monitoring Wait Interval")
	monitorCmd.Flags().DurationVarP(&idle, "idle", "i", idle, "Idle Time")

	//Commands
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(monitorCmd)
	RootCmd.AddCommand(autoCmd)
}
