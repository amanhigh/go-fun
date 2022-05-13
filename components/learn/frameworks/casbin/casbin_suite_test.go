package casbin_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCasbin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Casbin Suite")
}
