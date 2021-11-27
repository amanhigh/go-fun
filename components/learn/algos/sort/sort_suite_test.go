package sort_test

import (
	"github.com/amanhigh/go-fun/common/util"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSort(t *testing.T) {
	RegisterFailHandler(Fail)
	util.SeedRandom()
	RunSpecs(t, "Algo Sort Suite")
}
