package handler_test

import (
	"testing"

	"github.com/amanhigh/go-fun/components/kohan/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	core.RegisterJournalValidators()
	RunSpecs(t, "Handler Suite")
}
