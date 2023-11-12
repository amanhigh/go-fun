package it_test

import (
	"time"

	"github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
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
	err        error
	cancelFunc util.CancelFunc
	client     = clients.NewFunAppClient(BASE_URL)
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite", Label(models.GINKGO_INEGRATION_TEST))
}

var _ = BeforeSuite(func() {
	cancelFunc, err = tools.RunBackgroundCommand(COVER_CMD)
	Expect(err).ShouldNot(HaveOccurred())

	//Health Check every 1 second until Healthy or 5 Second Timeout
	for {
		select {
		case <-time.After(5 * time.Second):
			Fail("Unable to Start Funapp")
		case <-time.NewTicker(time.Second).C:
			color.HiBlue("FunApp: Health Check")
			if err = client.AdminService.HealthCheck(); err == nil {
				return
			}
		}
	}
})

var _ = AfterSuite(func() {
	//Send Stop Signal
	err = client.AdminService.Stop()
	Expect(err).ShouldNot(HaveOccurred())
})
