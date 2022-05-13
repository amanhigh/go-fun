package util

import (
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
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
	if ok := funk.Contains(enum, arg); !ok {
		err = errors.New(fmt.Sprintf("%v is not a Valid Argument. Valid Values: %v", arg, enum))
	}
	return
}
