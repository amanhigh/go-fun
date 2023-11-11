package util

import (
	"net/http"

	. "github.com/amanhigh/go-fun/models/common"
	"github.com/go-resty/resty/v2"
)

// Error Proccessor Mapping Http Code to Http Error
func ResponseProcessor(response *resty.Response, restyErr error) (err HttpError) {
	if restyErr != nil {
		err = NewServerError(restyErr)
	} else {
		switch response.StatusCode() {
		case http.StatusBadRequest:
			err = ErrBadRequest
		case http.StatusNotFound:
			err = ErrNotFound
		case http.StatusUnauthorized:
			err = ErrNotAuthorized
		case http.StatusForbidden:
			err = ErrNotAuthenticated
		default:
			err = nil
		}
	}
	return
}
