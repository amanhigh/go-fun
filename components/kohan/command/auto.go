package command

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var servePort = 9010

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

var serveCmd = &cobra.Command{
	Use:   "serve [CapturePath]",
	Short: "Start Kohan Server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Dur("Wait", wait).Str("Screenshots", args[0]).Msg("Starting Kohan Server")
		// TODO: Retry When Disk not Mounted, Watermill Exponential Backoff ?
		autoManager := core.GetKohanInterface().GetAutoManager(wait, args[0])
		server, err := core.GetKohanInterface().GetKohanServer(servePort, args[0], wait)
		if err != nil {
			return fmt.Errorf("failed to build kohan server: %w", err)
		}

		// TODO: #C Should use new Cron Libraries in learn Module
		go autoManager.MonitorInternetConnection(cmd.Context())

		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("Failed to start Kohan server")
			return fmt.Errorf("serve server startup failed: %w", err)
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
	serveCmd.Flags().DurationVarP(&wait, "wait", "w", wait, "OS Wait Interval")
	serveCmd.Flags().IntVarP(&servePort, "port", "p", servePort, "Kohan server port")

	// Commands
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(serveCmd)
	autoCmd.AddCommand(openTickerCmd)
	RootCmd.AddCommand(autoCmd)
}
