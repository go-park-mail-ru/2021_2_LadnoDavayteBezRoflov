package repositories

import "backendServer/models"

type SessionRepository interface {
	Create(user models.User) (SID string, err error)
	Get(sessionValue string) (user models.User, err error)
	Delete(sessionValue string) (err error)
}
