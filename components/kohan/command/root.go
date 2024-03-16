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
	RootCmd.PersistentFlags().BoolVarP(&config.KOHAN_DEBUG, "debug", "d", config.KOHAN_DEBUG, "Enable Debug")
	// BUG: Link to debug flag
	telemetry.InitLogger(zerolog.InfoLevel)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Battle Lost")
	}
}
