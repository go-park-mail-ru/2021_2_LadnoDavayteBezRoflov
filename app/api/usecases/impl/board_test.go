package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createBoardRepoMocks(controller *gomock.Controller) (
	*mocks.MockBoardRepository,
	*mocks.MockUserRepository,
	*mocks.MockTeamRepository,
	*mocks.MockCardListRepository,
	*mocks.MockCardRepository,
	*mocks.MockCheckListRepository,
) {
	boardRepoMock := mocks.NewMockBoardRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	teamRepoMock := mocks.NewMockTeamRepository(controller)
	cardListRepoMock := mocks.NewMockCardListRepository(controller)
	cardRepoMock := mocks.NewMockCardRepository(controller)
	checkListRepoMock := mocks.NewMockCheckListRepository(controller)
	return boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock
}

func TestCreateBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)

	// good
	boardRepoMock.EXPECT().Create(testBoard).Return(nil)
	resBID, err := boardUseCase.CreateBoard(testBoard)
	assert.NoError(t, err)
	assert.Equal(t, testBoard.BID, resBID)

	// error can't create
	boardRepoMock.EXPECT().Create(testBoard).Return(customErrors.ErrInternal)
	resBID, err = boardUseCase.CreateBoard(testBoard)
	assert.Equal(t, uint(0), resBID)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetUserBoards(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	testTeams := new([3]models.Team)
	err := faker.FakeData(testTeams)
	assert.NoError(t, err)
	testToggledBoards := new([3]models.Board)
	err = faker.FakeData(testToggledBoards)
	assert.NoError(t, err)
	testTeamsSlice := testTeams[:]
	testToggledBoardsSlice := testToggledBoards[:]

	// success (and no toggled boards)
	userRepoMock.EXPECT().GetUserTeams(uid).Return(&testTeamsSlice, nil)
	for _, testTeam := range testTeamsSlice {
		teamRepoMock.EXPECT().GetTeamBoards(testTeam.TID).Return(&testTeam.Boards, nil)
		teamRepoMock.EXPECT().GetTeamMembers(testTeam.TID).Return(&testTeam.Users, nil)
	}
	userRepoMock.EXPECT().GetUserToggledBoards(uid).Return(&[]models.Board{}, nil)
	resBoards, err := boardUseCase.GetUserBoards(uid)
	assert.NoError(t, err)
	assert.Equal(t, testTeamsSlice, *resBoards)

	// success (and toggled boards exists)
	userRepoMock.EXPECT().GetUserTeams(uid).Return(&testTeamsSlice, nil)
	for _, testTeam := range testTeamsSlice {
		teamRepoMock.EXPECT().GetTeamBoards(testTeam.TID).Return(&testTeam.Boards, nil)
		teamRepoMock.EXPECT().GetTeamMembers(testTeam.TID).Return(&testTeam.Users, nil)
	}
	userRepoMock.EXPECT().GetUserToggledBoards(uid).Return(&testToggledBoardsSlice, nil)
	resBoards, err = boardUseCase.GetUserBoards(uid)
	assert.NoError(t, err)
	assert.Equal(t, testTeamsSlice, *resBoards)

	// can't get user teams
	userRepoMock.EXPECT().GetUserTeams(uid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetUserBoards(uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get team boards
	userRepoMock.EXPECT().GetUserTeams(uid).Return(&testTeamsSlice, nil)
	teamRepoMock.EXPECT().GetTeamBoards(testTeamsSlice[0].TID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetUserBoards(uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get team members
	userRepoMock.EXPECT().GetUserTeams(uid).Return(&testTeamsSlice, nil)
	teamRepoMock.EXPECT().GetTeamBoards(testTeamsSlice[0].TID).Return(&testTeamsSlice[0].Boards, nil)
	teamRepoMock.EXPECT().GetTeamMembers(testTeamsSlice[0].TID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetUserBoards(uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get user toggled boards
	userRepoMock.EXPECT().GetUserTeams(uid).Return(&testTeamsSlice, nil)
	for _, testTeam := range testTeamsSlice {
		teamRepoMock.EXPECT().GetTeamBoards(testTeam.TID).Return(&testTeam.Boards, nil)
		teamRepoMock.EXPECT().GetTeamMembers(testTeam.TID).Return(&testTeam.Users, nil)
	}
	userRepoMock.EXPECT().GetUserToggledBoards(uid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetUserBoards(uid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	bid := uint(1)

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)
	testCardList := new(models.CardList)
	err = faker.FakeData(testBoard)
	assert.NoError(t, err)
	testCard := new(models.Card)
	err = faker.FakeData(testCard)
	assert.NoError(t, err)
	testTag := new(models.Tag)
	err = faker.FakeData(testTag)
	assert.NoError(t, err)
	testTag.ColorID = uint(2)
	testComment := new(models.Comment)
	err = faker.FakeData(testComment)
	assert.NoError(t, err)
	testAttachment := new(models.Attachment)
	err = faker.FakeData(testAttachment)
	assert.NoError(t, err)
	testCheckList := new(models.CheckList)
	err = faker.FakeData(testCheckList)
	assert.NoError(t, err)
	testCheckListItem := new(models.CheckListItem)
	err = faker.FakeData(testCheckListItem)
	assert.NoError(t, err)
	testCheckList.CheckListItems = append(testCheckList.CheckListItems, *testCheckListItem)
	testCard.CheckLists = append(testCard.CheckLists, *testCheckList)
	testCard.Comments = append(testCard.Comments, *testComment)
	testCard.Attachments = append(testCard.Attachments, *testAttachment)
	testCard.Tags = append(testCard.Tags, *testTag)
	testCardList.Cards = append(testCardList.Cards, *testCard)
	testBoard.CardLists = append(testBoard.CardLists, *testCardList)

	for i, cardList := range testBoard.CardLists {
		for j, card := range cardList.Cards {
			for index, comment := range card.Comments {
				testBoard.CardLists[i].Cards[j].Comments[index].DateParsed = comment.Date.Round(time.Second).String()
			}
		}
	}

	// success
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	for _, cardList := range testBoard.CardLists {
		cardListRepoMock.EXPECT().GetCardListCards(cardList.CLID).Return(&cardList.Cards, nil)
		for _, card := range cardList.Cards {
			cardRepoMock.EXPECT().GetCardComments(card.CID).Return(&card.Comments, nil)
			cardRepoMock.EXPECT().GetCardAttachments(card.CID).Return(&card.Attachments, nil)
			cardRepoMock.EXPECT().GetAssignedUsers(card.CID).Return(&card.Assignees, nil)
			for _, comment := range card.Comments {
				userRepoMock.EXPECT().GetPublicData(comment.UID).Return(&comment.User, nil)
			}
			cardRepoMock.EXPECT().GetCardTags(card.CID).Return(&card.Tags, nil)
			cardRepoMock.EXPECT().GetCardCheckLists(card.CID).Return(&card.CheckLists, nil)
			for _, checkList := range card.CheckLists {
				checkListRepoMock.EXPECT().GetCheckListItems(checkList.CHLID).Return(&checkList.CheckListItems, nil)
			}
		}
	}
	resBoard, err := boardUseCase.GetBoard(uid, bid)
	assert.NoError(t, err)
	assert.Equal(t, testBoard, resBoard)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(nil, customErrors.ErrTeamNotFound)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrTeamNotFound, err)

	// can't get board members
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get board invited members
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get board card lists
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get board tags
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get board cards
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get comments
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get attachments
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get assigned users
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Attachments, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(testBoard.CardLists[0].Cards[0].CID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get user's public data
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Attachments, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Assignees, nil)
	userRepoMock.EXPECT().GetPublicData(testBoard.CardLists[0].Cards[0].Comments[0].UID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get card's tags
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Attachments, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Assignees, nil)
	userRepoMock.EXPECT().GetPublicData(testBoard.CardLists[0].Cards[0].Comments[0].UID).Return(&testBoard.CardLists[0].Cards[0].Comments[0].User, nil)
	cardRepoMock.EXPECT().GetCardTags(testBoard.CardLists[0].Cards[0].CID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get check lists
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Attachments, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Assignees, nil)
	userRepoMock.EXPECT().GetPublicData(testBoard.CardLists[0].Cards[0].Comments[0].UID).Return(&testBoard.CardLists[0].Cards[0].Comments[0].User, nil)
	cardRepoMock.EXPECT().GetCardTags(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Tags, nil)
	cardRepoMock.EXPECT().GetCardCheckLists(testBoard.CardLists[0].Cards[0].CID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get check list's items
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	cardListRepoMock.EXPECT().GetCardListCards(testBoard.CardLists[0].CLID).Return(&testBoard.CardLists[0].Cards, nil)
	cardRepoMock.EXPECT().GetCardComments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Comments, nil)
	cardRepoMock.EXPECT().GetCardAttachments(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Attachments, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Assignees, nil)
	userRepoMock.EXPECT().GetPublicData(testBoard.CardLists[0].Cards[0].Comments[0].UID).Return(&testBoard.CardLists[0].Cards[0].Comments[0].User, nil)
	cardRepoMock.EXPECT().GetCardTags(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].Tags, nil)
	cardRepoMock.EXPECT().GetCardCheckLists(testBoard.CardLists[0].Cards[0].CID).Return(&testBoard.CardLists[0].Cards[0].CheckLists, nil)
	checkListRepoMock.EXPECT().GetCheckListItems(testBoard.CardLists[0].Cards[0].CheckLists[0].CHLID).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.GetBoard(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsBoardAccessed(uid, testBoard.BID).Return(true, nil)
	boardRepoMock.EXPECT().Update(testBoard).Return(nil)
	err = boardUseCase.UpdateBoard(uid, testBoard)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, testBoard.BID).Return(false, customErrors.ErrInternal)
	err = boardUseCase.UpdateBoard(uid, testBoard)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, testBoard.BID).Return(false, nil)
	err = boardUseCase.UpdateBoard(uid, testBoard)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsBoardAccessed(uid, testBoard.BID).Return(true, nil)
	boardRepoMock.EXPECT().Update(testBoard).Return(customErrors.ErrInternal)
	err = boardUseCase.UpdateBoard(uid, testBoard)
	assert.Equal(t, customErrors.ErrInternal, err)
}

//
//func TestDeleteBoard(t *testing.T) {
//	t.Parallel()
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
//	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)
//
//	uid := uint(1)
//	bid := uint(1)
//
//	// good
//	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
//	boardRepoMock.EXPECT().Delete(bid).Return(nil)
//	err := boardUseCase.DeleteBoard(uid, bid)
//	assert.NoError(t, err)
//
//	// error while checking access
//	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
//	err = boardUseCase.DeleteBoard(uid, bid)
//	assert.Equal(t, customErrors.ErrInternal, err)
//
//	// no access
//	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
//	err = boardUseCase.DeleteBoard(uid, bid)
//	assert.Equal(t, customErrors.ErrNoAccess, err)
//
//	// can't delete
//	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
//	boardRepoMock.EXPECT().Delete(bid).Return(customErrors.ErrInternal)
//	err = boardUseCase.DeleteBoard(uid, bid)
//	assert.Equal(t, customErrors.ErrInternal, err)
//}

func TestToggleBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	bid := uint(1)
	toggledUserId := uint(1)
	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToBoard(toggledUserId, bid).Return(nil)

	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	for _, cardList := range testBoard.CardLists {
		cardListRepoMock.EXPECT().GetCardListCards(cardList.CLID).Return(&cardList.Cards, nil)
		for _, card := range cardList.Cards {
			cardRepoMock.EXPECT().GetCardComments(card.CID).Return(&card.Comments, nil)
			cardRepoMock.EXPECT().GetAssignedUsers(card.CID).Return(&card.Attachments, nil)
			cardRepoMock.EXPECT().GetAssignedUsers(card.CID).Return(&card.Assignees, nil)
			for _, comment := range card.Comments {
				userRepoMock.EXPECT().GetPublicData(comment.UID).Return(&comment.User, nil)
			}
			cardRepoMock.EXPECT().GetCardTags(card.CID).Return(&card.Tags, nil)
			cardRepoMock.EXPECT().GetCardCheckLists(card.CID).Return(&card.CheckLists, nil)
			for _, checkList := range card.CheckLists {
				checkListRepoMock.EXPECT().GetCheckListItems(checkList.CHLID).Return(&checkList.CheckListItems, nil)
			}
		}
	}
	resBoard, err := boardUseCase.ToggleUser(uid, bid, toggledUserId)
	assert.NoError(t, err)
	assert.Equal(t, testBoard, resBoard)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
	_, err = boardUseCase.ToggleUser(uid, bid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	_, err = boardUseCase.ToggleUser(uid, bid, toggledUserId)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't toggle
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToBoard(toggledUserId, bid).Return(customErrors.ErrInternal)
	_, err = boardUseCase.ToggleUser(uid, bid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// get team error
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToBoard(toggledUserId, bid).Return(nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.ToggleUser(uid, bid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateAccessPathBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	bid := uint(1)
	newAccessPath := "new path"

	// good
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().UpdateAccessPath(bid).Return(newAccessPath, nil)
	resAccessPath, err := boardUseCase.UpdateAccessPath(uid, bid)
	assert.NoError(t, err)
	assert.Equal(t, newAccessPath, resAccessPath)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
	_, err = boardUseCase.UpdateAccessPath(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	_, err = boardUseCase.UpdateAccessPath(uid, bid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().UpdateAccessPath(bid).Return("", customErrors.ErrInternal)
	_, err = boardUseCase.UpdateAccessPath(uid, bid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestAddUserViaLinkBoard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock := createBoardRepoMocks(ctrl)
	boardUseCase := CreateBoardUseCase(boardRepoMock, userRepoMock, teamRepoMock, cardListRepoMock, cardRepoMock, checkListRepoMock)

	uid := uint(1)
	bid := uint(1)
	accessPath := "accessPath"

	testBoard := new(models.Board)
	err := faker.FakeData(testBoard)
	assert.NoError(t, err)

	// success
	boardRepoMock.EXPECT().FindBoardIDByPath(accessPath).Return(bid, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	userRepoMock.EXPECT().AddUserToBoard(uid, bid).Return(nil)

	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(testBoard, nil)
	boardRepoMock.EXPECT().GetBoardMembers(testBoard).Return(&testBoard.Members, nil)
	boardRepoMock.EXPECT().GetBoardInvitedMembers(bid).Return(&testBoard.InvitedMembers, nil)
	boardRepoMock.EXPECT().GetBoardCardLists(bid).Return(&testBoard.CardLists, nil)
	boardRepoMock.EXPECT().GetBoardTags(bid).Return(&testBoard.Tags, nil)
	for _, cardList := range testBoard.CardLists {
		cardListRepoMock.EXPECT().GetCardListCards(cardList.CLID).Return(&cardList.Cards, nil)
		for _, card := range cardList.Cards {
			cardRepoMock.EXPECT().GetCardComments(card.CID).Return(&card.Comments, nil)
			cardRepoMock.EXPECT().GetAssignedUsers(card.CID).Return(&card.Attachments, nil)
			cardRepoMock.EXPECT().GetAssignedUsers(card.CID).Return(&card.Assignees, nil)
			for _, comment := range card.Comments {
				userRepoMock.EXPECT().GetPublicData(comment.UID).Return(&comment.User, nil)
			}
			cardRepoMock.EXPECT().GetCardTags(card.CID).Return(&card.Tags, nil)
			cardRepoMock.EXPECT().GetCardCheckLists(card.CID).Return(&card.CheckLists, nil)
			for _, checkList := range card.CheckLists {
				checkListRepoMock.EXPECT().GetCheckListItems(checkList.CHLID).Return(&checkList.CheckListItems, nil)
			}
		}
	}
	resBoard, err := boardUseCase.AddUserViaLink(uid, accessPath)
	assert.NoError(t, err)
	assert.Equal(t, testBoard, resBoard)

	// error can't find board
	boardRepoMock.EXPECT().FindBoardIDByPath(accessPath).Return(uint(0), customErrors.ErrBoardNotFound)
	_, err = boardUseCase.AddUserViaLink(uid, accessPath)
	assert.Equal(t, customErrors.ErrBoardNotFound, err)

	// error while checking access
	boardRepoMock.EXPECT().FindBoardIDByPath(accessPath).Return(bid, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, customErrors.ErrInternal)
	_, err = boardUseCase.AddUserViaLink(uid, accessPath)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't toggle
	boardRepoMock.EXPECT().FindBoardIDByPath(accessPath).Return(bid, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	userRepoMock.EXPECT().AddUserToBoard(uid, bid).Return(customErrors.ErrInternal)
	_, err = boardUseCase.AddUserViaLink(uid, accessPath)
	assert.Equal(t, customErrors.ErrInternal, err)

	// get team error
	boardRepoMock.EXPECT().FindBoardIDByPath(accessPath).Return(bid, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(false, nil)
	userRepoMock.EXPECT().AddUserToBoard(uid, bid).Return(nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, bid).Return(true, nil)
	boardRepoMock.EXPECT().GetByID(bid).Return(nil, customErrors.ErrInternal)
	_, err = boardUseCase.AddUserViaLink(uid, accessPath)
	assert.Equal(t, customErrors.ErrInternal, err)
}
