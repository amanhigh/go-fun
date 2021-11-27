package util

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
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

func ParseFloat(value string) (result float64) {
	var err error
	result = -1

	floatVal := strings.TrimSpace(value)
	if floatVal != "" {
		if result, err = strconv.ParseFloat(floatVal, 64); err != nil {
			log.WithFields(log.Fields{"Value": value, "Error": err}).Error("Error Parsing Float value")
		}
	}

	return
}
