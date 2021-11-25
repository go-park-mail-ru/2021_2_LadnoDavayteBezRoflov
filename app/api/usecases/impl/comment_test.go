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

func createCommentRepoMocks(controller *gomock.Controller) (*mocks.MockCommentRepository, *mocks.MockUserRepository) {
	commentRepoMock := mocks.NewMockCommentRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return commentRepoMock, userRepoMock
}

func TestCreateComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepoMock, userRepoMock := createCommentRepoMocks(ctrl)
	commentUseCase := CreateCommentUseCase(commentRepoMock, userRepoMock)

	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	assert.NoError(t, err)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testComment.User = *testUser

	// good
	commentRepoMock.EXPECT().Create(testComment).Return(nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	resComment, err := commentUseCase.CreateComment(testComment)
	assert.NoError(t, err)
	assert.Equal(t, testComment, resComment)

	// error can't create
	commentRepoMock.EXPECT().Create(testComment).Return(customErrors.ErrInternal)
	_, err = commentUseCase.CreateComment(testComment)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error can't get public data
	commentRepoMock.EXPECT().Create(testComment).Return(nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(nil, customErrors.ErrInternal)
	_, err = commentUseCase.CreateComment(testComment)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepoMock, userRepoMock := createCommentRepoMocks(ctrl)
	commentUseCase := CreateCommentUseCase(commentRepoMock, userRepoMock)

	uid := uint(1)
	cmid := uint(1)

	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	assert.NoError(t, err)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testComment.User = *testUser

	// success
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(true, nil)
	commentRepoMock.EXPECT().GetByID(cmid).Return(testComment, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(testUser, nil)
	resComment, err := commentUseCase.GetComment(uid, cmid)
	assert.NoError(t, err)
	assert.Equal(t, testComment, resComment)

	// error while checking access
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(false, customErrors.ErrInternal)
	_, err = commentUseCase.GetComment(uid, cmid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(false, nil)
	_, err = commentUseCase.GetComment(uid, cmid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(true, nil)
	commentRepoMock.EXPECT().GetByID(cmid).Return(nil, customErrors.ErrCommentNotFound)
	_, err = commentUseCase.GetComment(uid, cmid)
	assert.Equal(t, customErrors.ErrCommentNotFound, err)

	// error can't get public data
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(true, nil)
	commentRepoMock.EXPECT().GetByID(cmid).Return(testComment, nil)
	userRepoMock.EXPECT().GetPublicData(testComment.UID).Return(nil, customErrors.ErrInternal)
	_, err = commentUseCase.GetComment(uid, cmid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepoMock, userRepoMock := createCommentRepoMocks(ctrl)
	commentUseCase := CreateCommentUseCase(commentRepoMock, userRepoMock)

	uid := uint(1)
	testComment := new(models.Comment)
	err := faker.FakeData(testComment)
	assert.NoError(t, err)
	testUser := new(models.PublicUserInfo)
	err = faker.FakeData(testUser)
	assert.NoError(t, err)
	testComment.User = *testUser

	// success
	userRepoMock.EXPECT().IsCommentAccessed(uid, testComment.CMID).Return(true, nil)
	commentRepoMock.EXPECT().Update(testComment).Return(nil)
	err = commentUseCase.UpdateComment(uid, testComment)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCommentAccessed(uid, testComment.CMID).Return(false, customErrors.ErrInternal)
	err = commentUseCase.UpdateComment(uid, testComment)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCommentAccessed(uid, testComment.CMID).Return(false, nil)
	err = commentUseCase.UpdateComment(uid, testComment)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCommentAccessed(uid, testComment.CMID).Return(true, nil)
	commentRepoMock.EXPECT().Update(testComment).Return(customErrors.ErrInternal)
	err = commentUseCase.UpdateComment(uid, testComment)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepoMock, userRepoMock := createCommentRepoMocks(ctrl)
	commentUseCase := CreateCommentUseCase(commentRepoMock, userRepoMock)

	uid := uint(1)
	cmid := uint(1)

	// good
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(true, nil)
	commentRepoMock.EXPECT().Delete(cmid).Return(nil)
	err := commentUseCase.DeleteComment(uid, cmid)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(false, customErrors.ErrInternal)
	err = commentUseCase.DeleteComment(uid, cmid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(false, nil)
	err = commentUseCase.DeleteComment(uid, cmid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	userRepoMock.EXPECT().IsCommentAccessed(uid, cmid).Return(true, nil)
	commentRepoMock.EXPECT().Delete(cmid).Return(customErrors.ErrInternal)
	err = commentUseCase.DeleteComment(uid, cmid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
