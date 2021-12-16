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

	attachment, err = attachmentUseCase.attachmentRepository.Create(file, cid)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) GetAttachment(atid, uid uint) (attachment *models.Attachment, err error) {
	attachment, err = attachmentUseCase.attachmentRepository.Get(atid)
	if err != nil {
		return nil, err
	}

	isAccessed, err := attachmentUseCase.userRepository.IsCardAccessed(uid, attachment.CID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) DeleteAttachment(atid, uid uint) (err error) {
	attachment, err := attachmentUseCase.attachmentRepository.Get(atid)
	if err != nil {
		return err
	}

	isAccessed, err := attachmentUseCase.userRepository.IsCardAccessed(uid, attachment.CID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return attachmentUseCase.attachmentRepository.Delete(atid)
}
