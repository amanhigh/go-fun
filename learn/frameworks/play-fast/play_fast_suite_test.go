package play_fast_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlayFast(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PlayFast Suite")
}
