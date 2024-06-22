package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/amanhigh/go-fun/common/telemetry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var logActual = new(bytes.Buffer)

func TestCommand(t *testing.T) {
	// https://github.com/onsi/ginkgo/issues/285
	// Trim os.Args to only the first arg
	os.Args = os.Args[:1] // trim to only the first arg, which is the command itself

	telemetry.InitTestLogger(logActual)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Suite")
}
