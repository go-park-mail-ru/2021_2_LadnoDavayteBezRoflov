package usecases

import (
	"backendServer/app/api/models"
)

type SessionUseCase interface {
	Create(user *models.User) (sid string, err error)
	Get(sid string) (string, error)
	GetUID(sid string) (uint, error)
	Delete(sid string) error
}
