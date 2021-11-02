package usecases

import (
	"backendServer/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (sid string, err error)
	Get(uid uint, login string) (user *models.User, err error)
	Update(user *models.User) (err error)
	UpdateAvatar(user *models.User) (err error)
}
