package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/utils"
)

type UserUseCaseImpl struct {
	sessionRepository repositories.SessionRepository
	userRepository    repositories.UserRepository
}

func CreateUserUseCase(sessionRepository repositories.SessionRepository,
	userRepository repositories.UserRepository) usecases.UserUseCase {
	return &UserUseCaseImpl{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (userUseCase *UserUseCaseImpl) Create(user *models.User) (sid string, err error) {
	if !utils.ValidateUserData(user, true) {
		err = customErrors.ErrBadInputData
		return
	}

	err = userUseCase.userRepository.Create(user)
	if err != nil {
		return
	}

	sid, err = userUseCase.sessionRepository.Create(user.UID)

	return
}
