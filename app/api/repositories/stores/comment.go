package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	"gorm.io/gorm"
)

type CommentStore struct {
	db *gorm.DB
}

func CreateCommentRepository(db *gorm.DB) repositories.CommentRepository {
	return &CommentStore{db: db}
}

func (commentStore *CommentStore) Create(comment *models.Comment) (err error) {
	return commentStore.db.Create(comment).Error
}

func (commentStore *CommentStore) Update(comment *models.Comment) (err error) {
	oldComment, err := commentStore.GetByID(comment.CMID)
	if err != nil {
		return
	}

	if comment.Text != "" && comment.Text != oldComment.Text {
		oldComment.Text = comment.Text
	}

	return commentStore.db.Save(oldComment).Error
}

func (commentStore *CommentStore) Delete(cmid uint) (err error) {
	return commentStore.db.Delete(&models.Comment{}, cmid).Error
}

func (commentStore *CommentStore) GetByID(cmid uint) (*models.Comment, error) {
	comment := new(models.Comment)
	if res := commentStore.db.Find(comment, cmid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCommentNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}
