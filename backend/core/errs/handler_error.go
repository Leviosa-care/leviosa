package errs

import (
	"errors"
)

// Standard HTTP error messages / categories
var (
	ErrNotFoundHTTP        = errors.New("requested resource not found")
	ErrInternalHTTP        = errors.New("an internal server error occurred")
	ErrBadRequestHTTP      = errors.New("bad request: invalid input provided")
	ErrServiceUnavailable  = errors.New("service is temporarily unavailable")
	ErrForbiddenHTTP       = errors.New("access forbidden: you do not have permission")
	ErrUnauthorizedHTTP    = errors.New("authentication required: please log in")
	ErrConflictHTTP        = errors.New("resource conflict: resource already exists or state is invalid")
	ErrUnprocessableEntity = errors.New("unprocessable entity: request was well-formed but could not be processed")
)
