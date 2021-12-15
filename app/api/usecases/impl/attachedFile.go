package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"time"
)

type AttachmentUseCaseImpl struct {
	attachmentRepository repositories.AttachmentRepository
	userRepository       repositories.UserRepository
}

func CreateAttachmentUseCase(
	attachmentRepository repositories.AttachmentRepository,
) usecases.AttachmentUseCase {
	return &AttachmentUseCaseImpl{
		attachmentRepository: attachmentRepository,
	}
}

func (attachmentUseCase *AttachmentUseCaseImpl) CreateAttachment(file *multipart.FileHeader, attachment *models.Attachment) (attachment *models.Attachment, err error) {
	err = attachmentUseCase.attachmentRepository.Create(file, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) GetAttachment(atid uint) (attachment *models.Attachment, err error) {
	attachment, err = attachmentUseCase.attachmentRepository.GetByID(atid)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (attachmentUseCase *AttachmentUseCaseImpl) DeleteAttachment(atid uint) (err error) {
	return attachmentUseCase.attachmentRepository.Delete(atid)
}
