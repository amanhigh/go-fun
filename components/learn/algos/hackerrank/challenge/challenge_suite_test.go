package challenge_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChallenge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Algo Challenge Suite")
}
