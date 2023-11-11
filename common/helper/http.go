package helper

import (
	"net/http"

	. "github.com/amanhigh/go-fun/models/common"
	"github.com/go-resty/resty/v2"
)

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
