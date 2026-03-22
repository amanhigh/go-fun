package layout_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLayout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Layout Suite")
}
