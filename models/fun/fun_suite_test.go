package fun_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Fun Suite")
}
