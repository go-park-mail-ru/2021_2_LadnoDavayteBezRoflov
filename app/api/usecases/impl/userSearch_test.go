package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func createUserSearchRepoMocks(controller *gomock.Controller) (
	*mocks.MockUserRepository,
	*mocks.MockCardRepository,
	*mocks.MockTeamRepository,
	*mocks.MockBoardRepository,
) {
	userRepoMock := mocks.NewMockUserRepository(controller)
	cardRepoMock := mocks.NewMockCardRepository(controller)
	teamRepoMock := mocks.NewMockTeamRepository(controller)
	boardRepoMock := mocks.NewMockBoardRepository(controller)
	return userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock
}

func TestFindForCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock := createUserSearchRepoMocks(ctrl)
	userSearchUseCase := CreateUserSearchUseCase(userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock)

	uid := uint(1)
	cid := uint(1)
	testText := "some user"

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	testMatchedUsers := new([3]models.PublicUserInfo)
	err = faker.FakeData(testMatchedUsers)
	assert.NoError(t, err)
	testMatchedUsersSlice := testMatchedUsers[:]
	testAssignedUsers := new([]models.PublicUserInfo)
	*testAssignedUsers = append(*testAssignedUsers, testMatchedUsersSlice[0])

	expectedUsers := new([]models.UserSearchInfo)
	for _, matchedUser := range testMatchedUsersSlice {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, assignedUser := range *testAssignedUsers {
			if assignedUser.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*expectedUsers = append(*expectedUsers, user)
	}

	// good
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	userRepoMock.EXPECT().FindBoardMembersByLogin(testCard.BID, testText, 15).Return(&testMatchedUsersSlice, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(testAssignedUsers, nil)
	resUsers, err := userSearchUseCase.FindForCard(uid, cid, testText)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, resUsers)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForCard(uid, cid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	_, err = userSearchUseCase.FindForCard(uid, cid, testText)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(nil, customErrors.ErrCardNotFound)
	_, err = userSearchUseCase.FindForCard(uid, cid, testText)
	assert.Equal(t, customErrors.ErrCardNotFound, err)

	// can't get board members
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	userRepoMock.EXPECT().FindBoardMembersByLogin(testCard.BID, testText, 15).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForCard(uid, cid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get assigned users
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	userRepoMock.EXPECT().FindBoardMembersByLogin(testCard.BID, testText, 15).Return(&testMatchedUsersSlice, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForCard(uid, cid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestFindForTeam(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock := createUserSearchRepoMocks(ctrl)
	userSearchUseCase := CreateUserSearchUseCase(userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock)

	uid := uint(1)
	tid := uint(1)
	testText := "some user"

	testMatchedUsers := new([3]models.PublicUserInfo)
	err := faker.FakeData(testMatchedUsers)
	assert.NoError(t, err)
	testMatchedUsersSlice := testMatchedUsers[:]
	testMembers := new([]models.User)
	testExistingMember := models.User{
		UID:    testMatchedUsersSlice[0].UID,
		Login:  testMatchedUsersSlice[0].Login,
		Avatar: testMatchedUsersSlice[0].Avatar,
	}
	*testMembers = append(*testMembers, testExistingMember)

	expectedUsers := new([]models.UserSearchInfo)
	for _, matchedUser := range testMatchedUsersSlice {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, member := range *testMembers {
			if member.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*expectedUsers = append(*expectedUsers, user)
	}

	// good
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().FindAllByLogin(testText, 15).Return(&testMatchedUsersSlice, nil)
	teamRepoMock.EXPECT().GetTeamMembers(tid).Return(testMembers, nil)
	resUsers, err := userSearchUseCase.FindForTeam(uid, tid, testText)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, resUsers)

	// error while checking access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForTeam(uid, tid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(false, nil)
	_, err = userSearchUseCase.FindForTeam(uid, tid, testText)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// error while finding
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().FindAllByLogin(testText, 15).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForTeam(uid, tid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get team members
	userRepoMock.EXPECT().IsUserInTeam(uid, tid).Return(true, nil)
	userRepoMock.EXPECT().FindAllByLogin(testText, 15).Return(&testMatchedUsersSlice, nil)
	teamRepoMock.EXPECT().GetTeamMembers(tid).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForTeam(uid, tid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestFindForBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock := createUserSearchRepoMocks(ctrl)
	userSearchUseCase := CreateUserSearchUseCase(userRepoMock, cardRepoMock, teamRepoMock, boardRepoMock)

	uid := uint(1)
	bid := uint(1)
	testText := "some user"

	testMatchedUsers := new([3]models.PublicUserInfo)
	err := faker.FakeData(testMatchedUsers)
	assert.NoError(t, err)
	testMatchedUsersSlice := testMatchedUsers[:]
	testBoard := new(models.Board)
	err = faker.FakeData(testBoard)
	assert.NoError(t, err)
	testBoard.InvitedMembers = append(testBoard.InvitedMembers, testMatchedUsers[0])

	expectedUsers := new([]models.UserSearchInfo)
	for _, matchedUser := range testMatchedUsersSlice {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, member := range testBoard.InvitedMembers {
			if member.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*expectedUsers = append(*expectedUsers, user)
	}

	// good
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().FindBoardInvitedMembersByLogin(bid, testText, 15).Return(&testMatchedUsersSlice, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	resUsers, err := userSearchUseCase.FindForBoard(uid, bid, testText)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, resUsers)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForBoard(uid, bid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	_, err = userSearchUseCase.FindForBoard(uid, bid, testText)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// error while finding
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().FindBoardInvitedMembersByLogin(bid, testText, 15).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForBoard(uid, bid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get board invited members
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().FindBoardInvitedMembersByLogin(bid, testText, 15).Return(&testMatchedUsersSlice, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(nil, customErrors.ErrInternal)
	_, err = userSearchUseCase.FindForBoard(uid, bid, testText)
	assert.Equal(t, customErrors.ErrInternal, err)
}
