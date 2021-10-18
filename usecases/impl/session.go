package impl

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories"
	"backendServer/usecases"
	"backendServer/utils"
	"strconv"
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

func (sessionUseCase *SessionUseCaseImpl) Create(user *models.User) (string, error) {
	if !utils.ValidateUserData(user, false) {
		return "", errors.ErrBadInputData
	}

	SID, err := sessionUseCase.sessionRepository.Create(user)
	if err != nil {
		return "", err
	}

	return SID, nil
}

func (sessionUseCase *SessionUseCaseImpl) Get(sid string) (string, error) {
	uid, err := sessionUseCase.sessionRepository.Get(sid)
	if err != nil {
		return "", err
	}
	user, err := sessionUseCase.userRepository.GetById(uid)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

func (sessionUseCase *SessionUseCaseImpl) GetUID(sid string) (string, error) {
	uid, err := sessionUseCase.sessionRepository.Get(sid)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(uid), 10), nil
}

func (sessionUseCase *SessionUseCaseImpl) Delete(sid string) (err error) {
	err = sessionUseCase.sessionRepository.Delete(sid)
	return
}
