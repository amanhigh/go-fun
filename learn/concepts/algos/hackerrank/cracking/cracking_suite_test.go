package cracking_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCracking(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cracking Suite")
}
