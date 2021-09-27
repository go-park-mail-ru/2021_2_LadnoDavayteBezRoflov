package errors

import "errors"

var (
	ErrBadRequest    = errors.New("Bad request")
	ErrBadInputData  = errors.New("Bad input data")
	ErrNotAuthorized = errors.New("Not authorized")

	ErrUserAlreadyCreated = errors.New("User already created")
	ErrEmailAlreadyUsed   = errors.New("Email already used")
)
