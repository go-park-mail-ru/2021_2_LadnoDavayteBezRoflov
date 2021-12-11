package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"mime/multipart"
	"os"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createUserRepoMocks(controller *gomock.Controller) (*mocks.MockSessionRepository, *mocks.MockUserRepository, *mocks.MockTeamRepository) {
	sessionRepoMock := mocks.NewMockSessionRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	teamRepoMock := mocks.NewMockTeamRepository(controller)
	return sessionRepoMock, userRepoMock, teamRepoMock
}

func TestMain(m *testing.M) {
	_ = faker.SetRandomMapAndSliceSize(1)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock, teamRepoMock := createUserRepoMocks(ctrl)
	userUseCase := CreateUserUseCase(sessionRepoMock, userRepoMock, teamRepoMock)

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	testUserTeam := new(models.Team)
	testUserTeam.Title = "Личное пространство " + testUser.Login
	testUserTeam.Type = models.PrivateSpaceTeam
	SID := "test SID"

	// good
	userRepoMock.EXPECT().Create(testUser).Return(nil)
	teamRepoMock.EXPECT().Create(testUserTeam).Return(nil)
	userRepoMock.EXPECT().AddUserToTeam(testUser.UID, testUserTeam.TID).Return(nil)
	sessionRepoMock.EXPECT().Create(testUser.UID).Return(SID, nil)
	resSID, err := userUseCase.Create(testUser)
	assert.NoError(t, err)
	assert.Equal(t, SID, resSID)

	// error can't validate data
	badUser := new(models.User)
	badUser.Email = "not an email"
	_, err = userUseCase.Create(badUser)
	assert.Equal(t, customErrors.ErrBadInputData, err)

	// error can't create user
	userRepoMock.EXPECT().Create(testUser).Return(customErrors.ErrInternal)
	_, err = userUseCase.Create(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error can't create user's team
	userRepoMock.EXPECT().Create(testUser).Return(nil)
	teamRepoMock.EXPECT().Create(testUserTeam).Return(customErrors.ErrInternal)
	_, err = userUseCase.Create(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error can't add user to team
	userRepoMock.EXPECT().Create(testUser).Return(nil)
	teamRepoMock.EXPECT().Create(testUserTeam).Return(nil)
	userRepoMock.EXPECT().AddUserToTeam(testUser.UID, testUserTeam.TID).Return(customErrors.ErrInternal)
	_, err = userUseCase.Create(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error can't create session
	userRepoMock.EXPECT().Create(testUser).Return(nil)
	teamRepoMock.EXPECT().Create(testUserTeam).Return(nil)
	userRepoMock.EXPECT().AddUserToTeam(testUser.UID, testUserTeam.TID).Return(nil)
	sessionRepoMock.EXPECT().Create(testUser.UID).Return("", customErrors.ErrInternal)
	_, err = userUseCase.Create(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock, teamRepoMock := createUserRepoMocks(ctrl)
	userUseCase := CreateUserUseCase(sessionRepoMock, userRepoMock, teamRepoMock)

	uid := uint(1)
	login := "test login"

	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	testUser.UID = uid
	testUser.Login = login

	// success
	userRepoMock.EXPECT().GetByLogin(login).Return(testUser, nil)
	resUser, err := userUseCase.Get(uid, login)
	assert.NoError(t, err)
	assert.Equal(t, testUser, resUser)

	// can't found
	userRepoMock.EXPECT().GetByLogin(login).Return(nil, customErrors.ErrUserNotFound)
	_, err = userUseCase.Get(uid, login)
	assert.Equal(t, customErrors.ErrUserNotFound, err)

	// error can't get user
	userRepoMock.EXPECT().GetByLogin(login).Return(nil, customErrors.ErrInternal)
	_, err = userUseCase.Get(uid, login)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().GetByLogin(login).Return(testUser, nil)
	_, err = userUseCase.Get(uint(0), login)
	assert.Equal(t, customErrors.ErrNoAccess, err)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock, teamRepoMock := createUserRepoMocks(ctrl)
	userUseCase := CreateUserUseCase(sessionRepoMock, userRepoMock, teamRepoMock)

	uid := uint(1)
	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	testUser.UID = uid
	testUser.HashedPassword, err = hasher.HashPassword(testUser.OldPassword)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	userRepoMock.EXPECT().Update(testUser).Return(nil)
	err = userUseCase.Update(testUser)
	assert.NoError(t, err)

	// can't found
	userRepoMock.EXPECT().GetByID(uid).Return(nil, customErrors.ErrUserNotFound)
	err = userUseCase.Update(testUser)
	assert.Equal(t, customErrors.ErrUserNotFound, err)

	// error can't get user
	userRepoMock.EXPECT().GetByID(uid).Return(nil, customErrors.ErrInternal)
	err = userUseCase.Update(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't update
	testUser.OldPassword = "some old password"
	testUser.HashedPassword, err = hasher.HashPassword(testUser.OldPassword)
	assert.NoError(t, err)
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	userRepoMock.EXPECT().Update(testUser).Return(customErrors.ErrInternal)
	err = userUseCase.Update(testUser)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error bad password
	testUser.OldPassword = "bad password"
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	err = userUseCase.Update(testUser)
	assert.Equal(t, customErrors.ErrBadRequest, err)
}

func TestUpdateUserAvatar(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock, userRepoMock, teamRepoMock := createUserRepoMocks(ctrl)
	userUseCase := CreateUserUseCase(sessionRepoMock, userRepoMock, teamRepoMock)

	uid := uint(1)
	testUser := new(models.User)
	err := faker.FakeData(testUser)
	assert.NoError(t, err)
	testUser.UID = uid
	testAvatar := new(multipart.FileHeader)

	// success
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	userRepoMock.EXPECT().UpdateAvatar(testUser, testAvatar).Return(nil)
	err = userUseCase.UpdateAvatar(testUser, testAvatar)
	assert.NoError(t, err)

	// can't found
	userRepoMock.EXPECT().GetByID(uid).Return(nil, customErrors.ErrUserNotFound)
	err = userUseCase.UpdateAvatar(testUser, testAvatar)
	assert.Equal(t, customErrors.ErrUserNotFound, err)

	// error can't get user
	userRepoMock.EXPECT().GetByID(uid).Return(nil, customErrors.ErrInternal)
	err = userUseCase.UpdateAvatar(testUser, testAvatar)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't update avatar
	userRepoMock.EXPECT().GetByID(uid).Return(testUser, nil)
	userRepoMock.EXPECT().UpdateAvatar(testUser, testAvatar).Return(customErrors.ErrInternal)
	err = userUseCase.UpdateAvatar(testUser, testAvatar)
	assert.Equal(t, customErrors.ErrInternal, err)
}
