package util

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

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

func ParseInt(value string) (i int,err error) {
	if i, err = strconv.Atoi(value); err != nil {
		err = errors.New(fmt.Sprintf("%v is not a Valid Integer", value))
	}
	return
}
