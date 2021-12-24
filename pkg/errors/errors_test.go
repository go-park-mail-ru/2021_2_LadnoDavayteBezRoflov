package customErrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveErrorToCode(t *testing.T) {
	t.Parallel()

	testErr := ErrBoardNotFound
	testCode := http.StatusNotFound

	testUnexpectedErr := errors.New("some unexpected error")
	testUnexpectedCode := http.StatusInternalServerError

	// success
	resCode := ResolveErrorToCode(testErr)
	assert.Equal(t, testCode, resCode)

	// error not found
	resCode = ResolveErrorToCode(testUnexpectedErr)
	assert.Equal(t, testUnexpectedCode, resCode)
}
