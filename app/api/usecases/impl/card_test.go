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

func createCardRepoMocks(controller *gomock.Controller) (*mocks.MockCardRepository, *mocks.MockUserRepository, *mocks.MockTagRepository) {
	cardRepoMock := mocks.NewMockCardRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	tagRepoMock := mocks.NewMockTagRepository(controller)
	return cardRepoMock, userRepoMock, tagRepoMock
}

func TestCreateCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)

	// good
	cardRepoMock.EXPECT().Create(testCard).Return(nil)
	resCid, err := cardUseCase.CreateCard(testCard)
	assert.NoError(t, err)
	assert.Equal(t, testCard.CID, resCid)

	// error can't create
	cardRepoMock.EXPECT().Create(testCard).Return(customErrors.ErrInternal)
	resCid, err = cardUseCase.CreateCard(testCard)
	assert.Equal(t, uint(0), resCid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	cid := uint(1)

	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	testComments := new([]models.Comment)
	testComment := new(models.Comment)
	err = faker.FakeData(testComment)
	assert.NoError(t, err)
	*testComments = append(*testComments, *testComment)
	testUsers := new([]models.PublicUserInfo)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testCard.Comments = *testComments
	testTag := new(models.Tag)
	err = faker.FakeData(testTag)
	assert.NoError(t, err)
	testTag.ColorID = uint(2)
	testCard.Tags = append(testCard.Tags, *testTag)

	// success
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(testComments, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(testUsers, nil)
	resCard, err := cardUseCase.GetCard(uid, cid)
	assert.NoError(t, err)
	assert.Equal(t, testCard, resCard)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(nil, customErrors.ErrCardNotFound)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrCardNotFound, err)

	// can't get card's tags
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(nil, customErrors.ErrCardNotFound)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrCardNotFound, err)

	// can't get cards comments
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(nil, customErrors.ErrInternal)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get cards comment's user info
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(testComments, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(nil, customErrors.ErrInternal)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// can't get cards comments
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(testComments, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(nil, customErrors.ErrInternal)
	_, err = cardUseCase.GetCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsCardAccessed(uid, testCard.CID).Return(true, nil)
	cardRepoMock.EXPECT().Update(testCard).Return(nil)
	err = cardUseCase.UpdateCard(uid, testCard)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, testCard.CID).Return(false, customErrors.ErrInternal)
	err = cardUseCase.UpdateCard(uid, testCard)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, testCard.CID).Return(false, nil)
	err = cardUseCase.UpdateCard(uid, testCard)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCardAccessed(uid, testCard.CID).Return(true, nil)
	cardRepoMock.EXPECT().Update(testCard).Return(customErrors.ErrInternal)
	err = cardUseCase.UpdateCard(uid, testCard)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	cid := uint(1)

	// good
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().Delete(cid).Return(nil)
	err := cardUseCase.DeleteCard(uid, cid)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	err = cardUseCase.DeleteCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	err = cardUseCase.DeleteCard(uid, cid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().Delete(cid).Return(customErrors.ErrInternal)
	err = cardUseCase.DeleteCard(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestToggleUserCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	cid := uint(1)
	toggledUserId := uint(1)
	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	testComments := new([]models.Comment)
	testComment := new(models.Comment)
	err = faker.FakeData(testComment)
	assert.NoError(t, err)
	*testComments = append(*testComments, *testComment)
	testUsers := new([]models.PublicUserInfo)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testCard.Comments = *testComments
	testTag := new(models.Tag)
	err = faker.FakeData(testTag)
	assert.NoError(t, err)
	testTag.ColorID = uint(2)
	testCard.Tags = append(testCard.Tags, *testTag)

	// success
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToCard(toggledUserId, cid).Return(nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(testComments, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(testUsers, nil)
	resCard, err := cardUseCase.ToggleUser(uid, cid, toggledUserId)
	assert.NoError(t, err)
	assert.Equal(t, testCard, resCard)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	_, err = cardUseCase.ToggleUser(uid, cid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	_, err = cardUseCase.ToggleUser(uid, cid, toggledUserId)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't toggle
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToCard(toggledUserId, cid).Return(customErrors.ErrInternal)
	_, err = cardUseCase.ToggleUser(uid, cid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// get card fail
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	userRepoMock.EXPECT().AddUserToCard(toggledUserId, cid).Return(nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(nil, customErrors.ErrInternal)
	_, err = cardUseCase.ToggleUser(uid, cid, toggledUserId)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestToggleTagCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	cid := uint(1)
	toggledCardId := uint(1)
	testCard := new(models.Card)
	err := faker.FakeData(testCard)
	assert.NoError(t, err)
	testComments := new([]models.Comment)
	testComment := new(models.Comment)
	err = faker.FakeData(testComment)
	assert.NoError(t, err)
	*testComments = append(*testComments, *testComment)
	testUsers := new([]models.PublicUserInfo)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testCard.Comments = *testComments
	testTag := new(models.Tag)
	err = faker.FakeData(testTag)
	assert.NoError(t, err)
	testTag.ColorID = uint(2)
	testCard.Tags = append(testCard.Tags, *testTag)

	// success
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	tagRepoMock.EXPECT().AddTagToCard(toggledCardId, cid).Return(nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(testCard, nil)
	cardRepoMock.EXPECT().GetCardTags(cid).Return(&testCard.Tags, nil)
	cardRepoMock.EXPECT().GetCardComments(cid).Return(testComments, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	cardRepoMock.EXPECT().GetAssignedUsers(cid).Return(testUsers, nil)
	resCard, err := cardUseCase.ToggleTag(uid, cid, toggledCardId)
	assert.NoError(t, err)
	assert.Equal(t, testCard, resCard)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	_, err = cardUseCase.ToggleTag(uid, cid, toggledCardId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	_, err = cardUseCase.ToggleTag(uid, cid, toggledCardId)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't toggle
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	tagRepoMock.EXPECT().AddTagToCard(toggledCardId, cid).Return(customErrors.ErrInternal)
	_, err = cardUseCase.ToggleTag(uid, cid, toggledCardId)
	assert.Equal(t, customErrors.ErrInternal, err)

	// get card fail
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	tagRepoMock.EXPECT().AddTagToCard(toggledCardId, cid).Return(nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().GetByID(cid).Return(nil, customErrors.ErrInternal)
	_, err = cardUseCase.ToggleTag(uid, cid, toggledCardId)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateAccessPathCard(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardRepoMock, userRepoMock, tagRepoMock := createCardRepoMocks(ctrl)
	cardUseCase := CreateCardUseCase(cardRepoMock, userRepoMock, tagRepoMock)

	uid := uint(1)
	cid := uint(1)
	newAccessPath := "new path"

	// good
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().UpdateAccessPath(cid).Return(newAccessPath, nil)
	resAccessPath, err := cardUseCase.UpdateAccessPath(uid, cid)
	assert.NoError(t, err)
	assert.Equal(t, newAccessPath, resAccessPath)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, customErrors.ErrInternal)
	_, err = cardUseCase.UpdateAccessPath(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(false, nil)
	_, err = cardUseCase.UpdateAccessPath(uid, cid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCardAccessed(uid, cid).Return(true, nil)
	cardRepoMock.EXPECT().UpdateAccessPath(cid).Return("", customErrors.ErrInternal)
	_, err = cardUseCase.UpdateAccessPath(uid, cid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
