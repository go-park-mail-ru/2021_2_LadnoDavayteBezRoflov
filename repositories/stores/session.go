package stores

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories"

	"github.com/google/uuid"
)

type SessionStore struct {
	data *models.Data
}

func CreateSessionRepository(data *models.Data) repositories.SessionRepository {
	return &SessionStore{data: data}
}

func (sessionStore *SessionStore) Create(user models.User) (SID string, err error) {
	sessionStore.data.Mu.RLock()
	curUser, ok := sessionStore.data.Users[user.Login]
	sessionStore.data.Mu.RUnlock()

	if !ok || curUser.Password != user.Password {
		err = errors.ErrBadInputData
		return
	}

	SID = uuid.NewString()

	sessionStore.data.Mu.Lock()
	sessionStore.data.Sessions[SID] = user.ID
	sessionStore.data.Mu.Unlock()

	return
}

func (sessionStore *SessionStore) Get(sessionValue string) (user models.User, err error) {
	sessionStore.data.Mu.RLock()
	defer sessionStore.data.Mu.RUnlock()

	userID, ok := sessionStore.data.Sessions[sessionValue]
	if !ok {
		err = errors.ErrBadInputData
		return
	}

	for _, curUser := range sessionStore.data.Users {
		if curUser.ID == userID {
			user = curUser
			return
		}
	}

	return
}

func (sessionStore *SessionStore) Delete(sessionValue string) (err error) {
	sessionStore.data.Mu.RLock()
	_, ok := sessionStore.data.Sessions[sessionValue]
	sessionStore.data.Mu.RUnlock()

	if !ok {
		err = errors.ErrNotAuthorized
		return
	}

	sessionStore.data.Mu.Lock()
	delete(sessionStore.data.Sessions, sessionValue)
	sessionStore.data.Mu.Unlock()

	return
}
