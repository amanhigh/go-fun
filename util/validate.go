package util

import (
	"errors"
	"fmt"
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
	if ok := SliceContains(arg, enum); !ok {
		err = errors.New(fmt.Sprintf("%v is not a Valid Argument. Valid Values: %v", arg, enum))
	}
	return
}
