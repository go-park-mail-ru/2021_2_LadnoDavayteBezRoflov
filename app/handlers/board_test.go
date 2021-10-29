package handlers

import (
	"backendServer/app/models"
	"backendServer/app/repositories/stores"
	"backendServer/pkg/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	SID := strconv.Itoa(int(user.UID + 1))

	data.Mu.Lock()
	data.Sessions[SID] = user.UID
	data.Mu.Unlock()

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

	allReturnedTeams := []models.Team{}
	err := json.Unmarshal(writer.Body.Bytes(), &allReturnedTeams)
	if err != nil {
		t.Error(err)
	}

	data.Mu.RLock()
	allExpectedTeams := data.Teams

	isEqual := true
	if len(allReturnedTeams) != len(allExpectedTeams) {
		isEqual = false
	}

	for _, team := range allReturnedTeams {
		if _, teamExpected := allExpectedTeams[team.TID]; !teamExpected {
			isEqual = false
		}
	}
	data.Mu.RUnlock()

	require.True(t, isEqual)
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
	SID := strconv.Itoa(int(user.UID + 1))

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
