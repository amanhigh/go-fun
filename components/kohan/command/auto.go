package command

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	MonitorServerPort = 9010
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automation Related Commands",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(_ *cobra.Command, args []string) {
		cluster = args[0]
		setLogLevel()
	},
}

var runOrFocusCmd = &cobra.Command{
	Use:   "run-or-focus [Title]",
	Short: "Focus Window or Start Program",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) (err error) {
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

		// BUG: Retry When Disk not Mounter, Watermill Exponential Backoff ?
		autoManager := core.GetKohanInterface().GetAutoManager(wait, args[0])
		server := core.NewKohanServer(args[0], autoManager, nil)

		go autoManager.MonitorInternetConnection(cmd.Context())

		shutdown := util.NewGracefulShutdown()
		if err := server.Start(MonitorServerPort, shutdown); err != nil {
			log.Error().Err(err).Msg("Failed to start monitor server")
			return fmt.Errorf("monitor server startup failed: %w", err)
		}
		return
	},
}

var openTickerCmd = &cobra.Command{
	Use:   "open-ticker [Ticker]",
	Short: "Opens Ticker",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		autoManager := core.GetKohanInterface().GetAutoManager(wait, "")
		autoManager.TryOpenTicker(cmd.Context(), args[0])
		return
	},
}

func init() {
	// Flags
	monitorCmd.Flags().DurationVarP(&wait, "wait", "w", wait, "Monitoring Wait Interval")

	// Commands
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(monitorCmd)
	autoCmd.AddCommand(openTickerCmd)
	RootCmd.AddCommand(autoCmd)
}
