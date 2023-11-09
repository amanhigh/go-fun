package helper

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Standard Http Errors
var BadRequestErr = errors.New("BadRequest")
var NotFoundErr = errors.New("NotFound")
var NotAuthorizedErr = errors.New("NotAuthorized")
var NotAuthenticatedErr = errors.New("NotAuthenticated")

// Error Proccessor Mapping Http Code to Error
func ResponseProcessor(response *resty.Response, restyErr error) (err error) {
	if restyErr != nil {
		err = restyErr
	} else {
		switch response.StatusCode() {
		case http.StatusBadRequest:
			err = BadRequestErr
		case http.StatusNotFound:
			err = NotFoundErr
		case http.StatusUnauthorized:
			err = NotAuthorizedErr
		case http.StatusForbidden:
			err = NotAuthenticatedErr
		default:
			err = nil
		}
	}
	return
}
