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
	msg  string
	code int
}

func NewHttpError(msg string, code int) HttpError {
	return &HttpErrorImpl{msg: msg, code: code}
}

func (self *HttpErrorImpl) Error() string {
	return self.msg
}

func (self *HttpErrorImpl) Code() int {
	return self.code
}

func NewServerError(err error) HttpError {
	return NewHttpError(err.Error(), 500)
}
