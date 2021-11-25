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

func TestCreateComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCommentUseCase(ctrl)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCommentHandler(api, "/comments", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	testComment.UID = testUID
	body, err := json.Marshal(testComment)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateComment(testComment).Return(testComment, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/comments", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/comments", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/comments", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCommentUseCase(ctrl)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCommentHandler(api, "/comments", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	testComment.UID = testUID
	testComment.CMID = uint(1)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().GetComment(testUID, testComment.CMID).Return(testComment, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/comments/1", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/comments/tete", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

 	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/comments/1", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCommentUseCase(ctrl)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCommentHandler(api, "/comments", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	testComment.UID = testUID
	testComment.CMID = uint(1)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteComment(testUID, testComment.CMID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/comments/1", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/tete", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

 	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/1", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteComment(testUID, testComment.CMID).Return(customErrors.ErrNoAccess)
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/1", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}