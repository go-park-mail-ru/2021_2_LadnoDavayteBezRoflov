package usecases

import (
	"backendServer/app/api/models"
)

type AttachmentUseCase interface {
	CreateAttachment(file *multipart.FileHeader, attachment *models.Attachment) (attachment *models.Attachment, err error)
	GetAttachment(atid uint) (attachment *models.Attachment, err error)
	DeleteAttachment(atid uint) (err error)
}
