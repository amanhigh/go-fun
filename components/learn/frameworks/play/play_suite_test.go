package play_test

import (
	"testing"

	"github.com/amanhigh/go-fun/common/telemetry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestPlay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Frameworks Play Suite")
}

var _ = BeforeSuite(func() {
	telemetry.InitLogger(zerolog.InfoLevel)
})
