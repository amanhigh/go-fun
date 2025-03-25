package util

import (
	"errors"
	"fmt"
	"net/http"

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

func ProcessValidationError(validationErr error) (err common.HttpError) {
	var errs validator.ValidationErrors
	if errors.As(validationErr, &errs) {
		for _, e := range errs {
			err = common.NewHttpError(fmt.Sprintf("'%s' with Value '%v' Violates '%s (%s)'", e.Field(), e.Value(), e.Tag(), e.Param()), http.StatusBadRequest)
			break
		}
	} else {
		var httpErr common.HttpError
		if errors.As(validationErr, &httpErr) {
			return httpErr
		}
		log.Warn().
			Str("ActualType", fmt.Sprintf("%T", validationErr)).
			Msg("Failed to convert validation error to HttpError")
		return common.NewHttpError("Invalid validation error format", http.StatusInternalServerError)
	}
	return
}
