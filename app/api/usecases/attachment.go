package usecases

import (
	"backendServer/app/api/models"
	"mime/multipart"
)

type AttachmentUseCase interface {
	CreateAttachment(file *multipart.FileHeader, cid, uid uint) (attachment *models.Attachment, err error)
	GetAttachment(atid, uid uint) (attachment *models.Attachment, err error)
	DeleteAttachment(atid, uid uint) (err error)
}
