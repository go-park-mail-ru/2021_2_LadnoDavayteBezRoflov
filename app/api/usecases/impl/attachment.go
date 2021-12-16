package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"

	"mime/multipart"
)

type AttachmentUseCaseImpl struct {
	attachmentRepository repositories.AttachmentRepository
	userRepository       repositories.UserRepository
}

func CreateAttachmentUseCase(
	attachmentRepository repositories.AttachmentRepository,
	userRepository repositories.UserRepository,
) usecases.AttachmentUseCase {
	return &AttachmentUseCaseImpl{
		attachmentRepository: attachmentRepository,
		userRepository:       userRepository,
	}
}

func (attachmentUseCase *AttachmentUseCaseImpl) CreateAttachment(file *multipart.FileHeader, cid, uid uint) (attachment *models.Attachment, err error) {
	isAccessed, err := attachmentUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	attachment, err = attachmentUseCase.attachmentRepository.Create(file, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) GetAttachment(atid, cid, uid uint) (attachment *models.Attachment, err error) {
	isAccessed, err := attachmentUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	attachment, err = attachmentUseCase.attachmentRepository.Get(atid)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) DeleteAttachment(atid, cid, uid uint) (err error) {
	isAccessed, err := attachmentUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return attachmentUseCase.attachmentRepository.Delete(atid)
}
