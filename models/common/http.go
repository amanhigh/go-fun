package common

import "net/http"

// Standard Http Errors
var ErrBadRequest = NewHttpError("BadRequest", http.StatusBadRequest)
var ErrNotFound = NewHttpError("NotFound", http.StatusNotFound)
var ErrNotAuthorized = NewHttpError("NotAuthorized", http.StatusUnauthorized)
var ErrNotAuthenticated = NewHttpError("NotAuthenticated", http.StatusForbidden)
var ErrEntityExists = NewHttpError("EntityExists", http.StatusConflict)
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

func (self *HttpErrorImpl) Error() string {
	return self.Msg
}

func (self *HttpErrorImpl) Code() int {
	return self.ErrCode
}

func NewServerError(err error) HttpError {
	return NewHttpError(err.Error(), http.StatusInternalServerError)
}

type Pagination struct {
	Offset int `form:"offset" binding:"min=0"`
	Limit  int `form:"limit" binding:"required,min=1,max=10"`
}

type Sort struct {
	SortBy string `form:"sort_by" binding:"omitempty,eq=name|eq=age|eq=gender"`
	Order  string `form:"order" binding:"omitempty,eq=asc|eq=desc"`
}

type PaginatedResponse struct {
	Total int64 `json:"total"`
}
