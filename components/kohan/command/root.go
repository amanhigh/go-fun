package command

import (
	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{}
)

func init() {
	//BUG: Connect logger to Debug in Root Cmd
	telemetry.InitLogger(zerolog.InfoLevel)
	RootCmd.PersistentFlags().BoolVarP(&config.KOHAN_DEBUG, "debug", "d", config.KOHAN_DEBUG, "Enable Debug")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Battle Lost")
	}
}

func setLogLevel() {
	if config.KOHAN_DEBUG {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Debug Mode Enabled")
	}
}
