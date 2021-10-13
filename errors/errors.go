package errors

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrBadInputData  = errors.New("bad input data")
	ErrNotAuthorized = errors.New("not authorized")

	ErrUserAlreadyCreated = errors.New("user already created")
	ErrEmailAlreadyUsed   = errors.New("email already used")

	ErrInternal = errors.New("internal error")
)

var errorToCodeMap = map[error]int{
	ErrBadRequest:    http.StatusBadRequest,
	ErrBadInputData:  http.StatusUnauthorized,
	ErrNotAuthorized: http.StatusUnauthorized,

	ErrUserAlreadyCreated: http.StatusUnauthorized,
	ErrEmailAlreadyUsed:   http.StatusUnauthorized,

	ErrInternal: http.StatusInternalServerError,
}

func ResolveErrorToCode(err error) (code int) {
	code, isErrorFound := errorToCodeMap[err]
	if !isErrorFound {
		code = http.StatusInternalServerError
	}
	return
}
