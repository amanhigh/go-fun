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
		setLogLevel()
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
	Use:   "monitor [IdleCmd] [CapturePath]",
	Short: "System Monitoring",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Dur("Wait", wait).Dur("Idle", idle).Time("Time", time.Now()).Msg("Monitoring System")
		server := core.NewMonitorServer(args[1])
		go core.MonitorIdle(args[0], wait, idle)
		go server.Start(9010)
		go core.MonitorSubmap()
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
