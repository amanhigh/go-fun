package command

import (
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/rs/zerolog/log"
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
	Use:   "monitor [IdleCmd]",
	Short: "System Monitoring",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Dur("Wait", wait).Dur("Idle", idle).Time("Time", time.Now()).Msg("Monitoring System")
		go core.MonitorIdle(args[0], wait, idle)
		// go core.MonitorClipboard(ctx, args[1])
		go core.MonitorSubmap()
		core.MonitorInternetConnection(wait)
		return
	},
}

var processTicker = &cobra.Command{
	Use:   "ticker [Ticker] [CapturePath]",
	Short: "Process Ticker",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		core.ProcessTicker(args[0], args[1])
	},
}

func init() {
	//Flags
	monitorCmd.Flags().DurationVarP(&wait, "wait", "w", wait, "Monitoring Wait Interval")
	monitorCmd.Flags().DurationVarP(&idle, "idle", "i", idle, "Idle Time")

	//Commands
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(monitorCmd)
	autoCmd.AddCommand(processTicker)
	RootCmd.AddCommand(autoCmd)
}
