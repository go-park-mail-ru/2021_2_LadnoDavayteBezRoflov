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

func (userUseCase *UserUseCaseImpl) Create(user *models.User) (string, error) {
	if !utils.ValidateUserData(user, true) {
		return "", errors.ErrBadInputData
	}

	addedUser, err := userUseCase.userRepository.Create(user)
	if err != nil {
		return "", err
	}

	SID, err := userUseCase.sessionRepository.Create(addedUser)
	if err != nil {
		return "", err
	}
	return SID, nil
}
