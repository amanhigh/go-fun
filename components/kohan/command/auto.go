package command

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automation Related Commands",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
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

var kohanServerCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Kohan Server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		// TODO: Retry When Disk not Mounted, Watermill Exponential Backoff ?
		osManager := core.GetKohanInterface().GetOSManager()
		server, err := core.GetKohanInterface().GetKohanServer()
		if err != nil {
			return fmt.Errorf("failed to build kohan server: %w", err)
		}

		monitorCtx, stopMonitor := context.WithCancel(cmd.Context())
		defer stopMonitor()
		go osManager.MonitorInternetConnection(monitorCtx)

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
		osManager := core.GetKohanInterface().GetOSManager()
		osManager.TryOpenTicker(cmd.Context(), args[0])
		return
	},
}

func init() {
	autoCmd.AddCommand(runOrFocusCmd)
	autoCmd.AddCommand(kohanServerCmd)
	autoCmd.AddCommand(openTickerCmd)
	RootCmd.AddCommand(autoCmd)
}
