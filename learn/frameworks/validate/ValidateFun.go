package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidateTestStruct struct {
	Name string `validate:"required,aman=preet"`
}

func ValidateFun() {
	validate := validator.New()

	// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
	validate.RegisterValidation("aman", func(fl validator.FieldLevel) bool {
		param := fl.Param()
		return fl.Field().String() == "Aman" && param == "preet"
	})

	fmt.Println(validate.Struct(ValidateTestStruct{Name: "Xyz"}))
	fmt.Println(validate.Struct(ValidateTestStruct{Name: "aman"}))
	fmt.Println(validate.Struct(ValidateTestStruct{Name: "Aman"}))
}
