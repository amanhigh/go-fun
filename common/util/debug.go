package util

import (
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"
)

func DebugControl(flag bool) {
	if flag {
		log.Info().Msg("Enabling Debug Mode")
		if err := os.WriteFile(config.DEBUG_FILE, []byte{}, DEFAULT_PERM); err != nil {
			log.Error().Err(err).Str("Path", config.DEBUG_FILE).Msg("Failed to enable debug mode")
		}
	} else {
		log.Info().Msg("Disabling Debug Mode")
		os.Remove(config.DEBUG_FILE)
	}
	log.Info().Bool("State", IsDebugMode()).Msg("Debug Mode")
}

func IsDebugMode() bool {
	return config.KOHAN_DEBUG || PathExists(config.DEBUG_FILE)
}
