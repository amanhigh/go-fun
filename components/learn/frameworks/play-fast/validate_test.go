package play_fast

import (
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ValidateTestStruct struct {
	Name string `validate:"required,aman=preet"`
}

var _ = Describe("Validate", func() {
	var (
		validate *validator.Validate
		err      error
	)

	BeforeEach(func() {
		validate = validator.New()
	})

	It("should build", func() {
		Expect(validate).To(Not(BeNil()))
	})

	Context("Custom Validator", func() {
		BeforeEach(func() {
			// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
			err = validate.RegisterValidation("aman", func(fl validator.FieldLevel) bool {
				param := fl.Param()
				return fl.Field().String() == "Aman" && param == "preet"
			})
			Expect(err).To(BeNil())
		})

		It("should pass", func() {
			Expect(validate.Struct(ValidateTestStruct{Name: "Aman"})).To(BeNil())
		})

		It("should fail", func() {
			Expect(validate.Struct(ValidateTestStruct{Name: "Xyz"})).To(Not(BeNil()))
		})

		It("should be case sensitive", func() {
			Expect(validate.Struct(ValidateTestStruct{Name: "aman"})).To(Not(BeNil()))
		})

	})

})
