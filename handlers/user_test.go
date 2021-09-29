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
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

var (
	userURL = "/profile"

	userRepo        = stores.CreateUserRepository(data)
	jsonSomeUser, _ = json.Marshal(utils.GetSomeUser(data))

	testsCreateUserFail = []struct {
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
			testName:      "user already exist",
			body:          bytes.NewReader(jsonSomeUser),
			expectedError: errors.ErrUserAlreadyCreated.Error(),
		},
	}
)

func TestCreateUserHandler(t *testing.T) {
	t.Parallel()

	request, _ := http.NewRequest("POST", rootURL+userURL, nil)
	writer := httptest.NewRecorder()

	router.ServeHTTP(writer, request)

	require.NotEqual(t, notExpectedErrStatus, writer.Code)
}

func TestUserHandlerCreateSuccess(t *testing.T) {
	t.Parallel()

	newUser := &models.User{}
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}

	jsonNewUser, _ := json.Marshal(newUser)
	body := bytes.NewReader(jsonNewUser)

	request, _ := http.NewRequest("POST", rootURL+userURL, body)
	writer := httptest.NewRecorder()
	expectedStatus := status{StatusDescription: "you are logged in"}

	router.ServeHTTP(writer, request)

	if writer.Code != http.StatusCreated {
		t.Error("status is not ok")
	}

	returnedStatus := &status{}
	err = json.Unmarshal(writer.Body.Bytes(), returnedStatus)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, expectedStatus.StatusDescription, returnedStatus.StatusDescription)
}

func TestUserHandlerCreateFail(t *testing.T) {
	t.Parallel()

	for _, tt := range testsCreateUserFail {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			request, _ := http.NewRequest("POST", rootURL+userURL, tt.body)
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
