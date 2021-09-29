package handlers

import (
	"backendServer/repositories/stores"
	"backendServer/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	boardURL = "/boards"

	boardRepo = stores.CreateBoardRepository(data)
)

func TestCreateBoardHandler(t *testing.T) {
	t.Parallel()

	request, _ := http.NewRequest("GET", rootURL+boardURL, nil)
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, request)

	require.NotEqual(t, notExpectedErrStatus, writer.Code)
}

func TestBoardHandlerGetAllSuccess(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))
	data.Sessions[SID] = user.ID

	request, _ := http.NewRequest("GET", rootURL+boardURL, nil)
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: SID,
	}
	request.AddCookie(cookie)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	reflect.DeepEqual(data.Teams, writer.Body.String())
}

func TestBoardHandlerGetAllFailNoCookie(t *testing.T) {
	t.Parallel()

	request, _ := http.NewRequest("GET", rootURL+boardURL, nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestBoardHandlerGetAllFailNoSession(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))

	request, _ := http.NewRequest("GET", rootURL+boardURL, nil)
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: SID,
	}
	request.AddCookie(cookie)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}
