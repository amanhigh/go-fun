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
	COVER_CMD        = "./cover.zsh run"
	TEST_PORT        = "8085"
	DEFAULT_BASE_URL = "http://localhost:"
)

var (
	err        error
	spawned    = true
	port       = os.Getenv("PORT")
	baseUrl    = os.Getenv("URL")
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
	telemetry.InitLogger(zerolog.InfoLevel)

	// Check if Port is Set
	if os.Getenv("PORT") == "" {
		port = TEST_PORT
	}

	if baseUrl == "" {
		baseUrl = DEFAULT_BASE_URL
	}

	serviceUrl = baseUrl + port
	client = clients.NewFunAppClient(serviceUrl, config.DefaultHttpConfig)

	//Run FunApp If not already running
	if err = client.AdminService.HealthCheck(ctx); err == nil {
		log.Info().Str("URL", serviceUrl).Str("Port", port).Msg("FunApp: Running Already")
		spawned = false
	} else {
		os.Setenv("PORT", port)
		go common.RunFunApp()

		//Health Check every 1 second until Healthy or 5 Second Timeout
		timeout := time.After(10 * time.Second)
		for {
			select {
			case <-timeout:
				Fail("Unable to Start Funapp")
			case <-time.NewTicker(time.Second).C:
				log.Warn().Str("Port", port).Msg("FunApp: Starting (Not Running)")
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
