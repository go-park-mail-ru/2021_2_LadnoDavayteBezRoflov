package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases/mocks"
	"backendServer/pkg/closer"
	customErrors "backendServer/pkg/errors"
	zapLogger "backendServer/pkg/logger"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockUserUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateUserHandler(api, "/profile", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	body, err := json.Marshal(testUser)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil).AnyTimes()
	useCaseMock.EXPECT().Create(testUser).Return("sid", nil).AnyTimes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/profile", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil).AnyTimes()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/profile", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockUserUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateUserHandler(api, "/profile", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "2",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "2",
	}

	testUID := uint(2)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	testUser.Login = "2"
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Get(testUID, "2").Return(testUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profile/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/profile/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Get(testUID, testUser.Login).Return(nil, customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/profile/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockUserUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateUserHandler(api, "/profile", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "3",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "3",
	}

	testUID := uint(3)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	body, err := json.Marshal(testUser)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testUser.UID = uint(3)
	useCaseMock.EXPECT().Update(testUser).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/profile/3", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/profile/3", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testUser.UID = uint(3)
	useCaseMock.EXPECT().Update(testUser).Return(customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/profile/3", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/profile/3", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
