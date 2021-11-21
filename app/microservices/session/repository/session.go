package repository

type SessionRepository interface {
	Create(uid uint64) (sessionID string, err error)
	Get(sessionID string) (uid uint64, err error)
	Delete(sessionID string) (err error)
}
