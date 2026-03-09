package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// ResponseProcessor processes the resty response and error to generate an HttpError.
//
// Parameters:
// - response: a pointer to resty.Response
// - restyErr: an error
// Return type(s):
// - err: an HttpError
func handleStatusCode(statusCode int) common.HttpError {
	switch statusCode {
	case http.StatusBadRequest:
		return common.ErrBadRequest
	case http.StatusNotFound:
		return common.ErrNotFound
	case http.StatusUnauthorized:
		return common.ErrNotAuthorized
	case http.StatusForbidden:
		return common.ErrNotAuthenticated
	case http.StatusConflict:
		return common.ErrEntityExists
	case http.StatusInternalServerError:
		return common.ErrInternalServerError
	default:
		return nil
	}
}

func ResponseProcessor(response *resty.Response, restyErr error) common.HttpError {
	if restyErr != nil {
		// Rest Client Error hence No Respones
		return common.NewServerError(restyErr)
	}

	// If Error is Http Error & has Data, Use directly.
	if err, ok := response.Error().(common.HttpError); ok && err.Code() > 0 {
		return err
	}

	// TASK: Error From Response
	// Incase we have No Error Honor Status Codes of Http
	return handleStatusCode(response.StatusCode())
}

// FIXME: #A Review this file completely.

func handleValidationError(err error) common.HttpError {
	if httpErr := handleStructuralErrors(err); httpErr != nil {
		return httpErr
	}
	return handleFormatErrors(err)
}

func handleStructuralErrors(err error) common.HttpError {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			return common.NewHttpError(fmt.Sprintf("'%s' with Value '%v' Violates '%s (%s)'", e.Field(), e.Value(), e.Tag(), e.Param()), http.StatusBadRequest)
		}
	}

	var httpErr common.HttpError
	if errors.As(err, &httpErr) {
		return httpErr
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return common.NewHttpError(fmt.Sprintf("Invalid JSON at position %d", syntaxErr.Offset), http.StatusBadRequest)
	}

	return nil
}

func handleFormatErrors(err error) common.HttpError {
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		field := unmarshalTypeErr.Field
		if field == "" {
			field = unmarshalTypeErr.Struct
		}
		return common.NewHttpError(fmt.Sprintf("Field '%s' expects %s", field, unmarshalTypeErr.Type.String()), http.StatusBadRequest)
	}

	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return common.NewHttpError("Request body cannot be empty or malformed JSON", http.StatusBadRequest)
	}

	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return common.NewHttpError(fmt.Sprintf("Query parameter '%s' must be numeric", numErr.Func), http.StatusBadRequest)
	}

	return nil
}

func ProcessValidationError(validationErr error) common.HttpError {
	if err := handleValidationError(validationErr); err != nil {
		return err
	}

	log.Warn().
		Str("ActualType", fmt.Sprintf("%T", validationErr)).
		Msg("Failed to convert validation error to HttpError")
	return common.NewHttpError("Invalid validation error format", http.StatusInternalServerError)
}
