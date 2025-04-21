package errs

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrIDNotFound       = errors.New("id not found")
)
