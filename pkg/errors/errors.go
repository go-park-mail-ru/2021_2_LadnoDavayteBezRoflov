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

	ErrTeamNotFound = errors.New("team not found")

	ErrBoardNotFound = errors.New("board not found")

	ErrCardListNotFound = errors.New("card list not found")

	ErrCardNotFound = errors.New("card not found")

	ErrCommentNotFound = errors.New("comment not found")

	ErrCheckListNotFound = errors.New("check list not found")

	ErrCheckListItemNotFound = errors.New("check list item not found")

	ErrAttachmentNotFound = errors.New("attachment not found")

	ErrNoAccess = errors.New("no access")

	ErrNotImplemented = errors.New("not implemented")
	ErrInternal       = errors.New("internal error")
)

var errorToCodeMap = map[error]int{
	ErrBadRequest:    http.StatusBadRequest,
	ErrBadInputData:  http.StatusUnauthorized,
	ErrNotAuthorized: http.StatusUnauthorized,

	ErrUserAlreadyCreated: http.StatusUnauthorized,
	ErrEmailAlreadyUsed:   http.StatusUnauthorized,
	ErrUserNotFound:       http.StatusNotFound,

	ErrTeamNotFound: http.StatusNotFound,

	ErrBoardNotFound: http.StatusNotFound,

	ErrCardListNotFound: http.StatusNotFound,

	ErrCardNotFound: http.StatusNotFound,

	ErrCommentNotFound: http.StatusNotFound,

	ErrCheckListNotFound: http.StatusNotFound,

	ErrCheckListItemNotFound: http.StatusNotFound,

	ErrAttachmentNotFound: http.StatusNotFound,

	ErrNoAccess: http.StatusForbidden,

	ErrNotImplemented: http.StatusNotImplemented,
	ErrInternal:       http.StatusInternalServerError,
}

func FindError(err error) error {
	err = errors.Unwrap(err)
	_, isErrorFound := errorToCodeMap[err]
	if isErrorFound {
		return err
	} else {
		return ErrInternal
	}
}

func ResolveErrorToCode(err error) (code int) {
	code, isErrorFound := errorToCodeMap[err]
	if !isErrorFound {
		code = http.StatusInternalServerError
	}
	return
}
