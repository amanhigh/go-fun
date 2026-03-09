package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Standard Http Errors
var ErrBadRequest = NewHttpError("BadRequest", http.StatusBadRequest)
var ErrNotFound = NewHttpError("NotFound", http.StatusNotFound)
var ErrNotAuthorized = NewHttpError("NotAuthorized", http.StatusUnauthorized)
var ErrNotAuthenticated = NewHttpError("NotAuthenticated", http.StatusForbidden)
var ErrEntityExists = NewHttpError("EntityExists", http.StatusConflict)
var ErrPayloadTooLarge = NewHttpError("PayloadTooLarge", http.StatusRequestEntityTooLarge)
var ErrInternalServerError = NewHttpError("InternalServerError", http.StatusInternalServerError)

/* Error Reperesenting Http Error and Status Code  */
type HttpError interface {
	error
	Code() int
}

/* Error representing field-specific validation errors */
type FieldHttpError interface {
	HttpError
	Field() string
}

type HttpErrorImpl struct {
	Msg     string `json:"message"`
	ErrCode int    `json:"code"`
}

type FieldHttpErrorImpl struct {
	HttpErrorImpl
	fieldName string `json:"-"` // Field name for validation errors (not serialized directly)
}

func NewHttpError(msg string, code int) HttpError {
	return &HttpErrorImpl{Msg: msg, ErrCode: code}
}

// NewFieldHttpError creates a FieldHttpError with a field name for validation errors.
// Always uses 400 Bad Request as the status code for field validation failures.
func NewFieldHttpError(field, msg string) FieldHttpError {
	return &FieldHttpErrorImpl{
		HttpErrorImpl: HttpErrorImpl{Msg: msg, ErrCode: http.StatusBadRequest},
		fieldName:     field,
	}
}

func (e *HttpErrorImpl) Error() string {
	return e.Msg
}

func (e *HttpErrorImpl) Code() int {
	return e.ErrCode
}

func (e *FieldHttpErrorImpl) Field() string {
	return e.fieldName
}

// MarshalJSON implements json.Marshaler to output JSend envelope format.
// 4xx errors: { "status": "fail", "data": { "message": "error message" } }
// 5xx errors: { "status": "error", "message": "...", "code": ... }
func (e *HttpErrorImpl) MarshalJSON() ([]byte, error) {
	if e.ErrCode >= http.StatusInternalServerError {
		// 5xx: JSend "error" format
		data, err := json.Marshal(map[string]any{
			"status":  EnvelopeError,
			"message": e.Msg,
			"code":    e.ErrCode,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal server error (code %d): %w", e.ErrCode, err)
		}
		return data, nil
	}
	// 4xx: JSend "fail" format - use message key for regular HttpError
	data, err := json.Marshal(map[string]any{
		"status": EnvelopeFail,
		"data": map[string]string{
			"message": e.Msg,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client error (code %d): %w", e.ErrCode, err)
	}
	return data, nil
}

// MarshalJSON for FieldHttpErrorImpl uses field name as key for 4xx errors
func (e *FieldHttpErrorImpl) MarshalJSON() ([]byte, error) {
	// FieldHttpError should only be used for 4xx errors, so always use "fail" format
	fieldName := e.Field()
	if fieldName == "" {
		fieldName = "message" // fallback for non-field-specific errors
	}
	data, err := json.Marshal(map[string]any{
		"status": EnvelopeFail,
		"data": map[string]string{
			fieldName: e.Msg,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field validation error (field: %s, code %d): %w", fieldName, e.ErrCode, err)
	}
	return data, nil
}

func NewServerError(err error) HttpError {
	return NewHttpError(err.Error(), http.StatusInternalServerError)
}

const (
	// Base API routes
	APIV1 = "/v1"

	// Monitor routes
	MonitorBase = APIV1 + "/monitor"
)

type Pagination struct {
	Offset int `form:"offset,default=0" binding:"min=0"`
	Limit  int `form:"limit,default=20" binding:"min=1,max=100"`
}

type Sort struct {
	SortBy string `form:"sort_by" binding:"omitempty,eq=name|eq=age|eq=gender"`
	Order  string `form:"order,default=asc" binding:"omitempty,eq=asc|eq=desc"`
}

type PaginatedResponse struct {
	Total  int64 `json:"total"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}
