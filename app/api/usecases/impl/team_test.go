package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createTeamRepoMocks(controller *gomock.Controller) (*mocks.MockTeamRepository, *mocks.MockUserRepository, *mocks.MockBoardRepository) {
	teamRepoMock := mocks.NewMockTeamRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	boardRepoMock := mocks.NewMockBoardRepository(controller)
	return teamRepoMock, userRepoMock, boardRepoMock
}

func TestCreateTeam(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	teamRepoMock, userRepoMock, boardRepoMock := createTeamRepoMocks(ctrl)
	teamUseCase := CreateTeamUseCase(teamRepoMock, userRepoMock, boardRepoMock)

	uid := uint(1)
	testTeam := new(models.Team)
	err := faker.FakeData(testTeam)
	assert.NoError(t, err)

	// good
	teamRepoMock.EXPECT().Create(testTeam).Return(nil)
	userRepoMock.EXPECT().AddUserToTeam(uid, testTeam.TID).Return(nil)
	resTid, err := teamUseCase.CreateTeam(uid, testTeam)
	assert.NoError(t, err)
	assert.Equal(t, testTeam.TID, resTid)

	// error can't create
	teamRepoMock.EXPECT().Create(testTeam).Return(customErrors.ErrInternal)
	resTid, err = teamUseCase.CreateTeam(uid, testTeam)
	assert.Equal(t, uint(0), resTid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error can't add user
	teamRepoMock.EXPECT().Create(testTeam).Return(nil)
	userRepoMock.EXPECT().AddUserToTeam(uid, testTeam.TID).Return(customErrors.ErrInternal)
	resTid, err = teamUseCase.CreateTeam(uid, testTeam)
	assert.Equal(t, uint(0), resTid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetTeam(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	teamRepoMock, userRepoMock, boardRepoMock := createTeamRepoMocks(ctrl)
	teamUseCase := CreateTeamUseCase(teamRepoMock, userRepoMock, boardRepoMock)

	uid := uint(1)
	tid := uint(1)

	testTeam := new(models.Team)
	err := faker.FakeData(testTeam)
	assert.NoError(t, err)
	testBoards := new([]models.Board)
	testMembers := new([]models.User)
	testTeam.Boards = *testBoards
	testTeam.Users = *testMembers

	// success
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(testTeam, nil)
	teamRepoMock.EXPECT().GetTeamBoards(tid).Return(testBoards, nil)
	teamRepoMock.EXPECT().GetTeamMembers(tid).Return(testMembers, nil)
	resTeam, err := teamUseCase.GetTeam(uid, tid)
	assert.NoError(t, err)
	assert.Equal(t, testTeam, resTeam)

	// error while checking access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, customErrors.ErrInternal)
	_, err = teamUseCase.GetTeam(uid, tid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, nil)
	_, err = teamUseCase.GetTeam(uid, tid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(nil, customErrors.ErrTeamNotFound)
	_, err = teamUseCase.GetTeam(uid, tid)
	assert.Equal(t, customErrors.ErrTeamNotFound, err)

	// can't get team boards
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(testTeam, nil)
	teamRepoMock.EXPECT().GetTeamBoards(tid).Return(nil, customErrors.ErrInternal)
	_, err = teamUseCase.GetTeam(uid, tid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get team members
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(testTeam, nil)
	teamRepoMock.EXPECT().GetTeamBoards(tid).Return(testBoards, nil)
	teamRepoMock.EXPECT().GetTeamMembers(tid).Return(nil, customErrors.ErrInternal)
	_, err = teamUseCase.GetTeam(uid, tid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	teamRepoMock, userRepoMock, boardRepoMock := createTeamRepoMocks(ctrl)
	teamUseCase := CreateTeamUseCase(teamRepoMock, userRepoMock, boardRepoMock)

	uid := uint(1)
	testTeam := new(models.Team)
	err := faker.FakeData(testTeam)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsUserInTeam(uid, testTeam.TID).Return(true, nil)
	teamRepoMock.EXPECT().Update(testTeam).Return(nil)
	err = teamUseCase.UpdateTeam(uid, testTeam)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsUserInTeam(uid, testTeam.TID).Return(false, customErrors.ErrInternal)
	err = teamUseCase.UpdateTeam(uid, testTeam)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsUserInTeam(uid, testTeam.TID).Return(false, nil)
	err = teamUseCase.UpdateTeam(uid, testTeam)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsUserInTeam(uid, testTeam.TID).Return(true, nil)
	teamRepoMock.EXPECT().Update(testTeam).Return(customErrors.ErrInternal)
	err = teamUseCase.UpdateTeam(uid, testTeam)
	assert.Equal(t, customErrors.ErrInternal, err)
}

//
//func TestDeleteTeam(t *testing.T) {
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	teamRepoMock, userRepoMock := createTeamRepoMocks(ctrl)
//	teamUseCase := CreateTeamUseCase(teamRepoMock, userRepoMock)
//
//	uid := uint(1)
//	tid := uint(1)
//
//	// good
//	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
//	teamRepoMock.EXPECT().Delete(tid).Return(nil)
//	err := teamUseCase.DeleteTeam(uid, tid)
//	assert.NoError(t, err)
//
//	// error while checking access
//	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, customErrors.ErrInternal)
//	err = teamUseCase.DeleteTeam(uid, tid)
//	assert.Equal(t, customErrors.ErrInternal, err)
//
//	// no access
//	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, nil)
//	err = teamUseCase.DeleteTeam(uid, tid)
//	assert.Equal(t, customErrors.ErrNoAccess, err)
//
//	// can't delete
//	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
//	teamRepoMock.EXPECT().Delete(tid).Return(customErrors.ErrInternal)
//	err = teamUseCase.DeleteTeam(uid, tid)
//	assert.Equal(t, customErrors.ErrInternal, err)
//}

func TestToggleTeam(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	teamRepoMock, userRepoMock, boardRepoMock := createTeamRepoMocks(ctrl)
	teamUseCase := CreateTeamUseCase(teamRepoMock, userRepoMock, boardRepoMock)

	uid := uint(1)
	tid := uint(1)
	toggledUserId := uint(1)
	testTeam := new(models.Team)
	err := faker.FakeData(testTeam)
	assert.NoError(t, err)
	testBoards := new([]models.Board)
	testMembers := new([]models.User)
	testTeam.Boards = *testBoards
	testTeam.Users = *testMembers

	// success
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToTeam(toggledUserId, tid).Return(nil)

	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(testTeam, nil)
	teamRepoMock.EXPECT().GetTeamBoards(tid).Return(testBoards, nil)
	teamRepoMock.EXPECT().GetTeamMembers(tid).Return(testMembers, nil)
	resTeam, err := teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.NoError(t, err)
	assert.Equal(t, testTeam, resTeam)

	// error while checking access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, customErrors.ErrInternal)
	_, err = teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, nil)
	_, err = teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't toggle
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToTeam(toggledUserId, tid).Return(customErrors.ErrInternal)
	_, err = teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// user was successfully deleted from team
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToTeam(toggledUserId, tid).Return(nil)
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, customErrors.ErrNoAccess)
	_, err = teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.Equal(t, nil, err)

	// get team error
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToTeam(toggledUserId, tid).Return(nil)
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	teamRepoMock.EXPECT().GetByID(tid).Return(nil, customErrors.ErrInternal)
	_, err = teamUseCase.ToggleUser(uid, tid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)
}
