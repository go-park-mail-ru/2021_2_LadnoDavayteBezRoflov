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

func TestCreateBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockBoardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateBoardHandler(api, "/boards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)
	body, err := json.Marshal(testBoard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateBoard(gomock.Any()).Return(testBoard.BID, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateBoard(gomock.Any()).Return(uint(0), customErrors.ErrInternal)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockBoardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateBoardHandler(api, "/boards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "2",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "2",
	}

	testUID := uint(2)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	testBoard.BID = uint(2)
	testBoard.BID = uint(2)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().GetBoard(testUID, testBoard.BID).Return(testBoard, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/boards/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/boards/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/boards/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockBoardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateBoardHandler(api, "/boards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "3",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "3",
	}

	testUID := uint(3)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	testBoard.BID = uint(3)
	testBoard.BID = uint(3)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteBoard(testUID, testBoard.BID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/boards/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/boards/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/boards/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockBoardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateBoardHandler(api, "/boards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "4",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "4",
	}

	testUID := uint(4)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)
	body, err := json.Marshal(testBoard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testBoard.BID = uint(4)
	useCaseMock.EXPECT().UpdateBoard(testUID, gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/boards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/test", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testBoard.BID = uint(4)
	useCaseMock.EXPECT().UpdateBoard(testUID, gomock.Any()).Return(customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/4", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestToggleUserBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockBoardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateBoardHandler(api, "/boards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "5",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "5",
	}

	testUID := uint(4)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testBoard.BID = uint(5)
	toggledUserID := uint(5)
	useCaseMock.EXPECT().ToggleUser(testUID, testBoard.BID, toggledUserID).Return(testBoard, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/boards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/5/toggleuser/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/test/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testBoard.BID = uint(5)
	useCaseMock.EXPECT().ToggleUser(testUID, testBoard.BID, toggledUserID).Return(nil, customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/boards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
