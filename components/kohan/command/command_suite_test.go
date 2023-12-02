package command_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommand(t *testing.T) {
	// https://github.com/onsi/ginkgo/issues/285
	// Trim os.Args to only the first arg
	os.Args = os.Args[:1] // trim to only the first arg, which is the command itself

	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Suite")
}
