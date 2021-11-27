package sort_test

import (
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSort(t *testing.T) {
	RegisterFailHandler(Fail)
	util2.SeedRandom()
	RunSpecs(t, "Algo Sort Suite")
}
