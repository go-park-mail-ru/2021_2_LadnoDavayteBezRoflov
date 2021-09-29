package handlers

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories/stores"
	"backendServer/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

type status struct {
	StatusDescription string `json:"status"`
}

type errorResponse struct {
	ErrDescription string `json:"error"`
}

var (
	rootURL    = "/api"
	sessionURL = "/sessions"

	data, _     = utils.FillTestData(10, 3, 100)
	router      = gin.New()
	routerGroup = router.Group(rootURL)
	sessionRepo = stores.CreateSessionRepository(data)

	notExpectedErrStatus      = http.StatusNotFound
	testsCreateSessionHandler = []struct {
		testName   string
		methodType string
	}{
		{
			testName:   "POST request",
			methodType: "POST",
		},
		{
			testName:   "GET request",
			methodType: "GET",
		},
		{
			testName:   "DELETE request",
			methodType: "DELETE",
		},
	}

	testsCreateSessionFail = []struct {
		testName      string
		body          *bytes.Reader
		expectedError string
	}{
		{
			testName:      "not json",
			body:          bytes.NewReader([]byte(`{"someWrongData"}`)),
			expectedError: errors.ErrBadRequest.Error(),
		},
		{
			testName:      "not user json",
			body:          bytes.NewReader([]byte(`{"name": "Anthony", "title": "qwerty"}`)),
			expectedError: errors.ErrBadInputData.Error(),
		},
		{
			testName:      "invalid user data",
			body:          bytes.NewReader([]byte(`{"login": "РусскийЛогин", "password": "qwerty"}`)),
			expectedError: errors.ErrBadInputData.Error(),
		},
		{
			testName:      "user don't exist",
			body:          bytes.NewReader([]byte(`{"login": "AnthonyChum", "password": "qwerty"}`)),
			expectedError: errors.ErrBadInputData.Error(),
		},
	}
)

func TestMain(m *testing.M) {
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)
	CreateUserHandler(routerGroup, userURL, userRepo, sessionRepo)
	CreateBoardHandler(routerGroup, boardURL, boardRepo, sessionRepo)

	code := m.Run()
	os.Exit(code)
}

func TestCreateSessionHandler(t *testing.T) {
	t.Parallel()

	for _, tt := range testsCreateSessionHandler {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			request, _ := http.NewRequest(tt.methodType, rootURL+sessionURL, nil)
			writer := httptest.NewRecorder()

			router.ServeHTTP(writer, request)

			require.NotEqual(t, notExpectedErrStatus, writer.Code)
		})
	}
}

func TestSessionHandlerCreateSuccess(t *testing.T) {
	t.Parallel()

	newUser := &models.User{}
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}

	data.Mu.RLock()
	usersAmount := len(data.Users)
	data.Mu.RUnlock()

	newUser.ID = uint(usersAmount + 1)

	data.Mu.Lock()
	data.Users[newUser.Login] = *newUser
	data.Mu.Unlock()

	jsonNewUser, _ := json.Marshal(newUser)

	body := bytes.NewReader(jsonNewUser)
	request, _ := http.NewRequest("POST", rootURL+sessionURL, body)
	writer := httptest.NewRecorder()
	expectedStatus := status{StatusDescription: "you are logged in"}

	router.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	returnedStatus := &status{}
	err = json.Unmarshal(writer.Body.Bytes(), returnedStatus)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, expectedStatus.StatusDescription, returnedStatus.StatusDescription)
}

func TestSessionHandlerCreateFail(t *testing.T) {
	t.Parallel()

	for _, tt := range testsCreateSessionFail {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			request, _ := http.NewRequest("POST", rootURL+sessionURL, tt.body)
			writer := httptest.NewRecorder()

			router.ServeHTTP(writer, request)

			returnedErr := &errorResponse{}
			err := json.Unmarshal(writer.Body.Bytes(), returnedErr)
			if err != nil {
				t.Error(err)
			}

			require.Equal(t, tt.expectedError, returnedErr.ErrDescription)
		})
	}
}

func TestSessionHandlerGetSuccess(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))

	data.Mu.Lock()
	data.Sessions[SID] = user.ID
	data.Mu.Unlock()

	request, _ := http.NewRequest("GET", rootURL+sessionURL, nil)
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

	reflect.DeepEqual(user.Login, writer.Body.String())
}

func TestSessionHandlerGetFailNoCookie(t *testing.T) {
	t.Parallel()

	request, _ := http.NewRequest("GET", rootURL+sessionURL, nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestSessionHandlerGetFailNoLogin(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))

	request, _ := http.NewRequest("GET", rootURL+sessionURL, nil)
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: SID,
	}
	request.AddCookie(cookie)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestSessionHandlerDeleteSuccess(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))

	data.Mu.Lock()
	data.Sessions[SID] = user.ID
	data.Mu.Unlock()

	request, _ := http.NewRequest("GET", rootURL+sessionURL, nil)
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: SID,
	}
	request.AddCookie(cookie)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusOK, writer.Code)
}

func TestSessionHandlerDeleteFailNoCookie(t *testing.T) {
	t.Parallel()

	request, _ := http.NewRequest("DELETE", rootURL+sessionURL, nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestSessionHandlerDeleteFailNoSession(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))

	request, _ := http.NewRequest("DELETE", rootURL+sessionURL, nil)
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: SID,
	}
	request.AddCookie(cookie)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}
