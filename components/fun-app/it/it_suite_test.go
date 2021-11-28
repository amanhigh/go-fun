package it_test

import (
	"github.com/amanhigh/go-fun/components/fun-app/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"

	"testing"
)

func TestIt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunApp Suite")
}

var _ = BeforeSuite(func() {
	go common.RunFunApp()
	//TODO:Avoid Sleep, Better Way ?
	time.Sleep(time.Second)
})
