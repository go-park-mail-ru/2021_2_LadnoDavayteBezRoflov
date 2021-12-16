package repositories

import (
	"backendServer/app/api/models"
)

type AttachmentRepository interface {
	Create(file *multipart.FileHeader, cid uint) (attachment *models.AttachedFile, err error)
	Delete(atid uint) (err error)
	Get(atid uint) (attachment *models.AttachedFile, err error)
}
