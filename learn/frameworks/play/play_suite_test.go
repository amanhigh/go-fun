package play_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Play Suite")
}
