package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"mime/multipart"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func createAttachmentRepoMocks(controller *gomock.Controller) (*mocks.MockAttachmentRepository, *mocks.MockUserRepository) {
	attachmentRepoMock := mocks.NewMockAttachmentRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return attachmentRepoMock, userRepoMock
}

func TestCreateAttachment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	attachmentRepoMock, userRepoMock := createAttachmentRepoMocks(ctrl)
	attachmentUseCase := CreateAttachmentUseCase(attachmentRepoMock, userRepoMock)

	uid := uint(1)
	testAttachment := new(models.Attachment)
	err := faker.FakeData(testAttachment)
	assert.NoError(t, err)

	file := new(multipart.FileHeader)

	// good
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(true, nil)
	attachmentRepoMock.EXPECT().Create(file, testAttachment.CID).Return(testAttachment, nil)
	resAttachment, err := attachmentUseCase.CreateAttachment(file, testAttachment.CID, uid)
	assert.NoError(t, err)
	assert.Equal(t, testAttachment, resAttachment)

	// error while checking access
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, customErrors.ErrInternal)
	_, err = attachmentUseCase.CreateAttachment(file, testAttachment.CID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, nil)
	_, err = attachmentUseCase.CreateAttachment(file, testAttachment.CID, uid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't create
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(true, nil)
	attachmentRepoMock.EXPECT().Create(file, testAttachment.CID).Return(nil, customErrors.ErrInternal)
	_, err = attachmentUseCase.CreateAttachment(file, testAttachment.CID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetAttachment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	attachmentRepoMock, userRepoMock := createAttachmentRepoMocks(ctrl)
	attachmentUseCase := CreateAttachmentUseCase(attachmentRepoMock, userRepoMock)

	uid := uint(1)
	testAttachment := new(models.Attachment)
	err := faker.FakeData(testAttachment)
	assert.NoError(t, err)

	// good
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(true, nil)
	resAttachment, err := attachmentUseCase.GetAttachment(testAttachment.ATID, uid)
	assert.NoError(t, err)
	assert.Equal(t, testAttachment, resAttachment)

	// error can't get
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(nil, customErrors.ErrInternal)
	_, err = attachmentUseCase.GetAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error while checking access
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, customErrors.ErrInternal)
	_, err = attachmentUseCase.GetAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, nil)
	_, err = attachmentUseCase.GetAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrNoAccess, err)
}

func TestDeleteAttachment(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	attachmentRepoMock, userRepoMock := createAttachmentRepoMocks(ctrl)
	attachmentUseCase := CreateAttachmentUseCase(attachmentRepoMock, userRepoMock)

	uid := uint(1)
	testAttachment := new(models.Attachment)
	err := faker.FakeData(testAttachment)
	assert.NoError(t, err)

	// good
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(true, nil)
	attachmentRepoMock.EXPECT().Delete(testAttachment.ATID).Return(nil)
	err = attachmentUseCase.DeleteAttachment(testAttachment.ATID, uid)
	assert.NoError(t, err)

	// error can't get
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(nil, customErrors.ErrInternal)
	err = attachmentUseCase.DeleteAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// error while checking access
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, customErrors.ErrInternal)
	err = attachmentUseCase.DeleteAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(false, nil)
	err = attachmentUseCase.DeleteAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	attachmentRepoMock.EXPECT().Get(testAttachment.ATID).Return(testAttachment, nil)
	userRepoMock.EXPECT().IsCardAccessed(uid, testAttachment.CID).Return(true, nil)
	attachmentRepoMock.EXPECT().Delete(testAttachment.ATID).Return(customErrors.ErrInternal)
	err = attachmentUseCase.DeleteAttachment(testAttachment.ATID, uid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
