package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type AttachmentStore struct {
	db              *gorm.DB
	attachmentsPath string
}

func CreateAttachmentRepository(db *gorm.DB, attachmentsPath string) repositories.AttachmentRepository {
	return &AttachmentStore{db: db, attachmentsPath: attachmentsPath}
}

func (attachmentStore *AttachmentStore) Create(file *multipart.FileHeader, cid uint) (attachment *models.Attachment, err error) {
	attachment = new(models.Attachment)
	attachment.AttachmentPubName = filepath.Base(file.Filename)
	attachment.CID = cid

	fileNameID := uuid.NewString()
	fileName := strings.Join([]string{attachmentStore.attachmentsPath, "/", fileNameID}, "")
	attachment.AttachmentTechName = fileName

	src, err := file.Open()
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

	return attachment, attachmentStore.db.Create(attachment).Error
}

func (attachmentStore *AttachmentStore) Delete(atid uint) (err error) {
	attachment, err := attachmentStore.Get(atid)
	if err != nil {
		return err
	}

	fileToDelete := attachment.AttachmentTechName

	err = os.Remove(fileToDelete)
	if err != nil {
		return err
	}

	return attachmentStore.db.Delete(&models.Attachment{}, atid).Error
}

func (attachmentStore *AttachmentStore) Get(atid uint) (*models.Attachment, error) {
	attachment := new(models.Attachment)
	if res := attachmentStore.db.Find(attachment, atid); res.RowsAffected == 0 {
		return nil, customErrors.ErrAttachmentNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return attachment, nil
}
