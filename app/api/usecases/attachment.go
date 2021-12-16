package usecases

import (
	"backendServer/app/api/models"
)

type AttachmentUseCase interface {
	CreateAttachment(file *multipart.FileHeader, cid, uid uint) (attachment *models.Attachment, err error)
	GetAttachment(atid, cid, uid uint) (attachment *models.Attachment, err error)
	DeleteAttachment(atid, cid, uid uint) (err error)
}
