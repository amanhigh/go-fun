package cracking_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCracking(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Algo Cracking Suite")
}
