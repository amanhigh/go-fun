package learn_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLearn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Learn Suite")
}
