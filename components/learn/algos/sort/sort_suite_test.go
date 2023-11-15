package sort_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSort(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Algo Sort Suite")
}
