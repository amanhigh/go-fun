package util

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type CancelFunc func() (err error)

const decimalBase = 10 // Base for decimal place calculations

func IsInt(value string) (err error) {
	if _, err = strconv.Atoi(value); err != nil {
		err = fmt.Errorf("%v is not a Valid Integer", value)
	}
	return
}

func ParseInt(value string) (i int, err error) {
	if i, err = strconv.Atoi(value); err != nil {
		err = fmt.Errorf("%v is not a Valid Integer", value)
	}
	return
}

func ParseBool(value string) (b bool, err error) {
	if b, err = strconv.ParseBool(value); err != nil {
		err = fmt.Errorf("%v is not a Valid Boolean", value)
	}
	return
}

func ReverseArray(input any) {
	n := reflect.ValueOf(input).Len()
	swap := reflect.Swapper(input)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func ParseFloat(value string) float64 {
	const defaultValue = -1.0

	floatVal := strings.TrimSpace(value)
	if floatVal == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseFloat(floatVal, 64)
	if err != nil {
		log.Error().Str("Value", value).Err(err).Msg("Error Parsing Float value")
		return defaultValue
	}

	return parsed
}

// RoundToDecimals rounds a float64 value to specified decimal places
// For example: RoundToDecimals(0.3899999999999999, 2) returns 0.39
func RoundToDecimals(value float64, decimals int) float64 {
	multiplier := math.Pow(decimalBase, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}
