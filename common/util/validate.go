package util

import (
	"fmt"

	"github.com/samber/lo"
)

func Verify(err ...error) (e error) {
	for _, e = range err {
		if e != nil {
			break
		}
	}
	return
}

func ValidateEnumArg(arg string, enum []string) (err error) {
	if ok := lo.Contains(enum, arg); !ok {
		err = fmt.Errorf("%v is not a Valid Argument. Valid Values: %v", arg, enum)
	}
	return
}
