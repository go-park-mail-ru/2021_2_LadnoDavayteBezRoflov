package impl

import (
	"backendServer/app/microservices/session/repository"
	"backendServer/app/microservices/session/usecase"
)

type SessionUseCaseImpl struct {
	sessionRepository repository.SessionRepository
}

func CreateSessionUseCase(sessionRepository repository.SessionRepository) usecase.SessionUseCase {
	return &SessionUseCaseImpl{sessionRepository: sessionRepository}
}

func (sessionUseCase *SessionUseCaseImpl) Create(uid uint64) (sessionID string, err error) {
	return sessionUseCase.sessionRepository.Create(uid)
}

func (sessionUseCase *SessionUseCaseImpl) Get(sessionID string) (uid uint64, err error) {
	return sessionUseCase.sessionRepository.Get(sessionID)
}

func (sessionUseCase *SessionUseCaseImpl) Delete(sessionID string) (err error) {
	return sessionUseCase.sessionRepository.Delete(sessionID)
}
