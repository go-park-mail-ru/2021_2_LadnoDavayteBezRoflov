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

type error struct {
	ErrDescription string `json:"error"`
}

var (
	rootURL    = "/api"
	sessionURL = "/sessions"

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
			body:          bytes.NewReader([]byte(`{"login": "1234Anthony", "password": "qwerty"}`)),
			expectedError: errors.ErrBadInputData.Error(),
		},
		{
			testName:      "user don't exist",
			body:          bytes.NewReader([]byte(`{"login": "Anthony", "password": "qwerty"}`)),
			expectedError: errors.ErrBadInputData.Error(),
		},
	}
)

func TestCreateSessionHandler(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

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

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	data.Users["Anthony"] = models.User{
		ID:       101,
		Login:    "Anthony",
		Email:    "ant@mail",
		Password: "qwerty",
	}

	body := bytes.NewReader([]byte(`{"login": "Anthony", "password": "qwerty"}`))
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

	reflect.DeepEqual(expectedStatus, returnedStatus)
}

func TestSessionHandlerCreateFail(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	for _, tt := range testsCreateSessionFail {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			request, _ := http.NewRequest("POST", rootURL+sessionURL, tt.body)
			writer := httptest.NewRecorder()

			router.ServeHTTP(writer, request)

			returnedErr := &error{}
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

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))
	data.Sessions[SID] = user.ID

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

	reflect.DeepEqual(data.Users["Anthony"].Login, writer.Body.String())
}

func TestSessionHandlerGetFailNoCookie(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	request, _ := http.NewRequest("GET", rootURL+sessionURL, nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestSessionHandlerGetFailNoLogin(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

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

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	user := utils.GetSomeUser(data)
	SID := strconv.Itoa(int(user.ID + 1))
	data.Sessions[SID] = user.ID

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

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

	request, _ := http.NewRequest("DELETE", rootURL+sessionURL, nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	require.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestSessionHandlerDeleteFailNoSession(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	router := gin.New()
	routerGroup := router.Group(rootURL)
	sessionRepo := stores.CreateSessionRepository(data)
	CreateSessionHandler(routerGroup, sessionURL, sessionRepo)

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
