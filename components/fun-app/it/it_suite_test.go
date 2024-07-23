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
	"github.com/rs/zerolog/log"

	"testing"
)

const (
	COVER_CMD        = "./cover.zsh run"
	DEFAULT_PORT     = "8085"
	DEFAULT_BASE_URL = "http://localhost:" + DEFAULT_PORT
)

var (
	err        error
	spawned    = true
	ctx        = context.Background()
	client     *clients.FunClient
	serviceUrl string
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite", Label(models.GINKGO_INEGRATION))
}

var _ = BeforeSuite(func() {
	// Init Logger
	telemetry.InitLogger(config.DefaultLogConfig)

	baseUrl := os.Getenv("URL")
	if baseUrl == "" {
		serviceUrl = DEFAULT_BASE_URL
	} else {
		serviceUrl = baseUrl
	}

	client = clients.NewFunAppClient(serviceUrl, config.DefaultHttpConfig)
	logger := log.With().Str("URL", serviceUrl).Logger()

	//Run FunApp If not already running
	if err = client.AdminService.HealthCheck(ctx); err == nil {
		logger.Info().Msg("FunApp: Already Running")
		spawned = false
	} else {
		logger.Info().Msg("FunApp: Starting (Not Running)")
		os.Setenv("PORT", DEFAULT_PORT)
		go common.RunFunApp()

		//Health Check every 1 second until Healthy or 5 Second Timeout
		timeout := time.After(10 * time.Second)
		for {
			select {
			case <-timeout:
				Fail("Unable to Start Funapp")
			case <-time.NewTicker(time.Second).C:
				logger.Warn().Msg("FunApp: Starting (Not Running)")
				if err = client.AdminService.HealthCheck(ctx); err == nil {
					return
				}
			}
		}
	}
})

var _ = AfterSuite(func() {
	//Send Stop Signal if Spawned
	if spawned {
		err = client.AdminService.Stop(ctx)
		Expect(err).ShouldNot(HaveOccurred())
	}
})
