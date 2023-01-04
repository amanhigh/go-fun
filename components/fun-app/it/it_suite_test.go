package it_test

import (
	"os"
	"time"

	"github.com/amanhigh/go-fun/components/fun-app/common"
	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite", Label(models.GINKGO_INEGRATION_TEST))
}

var _ = BeforeSuite(func() {
	// Override Port to Avoid Collision with Default APp
	os.Setenv("PORT", "8085")
	go common.RunFunApp()

	//TODO:Avoid Sleep, Better Way ?
	time.Sleep(time.Second)
})
