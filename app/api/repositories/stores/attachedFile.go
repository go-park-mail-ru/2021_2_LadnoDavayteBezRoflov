package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	"github.com/google/uuid"
	"mime/multipart"
	"os"
	"strings"

	"gorm.io/gorm"
)

type AttachmentStore struct {
	db              *gorm.DB
	attachmentsPath string
}

func CreateAttachmentRepository(db *gorm.DB, attachmentPath) repositories.AttachmentRepository {
	return &AttachmentStore{db: db, attachmentPath: attachmentPath}
}

func (attachmentStore *AttachmentStore) Create(file *multipart.FileHeader, attachment *models.AttachedFile) (err error) {
	fileNameID := uuid.NewString()
	fileName := strings.Join([]string{attachmentStore.attachmentsPath, "/", fileNameID}, "")
	attachment.AttachmentTech = fileName

	// save file
	src, err := attachment.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	return attachmentStore.db.Create(attachment).Error
}

func (attachmentStore *AttachmentStore) Delete(atid uint) (err error) {
	return attachmentStore.db.Delete(&models.AttachedFile{}, atid).Error
}

func (attachmentStore *AttachmentStore) GetAttachment(atid uint) (*models.AttachedFile, error) {
	attachment := new(models.AttachedFile)
	if res := attachmentStore.db.Find(attachment, atid); res.RowsAffected == 0 {
		return nil, customErrors.ErrAttachmentNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return attachment, nil
}
