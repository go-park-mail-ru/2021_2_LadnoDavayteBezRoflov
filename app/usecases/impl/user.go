package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"backendServer/pkg/utils"
	"mime/multipart"
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

	// TODO TEMP каждый пользователь в своей команде
	privateTeam := &models.Team{Title: "Личное пространство " + user.Login}

	err = userUseCase.teamRepository.Create(privateTeam)
	if err != nil {
		return
	}

	err = userUseCase.userRepository.AddUserToTeam(user.UID, privateTeam.TID)
	if err != nil {
		return
	}
	// TODO END TEMP

	sid, err = userUseCase.sessionRepository.Create(user.UID)

	return
}

func (userUseCase *UserUseCaseImpl) Get(uid uint, login string) (user *models.User, err error) {
	user, err = userUseCase.userRepository.GetByLogin(login)
	if err != nil {
		return
	}

	if user.UID != uid {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userUseCase *UserUseCaseImpl) Update(user *models.User) (err error) {
	oldUser, err := userUseCase.userRepository.GetByID(user.UID)
	if err != nil {
		return
	}

	if !hasher.IsPasswordsEqual(user.OldPassword, oldUser.HashedPassword) {
		err = customErrors.ErrBadRequest
		return
	}

	err = userUseCase.userRepository.Update(user)
	return
}

func (userUseCase *UserUseCaseImpl) UpdateAvatar(user *models.User, avatar *multipart.FileHeader) (err error) {
	_, err = userUseCase.userRepository.GetByID(user.UID)
	if err != nil {
		return
	}

	err = userUseCase.userRepository.UpdateAvatar(user, avatar)
	return
}
