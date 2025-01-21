package command

import (
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
	Use:   "monitor [CapturePath]",
	Short: "System Monitoring",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Dur("Wait", wait).Str("Screenshots", args[0]).Msg("Monitoring Systems")
		// XXX: Move to Injector once Created
		server := core.NewMonitorServer(args[0])
		core.MonitorInternetConnection(wait)
		go server.Start(9010)
		core.MonitorSubmap()
		return
	},
}

var openTickerCmd = &cobra.Command{
	Use:   "open-ticker [Ticker]",
	Short: "Opens Ticker",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		core.TryOpenTicker(args[0])
		return
	},
}

func init() {
	//Flags
	monitorCmd.Flags().DurationVarP(&wait, "wait", "w", wait, "Monitoring Wait Interval")

	//Commands
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(monitorCmd)
	autoCmd.AddCommand(openTickerCmd)
	RootCmd.AddCommand(autoCmd)
}
