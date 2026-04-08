package ui_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSpinner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spinner Suite")
}
