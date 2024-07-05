package common

// Standard Http Errors
var ErrBadRequest = NewHttpError("BadRequest", 400)
var ErrNotFound = NewHttpError("NotFound", 404)
var ErrNotAuthorized = NewHttpError("NotAuthorized", 401)
var ErrNotAuthenticated = NewHttpError("NotAuthenticated", 403)

/* Error Reperesenting Http Error and Status Code  */
type HttpError interface {
	error
	Code() int
}

type HttpErrorImpl struct {
	Msg     string `json:"message"`
	ErrCode int    `json:"code"`
	//FIXME: Add Custom Error Codes.
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
	return NewHttpError(err.Error(), 500)
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
