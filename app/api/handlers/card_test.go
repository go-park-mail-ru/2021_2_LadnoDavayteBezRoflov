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

func TestCreateCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCardHandler(api, "/cards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	body, err := json.Marshal(testCard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateCard(gomock.Any()).Return(testCard.CID, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/cards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/cards", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/cards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateCard(gomock.Any()).Return(uint(0), customErrors.ErrInternal)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/cards", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCardHandler(api, "/cards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "2",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "2",
	}

	testUID := uint(2)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	testCard.CID = uint(2)
	testCard.CID = uint(2)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().GetCard(testUID, testCard.CID).Return(testCard, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cards/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/cards/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/cards/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCard.CID = uint(2)
	useCaseMock.EXPECT().GetCard(testUID, testCard.CID).Return(nil, customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/cards/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCardHandler(api, "/cards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "3",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "3",
	}

	testUID := uint(3)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	testCard.CID = uint(3)
	testCard.BID = uint(3)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteCard(testUID, testCard.CID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/cards/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/cards/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/cards/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCardHandler(api, "/cards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "4",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "4",
	}

	testUID := uint(4)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	body, err := json.Marshal(testCard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCard.CID = uint(4)
	useCaseMock.EXPECT().UpdateCard(testUID, gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/cards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/test", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCard.CID = uint(4)
	useCaseMock.EXPECT().UpdateCard(testUID, gomock.Any()).Return(customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/4", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestToggleUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCardUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCardHandler(api, "/cards", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "5",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "5",
	}

	testUID := uint(4)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCard.CID = uint(5)
	toggledUserID := uint(5)
	useCaseMock.EXPECT().ToggleUser(testUID, testCard.CID, toggledUserID).Return(testCard, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/cards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/5/toggleuser/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/test/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCard.CID = uint(5)
	useCaseMock.EXPECT().ToggleUser(testUID, testCard.CID, toggledUserID).Return(nil, customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/cards/5/toggleuser/5", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
