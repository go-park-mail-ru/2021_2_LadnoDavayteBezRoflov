package usecases

import (
	"backendServer/app/api/models"
)

type CommentUseCase interface {
	CreateComment(comment *models.Comment) (cmid uint, err error)
	GetComment(uid, cmid uint) (comment *models.Comment, err error)
	UpdateComment(uid uint, comment *models.Comment) (err error)
	DeleteComment(uid, cmid uint) (err error)
}
