package it_test

import (
	"os"
	"time"

	"github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/components/fun-app/common"
	"github.com/amanhigh/go-fun/models"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

const (
	COVER_CMD = "./cover.zsh run"
	BASE_URL  = "http://localhost:8085"
)

var (
	err    error
	client = clients.NewFunAppClient(BASE_URL)
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite", Label(models.GINKGO_INEGRATION_TEST))
}

var _ = BeforeSuite(func() {
	//Run FunApp If not already running
	if err = client.AdminService.HealthCheck(); err == nil {
		color.HiGreen("FunApp: Running Already")
	} else {
		os.Setenv("PORT", "8085")
		go common.RunFunApp()

		//Health Check every 1 second until Healthy or 5 Second Timeout
		timeout := time.After(10 * time.Second)
		for {
			select {
			case <-timeout:
				Fail("Unable to Start Funapp")
			case <-time.NewTicker(time.Second).C:
				color.HiBlue("FunApp: Health Check")
				if err = client.AdminService.HealthCheck(); err == nil {
					return
				}
			}
		}
	}
})

var _ = AfterSuite(func() {
	//Send Stop Signal
	err = client.AdminService.Stop()
	Expect(err).ShouldNot(HaveOccurred())
})
