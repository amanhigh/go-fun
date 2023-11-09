package common

import "errors"

// Standard Http Errors
var BadRequestErr = errors.New("BadRequest")
var NotFoundErr = errors.New("NotFound")
var NotAuthorizedErr = errors.New("NotAuthorized")
var NotAuthenticatedErr = errors.New("NotAuthenticated")
