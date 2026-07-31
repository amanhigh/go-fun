package handlers_test

import (
	"testing"

	funcommon "github.com/amanhigh/go-fun/components/fun-app/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	validatorEngine, ok := binding.Validator.Engine().(*validator.Validate)
	Expect(ok).To(BeTrue())
	Expect(validatorEngine.RegisterValidation("name", funcommon.NameValidator)).To(Succeed())
})

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}
