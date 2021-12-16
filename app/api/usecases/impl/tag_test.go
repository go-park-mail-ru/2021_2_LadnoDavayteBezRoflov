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

func createTagRepoMocks(controller *gomock.Controller) (*mocks.MockTagRepository, *mocks.MockUserRepository) {
	tagRepoMock := mocks.NewMockTagRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return tagRepoMock, userRepoMock
}

func TestCreateTag(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagRepoMock, userRepoMock := createTagRepoMocks(ctrl)
	tagUseCase := CreateTagUseCase(tagRepoMock, userRepoMock)

	testTag := new(models.Tag)
	err := faker.FakeData(testTag)
	assert.NoError(t, err)

	// good
	tagRepoMock.EXPECT().Create(testTag).Return(nil)
	resTagID, err := tagUseCase.CreateTag(testTag)
	assert.NoError(t, err)
	assert.Equal(t, testTag.TGID, resTagID)

	// error can't create
	tagRepoMock.EXPECT().Create(testTag).Return(customErrors.ErrInternal)
	_, err = tagUseCase.CreateTag(testTag)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetTag(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagRepoMock, userRepoMock := createTagRepoMocks(ctrl)
	tagUseCase := CreateTagUseCase(tagRepoMock, userRepoMock)

	uid := uint(1)
	tgid := uint(1)

	testTag := new(models.Tag)
	err := faker.FakeData(testTag)
	assert.NoError(t, err)

	// success
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(true, nil)
	resTag, err := tagUseCase.GetTag(uid, tgid)
	assert.NoError(t, err)
	assert.Equal(t, testTag, resTag)

	// can't found
	tagRepoMock.EXPECT().GetByID(tgid).Return(nil, customErrors.ErrInternal)
	_, err = tagUseCase.GetTag(uid, tgid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error while checking access
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, customErrors.ErrInternal)
	_, err = tagUseCase.GetTag(uid, tgid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, nil)
	_, err = tagUseCase.GetTag(uid, tgid)
	assert.Equal(t, customErrors.ErrNoAccess, err)
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagRepoMock, userRepoMock := createTagRepoMocks(ctrl)
	tagUseCase := CreateTagUseCase(tagRepoMock, userRepoMock)

	uid := uint(1)
	testTag := new(models.Tag)
	err := faker.FakeData(testTag)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(true, nil)
	tagRepoMock.EXPECT().Update(testTag).Return(nil)
	err = tagUseCase.UpdateTag(uid, testTag)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, customErrors.ErrInternal)
	err = tagUseCase.UpdateTag(uid, testTag)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, nil)
	err = tagUseCase.UpdateTag(uid, testTag)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(true, nil)
	tagRepoMock.EXPECT().Update(testTag).Return(customErrors.ErrInternal)
	err = tagUseCase.UpdateTag(uid, testTag)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagRepoMock, userRepoMock := createTagRepoMocks(ctrl)
	tagUseCase := CreateTagUseCase(tagRepoMock, userRepoMock)

	uid := uint(1)
	tgid := uint(1)

	testTag := new(models.Tag)
	err := faker.FakeData(testTag)
	assert.NoError(t, err)

	// good
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(true, nil)
	tagRepoMock.EXPECT().Delete(tgid).Return(nil)
	err = tagUseCase.DeleteTag(uid, tgid)
	assert.NoError(t, err)

	// error can't get
	tagRepoMock.EXPECT().GetByID(tgid).Return(nil, customErrors.ErrInternal)
	err = tagUseCase.DeleteTag(uid, tgid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error while checking access
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, customErrors.ErrInternal)
	err = tagUseCase.DeleteTag(uid, tgid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(false, nil)
	err = tagUseCase.DeleteTag(uid, tgid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	tagRepoMock.EXPECT().GetByID(tgid).Return(testTag, nil)
	userRepoMock.EXPECT().IsBoardAccessed(uid, testTag.BID).Return(true, nil)
	tagRepoMock.EXPECT().Delete(tgid).Return(customErrors.ErrInternal)
	err = tagUseCase.DeleteTag(uid, tgid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
