package usecases

import (
	"backendServer/app/models"
)

type SessionUseCase interface {
	Create(user *models.User) (sid string, err error)
	AddTime(sid string, secs uint) (err error)
	Get(sid string) (string, error)
	GetUID(sid string) (uint, error)
	Delete(sid string) error
}
