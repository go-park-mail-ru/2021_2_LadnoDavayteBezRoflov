package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"backendServer/pkg/utils"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func createSessionRepoMocks(controller *gomock.Controller) (*mocks.MockSessionRepository, *mocks.MockUserRepository) {
	sessionRepoMock := mocks.NewMockSessionRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return sessionRepoMock, userRepoMock
}

func TestCreateSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock := createSessionRepoMocks(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock, userRepoMock)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	testUser.HashedPassword, err = hasher.HashPassword(testUser.Password)
	assert.NoError(t, err)

	// good
	sid := "test_sid"
	isCorrect := utils.ValidateUserData(testUser, false)
	assert.True(t, isCorrect)
	userRepoMock.EXPECT().GetByLogin(testUser.Login).Return(testUser, nil)
	isEqual := hasher.IsPasswordsEqual(testUser.Password, testUser.HashedPassword)
	assert.True(t, isEqual)
	sessionRepoMock.EXPECT().Create(testUser.UID).Return(sid, nil)
	resSid, err := sessionUseCase.Create(testUser)
	assert.NoError(t, err)
	assert.Equal(t, sid, resSid)

	// error user not found
	testUser.Login += "a"
	isCorrect = utils.ValidateUserData(testUser, false)
	assert.True(t, isCorrect)
	userRepoMock.EXPECT().GetByLogin(testUser.Login).Return(nil, customErrors.ErrUserNotFound)
	resSid, err = sessionUseCase.Create(testUser)
	assert.Equal(t, "", resSid)
	assert.Equal(t, customErrors.ErrUserNotFound, err)

	// error bad password not found
	testUser.Password += "a"
	isCorrect = utils.ValidateUserData(testUser, false)
	assert.True(t, isCorrect)
	userRepoMock.EXPECT().GetByLogin(testUser.Login).Return(testUser, nil)
	isEqual = hasher.IsPasswordsEqual(testUser.Password, testUser.HashedPassword)
	assert.False(t, isEqual)
	resSid, err = sessionUseCase.Create(testUser)
	assert.Equal(t, resSid, "")
	assert.Equal(t, customErrors.ErrBadInputData, err)

	// error bad input data
	testUser.Login = "1"
	isCorrect = utils.ValidateUserData(testUser, false)
	assert.False(t, isCorrect)
	resSid, err = sessionUseCase.Create(testUser)
	assert.Equal(t, "", resSid)
	assert.Equal(t, customErrors.ErrBadInputData, err)
}

func TestGetSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock := createSessionRepoMocks(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock, userRepoMock)

	sid := "test_sid"
	uid := uint(1)
	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)

	// has
	sessionRepoMock.EXPECT().Get(sid).Return(uid, nil)
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	login, err := sessionUseCase.Get(sid)
	assert.NoError(t, err)
	assert.Equal(t, testUser.Login, login)

	// not has
	sid += "a"
	sessionRepoMock.EXPECT().Get(sid).Return(uint(0), customErrors.ErrInternal)
	login, err = sessionUseCase.Get(sid)
	assert.Equal(t, "", login)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't find user by uid
	sessionRepoMock.EXPECT().Get(sid).Return(uid, nil)
	userRepoMock.EXPECT().GetByID(uid).Return(nil, customErrors.ErrUserNotFound)
	login, err = sessionUseCase.Get(sid)
	assert.Equal(t, "", login)
	assert.Equal(t, customErrors.ErrUserNotFound, err)
}

func TestGetUID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock := createSessionRepoMocks(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock, userRepoMock)

	sid := "test_sid"
	uid := uint(1)

	// has
	sessionRepoMock.EXPECT().Get(sid).Return(uid, nil)
	resUid, err := sessionUseCase.GetUID(sid)
	assert.NoError(t, err)
	assert.Equal(t, uid, resUid)

	// not has
	sid += "a"
	sessionRepoMock.EXPECT().Get(sid).Return(uint(0), customErrors.ErrInternal)
	resUid, err = sessionUseCase.GetUID(sid)
	assert.Equal(t, uint(0), resUid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock := createSessionRepoMocks(ctrl)
	sessionUseCase := CreateSessionUseCase(sessionRepoMock, userRepoMock)

	sid := "test_sid"

	// good
	sessionRepoMock.EXPECT().Delete(sid).Return(nil)
	err := sessionUseCase.Delete(sid)
	assert.NoError(t, err)

	// error
	sid += "a"
	sessionRepoMock.EXPECT().Delete(sid).Return(customErrors.ErrInternal)
	err = sessionUseCase.Delete(sid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
