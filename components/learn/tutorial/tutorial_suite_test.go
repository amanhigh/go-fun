package tutorial_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTutorial(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tutorial Suite")
}
