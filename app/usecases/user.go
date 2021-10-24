package usecases

import (
	"backendServer/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (string, error)
}
