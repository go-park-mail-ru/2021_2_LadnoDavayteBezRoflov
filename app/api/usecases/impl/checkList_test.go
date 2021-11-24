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

func createCheckListRepoMocks(controller *gomock.Controller) (*mocks.MockCheckListRepository, *mocks.MockUserRepository) {
	checkListRepoMock := mocks.NewMockCheckListRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return checkListRepoMock, userRepoMock
}

func TestCreateCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListRepoMock, userRepoMock := createCheckListRepoMocks(ctrl)
	checkListUseCase := CreateCheckListUseCase(checkListRepoMock, userRepoMock)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	assert.NoError(t, err)

	// good
	checkListRepoMock.EXPECT().Create(testCheckList).Return(nil)
	resChlid, err := checkListUseCase.CreateCheckList(testCheckList)
	assert.NoError(t, err)
	assert.Equal(t, testCheckList.CHLID, resChlid)

	// error can't create
	checkListRepoMock.EXPECT().Create(testCheckList).Return(customErrors.ErrInternal)
	resChlid, err = checkListUseCase.CreateCheckList(testCheckList)
	assert.Equal(t, uint(0), resChlid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListRepoMock, userRepoMock := createCheckListRepoMocks(ctrl)
	checkListUseCase := CreateCheckListUseCase(checkListRepoMock, userRepoMock)

	uid := uint(1)
	chlid := uint(1)

	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	assert.NoError(t, err)
	testCheckListItems := make([]models.CheckListItem, 3)
	for i := range testCheckListItems {
		err = faker.FakeData(&testCheckListItems[i])
		assert.NoError(t, err)
	}
	testCheckList.CheckListItems = testCheckListItems

	// success
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(true, nil)
	checkListRepoMock.EXPECT().GetByID(chlid).Return(testCheckList, nil)
	checkListRepoMock.EXPECT().GetCheckListItems(chlid).Return(&testCheckListItems, nil)
	resCheckList, err := checkListUseCase.GetCheckList(uid, chlid)
	assert.NoError(t, err)
	assert.Equal(t, testCheckList, resCheckList)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(false, customErrors.ErrInternal)
	_, err = checkListUseCase.GetCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(false, nil)
	_, err = checkListUseCase.GetCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(true, nil)
	checkListRepoMock.EXPECT().GetByID(chlid).Return(nil, customErrors.ErrCheckListNotFound)
	_, err = checkListUseCase.GetCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrCheckListNotFound, err)

	// can't get check list items
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(true, nil)
	checkListRepoMock.EXPECT().GetByID(chlid).Return(testCheckList, nil)
	checkListRepoMock.EXPECT().GetCheckListItems(chlid).Return(nil, customErrors.ErrInternal)
	_, err = checkListUseCase.GetCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListRepoMock, userRepoMock := createCheckListRepoMocks(ctrl)
	checkListUseCase := CreateCheckListUseCase(checkListRepoMock, userRepoMock)

	uid := uint(1)
	testCheckList := new(models.CheckList)
	err := faker.FakeData(testCheckList)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsCheckListAccessed(uid, testCheckList.CHLID).Return(true, nil)
	checkListRepoMock.EXPECT().Update(testCheckList).Return(nil)
	err = checkListUseCase.UpdateCheckList(uid, testCheckList)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, testCheckList.CHLID).Return(false, customErrors.ErrInternal)
	err = checkListUseCase.UpdateCheckList(uid, testCheckList)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, testCheckList.CHLID).Return(false, nil)
	err = checkListUseCase.UpdateCheckList(uid, testCheckList)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCheckListAccessed(uid, testCheckList.CHLID).Return(true, nil)
	checkListRepoMock.EXPECT().Update(testCheckList).Return(customErrors.ErrInternal)
	err = checkListUseCase.UpdateCheckList(uid, testCheckList)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteCheckList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	checkListRepoMock, userRepoMock := createCheckListRepoMocks(ctrl)
	checkListUseCase := CreateCheckListUseCase(checkListRepoMock, userRepoMock)

	uid := uint(1)
	chlid := uint(1)

	// good
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(true, nil)
	checkListRepoMock.EXPECT().Delete(chlid).Return(nil)
	err := checkListUseCase.DeleteCheckList(uid, chlid)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(false, customErrors.ErrInternal)
	err = checkListUseCase.DeleteCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(false, nil)
	err = checkListUseCase.DeleteCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	userRepoMock.EXPECT().IsCheckListAccessed(uid, chlid).Return(true, nil)
	checkListRepoMock.EXPECT().Delete(chlid).Return(customErrors.ErrInternal)
	err = checkListUseCase.DeleteCheckList(uid, chlid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
