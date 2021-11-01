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
	teamRepository    repositories.TeamRepository
}

func CreateUserUseCase(
	sessionRepository repositories.SessionRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
) usecases.UserUseCase {
	return &UserUseCaseImpl{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		teamRepository:    teamRepository,
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

	privateTeam := &models.Team{Title: "Личное пространство " + user.Login}

	err = userUseCase.teamRepository.Create(privateTeam)
	if err != nil {
		return
	}

	err = userUseCase.userRepository.AddUserToTeam(user.UID, privateTeam.TID) // TODO TEMP все пользователи в одной команде
	if err != nil {
		return
	}

	sid, err = userUseCase.sessionRepository.Create(user.UID)

	return
}
