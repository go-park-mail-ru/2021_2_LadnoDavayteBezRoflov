package repositories

import (
	"backendServer/app/api/models"
	"mime/multipart"
)

type AttachmentRepository interface {
	Create(file *multipart.FileHeader, cid uint) (attachment *models.Attachment, err error)
	Delete(atid uint) (err error)
	Get(atid uint) (attachment *models.Attachment, err error)
}
