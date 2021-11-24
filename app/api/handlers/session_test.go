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

func TestCreateSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	useCaseMock := mocks.NewMockSessionUseCase(ctrl)

	gin.SetMode(gin.ReleaseMode)
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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
