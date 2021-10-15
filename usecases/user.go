package usecases

import "backendServer/models"

type UserUseCase interface {
	Create(user *models.User) (string, error)
}
