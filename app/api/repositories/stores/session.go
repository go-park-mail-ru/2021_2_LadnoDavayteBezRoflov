package stores

import (
	"backendServer/app/api/repositories"
	"backendServer/app/microservices/session/handler"
	customErrors "backendServer/pkg/errors"
	"context"
	"fmt"
)

type SessionStore struct {
	sessionManager handler.SessionCheckerClient
	ctx            context.Context
}

func CreateSessionRepository(sessionManager handler.SessionCheckerClient) repositories.SessionRepository {
	return &SessionStore{
		sessionManager: sessionManager,
		ctx:            context.Background(),
	}
}

func (sessionStore *SessionStore) Create(uid uint) (sid string, err error) {
	session, err := sessionStore.sessionManager.Create(sessionStore.ctx, &handler.SessionInfo{UID: uint64(uid)})
	fmt.Println(session)
	fmt.Println(err)
	if err != nil {
		return "", err
	} else if session == nil {
		return "", customErrors.ErrInternal
	}
	return session.ID, nil
}

func (sessionStore *SessionStore) Get(sid string) (uid uint, err error) {
	sessionInfo, err := sessionStore.sessionManager.Get(sessionStore.ctx, &handler.SessionID{ID: sid})
	if err != nil {
		return
	} else if sessionInfo == nil {
		err = customErrors.ErrInternal
		return
	}
	uid = uint(sessionInfo.UID)
	return
}

func (sessionStore *SessionStore) Delete(sid string) (err error) {
	_, err = sessionStore.sessionManager.Delete(sessionStore.ctx, &handler.SessionID{ID: sid})
	return
}
