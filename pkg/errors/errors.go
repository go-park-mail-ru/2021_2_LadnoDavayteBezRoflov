package customErrors

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
	ErrUserNotFound       = errors.New("user not found")

	ErrInternal = errors.New("internal error")
)

var errorToCodeMap = map[error]int{
	ErrBadRequest:    http.StatusBadRequest,
	ErrBadInputData:  http.StatusUnauthorized,
	ErrNotAuthorized: http.StatusUnauthorized,

	ErrUserAlreadyCreated: http.StatusUnauthorized,
	ErrEmailAlreadyUsed:   http.StatusUnauthorized,
	ErrUserNotFound:       http.StatusNotFound,

	ErrInternal: http.StatusInternalServerError,
}

func ResolveErrorToCode(err error) (code int) {
	code, isErrorFound := errorToCodeMap[err]
	if !isErrorFound {
		code = http.StatusInternalServerError
	}
	return
}
