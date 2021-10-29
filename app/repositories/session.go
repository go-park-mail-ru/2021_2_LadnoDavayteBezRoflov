package repositories

type SessionRepository interface {
	Create(uid uint) (sid string, err error)
	AddTime(sid string, secs uint) (err error)
	Get(sid string) (uid uint, err error)
	Delete(sid string) (err error)
}
