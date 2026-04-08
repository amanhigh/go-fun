package widgets_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWidgets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Widget Suite")
}
