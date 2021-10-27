package repositories

import (
	"backendServer/app/models"
)

type SessionRepository interface {
	Create(user *models.User) (SID string, err error)
	Get(sessionValue string) (uid uint, err error)
	Delete(sessionValue string) (err error)
}
