package stores

import (
	"backendServer/app/repositories"
	"backendServer/pkg/closer"

	"github.com/gomodule/redigo/redis"

	"github.com/google/uuid"
)

type SessionStore struct {
	DB         *redis.Pool
	ExpiresSec uint64
	closer     *closer.Closer
}

func CreateSessionRepository(db *redis.Pool, expiresSec uint64, closer *closer.Closer) repositories.SessionRepository {
	return &SessionStore{DB: db, ExpiresSec: expiresSec, closer: closer}
}

func (sessionStore *SessionStore) Create(uid uint) (sid string, err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	sid = uuid.NewString()

	_, err = connection.Do("SETEX", sid, sessionStore.ExpiresSec, uid)

	return
}

func (sessionStore *SessionStore) AddTime(sid string, secs uint) (err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	_, err = connection.Do("EXPIRE", sid, secs)
	return
}

func (sessionStore *SessionStore) Get(sid string) (uid uint, err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	uid64, err := redis.Uint64(connection.Do("GET", sid))
	if err != nil {
		return
	}

	uid = uint(uid64)
	return
}

func (sessionStore *SessionStore) Delete(sid string) (err error) {
	connection := sessionStore.DB.Get()
	defer sessionStore.closer.Close(connection.Close)

	_, err = connection.Do("DEL", sid)

	return
}
