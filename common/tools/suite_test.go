package tools_test

import (
	"testing"

	"golang.design/x/clipboard"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTools(t *testing.T) {
	RegisterFailHandler(Fail)
	clipboard.Init()
	RunSpecs(t, "Tools Suite")
}
