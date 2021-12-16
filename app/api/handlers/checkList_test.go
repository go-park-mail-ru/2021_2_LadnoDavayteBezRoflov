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

func TestCreateCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCheckListUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCheckListHandler(api, "/CheckLists", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "1",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "1",
	}

	testUID := uint(1)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	assert.NoError(t, err)
	body, err := json.Marshal(testCheckList)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateCheckList(testCheckList).Return(testCheckList.CHLID, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/CheckLists", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/CheckLists", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/CheckLists", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().CreateCheckList(testCheckList).Return(uint(0), customErrors.ErrInternal)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/CheckLists", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCheckListUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCheckListHandler(api, "/CheckLists", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "2",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "2",
	}

	testUID := uint(2)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	testCheckList.CHLID = uint(2)
	testCheckList.CID = uint(2)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().GetCheckList(testUID, testCheckList.CHLID).Return(testCheckList, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/CheckLists/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/CheckLists/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/CheckLists/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCheckList.CHLID = uint(2)
	useCaseMock.EXPECT().GetCheckList(testUID, testCheckList.CHLID).Return(nil, customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/CheckLists/2", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCheckListUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCheckListHandler(api, "/CheckLists", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "3",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "3",
	}

	testUID := uint(3)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	testCheckList.CHLID = uint(3)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteCheckList(testUID, testCheckList.CHLID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/CheckLists/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/CheckLists/test", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/CheckLists/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	useCaseMock.EXPECT().DeleteCheckList(testUID, testCheckList.CHLID).Return(customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/CheckLists/3", nil)
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionMock := mocks.NewMockSessionUseCase(ctrl)
	useCaseMock := mocks.NewMockCheckListUseCase(ctrl)

	router := gin.Default()
	var logger zapLogger.Logger
	logger.InitLogger("./logs.log")
	everythingCloser := closer.CreateCloser(&logger)
	defer everythingCloser.Close(logger.Sync)
	commonMW := CreateCommonMiddleware(logger)
	router.Use(commonMW.Logger())
	mw := CreateSessionMiddleware(sessionMock)
	api := router.Group("/api")
	CreateCheckListHandler(api, "/CheckLists", useCaseMock, mw)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: "4",
	}

	csrfToken := &http.Cookie{
		Name:  "csrf_token",
		Value: "4",
	}

	testUID := uint(4)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	assert.NoError(t, err)
	body, err := json.Marshal(testCheckList)
	assert.NoError(t, err)

	// success
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCheckList.CHLID = uint(4)
	useCaseMock.EXPECT().UpdateCheckList(testUID, testCheckList).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/CheckLists/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ParseUint not parsing
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/CheckLists/test", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// not authorized
	sessionMock.EXPECT().GetUID(cookie.Value).Return(uint(0), customErrors.ErrNotAuthorized)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/CheckLists/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// fail
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)
	testCheckList.CHLID = uint(4)
	useCaseMock.EXPECT().UpdateCheckList(testUID, testCheckList).Return(customErrors.ErrNoAccess)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/CheckLists/4", bytes.NewBuffer(body))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// json not binding
	sessionMock.EXPECT().GetUID(cookie.Value).Return(testUID, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/CheckLists/4", bytes.NewBuffer(body[:1]))
	req.AddCookie(cookie)
	req.AddCookie(csrfToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
