package impl

import (
	"backendServer/app/microservices/session/repository/mock"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createSessionRepoMock(controller *gomock.Controller) *mock.MockSessionRepository {
	sessionRepoMock := mock.NewMockSessionRepository(controller)
	return sessionRepoMock
}

func TestCreateSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createSessionRepoMock(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock)

	testUID := uint64(1)
	testSessionID := "testSessionID"

	// good
	sessionRepoMock.EXPECT().Create(testUID).Return(testSessionID, nil)
	resSessionID, err := sessionUseCase.Create(testUID)
	assert.NoError(t, err)
	assert.Equal(t, testSessionID, resSessionID)

	// error can't create
	sessionRepoMock.EXPECT().Create(testUID).Return("", customErrors.ErrInternal)
	resSessionID, err = sessionUseCase.Create(testUID)
	assert.Equal(t, "", resSessionID)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createSessionRepoMock(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock)

	testUID := uint64(1)
	testSessionID := "testSessionID"

	// good
	sessionRepoMock.EXPECT().Get(testSessionID).Return(testUID, nil)
	resUID, err := sessionUseCase.Get(testSessionID)
	assert.NoError(t, err)
	assert.Equal(t, testUID, resUID)

	// error can't get
	sessionRepoMock.EXPECT().Get(testSessionID).Return(uint64(0), customErrors.ErrInternal)
	resUID, err = sessionUseCase.Get(testSessionID)
	assert.Equal(t, uint64(0), resUID)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createSessionRepoMock(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock)

	testSessionID := "testSessionID"

	// good
	sessionRepoMock.EXPECT().Delete(testSessionID).Return(nil)
	err := sessionUseCase.Delete(testSessionID)
	assert.NoError(t, err)

	// error can't delete
	sessionRepoMock.EXPECT().Delete(testSessionID).Return(customErrors.ErrInternal)
	err = sessionUseCase.Delete(testSessionID)
	assert.Equal(t, customErrors.ErrInternal, err)
}
