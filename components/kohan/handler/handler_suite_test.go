package handler_test

import (
	"testing"

	"github.com/amanhigh/go-fun/components/kohan/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.design/x/clipboard"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	core.RegisterJournalValidators()
	clipboard.Init()
	RunSpecs(t, "Handler Suite")
}
