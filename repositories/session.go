package repositories

import (
	"backendServer/errors"
	"backendServer/models"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(user models.User) (SID string, err error)
	Get(sessionValue string) (user models.User)
	Delete(sessionValue string) (err error)
}

type SessionStore struct {
	data *models.Data
}

func CreateSessionRepository(data *models.Data) SessionRepository {
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

func (sessionStore *SessionStore) Get(sessionValue string) (user models.User) {
	sessionStore.data.Mu.RLock()
	userID, ok := sessionStore.data.Sessions[sessionValue]
	sessionStore.data.Mu.RUnlock()

	if !ok {
		return
	}

	sessionStore.data.Mu.RLock()
	users := sessionStore.data.Users
	sessionStore.data.Mu.RUnlock()

	for _, curUser := range users {
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
