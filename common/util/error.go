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

// ============================================================================
// HTTP Response Processing
// ============================================================================

// statusCodeMap maps HTTP status codes to predefined HttpError instances.
var statusCodeMap = map[int]common.HttpError{
	http.StatusBadRequest:            common.ErrBadRequest,
	http.StatusNotFound:              common.ErrNotFound,
	http.StatusUnauthorized:          common.ErrNotAuthorized,
	http.StatusForbidden:             common.ErrNotAuthenticated,
	http.StatusConflict:              common.ErrEntityExists,
	http.StatusRequestEntityTooLarge: common.ErrPayloadTooLarge,
	http.StatusInternalServerError:   common.ErrInternalServerError,
}

// ResponseProcessor converts a resty response and error into a standardized HttpError.
//
// Processing Priority (highest to lowest):
// 1. Client/Network Error: restyErr != nil → wraps original error as HttpError(500)
// 2. Response Body Error: response.Error() returns HttpError → passes through unchanged
// 3. Status Code Only: statusCodeMap[response.StatusCode()] → maps to predefined errors
//
// Returns nil for unhandled status codes (e.g., 200 OK, 418 Teapot)
func ResponseProcessor(response *resty.Response, restyErr error) common.HttpError {
	// 1. Client/Network errors (connection timeout, DNS failure, etc.)
	if restyErr != nil {
		return common.NewServerError(restyErr)
	}

	// 2. Response body contains structured HttpError (rarely used in current codebase)
	if err, ok := response.Error().(common.HttpError); ok && err.Code() > 0 {
		return err
	}

	// 3. Status code mapping for standard HTTP errors
	return statusCodeMap[response.StatusCode()]
}

// ============================================================================
// Validation Error Processing
// ============================================================================

// ProcessValidationError converts various error types into HttpError.
// Handles: validator errors, JSON errors, strconv errors, and HttpError passthrough.
func ProcessValidationError(err error) common.HttpError {
	if err == nil {
		return unhandledError(err)
	}

	// Try each error handler in order of specificity
	if httpErr := handleValidatorError(err); httpErr != nil {
		return httpErr
	}
	if httpErr := handleHttpError(err); httpErr != nil {
		return httpErr
	}
	if httpErr := handleJSONError(err); httpErr != nil {
		return httpErr
	}
	if httpErr := handleNumericError(err); httpErr != nil {
		return httpErr
	}

	return unhandledError(err)
}

// handleValidatorError processes go-playground/validator errors.
func handleValidatorError(err error) common.HttpError {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) || len(errs) == 0 {
		return nil
	}

	e := errs[0]
	msg := fmt.Sprintf("'%s' with Value '%v' Violates '%s (%s)'", e.Field(), e.Value(), e.Tag(), e.Param())
	return common.NewHttpError(msg, http.StatusBadRequest)
}

// handleHttpError passes through existing HttpError instances.
func handleHttpError(err error) common.HttpError {
	var httpErr common.HttpError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return nil
}

// handleJSONError processes JSON parsing errors (syntax and type mismatch).
func handleJSONError(err error) common.HttpError {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return common.NewHttpError(fmt.Sprintf("Invalid JSON at position %d", syntaxErr.Offset), http.StatusBadRequest)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = typeErr.Struct
		}
		return common.NewHttpError(fmt.Sprintf("Field '%s' expects %s", field, typeErr.Type.String()), http.StatusBadRequest)
	}

	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return common.NewHttpError("Request body cannot be empty or malformed JSON", http.StatusBadRequest)
	}

	return nil
}

// handleNumericError processes strconv parsing errors.
func handleNumericError(err error) common.HttpError {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return common.NewHttpError(fmt.Sprintf("Query parameter '%s' must be numeric", numErr.Func), http.StatusBadRequest)
	}
	return nil
}

// unhandledError logs and returns a generic error for unknown error types.
func unhandledError(err error) common.HttpError {
	log.Warn().
		Str("ActualType", fmt.Sprintf("%T", err)).
		Msg("Failed to convert validation error to HttpError")
	return common.NewHttpError("Invalid validation error format", http.StatusInternalServerError)
}
