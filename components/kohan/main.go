package main

import (
	"github.com/amanhigh/go-fun/components/kohan/command"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.NewKohanConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Kohan config")
		return
	}

	core.SetupKohanInjector(cfg)
	command.Execute()
}
