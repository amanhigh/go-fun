package it_test

import (
	"context"
	"os"
	"time"

	"github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/components/fun-app/common"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"testing"
)

const (
	COVER_CMD = "./cover.zsh run"
	TEST_PORT = "8085"
	// HACK: Make Base URL configurable
	BASE_URL = "http://localhost:" + TEST_PORT
)

var (
	err    error
	client = clients.NewFunAppClient(BASE_URL, config.DefaultHttpConfig)
	ctx    = context.Background()
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite", Label(models.GINKGO_INEGRATION))
}

var _ = BeforeSuite(func() {
	// Init Logger
	telemetry.InitLogger(zerolog.InfoLevel)

	//Run FunApp If not already running
	if err = client.AdminService.HealthCheck(ctx); err == nil {
		log.Info().Str("Port", os.Getenv("PORT")).Msg("FunApp: Running Already")
	} else {
		os.Setenv("PORT", TEST_PORT)
		go common.RunFunApp()

		//Health Check every 1 second until Healthy or 5 Second Timeout
		timeout := time.After(10 * time.Second)
		for {
			select {
			case <-timeout:
				Fail("Unable to Start Funapp")
			case <-time.NewTicker(time.Second).C:
				log.Warn().Str("Port", os.Getenv("PORT")).Msg("FunApp: Starting (Not Running)")
				if err = client.AdminService.HealthCheck(ctx); err == nil {
					return
				}
			}
		}
	}
})

var _ = AfterSuite(func() {
	//Send Stop Signal
	err = client.AdminService.Stop(ctx)
	Expect(err).ShouldNot(HaveOccurred())
})
