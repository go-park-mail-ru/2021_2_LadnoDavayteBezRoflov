package errors

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrBadInputData  = errors.New("bad input data")
	ErrNotAuthorized = errors.New("not authorized")

	ErrUserAlreadyCreated = errors.New("user already created")
	ErrEmailAlreadyUsed   = errors.New("email already used")
)
