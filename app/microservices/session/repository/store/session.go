package store

import (
	"backendServer/app/microservices/session/repository"
	"backendServer/pkg/closer"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type SessionStore struct {
	DB         *redis.Pool
	ExpiresSec uint64
	closer     *closer.Closer
}

func CreateSessionRepository(db *redis.Pool, expiresSec uint64, closer *closer.Closer) repository.SessionRepository {
	return &SessionStore{DB: db, ExpiresSec: expiresSec, closer: closer}
}

func (sessionStore *SessionStore) Create(uid uint64) (sessionID string, err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	sessionID = uuid.NewString()

	_, err = connection.Do("SETEX", sessionID, sessionStore.ExpiresSec, uid)
	return
}

func (sessionStore *SessionStore) Get(sessionID string) (uid uint64, err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	return redis.Uint64(connection.Do("GET", sessionID))
}

func (sessionStore *SessionStore) Delete(sessionID string) (err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	_, err = connection.Do("DEL", sessionID)
	return
}
