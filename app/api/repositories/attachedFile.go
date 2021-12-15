package repositories

import (
	"backendServer/app/api/models"
)

type AttachmentRepository interface {
	Create(file *multipart.FileHeader, attachment *models.AttachedFile) (err error)
	Delete(atid uint) (err error)
	GetAttachment(atid uint) (attachment *models.AttachedFile, err error)
}
