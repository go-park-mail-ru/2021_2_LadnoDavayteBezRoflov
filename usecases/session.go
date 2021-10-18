package usecases

import "backendServer/models"

type SessionUseCase interface {
	Create(user *models.User) (string, error)
	Get(sid string) (string, error)
	GetUID(sid string) (string, error)
	Delete(sid string) error
}
