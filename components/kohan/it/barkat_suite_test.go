package it_test

import (
	"testing"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/models/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var baseURL string

func TestBarkat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Barkat Integration Suite")
}

var _ = BeforeSuite(func() {
	telemetry.InitLogger(config.DefaultLogConfig)
})
