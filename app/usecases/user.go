package usecases

import (
	"backendServer/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (sid string, err error)
}
