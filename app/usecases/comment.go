package usecases

import "backendServer/app/models"

type CommentUseCase interface {
	CreateComment(comment *models.Comment) (finalComment *models.Comment, err error)
	GetComment(uid, cmid uint) (comment *models.Comment, err error)
	UpdateComment(uid uint, comment *models.Comment) (err error)
	DeleteComment(uid, cmid uint) (err error)
}
