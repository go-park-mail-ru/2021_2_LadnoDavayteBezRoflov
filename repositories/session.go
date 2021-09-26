package repositories

import (
	"backendServer/models"
	"errors"

	"github.com/google/uuid"
)

type SessionRepository struct {
	data *models.Data
}

func CreateSessionRepository(data *models.Data) (sessionRepository SessionRepository) {
	return SessionRepository{data: data}
}

func (sessionRepository *SessionRepository) Create(user models.User) (SID string, err error) {
	sessionRepository.data.Mu.RLock()
	curUser, ok := sessionRepository.data.Users[user.Login]
	sessionRepository.data.Mu.RUnlock()

	if !ok || curUser.Password != user.Password {
		err = errors.New("Bad input data")
		return
	}

	SID = uuid.NewString()

	sessionRepository.data.Mu.Lock()
	sessionRepository.data.Sessions[SID] = user.ID
	sessionRepository.data.Mu.Unlock()

	return
}

func (sessionRepository *SessionRepository) Get(sessionValue string) (user models.User) {
	sessionRepository.data.Mu.RLock()
	userID, ok := sessionRepository.data.Sessions[sessionValue]
	sessionRepository.data.Mu.RUnlock()

	if !ok {
		return
	}

	sessionRepository.data.Mu.RLock()
	users := sessionRepository.data.Users
	sessionRepository.data.Mu.RUnlock()

	for _, curUser := range users {
		if curUser.ID == userID {
			user = curUser
			return
		}
	}

	return
}

func (sessionRepository *SessionRepository) Delete(sessionValue string) (err error) {
	sessionRepository.data.Mu.RLock()
	_, ok := sessionRepository.data.Sessions[sessionValue]
	sessionRepository.data.Mu.RUnlock()

	if !ok {
		err = errors.New("Not authorized")
		return
	}

	sessionRepository.data.Mu.Lock()
	delete(sessionRepository.data.Sessions, sessionValue)
	sessionRepository.data.Mu.Unlock()

	return
}
