package util

import (
	"errors"
	"fmt"
	"reflect"
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

func ReverseArray(input interface{}) {
	n := reflect.ValueOf(input).Len()
	swap := reflect.Swapper(input)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
