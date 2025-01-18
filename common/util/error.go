package util

import (
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	. "github.com/amanhigh/go-fun/models/common"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
)

// ResponseProcessor processes the resty response and error to generate an HttpError.
//
// Parameters:
// - response: a pointer to resty.Response
// - restyErr: an error
// Return type(s):
// - err: an HttpError
func ResponseProcessor(response *resty.Response, restyErr error) (err HttpError) {
	var ok bool
	if restyErr != nil {
		//Rest Client Error hence No Respones
		err = NewServerError(restyErr)
	} else if err, ok = response.Error().(HttpError); ok && err.Code() > 0 {
		//If Error is Http Error & has Data, Use directly.
	} else {
		//Incase we have No Error Honor Status Codes of Http
		switch response.StatusCode() {
		case http.StatusBadRequest:
			err = ErrBadRequest
		case http.StatusNotFound:
			err = ErrNotFound
		case http.StatusUnauthorized:
			err = ErrNotAuthorized
		case http.StatusForbidden:
			err = ErrNotAuthenticated
		case http.StatusConflict:
			err = ErrEntityExists
		case http.StatusInternalServerError:
			// TODO: Error From Response
			err = ErrInternalServerError
		default:
			err = nil
		}
	}
	return
}

func ProcessValidationError(validationErr error) (err HttpError) {
	if errs, ok := validationErr.(validator.ValidationErrors); ok {
		for _, e := range errs {
			err = common.NewHttpError(fmt.Sprintf("'%s' with Value '%v' Violates '%s (%s)'", e.Field(), e.Value(), e.Tag(), e.Param()), http.StatusBadRequest)
			break
		}
	} else {
		err = validationErr.(HttpError)
	}
	return
}
