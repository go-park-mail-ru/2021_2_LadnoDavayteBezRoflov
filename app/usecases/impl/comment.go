package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	customErrors "backendServer/pkg/errors"
	"time"
)

type CommentUseCaseImpl struct {
	commentRepository repositories.CommentRepository
	userRepository    repositories.UserRepository
}

func CreateCommentUseCase(
	commentRepository repositories.CommentRepository,
	userRepository repositories.UserRepository,
) usecases.CommentUseCase {
	return &CommentUseCaseImpl{
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (commentUseCase *CommentUseCaseImpl) CreateComment(comment *models.Comment) (finalComment *models.Comment, err error) {
	comment.Date = time.Now()
	err = commentUseCase.commentRepository.Create(comment)
	if err != nil {
		return nil, err
	}
	comment.DateParsed = comment.Date.Round(time.Second).String()
	return comment, nil
}

func (commentUseCase *CommentUseCaseImpl) GetComment(uid, cmid uint) (comment *models.Comment, err error) {
	isAccessed, err := commentUseCase.userRepository.IsCommentAccessed(uid, cmid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	comment, err = commentUseCase.commentRepository.GetByID(cmid)
	if err != nil {
		return
	}

	comment.DateParsed = comment.Date.Round(time.Second).String()
	return
}

func (commentUseCase *CommentUseCaseImpl) UpdateComment(uid uint, comment *models.Comment) (err error) {
	isAccessed, err := commentUseCase.userRepository.IsCommentAccessed(uid, comment.CMID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return commentUseCase.commentRepository.Update(comment)
}

func (commentUseCase *CommentUseCaseImpl) DeleteComment(uid, cmid uint) (err error) {
	isAccessed, err := commentUseCase.userRepository.IsCommentAccessed(uid, cmid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return commentUseCase.commentRepository.Delete(cmid)
}
