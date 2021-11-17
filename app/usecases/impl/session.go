package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"backendServer/pkg/utils"
)

type SessionUseCaseImpl struct {
	sessionRepository repositories.SessionRepository
	userRepository    repositories.UserRepository
}

func CreateSessionUseCase(sessionRepository repositories.SessionRepository,
	userRepository repositories.UserRepository,
) usecases.SessionUseCase {
	return &SessionUseCaseImpl{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (sessionUseCase *SessionUseCaseImpl) Create(user *models.User) (sid string, err error) {
	if !utils.ValidateUserData(user, false) {
		err = customErrors.ErrBadInputData
		return
	}

	existingUser, err := sessionUseCase.userRepository.GetByLogin(user.Login)
	if err != nil {
		return
	}
	if !hasher.IsPasswordsEqual(user.Password, existingUser.HashedPassword) { // TODO
		err = customErrors.ErrBadInputData
		return
	}

	sid, err = sessionUseCase.sessionRepository.Create(existingUser.UID)
	return
}

func (sessionUseCase *SessionUseCaseImpl) AddTime(sid string, secs uint) (err error) {
	return sessionUseCase.sessionRepository.AddTime(sid, secs)
}

func (sessionUseCase *SessionUseCaseImpl) Get(sid string) (string, error) {
	uid, err := sessionUseCase.sessionRepository.Get(sid)
	if err != nil {
		return "", err
	}
	user, err := sessionUseCase.userRepository.GetByID(uid)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

func (sessionUseCase *SessionUseCaseImpl) GetUID(sid string) (uint, error) {
	uid, err := sessionUseCase.sessionRepository.Get(sid)
	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (sessionUseCase *SessionUseCaseImpl) Delete(sid string) (err error) {
	err = sessionUseCase.sessionRepository.Delete(sid)
	return
}
