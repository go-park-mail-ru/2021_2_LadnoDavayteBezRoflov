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

func createCheckListItemRepoMocks(controller *gomock.Controller) (*mocks.MockCheckListItemRepository, *mocks.MockUserRepository) {
	checkListItemRepoMock := mocks.NewMockCheckListItemRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return checkListItemRepoMock, userRepoMock
}

func TestCreateCheckListItem(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListItemRepoMock, userRepoMock := createCheckListItemRepoMocks(ctrl)
	checkListItemUseCase := CreateCheckListItemUseCase(checkListItemRepoMock, userRepoMock)

	testCheckListItem := new(models.CheckListItem)
	err := faker.FakeData(testCheckListItem)
	assert.NoError(t, err)

	// good
	checkListItemRepoMock.EXPECT().Create(testCheckListItem).Return(nil)
	resChliid, err := checkListItemUseCase.CreateCheckListItem(testCheckListItem)
	assert.NoError(t, err)
	assert.Equal(t, testCheckListItem.CHLIID, resChliid)

	// error can't create
	checkListItemRepoMock.EXPECT().Create(testCheckListItem).Return(customErrors.ErrInternal)
	resChliid, err = checkListItemUseCase.CreateCheckListItem(testCheckListItem)
	assert.Equal(t, uint(0), resChliid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetCheckListItem(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListItemRepoMock, userRepoMock := createCheckListItemRepoMocks(ctrl)
	checkListItemUseCase := CreateCheckListItemUseCase(checkListItemRepoMock, userRepoMock)

	uid := uint(1)
	chliid := uint(1)

	testCheckListItem := new(models.CheckListItem)
	err := faker.FakeData(testCheckListItem)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(true, nil)
	checkListItemRepoMock.EXPECT().GetByID(chliid).Return(testCheckListItem, nil)
	resCheckListItem, err := checkListItemUseCase.GetCheckListItem(uid, chliid)
	assert.NoError(t, err)
	assert.Equal(t, testCheckListItem, resCheckListItem)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(false, customErrors.ErrInternal)
	_, err = checkListItemUseCase.GetCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(false, nil)
	_, err = checkListItemUseCase.GetCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(true, nil)
	checkListItemRepoMock.EXPECT().GetByID(chliid).Return(nil, customErrors.ErrCheckListItemNotFound)
	_, err = checkListItemUseCase.GetCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrCheckListItemNotFound, err)
}

func TestUpdateCheckListItem(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListItemRepoMock, userRepoMock := createCheckListItemRepoMocks(ctrl)
	checkListItemUseCase := CreateCheckListItemUseCase(checkListItemRepoMock, userRepoMock)

	uid := uint(1)
	testCheckListItem := new(models.CheckListItem)
	err := faker.FakeData(testCheckListItem)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, testCheckListItem.CHLIID).Return(true, nil)
	checkListItemRepoMock.EXPECT().Update(testCheckListItem).Return(nil)
	err = checkListItemUseCase.UpdateCheckListItem(uid, testCheckListItem)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, testCheckListItem.CHLIID).Return(false, customErrors.ErrInternal)
	err = checkListItemUseCase.UpdateCheckListItem(uid, testCheckListItem)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, testCheckListItem.CHLIID).Return(false, nil)
	err = checkListItemUseCase.UpdateCheckListItem(uid, testCheckListItem)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, testCheckListItem.CHLIID).Return(true, nil)
	checkListItemRepoMock.EXPECT().Update(testCheckListItem).Return(customErrors.ErrInternal)
	err = checkListItemUseCase.UpdateCheckListItem(uid, testCheckListItem)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteCheckListItem(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListItemRepoMock, userRepoMock := createCheckListItemRepoMocks(ctrl)
	checkListItemUseCase := CreateCheckListItemUseCase(checkListItemRepoMock, userRepoMock)

	uid := uint(1)
	chliid := uint(1)

	// good
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(true, nil)
	checkListItemRepoMock.EXPECT().Delete(chliid).Return(nil)
	err := checkListItemUseCase.DeleteCheckListItem(uid, chliid)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(false, customErrors.ErrInternal)
	err = checkListItemUseCase.DeleteCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(false, nil)
	err = checkListItemUseCase.DeleteCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	userRepoMock.EXPECT().IsCheckListItemAccessed(uid, chliid).Return(true, nil)
	checkListItemRepoMock.EXPECT().Delete(chliid).Return(customErrors.ErrInternal)
	err = checkListItemUseCase.DeleteCheckListItem(uid, chliid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
