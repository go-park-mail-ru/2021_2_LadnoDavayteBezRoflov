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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	_ = faker.SetRandomMapAndSliceSize(1)
	gin.SetMode(gin.ReleaseMode)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	useCaseMock := mocks.NewMockSessionUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(useCaseMock)
	api := router.Group("/api")
	CreateSessionHandler(api, "/sessions", useCaseMock, mw)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)

	body, err := json.Marshal(testUser)
	assert.NoError(t, err)

	// success
	useCaseMock.EXPECT().Create(testUser).Return("testSid", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// can't bind json
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body[:1]))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	useCaseMock.EXPECT().Create(testUser).Return("", customErrors.ErrUserNotFound)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	useCaseMock := mocks.NewMockSessionUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)

	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(useCaseMock)
	api := router.Group("/api")
	CreateSessionHandler(api, "/sessions", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	testUID := uint(1)
	testLogin := "login"

	// success
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Get(cookie.Value).Return(testLogin, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sessions", nil)
	req.AddCookie(cookie)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// not authorized
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Get(cookie.Value).Return("", customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/sessions", nil)
	req.AddCookie(cookie)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// session not found
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/sessions", nil)
	req.AddCookie(cookie)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	useCaseMock := mocks.NewMockSessionUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)

	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(useCaseMock)
	api := router.Group("/api")
	CreateSessionHandler(api, "/sessions", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	// success
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Delete(cookie.Value).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/sessions", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// not authorized
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().Delete(cookie.Value).Return(customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/sessions", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// session not found
	useCaseMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/sessions", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
