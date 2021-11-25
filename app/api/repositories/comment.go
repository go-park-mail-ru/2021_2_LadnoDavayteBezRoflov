package repositories

import (
	"backendServer/app/api/models"
)

type CommentRepository interface {
	Create(comment *models.Comment) (err error)
	Update(comment *models.Comment) (err error)
	Delete(cmid uint) (err error)
	GetByID(cmid uint) (comment *models.Comment, err error)
}
