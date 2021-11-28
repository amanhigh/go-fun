package play_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Frameworks Play Suite")
}
