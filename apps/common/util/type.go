package util

import (
	"errors"
	"fmt"
	"strconv"
)

func IsInt(value string) (err error) {
	if _, err = strconv.Atoi(value); err != nil {
		err = errors.New(fmt.Sprintf("%v is not a Valid Integer", value))
	}
	return
}

func ParseInt(value string) (i int, err error) {
	if i, err = strconv.Atoi(value); err != nil {
		err = errors.New(fmt.Sprintf("%v is not a Valid Integer", value))
	}
	return
}
