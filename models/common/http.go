package common

import "net/http"

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

type HttpErrorImpl struct {
	Msg     string `json:"message"`
	ErrCode int    `json:"code"`
}

func NewHttpError(msg string, code int) HttpError {
	return &HttpErrorImpl{Msg: msg, ErrCode: code}
}

func (e *HttpErrorImpl) Error() string {
	return e.Msg
}

func (e *HttpErrorImpl) Code() int {
	return e.ErrCode
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
